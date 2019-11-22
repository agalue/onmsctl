package api

import "github.com/OpenNMS/onmsctl/model"

// MonitoringLocationsAPI the API to manipulate Monitoring Locations
type MonitoringLocationsAPI interface {
	GetLocations() (*model.MonitoringLocationList, error)
	LocationExists(location string) (bool, error)
	GetLocation(location string) (*model.MonitoringLocation, error)
	SetLocation(location model.MonitoringLocation) error
}
