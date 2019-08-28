package api

// RestAPI the API for ReST Operations
type RestAPI interface {
	Get(path string) ([]byte, error)
	Post(path string, jsonBytes []byte) error
	Delete(path string) error
	Put(path string, jsonBytes []byte, contentType string) error
}
