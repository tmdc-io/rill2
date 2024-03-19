package worker

import (
	"context"
	"time"

	"github.com/rilldata/rill/admin/database"
	"github.com/rilldata/rill/admin/metrics"
	"go.uber.org/zap"
)

func (w *Worker) runAutoscaler(ctx context.Context) error {
	recs, ok, err := w.allRecommendations(ctx)
	if err != nil {
		return err
	}
	if !ok {
		w.logger.Debug("skipping autoscaler: no metrics project configured")
		return nil
	}

	for _, rec := range recs {
		duration := time.Since(rec.UpdatedOn)
		if duration < 24*time.Hour {
			w.logger.Debug("skipping autoscaler: the project has been scaled recently", zap.String("project_id", rec.ProjectID), zap.Time("project_updated_on", rec.UpdatedOn))
			break
		}

		opts := &database.UpdateProjectOptions{
			ProdSlots: rec.RecommendedSlots,
		}

		proj, err := w.admin.DB.UpdateProject(ctx, rec.ProjectID, opts)
		if err != nil {
			w.logger.Warn("failed to autoscale:", zap.String("project_id", rec.ProjectID), zap.Error(err))
			return err
		}

		w.logger.Debug("succeeded in autoscaling:", zap.String("project_id", proj.Name), zap.Int("project_slots", proj.ProdSlots))
	}

	return nil
}

func (w *Worker) allRecommendations(ctx context.Context) ([]metrics.AutoscalerSlotsRecommendation, bool, error) {
	client, ok, err := w.admin.OpenMetricsProject(ctx)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}

	var recs []metrics.AutoscalerSlotsRecommendation
	limit := 1000
	offset := 0
	for {
		batch, err := client.AutoscalerSlotsRecommendations(ctx, limit, offset)
		if err != nil {
			return nil, false, err
		}
		if len(batch) == 0 {
			break
		}
		recs = append(recs, batch...)
		offset += limit
	}

	return recs, true, nil
}
