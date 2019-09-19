package model

// OnmsCategory an entity that represents an OpenNMS category
type OnmsCategory struct {
	ID     int      `json:"id" yaml:"id"`
	Name   string   `json:"name" yaml:"name"`
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
}

// OnmsAssetRecord an entity that represents an OpenNMS asset record
type OnmsAssetRecord struct {
	ID int `json:"id,omitempty" yaml:"id,omitempty"`
	// Identification
	Description     string `json:"description,omitempty" yaml:"description,omitempty"`
	Category        string `json:"category,omitempty" yaml:"category,omitempty"`
	Manufacturer    string `json:"manufacturer,omitempty" yaml:"manufacturer,omitempty"`
	ModelNumber     string `json:"modelNumber,omitempty" yaml:"modelNumber,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty" yaml:"serialNumber,omitempty"`
	AssetNumber     string `json:"assetNumber,omitempty" yaml:"assetNumber,omitempty"`
	DateInstalled   *Time  `json:"dateInstalled,omitempty" yaml:"dateInstalled,omitempty"`
	OperatingSystem string `json:"operatingSystem,omitempty" yaml:"operatingSystem,omitempty"`
	// Location
	State          string  `json:"state,omitempty" yaml:"state,omitempty"`
	Region         string  `json:"region,omitempty" yaml:"region,omitempty"`
	Address1       string  `json:"address1,omitempty" yaml:"address1,omitempty"`
	Address2       string  `json:"address2,omitempty" yaml:"address2,omitempty"`
	City           string  `json:"city,omitempty" yaml:"city,omitempty"`
	ZIP            string  `json:"zip,omitempty" yaml:"zip,omitempty"`
	Country        string  `json:"country,omitempty" yaml:"country,omitempty"`
	Longitude      float32 `json:"longitude,omitempty" yaml:"longitude,omitempty"`
	Latitude       float32 `json:"latitude,omitempty" yaml:"latitude,omitempty"`
	Division       string  `json:"division,omitempty" yaml:"division,omitempty"`
	Department     string  `json:"department,omitempty" yaml:"department,omitempty"`
	Building       string  `json:"building,omitempty" yaml:"building,omitempty"`
	Floor          string  `json:"floor,omitempty" yaml:"floor,omitempty"`
	Room           string  `json:"room,omitempty" yaml:"room,omitempty"`
	Rack           string  `json:"rack,omitempty" yaml:"rack,omitempty"`
	RackUnitHeight string  `json:"rackunitheight,omitempty" yaml:"rackUnitHeight,omitempty"`
	Slot           string  `json:"slot,omitempty" yaml:"slot,omitempty"`
	Port           string  `json:"port,omitempty" yaml:"port,omitempty"`
	CircuitID      string  `json:"circuitId,omitempty" yaml:"circuitId,omitempty"`
	Admin          string  `json:"admin,omitempty" yaml:"admin,omitempty"`
	// Vendor
	Vendor                  string `json:"vendor,omitempty" yaml:"vendor,omitempty"`
	VendorPhone             string `json:"vendorPhone,omitempty" yaml:"vendorPhone,omitempty"`
	VendorFax               string `json:"vendorFax,omitempty" yaml:"vendorFax,omitempty"`
	VendorAssetNumber       string `json:"vendorAssetNumber,omitempty" yaml:"vendorAssetNumber,omitempty"`
	SupportPhone            string `json:"supportPhone,omitempty" yaml:"supportPhone,omitempty"`
	Lease                   string `json:"lease,omitempty" yaml:"lease,omitempty"`
	LeaseExpires            *Time  `json:"leaseExpires,omitempty" yaml:"leaseExpires,omitempty"`
	MaintContract           string `json:"maintcontract,omitempty" yaml:"maintcontract,omitempty"`
	MaintContractNumber     string `json:"maintContractNumber,omitempty" yaml:"maintContractNumber,omitempty"`
	MaintContractExpiration *Time  `json:"maintContractExpiration,omitempty" yaml:"maintContractExpiration,omitempty"`
	// Hardware
	CPU                string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	RAM                string `json:"ram,omitempty" yaml:"ram,omitempty"`
	AdditionalHardware string `json:"additionalhardware,omitempty" yaml:"additionalHardware,omitempty"`
	NumPowerSupplies   string `json:"numpowersupplies,omitempty" yaml:"numPowerSupplies,omitempty"`
	InputPower         string `json:"inputpower,omitempty" yaml:"inputPower,omitempty"`
	StorageCtrl        string `json:"storagectrl,omitempty" yaml:"storageCtrl,omitempty"`
	HDD1               string `json:"hdd1,omitempty" yaml:"hdd1,omitempty"`
	HDD2               string `json:"hdd2,omitempty" yaml:"hdd2,omitempty"`
	HDD3               string `json:"hdd3,omitempty" yaml:"hdd3,omitempty"`
	HDD4               string `json:"hdd4,omitempty" yaml:"hdd4,omitempty"`
	HDD5               string `json:"hdd5,omitempty" yaml:"hdd5,omitempty"`
	HDD6               string `json:"hdd6,omitempty" yaml:"hdd6,omitempty"`
	// Authentiation
	Username      string `json:"username,omitempty" yaml:"username,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty"`
	Enable        string `json:"enable,omitempty" yaml:"enable,omitempty"`
	AutoEnable    string `json:"autoenable,omitempty" yaml:"autoEnable,omitempty"`
	Connection    string `json:"connection,omitempty" yaml:"connection,omitempty"`
	SnmpCommunity string `json:"snmpcommunity,omitempty" yaml:"snmpCommunity,omitempty"`
	// Categories
	DisplayCategory   string `json:"displayCategory,omitempty" yaml:"displayCategory,omitempty"`
	NotifyCategory    string `json:"notifyCategory,omitempty" yaml:"notifyCategory,omitempty"`
	PollerCategory    string `json:"pollerCategory,omitempty" yaml:"pollerCategory,omitempty"`
	ThresholdCategory string `json:"thresholdCategory,omitempty" yaml:"thresholdCategory,omitempty"`
	// VMWare
	VmwareManagedObjectID   string `json:"vmwareManagedObjectId,omitempty" yaml:"vmwareManagedObjectId,omitempty"`
	VmwareManagedEntityType string `json:"vmwareManagedEntityType,omitempty" yaml:"vmwareManagedEntityType,omitempty"`
	VmwareManagementServer  string `json:"vmwareManagementServer,omitempty" yaml:"vmwareManagementServer,omitempty"`
	VmwareState             string `json:"vmwareState,omitempty" yaml:"vmwareState,omitempty"`
	VmwareTopologyInfo      string `json:"vmwareTopologyInfo,omitempty" yaml:"vmwareTopologyInfo,omitempty"`
	// General
	Comment               string `json:"comment,omitempty" yaml:"comment,omitempty"`
	LastModifiedBy        string `json:"lastModifiedBy,omitempty" yaml:"lastModifiedBy,omitempty"`
	LastModifiedDate      *Time  `json:"lastModifiedDate,omitempty" yaml:"lastModifiedDate,omitempty"`
	ManagedObjectType     string `json:"managedObjectType,omitempty" yaml:"managedObjectType,omitempty"`
	ManagedObjectInstance string `json:"managedObjectInstance,omitempty" yaml:"managedObjectInstance,omitempty"`
}

