package rest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Instance a global reference to the ReST Client instance
var Instance = Client{
	URL:                "http://localhost:8980/opennms",
	Username:           "admin",
	Password:           "admin",
	InsecureSkipVerify: true,
}

// Client OpenNMS ReST API configuration
type Client struct {
	URL                string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

func (cli Client) getHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cli.InsecureSkipVerify},
	}
	return &http.Client{Transport: tr}
}

// Get sends an HTTP GET request
func (cli Client) Get(path string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, cli.URL+path, nil)
	if err != nil {
		return nil, err
	}
	setCommonHeaders(cli, request)
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return nil, err
	}
	err = httpIsValid(response)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}

// Post sends an HTTP POST request
func (cli Client) Post(path string, jsonBytes []byte) error {
	request, err := http.NewRequest(http.MethodPost, cli.URL+path, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	setCommonHeaders(cli, request)
	request.Header.Set("Content-Type", "application/json")
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return httpIsValid(response)
}

// Delete sends an HTTP DELETE request
func (cli Client) Delete(path string) error {
	request, err := http.NewRequest(http.MethodDelete, cli.URL+path, nil)
	if err != nil {
		return err
	}
	setCommonHeaders(cli, request)
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return httpIsValid(response)
}

// Put sends an HTTP PUT request
func (cli Client) Put(path string, jsonBytes []byte) error {
	request, err := http.NewRequest(http.MethodPut, cli.URL+path, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	setCommonHeaders(cli, request)
	request.Header.Set("Content-Type", "application/json")
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return httpIsValid(response)
}

func setCommonHeaders(cli Client, request *http.Request) {
	request.Header.Set("Accept", "application/json")
	request.SetBasicAuth(cli.Username, cli.Password)
}

func httpIsValid(response *http.Response) error {
	code := response.StatusCode
	if code != http.StatusOK && code != http.StatusAccepted && code != http.StatusNoContent {
		return fmt.Errorf("Invalid Response: %s", response.Status)
	}
	return nil
}
