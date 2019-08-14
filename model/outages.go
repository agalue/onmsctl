package model

// OnmsOutage OpenNMS outage entity
type OnmsOutage struct {
	ID                   int                   `json:"id" yaml:"id"`
	ForeignSource        string                `json:"foreignSource,omitempty" yaml:"foreignSource,omitempty"`
	ForeignID            string                `json:"foreignId,omitempty" yaml:"foreignId,omitempty"`
	NodeID               int                   `json:"nodeId,omitempty" yaml:"nodeId,omitempty"`
	NodeLabel            string                `json:"nodeLabel,omitempty" yaml:"nodeLabel,omitempty"`
	IPAddress            string                `json:"ipAddress,omitempty" yaml:"ipAddress,omitempty"`
	ServiceID            int                   `json:"serviceId,omitempty" yaml:"serviceId,omitempty"`
	Location             string                `json:"locationName,omitempty" yaml:"locationName,omitempty"`
	MonitoredService     *OnmsMonitoredService `json:"monitoredService,omitempty" yaml:"monitoredService,omitempty"`
	SuppressedBy         string                `json:"suppressedBy,omitempty" yaml:"suppressedBy,omitempty"`
	SuppressedTime       *Time                 `json:"suppressedTime,omitempty" yaml:"suppressedTime,omitempty"`
	ServiceLostTime      *Time                 `json:"ifLostService,omitempty" yaml:"ifLostService,omitempty"`
	ServiceRegainedTime  *Time                 `json:"ifRegainedService,omitempty" yaml:"ifRegainedService,omitempty"`
	ServiceLostEvent     *OnmsEvent            `json:"serviceLostEvent,omitempty" yaml:"serviceLostEvent,omitempty"`
	ServiceRegainedEvent *OnmsEvent            `json:"serviceRegainedEvent,omitempty" yaml:"serviceRegainedEvent,omitempty"`
}

// OnmsOutageList a list of outages
type OnmsOutageList struct {
	Count      int          `json:"count" yaml:"count"`
	TotalCount int          `json:"totalCount" yaml:"totalCount"`
	Offset     int          `json:"offset" yaml:"offset"`
	Outages    []OnmsOutage `json:"outage" yaml:"outages"`
}
