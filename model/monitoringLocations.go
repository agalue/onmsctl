package model

// MonitoringLocation an OpenNMS Location
type MonitoringLocation struct {
	Tags                   []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	GeoLocation            string   `json:"geolocation,omitempty" yaml:"geoLocation,omitempty"`
	Latitude               float64  `json:"latitude,omitempty" yaml:"latitude,omitempty"`
	Longitude              float64  `json:"longitude,omitempty" yaml:"longitude,omitempty"`
	LocationName           string   `json:"location-name,omitempty" yaml:"name,omitempty"`
	Priority               int      `json:"priority,omitempty" yaml:"priority,omitempty"`
	MonitoringArea         string   `json:"monitoring-area,omitempty" yaml:"monitoringArea,omitempty"`
	PollingPackageNames    []string `json:"polling-package-names,omitempty" yaml:"pollingPackageNames,omitempty"`
	CollectionPackageNames []string `json:"collection-package-names,omitempty" yaml:"collectionPackageNames,omitempty"`
}

// MonitoringLocationList a list of monitoring locations
type MonitoringLocationList struct {
	Count      int                  `json:"count" yaml:"count"`
	TotalCount int                  `json:"totalCount" yaml:"totalCount"`
	Offset     int                  `json:"offset" yaml:"offset"`
	Locations  []MonitoringLocation `json:"location" yaml:"locations"`
}
