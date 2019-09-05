package services

import (
	"fmt"
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

type mockProvisioningRest struct {
	t *testing.T
}

func (api mockProvisioningRest) Get(path string) ([]byte, error) {
	switch path {
	case "/rest/requisitionNames":
		return []byte(`{"count":2,"foreign-source":["Test1","Test2"]}`), nil
	case "/rest/foreignSourcesConfig/assets":
		return []byte(`{"count":3,"element":["address1","city","zip"]}`), nil
	case "/rest/foreignSourcesConfig/policies":
		return []byte(test.PoliciesJSON), nil
	case "/rest/foreignSourcesConfig/detectors":
		return []byte(test.DetectorsJSON), nil
	}
	return nil, fmt.Errorf("should not be called")
}

func (api mockProvisioningRest) Post(path string, jsonBytes []byte) error {
	return fmt.Errorf("should not be called")
}

func (api mockProvisioningRest) Delete(path string) error {
	return fmt.Errorf("should not be called")
}

func (api mockProvisioningRest) Put(path string, jsonBytes []byte, contentType string) error {
	return fmt.Errorf("should not be called")
}

func TestGetRequisitionNames(t *testing.T) {
	api := GetProvisioningUtilsAPI(&mockProvisioningRest{t})
	list, err := api.GetRequisitionNames()
	assert.NilError(t, err)
	assert.Equal(t, 2, list.Count)
	assert.Equal(t, "Test1", list.ForeignSources[0])
}

func TestRequisitionExists(t *testing.T) {
	api := GetProvisioningUtilsAPI(&mockProvisioningRest{t})
	assert.Equal(t, true, api.RequisitionExists("Test1"))
	assert.Equal(t, false, api.RequisitionExists("Test3"))
}

func TestGetAvailableAssets(t *testing.T) {
	api := GetProvisioningUtilsAPI(&mockProvisioningRest{t})
	list, err := api.GetAvailableAssets()
	assert.NilError(t, err)
	assert.Equal(t, 3, list.Count)
	assert.Equal(t, "address1", list.Element[0])
}

func TestGetAvailableDetectors(t *testing.T) {
	api := GetProvisioningUtilsAPI(&mockProvisioningRest{t})
	list, err := api.GetAvailableDetectors()
	assert.NilError(t, err)
	assert.Equal(t, 2, list.Count)
	assert.Equal(t, "ICMP", list.Plugins[0].Name)
}

func TestGetAvailablePolicies(t *testing.T) {
	api := GetProvisioningUtilsAPI(&mockProvisioningRest{t})
	list, err := api.GetAvailablePolicies()
	assert.NilError(t, err)
	assert.Equal(t, 1, list.Count)
	assert.Equal(t, "Set Node Category", list.Plugins[0].Name)
}
