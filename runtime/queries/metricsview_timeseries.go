package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	runtimev1 "github.com/rilldata/rill/proto/gen/rill/runtime/v1"
	"github.com/rilldata/rill/runtime"
	"github.com/rilldata/rill/runtime/drivers"
	"github.com/rilldata/rill/runtime/pkg/pbutil"
	"github.com/rilldata/rill/runtime/queries/forecast/holtwinters_v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetricsViewTimeSeries struct {
	MetricsViewName string                       `json:"metrics_view_name,omitempty"`
	MeasureNames    []string                     `json:"measure_names,omitempty"`
	TimeStart       *timestamppb.Timestamp       `json:"time_start,omitempty"`
	TimeEnd         *timestamppb.Timestamp       `json:"time_end,omitempty"`
	Limit           int64                        `json:"limit,omitempty"`
	Offset          int64                        `json:"offset,omitempty"`
	Sort            []*runtimev1.MetricsViewSort `json:"sort,omitempty"`
	Filter          *runtimev1.MetricsViewFilter `json:"filter,omitempty"`
	TimeGranularity runtimev1.TimeGrain          `json:"time_granularity,omitempty"`
	ForecastPeriod  int                          `json:"offset,omitempty"`

	Result *runtimev1.MetricsViewTimeSeriesResponse `json:"-"`
}

var _ runtime.Query = &MetricsViewTimeSeries{}

func (q *MetricsViewTimeSeries) Key() string {
	r, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("MetricsViewTimeSeries:%s", string(r))
}

func (q *MetricsViewTimeSeries) Deps() []string {
	return []string{q.MetricsViewName}
}

func (q *MetricsViewTimeSeries) MarshalResult() any {
	return q.Result
}

func (q *MetricsViewTimeSeries) UnmarshalResult(v any) error {
	res, ok := v.(*runtimev1.MetricsViewTimeSeriesResponse)
	if !ok {
		return fmt.Errorf("MetricsViewTimeSeries: mismatched unmarshal input")
	}
	q.Result = res
	return nil
}

func (q *MetricsViewTimeSeries) Resolve(ctx context.Context, rt *runtime.Runtime, instanceID string, priority int) error {
	olap, err := rt.OLAP(ctx, instanceID)
	if err != nil {
		return err
	}

	if olap.Dialect() != drivers.DialectDuckDB {
		return fmt.Errorf("not available for dialect '%s'", olap.Dialect())
	}

	mv, err := lookupMetricsView(ctx, rt, instanceID, q.MetricsViewName)
	if err != nil {
		return err
	}

	if mv.TimeDimension == "" {
		return fmt.Errorf("metrics view '%s' does not have a time dimension", q.MetricsViewName)
	}

	measures, err := toMeasures(mv.Measures, q.MeasureNames)
	if err != nil {
		return err
	}

	tsq := &ColumnTimeseries{
		TableName:           mv.Model,
		TimestampColumnName: mv.TimeDimension,
		TimeRange: &runtimev1.TimeSeriesTimeRange{
			Start:    q.TimeStart,
			End:      q.TimeEnd,
			Interval: q.TimeGranularity,
		},
		Measures: measures,
		Filters:  q.Filter,
	}
	err = rt.Query(ctx, instanceID, tsq, priority)
	if err != nil {
		return err
	}

	r := tsq.Result

	fResults := getForecasted(&q.TimeGranularity, r.Results, q.ForecastPeriod)
	q.Result = &runtimev1.MetricsViewTimeSeriesResponse{
		Meta:         r.Meta,
		Data:         r.Results,
		ForecastData: fResults,
	}

	return nil
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func toTimeGrainNs(specifier runtimev1.TimeGrain, ts time.Time) int64 {
	ts.Month()

	switch specifier {
	case runtimev1.TimeGrain_TIME_GRAIN_MILLISECOND:
		return time.Millisecond.Nanoseconds()
	case runtimev1.TimeGrain_TIME_GRAIN_SECOND:
		return time.Second.Nanoseconds()
	case runtimev1.TimeGrain_TIME_GRAIN_MINUTE:
		return time.Minute.Nanoseconds()
	case runtimev1.TimeGrain_TIME_GRAIN_HOUR:
		return time.Hour.Nanoseconds()
	case runtimev1.TimeGrain_TIME_GRAIN_DAY:
		return 24 * time.Hour.Nanoseconds()
	case runtimev1.TimeGrain_TIME_GRAIN_WEEK:
		return 24 * 7 * time.Hour.Nanoseconds()
	case runtimev1.TimeGrain_TIME_GRAIN_MONTH:
		return int64(daysIn(ts.Month(), ts.Year())) * 24 * time.Hour.Nanoseconds()
	case runtimev1.TimeGrain_TIME_GRAIN_YEAR:
		return int64(daysIn(ts.Month(), ts.Year())) * 24 * time.Hour.Nanoseconds()
	}
	panic(fmt.Errorf("unconvertable time grain specifier: %v", specifier))
}

func getForecasted(t *runtimev1.TimeGrain, results []*runtimev1.TimeSeriesValue, timePeriod int) []*runtimev1.TimeSeriesValue {
	nForecastedValues := timePeriod
	originalTsValues := make(map[string][]float64)
	for _, r := range results {
		for k, v := range r.Records.Fields {
			if originalTsValues[k] == nil {
				originalTsValues[k] = make([]float64, 0)
			}
			originalTsValues[k] = append(originalTsValues[k], v.GetNumberValue())
		}
	}
	forecastedTsValues := make(map[string][]float64)
	for k, v := range originalTsValues {
		forecasted, err := holtwinters_v2.PredictAdditive(v, 2, 0.5, 0.4, 0.6, nForecastedValues)
		if err != nil {
			log.Fatal(err)
		}
		forecastedTsValues[k] = forecasted
	}
	var forecastedResult []*runtimev1.TimeSeriesValue

	result := results[len(results)-1]
	ts := result.Ts
	for i := 0; i < nForecastedValues; i++ {
		duration := toTimeGrainNs(*t, ts.AsTime())
		ts = timestamppb.New(ts.AsTime().Add(time.Duration(duration)))
		fields := make(map[string]any)
		for k, v := range forecastedTsValues {
			fields[k] = v[i+len(originalTsValues[k])]
		}
		toStruct, _ := pbutil.ToStruct(fields)
		forecastedResult = append(forecastedResult, &runtimev1.TimeSeriesValue{
			Ts:      ts,
			Bin:     result.Bin,
			Records: toStruct,
		})
	}
	return forecastedResult
}

func toMeasures(measures []*runtimev1.MetricsView_Measure, measureNames []string) ([]*runtimev1.ColumnTimeSeriesRequest_BasicMeasure, error) {
	var res []*runtimev1.ColumnTimeSeriesRequest_BasicMeasure
	for _, n := range measureNames {
		found := false
		for _, m := range measures {
			if m.Name == n {
				res = append(res, &runtimev1.ColumnTimeSeriesRequest_BasicMeasure{
					SqlName:    m.Name,
					Expression: m.Expression,
				})
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("measure does not exist: '%s'", n)
		}
	}
	return res, nil
}
