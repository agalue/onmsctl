package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"gotest.tools/assert"

	"github.com/OpenNMS/onmsctl/model"
)

var mockResources = &model.ResourceList{
	Count: 1,
	Resources: []model.Resource{
		{
			ID:    "node[1].nodeSnmp[]",
			Label: "Sample",
		},
	},
}

type mockResourceRest struct {
	test *testing.T
}

func (api mockResourceRest) Get(path string) ([]byte, error) {
	switch path {
	case "/rest/resources":
		bytes, _ := json.Marshal(mockResources)
		return bytes, nil
	case "/rest/resources/node[1].nodeSnmp[]":
		bytes, _ := json.Marshal(mockResources.Resources[0])
		return bytes, nil
	case "/rest/resources/fornode/1":
		bytes, _ := json.Marshal(mockResources.Resources[0])
		return bytes, nil
	default:
		return nil, fmt.Errorf("should not be called")
	}
}

func (api mockResourceRest) Post(path string, jsonBytes []byte) error {
	return fmt.Errorf("should not be called")
}

func (api mockResourceRest) PostRaw(path string, dataBytes []byte, contentType string) (*http.Response, error) {
	return nil, fmt.Errorf("should not be called")
}

func (api mockResourceRest) Delete(path string) error {
	if path == "/rest/resources/node[1].nodeSnmp[]" {
		return nil
	}
	return fmt.Errorf("should not be called")
}

func (api mockResourceRest) Put(path string, dataBytes []byte, contentType string) error {
	return fmt.Errorf("should not be called")
}

func (api mockResourceRest) IsValid(r *http.Response) error {
	return nil
}

func TestGetResources(t *testing.T) {
	rest := &mockResourceRest{test: t}
	api := GetResourcesAPI(rest)
	list, err := api.GetResources()
	assert.NilError(t, err)
	assert.Equal(t, 1, list.Count)
}

func TestGetResource(t *testing.T) {
	rest := &mockResourceRest{test: t}
	api := GetResourcesAPI(rest)
	r, err := api.GetResource("node[1].nodeSnmp[]")
	assert.NilError(t, err)
	assert.Equal(t, "Sample", r.Label)
}

func TestGetResourceForNode(t *testing.T) {
	rest := &mockResourceRest{test: t}
	api := GetResourcesAPI(rest)
	r, err := api.GetResourceForNode("1")
	assert.NilError(t, err)
	assert.Equal(t, "Sample", r.Label)
}

func TestDeleteResource(t *testing.T) {
	rest := &mockResourceRest{test: t}
	api := GetResourcesAPI(rest)
	err := api.DeleteResource("node[1].nodeSnmp[]")
	assert.NilError(t, err)
}
