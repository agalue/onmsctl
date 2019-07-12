package info

// OnmsInfoDatetimeFormat provides information about the time format
type OnmsInfoDatetimeFormat struct {
	ZoneID string `json:"zoneId" yaml:"zoneId"`
	Format string `json:"datetimeformat" yaml:"format"`
}

// OnmsInfo provides information about the OpenNMS server
type OnmsInfo struct {
	DisplayVersion     string                  `json:"displayVersion" yaml:"displayVersion"`
	Version            string                  `json:"version" yaml:"version"`
	PackageName        string                  `json:"packageName" yaml:"packageName"`
	PackageDescription string                  `json:"packageDescription" yaml:"packageDescription"`
	DatetimeFormat     *OnmsInfoDatetimeFormat `json:"datetimeformatConfig" yaml:"datetimeFormat"`
}