// OnmsServiceType an entity that represents an OpenNMS Monitored Service type
type OnmsServiceType struct {
	ID   int    `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// OnmsMonitoredService an entity that represents an OpenNMS Monitored Service
type OnmsMonitoredService struct {
	ID          int              `json:"id,omitempty" yaml:"id,omitempty"`
	ServiceType *OnmsServiceType `json:"serviceType" yaml:"serviceType"`
	Notify      string           `json:"notify,omitempty" yaml:"notify,omitempty"`
	Qualifier   string           `json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
	Status      string           `json:"status,omitempty" yaml:"status,omitempty"`
	StatusLong  string           `json:"statusLong,omitempty" yaml:"statusLong,omitempty"`
	LastGood    *Time            `json:"lastGood,omitempty" yaml:"lastGood,omitempty"`
	LastFail    *Time            `json:"lastFail,omitempty" yaml:"lastFail,omitempty"`
	Source      string           `json:"source,omitempty" yaml:"source,omitempty"`
	IsDown      bool             `json:"down" yaml:"isDown"`
}

// OnmsMonitoredServiceList a list of nodes
type OnmsMonitoredServiceList struct {
	Count      int                    `json:"count" yaml:"count"`
	TotalCount int                    `json:"totalCount" yaml:"totalCount"`
	Offset     int                    `json:"offset" yaml:"offset"`
	Services   []OnmsMonitoredService `json:"service" yaml:"services"`
}

// OnmsIPInterface an entity that represents an OpenNMS IP Interface
type OnmsIPInterface struct {
	ID                    int                `json:"id" yaml:"id"`
	NodeID                int                `json:"nodeId,omitempty" yaml:"-,omitempty"`
	IsManaged             string             `json:"isManaged,omitempty" yaml:"isManaged,omitempty"`
	IPAddress             string             `json:"ipAddress" yaml:"ipAddress"`
	MonitoredServiceCount int                `json:"monitoredServiceCount" yaml:"monitoredServiceCount"`
	IfIndex               int                `json:"ifIndex,omitempty" yaml:"ifIndex,omitempty"`
	HostName              string             `json:"hostName,omitempty" yaml:"hostName,omitempty"`
	SnmpPrimary           string             `json:"snmpPrimary,omitempty" yaml:"snmpPrimary,omitempty"`
	LastPoll              *Time              `json:"lastCapsdPoll,omitempty" yaml:"lastPoll,omitempty"`
	SNMPInterface         *OnmsSnmpInterface `json:"snmpInterface,omitempty" yaml:"snmpInterface,omitempty"`
	IsDown                bool               `json:"isDown" yaml:"isDown"`
	HasFlows              bool               `json:"hasFlows" yaml:"hasFlows"`
}

