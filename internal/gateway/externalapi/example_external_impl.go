package externalapi

import (
	"boiler-plate-clean/config"
	"encoding/base64"
	"github.com/RumbiaID/pkg-library/app/pkg/httpclient"
)

type ExampleExternalImpl struct {
	config     *config.Config
	HttpClient httpclient.Client
}

func NewExampleExternalImpl(
	config *config.Config, HttpClient httpclient.Client,
) ExampleSvcExternal {
	return &ExampleExternalImpl{
		config:     config,
		HttpClient: HttpClient,
	}
}
func (b *ExampleExternalImpl) Post() (interface{}, int, error) {
	var response interface{}
	urlPath := "localhost:8000/user"
	//var request map[string]interface{}
	request := map[string]interface{}{
		"key": false,
	}

	params := map[string]string{
		"authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("credential")),
		"Content-Type":  "application/json",
		"accept":        "application/json",
	}

	statusCode, err := b.HttpClient.PostJSON(urlPath, request, params, &response)
	if err != nil {
		return &response, statusCode, err
	}
	return &response, statusCode, nil
}
