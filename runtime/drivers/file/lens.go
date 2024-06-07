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

func (c *connection) fetchLensData(props map[string]any) (string, error) {
	conf, err := c.parseSourceProps(props)
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

	dirPath, err := c.createDirectory(filepath.Dir(conf.Path))
	if err != nil {
		return "", err
	}

	offset := conf.Lens.Body.Start

	c.logger.Info("starting writing data...")
	for {
		c.logger.Debug(fmt.Sprintf("query struct for Parquet Data Write: %v", query))
		complete, err := c.writeDataToParquet(conf, conf.Lens.Body.Batch, offset, query, dirPath)
		if err != nil {
			return "", err
		}
		if complete {
			c.logger.Info("Writing data to Parquet files finished...")
			break
		}
		offset = offset + conf.Lens.Body.Batch
		query = Query{
			Dimensions:     query.Dimensions,
			Ungrouped:      query.Ungrouped,
			Limit:          conf.Lens.Body.Batch,
			Offset:         offset,
			ResponseFormat: "compact",
		}
	}
	return fmt.Sprintf("%s/%s", dirPath, filepath.Base(conf.Path)), nil
}

// Parse Source Properties to struct ApiSourceProperties
func (c *connection) parseSourceProps(props map[string]any) (*ApiSourceProperties, error) {
	conf := &ApiSourceProperties{}
	err := mapstructure.Decode(props, conf)
	if err != nil {
		c.logger.Error(fmt.Sprintf("conversion of props to conf structure: %s", err.Error()))
		return nil, err
	}
	c.logger.Debug(fmt.Sprintf("conf: %v", *conf))
	return conf, nil
}

func (c *connection) getDataFromLens2(conf *ApiSourceProperties, query Query) (*ApiResponse, error) {
	envApiKey := os.Getenv("DATAOS_RUN_AS_APIKEY")
	if len(envApiKey) == 0 {
		c.logger.Error("no apikey provide, `DATAOS_RUN_AS_APIKEY` missing")
		return nil, errors.New("no apikey given, please provide `DATAOS_RUN_AS_APIKEY` as env variable")
	}
	apiConf := Lens{
		BaseUri: fmt.Sprintf("%s/%s/v2/load", strings.Trim(conf.Lens.BaseUri, "/"), conf.Lens.Name),
		Apikey:  fmt.Sprintf("Bearer %s", envApiKey),
	}

	// fetch data from lens api
	body, err := c.httpApiCall(apiConf, map[string]any{"query": query})
	if err != nil {
		return nil, err
	}

	var response ApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		c.logger.Error(fmt.Sprintf("Json Unmarshalling failed for ResponseBody to ApiResponse: %s", err.Error()))
		c.logger.Debug(fmt.Sprintf("Response from LENS API: %s", string(body)))
		return nil, err
	}

	if response.Error != nil {
		switch response.Error.(type) {
		case string:
			if response.Error.(string) == "Continue wait" {
				c.logger.Debug("Response returned Continue wait...")
				data, err := c.getDataFromLens2(conf, query)
				if err != nil {
					return nil, err
				}
				return data, nil
			} else {
				c.logger.Error(fmt.Sprintf("Response Body contains Error: %s", response.Error.(string)))
				return nil, errors.New(response.Error.(string))
			}
		case map[string]interface{}:
			if response.Error.(map[string]interface{}) != nil {
				c.logger.Error(fmt.Sprintf("Response Body contains Error: %v", response.Error))
				return nil, errors.New(fmt.Sprintf("%v", response.Error))
			} else {
				c.logger.Error("Response Body return Error with no message")
				return nil, errors.New("no response error fetched")
			}
		default:
			c.logger.Error(fmt.Sprintf("Invalid error response from API: %v", response.Error))
			return nil, errors.New(fmt.Sprintf("Invalid error response: %v", response.Error))
		}
	}
	return &response, nil
}

// Http Request Method return JSON data
func (c *connection) httpApiCall(conf Lens, payload map[string]any) ([]byte, error) {
	// prepare body for API call
	var err error = nil
	var bodyBytes []byte = nil

	c.logger.Debug(fmt.Sprintf("Payload for API Call: %v", payload))
	if payload != nil {
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Error marshaling request body: %s", err.Error()))
			return nil, err
		}
	}

	uri, err := url.Parse(conf.BaseUri)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error parsing URI: %s", err.Error()))
		return nil, err
	}
	c.logger.Debug(fmt.Sprintf("URI for API Call: %v", uri))

	// Create a new HTTP Request
	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed creating a new http request: %s", err.Error()))
		return nil, err
	}
	c.logger.Debug(fmt.Sprintf("Http Request without Headers: %v", req))

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", conf.Apikey)

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error making HTTP request: %s", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error reading response body: %s", err.Error()))
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
func (c *connection) getSchema(response *ApiResponse) (string, error) {
	var jsonSchemaStrings []string = nil
	for _, member := range response.Data.Members {
		// create field schema element
		pair := ""
		switch dtype := getParquetDatatype(response.Annotation.Dimensions[member].Type); dtype {
		case "BYTE_ARRAY":
			pair = fmt.Sprintf("name=%s, type=%s, convertedtype=UTF8, repetitiontype=OPTIONAL", strings.Split(member, ".")[1], dtype)
		case "TIMESTAMP_MICROS":
			pair = fmt.Sprintf("name=%s, type=INT64, convertedtype=%s, repetitiontype=OPTIONAL", strings.Split(member, ".")[1], dtype)
		default:
			pair = fmt.Sprintf("name=%s, type=%s, repetitiontype=OPTIONAL", strings.Split(member, ".")[1], dtype)
		}

		data := map[string]interface{}{
			"Tag": pair,
		}
		// Marshal the map into JSON
		jsonString, err := json.Marshal(data)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Error marshaling JSON of parquet schema data: %s", err.Error()))
			return "", err
		}
		// collect all schema elements
		jsonSchemaStrings = append(jsonSchemaStrings, string(jsonString))
	}
	// create json schema and place jsonSchemaString of all schema elements
	var jsonSchema = fmt.Sprintf("{\"Tag\": \"name=parquet_go_root, repetitiontype=REQUIRED\",\"Fields\": [%s]}", strings.Join(jsonSchemaStrings, ","))
	return jsonSchema, nil
}

