package services

import (
	"encoding/json"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

type monitoringLocationsAPI struct {
	rest api.RestAPI
}

// GetMonitoringLocationsAPI Obtain an implementation of the Monitoring Locations API
func GetMonitoringLocationsAPI(rest api.RestAPI) api.MonitoringLocationsAPI {
	return &monitoringLocationsAPI{rest}
}

func (api monitoringLocationsAPI) GetLocations() (*model.MonitoringLocationList, error) {
	jsonString, err := api.rest.Get("/api/v2/monitoringLocations")
	if err != nil {
		return nil, err
	}
	locations := &model.MonitoringLocationList{}
	if err := json.Unmarshal(jsonString, locations); err != nil {
		return nil, err
	}
	return locations, nil
}

func (api monitoringLocationsAPI) LocationExists(location string) (bool, error) {
	locations, err := api.GetLocations()
	if err != nil {
		return false, err
	}
	for _, loc := range locations.Locations {
		if loc.LocationName == location {
			return true, nil
		}
	}
	return false, nil
}

func (api monitoringLocationsAPI) GetLocation(location string) (*model.MonitoringLocation, error) {
	jsonString, err := api.rest.Get("/api/v2/monitoringLocations/" + location)
	if err != nil {
		return nil, err
	}
	loc := &model.MonitoringLocation{}
	if err := json.Unmarshal(jsonString, loc); err != nil {
		return nil, err
	}
	return loc, nil
}

func (api monitoringLocationsAPI) SetLocation(location model.MonitoringLocation) error {
	jsonBytes, err := json.Marshal(location)
	if err != nil {
		return err
	}
	return api.rest.Post("/api/v2/monitoringLocations", jsonBytes)
}
