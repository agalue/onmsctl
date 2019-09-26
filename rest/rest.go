package rest

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

// Instance a global reference to the ReST Client instance
var Instance = Client{
	URL:      "http://localhost:8980/opennms",
	Username: "admin",
	Password: "admin",
	Timeout:  5,
}

// Client OpenNMS ReST API configuration
type Client struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Insecure bool   `yaml:"insecure"`
	Timeout  int    `yaml:"timeout"`
	Debug    bool   `yaml:"debug"`
}

func (cli Client) getHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cli.Insecure},
	}
	timeout := time.Duration(cli.Timeout) * time.Second
	return &http.Client{Transport: tr, Timeout: timeout}
}

// Get sends an HTTP GET request
func (cli Client) Get(path string) ([]byte, error) {
	request, err := cli.buildRequest(http.MethodGet, cli.URL+path, nil)
	if err != nil {
		return nil, err
	}
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return nil, err
	}
	err = httpIsValid(response)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(response.Body)
	if cli.Debug && err == nil {
		log.Println("Data received", string(data))
	}
	return data, err
}

// Post sends an HTTP POST request
func (cli Client) Post(path string, jsonBytes []byte) error {
	if cli.Debug {
		log.Println("Data to be sent", string(jsonBytes))
	}
	request, err := cli.buildRequest(http.MethodPost, cli.URL+path, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return httpIsValid(response)
}

// Delete sends an HTTP DELETE request
func (cli Client) Delete(path string) error {
	request, err := cli.buildRequest(http.MethodDelete, cli.URL+path, nil)
	if err != nil {
		return err
	}
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return httpIsValid(response)
}

// Put sends an HTTP PUT request
func (cli Client) Put(path string, jsonBytes []byte, contentType string) error {
	if cli.Debug {
		log.Println("Data to be sent", string(jsonBytes))
	}
	request, err := cli.buildRequest(http.MethodPut, cli.URL+path, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", contentType)
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return httpIsValid(response)
}

func (cli Client) buildRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.SetBasicAuth(cli.Username, cli.Password)
	if cli.Debug {
		trace := &httptrace.ClientTrace{
			GotConn: func(connInfo httptrace.GotConnInfo) {
				log.Println("Got Connection", connInfo)
			},
			ConnectStart: func(network, addr string) {
				log.Println("Dial start", network, addr)
			},
			ConnectDone: func(network, addr string, err error) {
				log.Println("Dial done", network, addr)
			},
			GotFirstResponseByte: func() {
				log.Println("Got first response byte!")
			},
			WroteHeaderField: func(key string, value []string) {
				log.Println("Wrote header", key, value)
			},
			WroteRequest: func(wr httptrace.WroteRequestInfo) {
				log.Println("Wrote request error?", wr.Err)
			},
		}
		request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	}
	return request, nil
}

func httpIsValid(response *http.Response) error {
	code := response.StatusCode
	if code != http.StatusOK && code != http.StatusAccepted && code != http.StatusNoContent {
		return fmt.Errorf("Invalid Response: %s", response.Status)
	}
	return nil
}