func (c *connection) writeDataToParquet(conf *ApiSourceProperties, limit, offset int, query Query, dirPath string) (bool, error) {
	if offset != 0 && offset == limit*(conf.Lens.Body.End+1)+conf.Lens.Body.Start && conf.Lens.Body.End != -1 {
		return true, nil
	}
	var fileName string
	fileName = fmt.Sprintf("%s/data_%d_%d.parquet", dirPath, conf.Lens.Body.Start, offset)
	c.logger.Debug(fmt.Sprintf("File Name: %s", fileName))

	count, response, err := c.writeIfDataExists(conf, query)
	if err != nil {
		return false, err
	}

	if count > 0 {
		schema, err := c.getSchema(response)
		if err != nil {
			return false, err
		}
		c.logger.Debug(fmt.Sprintf("Parquet Schema Created: %s", schema))

		// write data to parquet file
		fw, err := local.NewLocalFileWriter(fileName)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Parquet Local File Writer failed: %s", err.Error()))
			return false, err
		}
		defer fw.Close()

		pw, err := writer.NewJSONWriter(schema, fw, 10)
		if err != nil {
			c.logger.Error(fmt.Sprintf("Parquet Json Writer failed: %s", err.Error()))
			return false, err
		}
		defer pw.WriteStop()

		// Write data to the Parquet file
		for _, row := range response.Data.Dataset {
			rec := `{%s}`
			var recordFields []string = nil
			for index, memberValue := range row {
				field := response.Data.Members[index]
				recordField := ""
				if memberValue != nil {
					switch valueType := getParquetDatatype(response.Annotation.Dimensions[field].Type); valueType {
					case "BYTE_ARRAY":
						recordField = fmt.Sprintf("\"%s\":\"%s\"", strings.Split(field, ".")[1], escapeSpecialChars(memberValue.(string)))
					case "TIMESTAMP_MICROS":
						layout := "2006-01-02T15:04:05.000"
						timestamp, err := time.Parse(layout, memberValue.(string))
						if err != nil {
							c.logger.Error(fmt.Sprintf("time parsing error while creating recordField: %s", err.Error()))
							return false, err
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
				} else {
					recordField = fmt.Sprintf("\"%s\":null", strings.Split(field, ".")[1])
				}
				recordFields = append(recordFields, recordField)
			}
			recordString := strings.Join(recordFields, ",")
			rec = fmt.Sprintf(rec, recordString)
			if err = pw.Write(rec); err != nil {
				c.logger.Error(fmt.Sprintf("Parquet Writer failed to write: %s", err.Error()))
				return false, err
			}
		}

		c.logger.Debug("Parquet Batch Write Finished")

		// Ensure all buffered data is flushed to disk
		if err := pw.WriteStop(); err != nil {
			c.logger.Error(fmt.Sprintf("Error closing Parquet Json file writer: %s", err.Error()))
			return false, err
		}

		// Close the file writer
		if err := fw.Close(); err != nil {
			c.logger.Error(fmt.Sprintf("Error closing Parquet Local file writer: %s", err.Error()))
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (c *connection) createDirectory(dirPath string) (string, error) {
	// Create the directory
	currentTime := time.Now().Format("2006-01-02-150405-MST")

	timestampDirPath := fmt.Sprintf("%s/%s", dirPath, currentTime)
	if len(timestampDirPath) == 0 {
		c.logger.Error("directory path not created properly")
		return "", errors.New("directory path not created properly")
	}
	if err := os.MkdirAll(timestampDirPath, 0755); err != nil {
		c.logger.Error(fmt.Sprintf("Error creating directory: %s", err.Error()))
		return "", err
	}
	c.logger.Debug(fmt.Sprintf("directory created with path: %s", timestampDirPath))
	return timestampDirPath, nil
}

func (c *connection) writeIfDataExists(conf *ApiSourceProperties, query Query) (int, *ApiResponse, error) {
	response, err := c.getDataFromLens2(conf, query)
	if err != nil {
		return 0, nil, err
	}
	count := len(response.Data.Dataset)
	c.logger.Debug(fmt.Sprintf("Data Row Count from Response: %d", count))
	return count, response, nil
}

func escapeSpecialChars(input string) string {
	replacer := strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
		`\n`, `\n`,
		`\t`, `\t`,
		`\r`, `\r`,
	)
	return replacer.Replace(input)
}
