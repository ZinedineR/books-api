package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

// New creates client Factory
func New() ClientFactory {
	return &clientFactory{}
}

// ClientFactory creates specific client implementation
type ClientFactory interface {
	CreateClient() Client
}

type clientFactory struct{}

func (c *clientFactory) CreateClient() Client {
	return &client{}
}

// Client abstracts third party request client
type Client interface {
	Get(string, map[string]string, interface{}) (int, error)
	PostJSON(string, interface{}, map[string]string, interface{}) (int, error)
	PutJSON(string, interface{}, map[string]string, interface{}) (int, error)
	DeleteJSON(string, map[string]string) (int, error)
	PostJSONCallback(string, interface{}, map[string]string, interface{}, string) (int, error)
}

type client struct{}

//func (c *client) PostJSON(url string, payload interface{}, headers map[string]string, dest interface{}) (int, error) {
//	payloadBytes, err := json.Marshal(payload)
//	if err != nil {
//		return http.StatusInternalServerError, errors.Wrap(err, "error while encoding JSON payload")
//	}
//
//	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
//	if err != nil {
//		return http.StatusInternalServerError, errors.Wrap(err, "error creating HTTP request")
//	}
//
//	for k, v := range headers {
//		req.Header.Set(k, v)
//	}
//
//	client := http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return http.StatusInternalServerError, errors.Wrap(err, "error sending HTTP request")
//	}
//	defer resp.Body.Close()
//
//	bodyBytes, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return resp.StatusCode, errors.Wrap(err, "error reading response body")
//	}
//
//	if len(bodyBytes) == 0 && resp.StatusCode == http.StatusOK {
//		return resp.StatusCode, nil
//	}
//
//	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
//		return resp.StatusCode, errors.Wrap(err, "error while decoding response content")
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return resp.StatusCode, errors.New("failed to request")
//	}
//
//	return resp.StatusCode, nil
//}

func (c client) PostJSON(url string, payload interface{}, headers map[string]string, dest interface{}) (
	statusCode int, err error,
) {
	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the HTTP request
	client := &http.Client{}
	//client.Transport = &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	// Read the response body as a string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	// Decode the response JSON into the provided destination struct
	if err := json.Unmarshal(body, &dest); err != nil {
		return resp.StatusCode, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New(string(body[:]))
	}
	return resp.StatusCode, nil
}

func (c client) PutJSON(url string, payload interface{}, headers map[string]string, dest interface{}) (
	statusCode int, err error,
) {
	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	// Create HTTP request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body as a string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	// Decode the response JSON into the provided destination struct
	if err := json.Unmarshal(body, &dest); err != nil {
		return resp.StatusCode, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New(string(body[:]))
	}

	return resp.StatusCode, nil
}

func (c client) DeleteJSON(url string, headers map[string]string) (
	statusCode int, err error,
) {
	// Create HTTP request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body as a string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New(string(body[:]))
	}

	return resp.StatusCode, nil
}

func (c *client) PostJSONCallback(
	url string, payload interface{}, headers map[string]string, dest interface{}, apiRequestId string,
) (int, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return http.StatusInternalServerError, errors.New("error while encoding JSON payload: " + err.Error())
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return http.StatusInternalServerError, errors.New("error creating HTTP request")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.New("error sending HTTP request")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, errors.New("error reading response body")
	}

	logrus.Infoln(fmt.Sprintf("REQUEST ID: %s , RESPONSE CODE CALLBACK: %v", apiRequestId, resp.StatusCode))
	if len(bodyBytes) == 0 && resp.StatusCode == http.StatusOK {
		return resp.StatusCode, nil
	}

	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return resp.StatusCode, errors.New("error while decoding response content")
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New("failed to request")
	}

	return resp.StatusCode, nil
}

func (c *client) Get(url string, headers map[string]string, dest interface{}) (
	int, error,
) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return http.StatusInternalServerError, errors.New("error creating HTTP request")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, errors.New("error sending HTTP request")
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, errors.New("error reading response body")
	}

	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return resp.StatusCode, errors.New("error while decoding response content :" + string(bodyBytes[:]))
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, errors.New("failed to request")
	}

	return resp.StatusCode, nil
}
