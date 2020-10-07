package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"gotest.tools/assert"
)

type mockMonitoringLocationRest struct {
	test     *testing.T
	lastPath string
}

func (api *mockMonitoringLocationRest) Get(path string) ([]byte, error) {
	api.lastPath = path
	if path == "/api/v2/monitoringLocations" {
		bytes, _ := json.Marshal(&model.MonitoringLocationList{
			Count: 1,
			Locations: []model.MonitoringLocation{
				{
					LocationName: "Apex",
				},
			},
		})
		return bytes, nil
	}
	if path == "/api/v2/monitoringLocations/Apex" {
		bytes, _ := json.Marshal(model.MonitoringLocation{
			LocationName: "Apex",
		})
		return bytes, nil
	}
	return nil, fmt.Errorf("Not Found")
}

func (api mockMonitoringLocationRest) Post(path string, jsonBytes []byte) error {
	loc := &model.MonitoringLocation{}
	err := json.Unmarshal(jsonBytes, loc)
	assert.NilError(api.test, err)
	assert.Equal(api.test, "Cary", loc.LocationName)
	return nil
}

func (api mockMonitoringLocationRest) PostRaw(path string, dataBytes []byte, contentType string) (*http.Response, error) {
	return nil, fmt.Errorf("should not be called")
}

func (api mockMonitoringLocationRest) Delete(path string) error {
	return fmt.Errorf("should not be called")
}

func (api mockMonitoringLocationRest) Put(path string, dataBytes []byte, contentType string) error {
	return fmt.Errorf("should not be called")
}

func (api mockMonitoringLocationRest) IsValid(r *http.Response) error {
	return nil
}

func TestLocationExists(t *testing.T) {
	rest := &mockMonitoringLocationRest{test: t}
	api := GetMonitoringLocationsAPI(rest)

	exists, err := api.LocationExists("Apex")
	assert.NilError(t, err)
	assert.Assert(t, exists)

	exists, err = api.LocationExists("Cary")
	assert.NilError(t, err)
	assert.Assert(t, !exists)
}

func TestGetLocations(t *testing.T) {
	rest := &mockMonitoringLocationRest{test: t}
	api := GetMonitoringLocationsAPI(rest)

	list, err := api.GetLocations()
	assert.NilError(t, err)
	assert.Equal(t, 1, list.Count)
	assert.Equal(t, "Apex", list.Locations[0].LocationName)
}

func TestGetLocation(t *testing.T) {
	rest := &mockMonitoringLocationRest{test: t}
	api := GetMonitoringLocationsAPI(rest)

	loc, err := api.GetLocation("Apex")
	assert.NilError(t, err)
	assert.Equal(t, "Apex", loc.LocationName)

	loc, err = api.GetLocation("Cary")
	assert.Assert(t, loc == nil)
	assert.Assert(t, err != nil)
}

func TestSetLocation(t *testing.T) {
	rest := &mockMonitoringLocationRest{test: t}
	api := GetMonitoringLocationsAPI(rest)

	err := api.SetLocation(model.MonitoringLocation{
		LocationName: "Cary",
	})
	assert.NilError(t, err)
}