// OnmsIPInterfaceList a list of nodes
type OnmsIPInterfaceList struct {
	Count      int               `json:"count" yaml:"count"`
	TotalCount int               `json:"totalCount" yaml:"totalCount"`
	Offset     int               `json:"offset" yaml:"offset"`
	Interfaces []OnmsIPInterface `json:"ipInterface" yaml:"interfaces"`
}

// OnmsSnmpInterface an entity that represents an OpenNMS SNMP Interface
type OnmsSnmpInterface struct {
	ID                      int    `json:"id" yaml:"id"`
	IfType                  int    `json:"ifType,omitempty" yaml:"ifType,omitempty"`
	IfAlias                 string `json:"ifAlias,omitempty" yaml:"ifAlias,omitempty"`
	IfIndex                 int    `json:"ifIndex,omitempty" yaml:"ifIndex,omitempty"`
	IfDescr                 string `json:"ifDescr,omitempty" yaml:"ifDescr,omitempty"`
	IfName                  string `json:"ifName,omitempty" yaml:"ifName,omitempty"`
	PhysAddress             string `json:"physAddress,omitempty" yaml:"physAddress,omitempty"`
	IfSpeed                 int    `json:"ifSpeed,omitempty" yaml:"ifSpeed,omitempty"`
	IfAdminStatus           int    `json:"ifAdminStatus,omitempty" yaml:"ifAdminStatus,omitempty"`
	IfOperStatus            int    `json:"ifOperStatus,omitempty" yaml:"ifOperStatus,omitempty"`
	Collect                 bool   `json:"collect" yaml:"collect"`
	CollectFlag             string `json:"collectFlag,omitempty" yaml:"collectFlag,omitempty"`
	CollectionUserSpecified bool   `json:"collectionUserSpecified,omitempty" yaml:"collectionUserSpecified,omitempty"`
	Poll                    bool   `json:"poll" yaml:"poll"`
	PollFlag                string `json:"pollFlag,omitempty" yaml:"pollFlag,omitempty"`
	LastPoll                *Time  `json:"lastCapsdPoll,omitempty" yaml:"lastPoll,omitempty"`
	HasFlows                bool   `json:"hasFlows" yaml:"hasFlows"`
}

// OnmsSnmpInterfaceList a list of nodes
type OnmsSnmpInterfaceList struct {
	Count      int                 `json:"count" yaml:"count"`
	TotalCount int                 `json:"totalCount" yaml:"totalCount"`
	Offset     int                 `json:"offset" yaml:"offset"`
	Interfaces []OnmsSnmpInterface `json:"snmpInterface" yaml:"interfaces"`
}

// OnmsNode an entity that represents an OpenNMS node
type OnmsNode struct {
	ID             string           `json:"id" yaml:"id"` // TODO the JSON returns a string instead of an integer
	Label          string           `json:"label,omitempty" yaml:"label,omitempty"`
	LabelSource    string           `json:"labelSource,omitempty" yaml:"labelSource,omitempty"`
	ForeignSource  string           `json:"foreignSource,omitempty" yaml:"foreignSource,omitempty"`
	ForeignID      string           `json:"foreignId,omitempty" yaml:"foreignId,omitempty"`
	Location       string           `json:"location,omitempty" yaml:"location,omitempty"`
	SysObjectID    string           `json:"sysObjectId,omitempty" yaml:"sysObjectId,omitempty"`
	SysName        string           `json:"sysName,omitempty" yaml:"sysName,omitempty"`
	SysLocation    string           `json:"sysLocation,omitempty" yaml:"sysLocation,omitempty"`
	SysDescription string           `json:"sysDescription,omitempty" yaml:"sysDescription,omitempty"`
	SysContact     string           `json:"sysContact,omitempty" yaml:"sysContact,omitempty"`
	HasFlows       bool             `json:"hasFlows" yaml:"hasFlows"`
	CreateTime     *Time            `json:"createTime,omitempty" yaml:"createTime,omitempty"`
	LastPoll       *Time            `json:"lastCapsdPoll,omitempty" yaml:"lastPoll,omitempty"`
	AssetRecord    *OnmsAssetRecord `json:"assetRecord,omitempty" yaml:"assetRecord,omitempty"`
	Categories     []OnmsCategory   `json:"categories,omitempty" yaml:"categories,omitempty"`
}

// OnmsNodeList a list of nodes
type OnmsNodeList struct {
	Count      int        `json:"count" yaml:"count"`
	TotalCount int        `json:"totalCount" yaml:"totalCount"`
	Offset     int        `json:"offset" yaml:"offset"`
	Nodes      []OnmsNode `json:"node" yaml:"nodes"`
}
