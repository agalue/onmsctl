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
	if cli.Debug {
		log.Printf("GET, Path: %s", cli.URL+path)
	}
	request, err := cli.buildRequest(http.MethodGet, cli.URL+path, nil)
	if err != nil {
		return nil, err
	}
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return nil, err
	}
	err = cli.IsValid(response)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(response.Body)
	if cli.Debug && err == nil {
		log.Printf("GET, Data: %s", string(data))
	}
	response.Body.Close()
	return data, err
}

// Post sends an HTTP POST request
func (cli Client) Post(path string, jsonBytes []byte) error {
	response, err := cli.PostRaw(path, jsonBytes, "application/json")
	if err != nil {
		return err
	}
	return cli.IsValid(response)
}

// PostRaw sends an HTTP POST request, returning the raw response
func (cli Client) PostRaw(path string, dataBytes []byte, contentType string) (*http.Response, error) {
	if cli.Debug {
		log.Printf("POST, Path: %s, Type: %s, Data: %s", cli.URL+path, contentType, string(dataBytes))
	}
	request, err := cli.buildRequest(http.MethodPost, cli.URL+path, bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return cli.getHTTPClient().Do(request)
}

// Delete sends an HTTP DELETE request
func (cli Client) Delete(path string) error {
	if cli.Debug {
		log.Printf("DELETE, Path: %s", cli.URL+path)
	}
	request, err := cli.buildRequest(http.MethodDelete, cli.URL+path, nil)
	if err != nil {
		return err
	}
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return cli.IsValid(response)
}

// Put sends an HTTP PUT request
func (cli Client) Put(path string, dataBytes []byte, contentType string) error {
	if cli.Debug {
		log.Println("Data to be sent", string(dataBytes))
	}
	request, err := cli.buildRequest(http.MethodPut, cli.URL+path, bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", contentType)
	response, err := cli.getHTTPClient().Do(request)
	if err != nil {
		return err
	}
	return cli.IsValid(response)
}

// IsValid verifies HTTP response, return an error if it is not valid
func (cli Client) IsValid(response *http.Response) error {
	if cli.Debug {
		log.Printf("Got response: %s", response.Status)
	}
	code := response.StatusCode
	if code == http.StatusOK ||
		code == http.StatusAccepted ||
		code == http.StatusNoContent ||
		code == http.StatusCreated {
		return nil
	}
	return fmt.Errorf("Invalid Response: %s", response.Status)
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
