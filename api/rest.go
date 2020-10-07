package api

import "net/http"

// RestAPI the API for ReST Operations
type RestAPI interface {
	Get(path string) ([]byte, error)
	Post(path string, jsonBytes []byte) error
	PostRaw(path string, dataBytes []byte, contentType string) (*http.Response, error)
	Delete(path string) error
	Put(path string, dataBytes []byte, contentType string) error
	IsValid(response *http.Response) error
}
