package file

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/types"
	"github.com/xitongsys/parquet-go/writer"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ApiSourceProperties struct {
	Path        string         `mapstructure:"path"`
	Format      string         `mapstructure:"format"`
	LensName    string         `mapstructure:"lensName"`
	Uri         string         `mapstructure:"uri"`
	Method      string         `mapstructure:"method"`
	Body        Body           `mapstructure:"body"`
	Headers     map[string]any `mapstructure:"headers"`
	QueryParams map[string]any `mapstructure:"queryParams"`
}

type Body struct {
	Dimensions []string `mapstructure:"dimensions"`
	Limit      int
	Offset     int
}

type Query struct {
	Dimensions     []string `json:"dimensions"`
	Ungrouped      bool     `json:"ungrouped"`
	Limit          int      `json:"limit"`
	Offset         int      `json:"offset"`
	ResponseFormat string   `json:"responseFormat"`
}

type Data struct {
	Members []string        `json:"members"`
	Dataset [][]interface{} `json:"dataset"`
}

type Annotation struct {
	Dimensions map[string]Dimension `json:"dimensions"`
}

type Dimension struct {
	Title       string `json:"title"`
	ShortTitle  string `json:"shortTitle"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type ApiResponse struct {
	Data       Data       `json:"data"`
	Annotation Annotation `json:"annotation"`
	Error      string     `json:"error"`
}

func storeApiData(props map[string]any) error {
	conf, err := parseSourceProps(props)
	if err != nil {
		return err
	}

	// create Query Object
	query := Query{
		Dimensions:     conf.Body.Dimensions,
		Ungrouped:      false,
		Limit:          conf.Body.Limit,
		Offset:         conf.Body.Offset,
		ResponseFormat: "compact",
	}

	createDirectory(filepath.Dir(conf.Path))
	err = writeDataToParquet(conf, conf.Body.Limit, conf.Body.Offset, query)
	if err != nil {
		return err
	}
	return nil
}

// Parse Source Properties to struct ApiSourceProperties
func parseSourceProps(props map[string]any) (*ApiSourceProperties, error) {
	// Todo: put proper validation for props
	conf := &ApiSourceProperties{}
	err := mapstructure.Decode(props, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func getDataFromLens2(conf *ApiSourceProperties, query Query) (*ApiResponse, error) {
	apiConf := ApiSourceProperties{
		Uri:     fmt.Sprintf("%s/%s/v2/load", strings.Trim(conf.Uri, "/"), conf.LensName),
		Headers: conf.Headers,
		Method:  "POST",
	}

	// fetch data from lens api
	body, err := httpApiCall(apiConf, map[string]any{"query": query})
	if err != nil {
		return nil, err
	}

	var response ApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		zap.L().Error(fmt.Sprintf("Error unmarshalling JSON: %s", err.Error()))
		return nil, err
	}

	if len(response.Error) > 0 {
		if response.Error == "Continue wait" {
			data, err := getDataFromLens2(conf, query)
			if err != nil {
				return nil, err
			}
			return data, nil
		} else {
			return nil, errors.New(response.Error)
		}
	}

	return &response, nil
}

// Http Request Method return JSON data
func httpApiCall(conf ApiSourceProperties, payload map[string]any) ([]byte, error) {
	// prepare body for API call
	var err error = nil
	var bodyBytes []byte = nil
	if payload != nil {
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Error marshaling request body: %s", err.Error()))
			return nil, err
		}
	}

	// Add query parameters to the URI
	uri, err := url.Parse(conf.Uri)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error parsing URI: %s", err.Error()))
		return nil, err
	}
	query := uri.Query()
	for key, value := range conf.QueryParams {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	uri.RawQuery = query.Encode()

	// Create a new HTTP Request
	req, err := http.NewRequest(conf.Method, uri.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	// Set the request headers
	for key, value := range conf.Headers {
		req.Header.Set(key, value.(string))
	}

	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error making HTTP request: %s", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error reading response body: %s", err.Error()))
		return nil, err
	}

	return body, nil
}

func getParquetDatatype(dtype string) string {
	datatype := map[string]string{
		"number": "FLOAT",
		"time":   "TIMESTAMP_MICROS",
		"string": "BYTE_ARRAY",
		"bool":   "BOOLEAN",
	}
	return datatype[dtype]
}

// Function to generate Parquet schema based on field details
func getSchema(response *ApiResponse) (string, error) {
	var jsonSchemaStrings []string = nil
	for _, member := range response.Data.Members {
		// create field schema element
		pair := ""
		switch dtype := getParquetDatatype(response.Annotation.Dimensions[member].Type); dtype {
		case "BYTE_ARRAY":
			pair = fmt.Sprintf("name=%s, type=%s, convertedtype=UTF8", member, dtype)
		case "TIMESTAMP_MICROS":
			pair = fmt.Sprintf("name=%s, type=INT64, convertedtype=%s", member, dtype)
		default:
			pair = fmt.Sprintf("name=%s, type=%s", member, dtype)
		}

		data := map[string]interface{}{
			"Tag": pair,
		}
		// Marshal the map into JSON
		jsonString, err := json.Marshal(data)
		if err != nil {
			zap.L().Error(fmt.Sprintf("Error marshaling JSON: %s", err.Error()))
			return "", err
		}
		// collect all schema elements
		jsonSchemaStrings = append(jsonSchemaStrings, string(jsonString))
	}
	// create json schema and place jsonSchemaString of all schema elements
	var jsonSchema = fmt.Sprintf("{\"Tag\": \"name=parquet_go_root, repetitiontype=REQUIRED\",\"Fields\": [%s]}", strings.Join(jsonSchemaStrings, ","))
	return jsonSchema, nil
}

func writeDataToParquet(conf *ApiSourceProperties, limit, offset int, query Query) error {
	var fileName string
	fileName = fmt.Sprintf("%s/data_%d_%d.parquet", filepath.Dir(conf.Path), conf.Body.Offset, offset)
	response, err := getDataFromLens2(conf, query)
	if err != nil {
		return err
	}

	if len(response.Data.Dataset) == 0 {
		return nil
	}

	schema, err := getSchema(response)
	if err != nil {
		return err
	}

	// write data to parquet file
	fw, err := local.NewLocalFileWriter(fileName)
	if err != nil {
		return err
	}
	defer fw.Close()
	pw, err := writer.NewJSONWriter(schema, fw, 4)
	if err != nil {
		return err
	}
	defer pw.WriteStop()

	// Write data to the Parquet file
	for _, row := range response.Data.Dataset {
		rec := `{%s}`
		var recordFields []string = nil
		for index, memberValue := range row {
			field := response.Data.Members[index]
			recordField := ""
			switch valueType := getParquetDatatype(response.Annotation.Dimensions[field].Type); valueType {
			case "BYTE_ARRAY":
				recordField = fmt.Sprintf("\"%s\":\"%s\"", field, memberValue.(string))
			case "TIMESTAMP_MICROS":
				// Parse the string as a time.Time object
				layout := "2006-01-02T15:04:05.000"
				timestamp, err := time.Parse(layout, memberValue.(string))
				if err != nil {
					return err
				}
				// Convert the time to a Unix timestamp (int64)
				timestampMicros := types.TimeToTIMESTAMP_MICROS(timestamp, false)
				recordField = fmt.Sprintf("\"%s\":\"%d\"", field, timestampMicros)
			case "BOOLEAN":
				recordField = fmt.Sprintf("\"%s\":\"%t\"", field, memberValue.(bool))
			case "FLOAT":
				recordField = fmt.Sprintf("\"%s\":\"%f\"", field, memberValue.(float64))
			default:
				recordField = fmt.Sprintf("\"%s\":\"%s\"", field, fmt.Sprintf("%s", memberValue))
			}
			//fmt.Println(recordField)
			recordFields = append(recordFields, recordField)
		}
		recordString := strings.Join(recordFields, ",")
		rec = fmt.Sprintf(rec, recordString)
		if err = pw.Write(rec); err != nil {
			zap.L().Error(fmt.Sprintf("Write error %s", err.Error()))
		}
	}

	zap.L().Info("Write Finished")
	offset = offset + limit
	query = Query{
		Dimensions:     query.Dimensions,
		Ungrouped:      query.Ungrouped,
		Limit:          limit,
		Offset:         offset,
		ResponseFormat: "compact",
	}
	err = writeDataToParquet(conf, limit, offset, query)
	if err != nil {
		return err
	}
	// Ensure all buffered data is flushed to disk
	if err := pw.WriteStop(); err != nil {
		fmt.Println("Error writing Parquet file:", err)
		return err
	}

	// Close the file writer
	if err := fw.Close(); err != nil {
		fmt.Println("Error closing Parquet file writer:", err)
		return err
	}
	return nil
}

func createDirectory(dirPath string) {
	// Check if the directory exists
	if _, err := os.Stat(dirPath); err == nil {
		// Directory exists, remove it
		if err := os.RemoveAll(dirPath); err != nil {
			zap.L().Error(fmt.Sprintf("Error removing directory: %s", err.Error()))
			return
		}
	}

	// Create the directory
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		zap.L().Error(fmt.Sprintf("Error creating directory: %s", err))
		return
	}

}
