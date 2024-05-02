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
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ApiSourceProperties struct {
	Path   string `mapstructure:"path"`
	Format string `mapstructure:"format"`
	Lens   Lens   `mapstructure:"lens"`
}

type Lens struct {
	BaseUri string `mapstructure:"baseUri"`
	Name    string `mapstructure:"name"`
	Body    Body   `mapstructure:"query"`
	Apikey  string `mapstructure:"apikey"`
}

type Body struct {
	Dimensions []string `mapstructure:"dimensions"`
	Batch      int      `mapstructure:"batch"`
	Start      int      `mapstructure:"start"`
	End        int      `mapstructure:"end"`
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
	Data       Data        `json:"data"`
	Annotation Annotation  `json:"annotation"`
	Error      interface{} `json:"error"`
}

func fetchLensData(props map[string]any) (string, error) {
	conf, err := parseSourceProps(props)
	if err != nil {
		return "", err
	}

	// create Query Object
	query := Query{
		Dimensions:     conf.Lens.Body.Dimensions,
		Ungrouped:      true,
		Limit:          conf.Lens.Body.Batch,
		Offset:         conf.Lens.Body.Start,
		ResponseFormat: "compact",
	}

	dirPath := createDirectory(filepath.Dir(conf.Path))
	if len(dirPath) == 0 {
		return "", errors.New("directory creation failed")
	}
	offset := conf.Lens.Body.Start
	err = writeDataToParquet(conf, conf.Lens.Body.Batch, offset, query, dirPath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", dirPath, filepath.Base(conf.Path)), nil
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
	apiConf := Lens{
		BaseUri: fmt.Sprintf("%s/%s/v2/load", strings.Trim(conf.Lens.BaseUri, "/"), conf.Lens.Name),
		Apikey:  fmt.Sprintf("Bearer %s", conf.Lens.Apikey),
	}

	// fetch data from lens api
	body, err := httpApiCall(apiConf, map[string]any{"query": query})
	if err != nil {
		return nil, err
	}

	var response ApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err.Error(), response.Error)
		//log.Printf("error unmarshalling JSON: %s\n", err.Error())
		return nil, err
	}

	if response.Error != nil {
		switch response.Error.(type) {
		case string:
			if response.Error.(string) == "Continue wait" {
				data, err := getDataFromLens2(conf, query)
				if err != nil {
					return nil, err
				}
				return data, nil
			} else {
				return nil, errors.New(response.Error.(string))
			}
		case map[string]interface{}:
			if response.Error.(map[string]interface{}) != nil {
				return nil, errors.New(fmt.Sprintf("%v", response.Error))
			} else {
				return nil, errors.New("no response error fetched")
			}
		default:
			return nil, errors.New(fmt.Sprintf("Invalid error response: %v", response.Error))
		}
	}

	return &response, nil
}

// Http Request Method return JSON data
func httpApiCall(conf Lens, payload map[string]any) ([]byte, error) {
	// prepare body for API call
	var err error = nil
	var bodyBytes []byte = nil
	if payload != nil {
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			log.Printf(fmt.Sprintf("Error marshaling request body: %s", err.Error()))
			return nil, err
		}
	}

	// Add query parameters to the URI
	uri, err := url.Parse(conf.BaseUri)
	if err != nil {
		log.Printf(fmt.Sprintf("Error parsing URI: %s", err.Error()))
		return nil, err
	}

	// Create a new HTTP Request
	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", conf.Apikey)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf(fmt.Sprintf("Error making HTTP request: %s", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf(fmt.Sprintf("Error reading response body: %s", err.Error()))
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
			pair = fmt.Sprintf("name=%s, type=%s, convertedtype=UTF8", strings.Split(member, ".")[1], dtype)
		case "TIMESTAMP_MICROS":
			pair = fmt.Sprintf("name=%s, type=INT64, convertedtype=%s", strings.Split(member, ".")[1], dtype)
		default:
			pair = fmt.Sprintf("name=%s, type=%s", strings.Split(member, ".")[1], dtype)
		}

		data := map[string]interface{}{
			"Tag": pair,
		}
		// Marshal the map into JSON
		jsonString, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling JSON: %s", err.Error())
			return "", err
		}
		// collect all schema elements
		jsonSchemaStrings = append(jsonSchemaStrings, string(jsonString))
	}
	// create json schema and place jsonSchemaString of all schema elements
	var jsonSchema = fmt.Sprintf("{\"Tag\": \"name=parquet_go_root, repetitiontype=REQUIRED\",\"Fields\": [%s]}", strings.Join(jsonSchemaStrings, ","))
	return jsonSchema, nil
}

func writeDataToParquet(conf *ApiSourceProperties, limit, offset int, query Query, dirPath string) error {
	if offset != 0 && offset == limit*(conf.Lens.Body.End+1) {
		return nil
	}
	var fileName string
	fileName = fmt.Sprintf("%s/data_%d_%d.parquet", dirPath, conf.Lens.Body.Start, offset)

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
				recordField = fmt.Sprintf("\"%s\":\"%s\"", strings.Split(field, ".")[1], memberValue.(string))
			case "TIMESTAMP_MICROS":
				// Parse the string as a time.Time object
				layout := "2006-01-02T15:04:05.000"
				timestamp, err := time.Parse(layout, memberValue.(string))
				if err != nil {
					return err
				}
				// Convert the time to a Unix timestamp (int64)
				timestampMicros := types.TimeToTIMESTAMP_MICROS(timestamp, false)
				recordField = fmt.Sprintf("\"%s\":\"%d\"", strings.Split(field, ".")[1], timestampMicros)
			case "BOOLEAN":
				recordField = fmt.Sprintf("\"%s\":\"%t\"", strings.Split(field, ".")[1], memberValue.(bool))
			case "FLOAT":
				recordField = fmt.Sprintf("\"%s\":\"%f\"", strings.Split(field, ".")[1], memberValue.(float64))
			default:
				recordField = fmt.Sprintf("\"%s\":\"%s\"", strings.Split(field, ".")[1], fmt.Sprintf("%s", memberValue))
			}
			//fmt.Println(recordField)
			recordFields = append(recordFields, recordField)
		}
		recordString := strings.Join(recordFields, ",")
		rec = fmt.Sprintf(rec, recordString)
		if err = pw.Write(rec); err != nil {
			log.Printf("Write error %s", err.Error())
		}
	}

	log.Printf("Write Finished")
	offset = offset + limit
	query = Query{
		Dimensions:     query.Dimensions,
		Ungrouped:      query.Ungrouped,
		Limit:          limit,
		Offset:         offset,
		ResponseFormat: "compact",
	}
	err = writeDataToParquet(conf, limit, offset, query, dirPath)
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

func createDirectory(dirPath string) string {
	// Create the directory
	currentTime := time.Now()
	timestamp := currentTime.Format("2006-01-02-15-04")

	timestampDirPath := fmt.Sprintf("%s/%s", dirPath, timestamp)
	if err := os.MkdirAll(timestampDirPath, 0755); err != nil {
		log.Printf("Error creating directory: %s", err)
		return ""
	}
	return timestampDirPath
}
