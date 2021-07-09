package model

import (
	"encoding/xml"
	"fmt"
	"net"
)

// MetaDataList a list of metadata
type MetaDataList struct {
	XMLName    xml.Name   `xml:"meta-data-list" json:"-" yaml:"-"`
	Count      int        `xml:"count,attr" json:"count" yaml:"count"`
	TotalCount int        `xml:"totalCount,attr" json:"totalCount" yaml:"totalCount"`
	Offset     int        `xml:"offset,attr" json:"offset" yaml:"offset"`
	Metadata   []MetaData `xml:"meta-data" json:"metaData" yaml:"metadata"`
}

// MetaData a meta-data entry
type MetaData struct {
	XMLName xml.Name `xml:"meta-data" json:"-" yaml:"-"`
	Key     string   `xml:"key" json:"key" yaml:"key"`
	Value   string   `xml:"value" json:"value" yaml:"value"`
	Context string   `xml:"context" json:"context" yaml:"context"`
}

// Validate verify structure and apply defaults when needed
func (obj *MetaData) Validate() error {
	if obj.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	if obj.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	if obj.Context == "" {
		return fmt.Errorf("context cannot be empty")
	}
	return nil
}

// OnmsCategoryList a list of metadata
type OnmsCategoryList struct {
	XMLName    xml.Name       `xml:"categories" json:"-" yaml:"-"`
	Count      int            `xml:"count,attr" json:"count" yaml:"count"`
	TotalCount int            `xml:"totalCount,attr" json:"totalCount" yaml:"totalCount"`
	Offset     int            `xml:"offset,attr" json:"offset" yaml:"offset"`
	Categories []OnmsCategory `xml:"category" json:"category" yaml:"categories"`
}

// OnmsCategory an entity that represents an OpenNMS category
type OnmsCategory struct {
	XMLName xml.Name `xml:"category" json:"-" yaml:"-"`
	ID      int      `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	Name    string   `xml:"name,attr" json:"name" yaml:"name"`
	Groups  []string `xml:"groups,omitempty" json:"groups,omitempty" yaml:"groups,omitempty"`
}

// OnmsAssetRecord an entity that represents an OpenNMS asset record
type OnmsAssetRecord struct {
	XMLName xml.Name `xml:"assetRecord" json:"-" yaml:"-"`
	ID      int      `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	// Identification
	Description     string `xml:"description,omitempty" json:"description,omitempty" yaml:"description,omitempty"`
	Category        string `xml:"category,omitempty" json:"category,omitempty" yaml:"category,omitempty"`
	Manufacturer    string `xml:"manufacturer,omitempty" json:"manufacturer,omitempty" yaml:"manufacturer,omitempty"`
	ModelNumber     string `xml:"modelNumber,omitempty" json:"modelNumber,omitempty" yaml:"modelNumber,omitempty"`
	SerialNumber    string `xml:"serialNumber,omitempty" json:"serialNumber,omitempty" yaml:"serialNumber,omitempty"`
	AssetNumber     string `xml:"assetNumber,omitempty" json:"assetNumber,omitempty" yaml:"assetNumber,omitempty"`
	DateInstalled   *Time  `xml:"dateInstalled,omitempty" json:"dateInstalled,omitempty" yaml:"dateInstalled,omitempty"`
	OperatingSystem string `xml:"operatingSystem,omitempty" json:"operatingSystem,omitempty" yaml:"operatingSystem,omitempty"`
	// Location
	State          string  `xml:"state,omitempty" json:"state,omitempty" yaml:"state,omitempty"`
	Region         string  `xml:"region,omitempty" json:"region,omitempty" yaml:"region,omitempty"`
	Address1       string  `xml:"address1,omitempty" json:"address1,omitempty" yaml:"address1,omitempty"`
	Address2       string  `xml:"address2,omitempty" json:"address2,omitempty" yaml:"address2,omitempty"`
	City           string  `xml:"city,omitempty" json:"city,omitempty" yaml:"city,omitempty"`
	ZIP            string  `xml:"zip,omitempty" json:"zip,omitempty" yaml:"zip,omitempty"`
	Country        string  `xml:"country,omitempty" json:"country,omitempty" yaml:"country,omitempty"`
	Longitude      float32 `xml:"longitude,omitempty" json:"longitude,omitempty" yaml:"longitude,omitempty"`
	Latitude       float32 `xml:"latitude,omitempty" json:"latitude,omitempty" yaml:"latitude,omitempty"`
	Division       string  `xml:"division,omitempty" json:"division,omitempty" yaml:"division,omitempty"`
	Department     string  `xml:"department,omitempty" json:"department,omitempty" yaml:"department,omitempty"`
	Building       string  `xml:"building,omitempty" json:"building,omitempty" yaml:"building,omitempty"`
	Floor          string  `xml:"floor,omitempty" json:"floor,omitempty" yaml:"floor,omitempty"`
	Room           string  `xml:"room,omitempty" json:"room,omitempty" yaml:"room,omitempty"`
	Rack           string  `xml:"rack,omitempty" json:"rack,omitempty" yaml:"rack,omitempty"`
	RackUnitHeight string  `xml:"rackunitheight,omitempty" json:"rackunitheight,omitempty" yaml:"rackUnitHeight,omitempty"`
	Slot           string  `xml:"slot,omitempty" json:"slot,omitempty" yaml:"slot,omitempty"`
	Port           string  `xml:"port,omitempty" json:"port,omitempty" yaml:"port,omitempty"`
	CircuitID      string  `xml:"circuitId,omitempty" json:"circuitId,omitempty" yaml:"circuitId,omitempty"`
	Admin          string  `xml:"admin,omitempty" json:"admin,omitempty" yaml:"admin,omitempty"`
	// Vendor
	Vendor                  string `xml:"vendor,omitempty" json:"vendor,omitempty" yaml:"vendor,omitempty"`
	VendorPhone             string `xml:"vendorPhone,omitempty" json:"vendorPhone,omitempty" yaml:"vendorPhone,omitempty"`
	VendorFax               string `xml:"vendorFax,omitempty" json:"vendorFax,omitempty" yaml:"vendorFax,omitempty"`
	VendorAssetNumber       string `xml:"vendorAssetNumber,omitempty" json:"vendorAssetNumber,omitempty" yaml:"vendorAssetNumber,omitempty"`
	SupportPhone            string `xml:"supportPhone,omitempty" json:"supportPhone,omitempty" yaml:"supportPhone,omitempty"`
	Lease                   string `xml:"lease,omitempty" json:"lease,omitempty" yaml:"lease,omitempty"`
	LeaseExpires            *Time  `xml:"leaseExpires,omitempty" json:"leaseExpires,omitempty" yaml:"leaseExpires,omitempty"`
	MaintContract           string `xml:"maintcontract,omitempty" json:"maintcontract,omitempty" yaml:"maintcontract,omitempty"`
	MaintContractNumber     string `xml:"maintContractNumber,omitempty" json:"maintContractNumber,omitempty" yaml:"maintContractNumber,omitempty"`
	MaintContractExpiration *Time  `xml:"maintContractExpiration,omitempty" json:"maintContractExpiration,omitempty" yaml:"maintContractExpiration,omitempty"`
	// Hardware
	CPU                string `xml:"cpu,omitempty" json:"cpu,omitempty" yaml:"cpu,omitempty"`
	RAM                string `xml:"ram,omitempty" json:"ram,omitempty" yaml:"ram,omitempty"`
	AdditionalHardware string `xml:"additionalhardware,omitempty" json:"additionalhardware,omitempty" yaml:"additionalHardware,omitempty"`
	NumPowerSupplies   string `xml:"numpowersupplies,omitempty" json:"numpowersupplies,omitempty" yaml:"numPowerSupplies,omitempty"`
	InputPower         string `xml:"inputpower,omitempty" json:"inputpower,omitempty" yaml:"inputPower,omitempty"`
	StorageCtrl        string `xml:"storagectrl,omitempty" json:"storagectrl,omitempty" yaml:"storageCtrl,omitempty"`
	HDD1               string `xml:"hdd1,omitempty" json:"hdd1,omitempty" yaml:"hdd1,omitempty"`
	HDD2               string `xml:"hdd2,omitempty" json:"hdd2,omitempty" yaml:"hdd2,omitempty"`
	HDD3               string `xml:"hdd3,omitempty" json:"hdd3,omitempty" yaml:"hdd3,omitempty"`
	HDD4               string `xml:"hdd4,omitempty" json:"hdd4,omitempty" yaml:"hdd4,omitempty"`
	HDD5               string `xml:"hdd5,omitempty" json:"hdd5,omitempty" yaml:"hdd5,omitempty"`
	HDD6               string `xml:"hdd6,omitempty" json:"hdd6,omitempty" yaml:"hdd6,omitempty"`
	// Authentiation
	Username      string `xml:"username,omitempty" json:"username,omitempty" yaml:"username,omitempty"`
	Password      string `xml:"password,omitempty" json:"password,omitempty" yaml:"password,omitempty"`
	Enable        string `xml:"enable,omitempty" json:"enable,omitempty" yaml:"enable,omitempty"`
	AutoEnable    string `xml:"autoenable,omitempty" json:"autoenable,omitempty" yaml:"autoEnable,omitempty"`
	Connection    string `xml:"connection,omitempty" json:"connection,omitempty" yaml:"connection,omitempty"`
	SnmpCommunity string `xml:"snmpcommunity,omitempty" json:"snmpcommunity,omitempty" yaml:"snmpCommunity,omitempty"`
	// Categories
	DisplayCategory   string `xml:"displayCategory,omitempty" json:"displayCategory,omitempty" yaml:"displayCategory,omitempty"`
	NotifyCategory    string `xml:"notifyCategory,omitempty" json:"notifyCategory,omitempty" yaml:"notifyCategory,omitempty"`
	PollerCategory    string `xml:"pollerCategory,omitempty" json:"pollerCategory,omitempty" yaml:"pollerCategory,omitempty"`
	ThresholdCategory string `xml:"thresholdCategory,omitempty" json:"thresholdCategory,omitempty" yaml:"thresholdCategory,omitempty"`
	// VMWare
	VmwareManagedObjectID   string `xml:"vmwareManagedObjectId,omitempty" json:"vmwareManagedObjectId,omitempty" yaml:"vmwareManagedObjectId,omitempty"`
	VmwareManagedEntityType string `xml:"vmwareManagedEntityType,omitempty" json:"vmwareManagedEntityType,omitempty" yaml:"vmwareManagedEntityType,omitempty"`
	VmwareManagementServer  string `xml:"vmwareManagementServer,omitempty" json:"vmwareManagementServer,omitempty" yaml:"vmwareManagementServer,omitempty"`
	VmwareState             string `xml:"vmwareState,omitempty" json:"vmwareState,omitempty" yaml:"vmwareState,omitempty"`
	VmwareTopologyInfo      string `xml:"vmwareTopologyInfo,omitempty" json:"vmwareTopologyInfo,omitempty" yaml:"vmwareTopologyInfo,omitempty"`
	// General
	Comment               string `xml:"comment,omitempty" json:"comment,omitempty" yaml:"comment,omitempty"`
	LastModifiedBy        string `xml:"lastModifiedBy,omitempty" json:"lastModifiedBy,omitempty" yaml:"lastModifiedBy,omitempty"`
	LastModifiedDate      *Time  `xml:"lastModifiedDate,omitempty" json:"lastModifiedDate,omitempty" yaml:"lastModifiedDate,omitempty"`
	ManagedObjectType     string `xml:"managedObjectType,omitempty" json:"managedObjectType,omitempty" yaml:"managedObjectType,omitempty"`
	ManagedObjectInstance string `xml:"managedObjectInstance,omitempty" json:"managedObjectInstance,omitempty" yaml:"managedObjectInstance,omitempty"`
}

// OnmsServiceType an entity that represents an OpenNMS Monitored Service type
type OnmsServiceType struct {
	XMLName xml.Name `xml:"serviceType" json:"-" yaml:"-"`
	ID      int      `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	Name    string   `xml:"name" json:"name" yaml:"name"`
}

// OnmsMonitoredService an entity that represents an OpenNMS Monitored Service
type OnmsMonitoredService struct {
	XMLName     xml.Name         `xml:"service" json:"-" yaml:"-"`
	ID          int              `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	ServiceType *OnmsServiceType `xml:"serviceType" json:"serviceType" yaml:"serviceType"`
	Notify      string           `xml:"notify,omitempty" json:"notify,omitempty" yaml:"notify,omitempty"`
	Qualifier   string           `xml:"qualifier,omitempty" json:"qualifier,omitempty" yaml:"qualifier,omitempty"`
	Status      string           `xml:"status,attr,omitempty" json:"status,omitempty" yaml:"status,omitempty"`
	StatusLong  string           `xml:"statusLong,attr,omitempty" json:"statusLong,omitempty" yaml:"statusLong,omitempty"`
	LastGood    *Time            `xml:"lastGood,omitempty" json:"lastGood,omitempty" yaml:"lastGood,omitempty"`
	LastFail    *Time            `xml:"lastFail,omitempty" json:"lastFail,omitempty" yaml:"lastFail,omitempty"`
	Source      string           `xml:"source,attr,omitempty" json:"source,omitempty" yaml:"source,omitempty"`
	IsDown      bool             `xml:"down,attr,omitempty" json:"down,omitempty" yaml:"isDown,omitempty"`
	Meta        []MetaData       `xml:"metaData,attr,omitempty" json:"metaData,omitempty" yaml:"metaData,omitempty"`
}

// Validate verify structure and apply defaults when needed
func (obj *OnmsMonitoredService) Validate() error {
	if obj.ServiceType == nil || obj.ServiceType.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if obj.Status == "" {
		obj.Status = "A"
	}
	for m := range obj.Meta {
		meta := &obj.Meta[m]
		if err := meta.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// OnmsMonitoredServiceList a list of nodes
type OnmsMonitoredServiceList struct {
	XMLName    xml.Name               `xml:"services" json:"-" yaml:"-"`
	Count      int                    `xml:"count,attr" json:"count" yaml:"count"`
	TotalCount int                    `xml:"totalCount,attr" json:"totalCount" yaml:"totalCount"`
	Offset     int                    `xml:"offset,attr" json:"offset" yaml:"offset"`
	Services   []OnmsMonitoredService `xml:"service" json:"service" yaml:"services"`
}

// OnmsIPInterface an entity that represents an OpenNMS IP Interface
type OnmsIPInterface struct {
	XMLName               xml.Name               `xml:"ipInterface" json:"-" yaml:"-"`
	ID                    string                 `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"` // The JSON returns a string instead of an integer
	NodeID                int                    `xml:"nodeId,omitempty" json:"nodeId,omitempty" yaml:"-,omitempty"`
	IsManaged             string                 `xml:"isManaged,attr,omitempty" json:"isManaged,omitempty" yaml:"isManaged,omitempty"`
	IPAddress             string                 `xml:"ipAddress" json:"ipAddress" yaml:"ipAddress"`
	MonitoredServiceCount int                    `xml:"monitoredServiceCount,attr,omitempty" json:"monitoredServiceCount,omitempty" yaml:"monitoredServiceCount,omitempty"`
	IfIndex               int                    `xml:"ifIndex,attr,omitempty" json:"ifIndex,omitempty" yaml:"ifIndex,omitempty"`
	HostName              string                 `xml:"hostName,omitempty" json:"hostName,omitempty" yaml:"hostName,omitempty"`
	SnmpPrimary           string                 `xml:"snmpPrimary,attr,omitempty" json:"snmpPrimary,omitempty" yaml:"snmpPrimary,omitempty"`
	LastPoll              *Time                  `xml:"lastCapsdPoll,omitempty" json:"lastCapsdPoll,omitempty" yaml:"lastPoll,omitempty"`
	SNMPInterface         *OnmsSnmpInterface     `xml:"snmpInterface,omitempty" json:"snmpInterface,omitempty" yaml:"snmpInterface,omitempty"`
	IsDown                bool                   `xml:"isDown,attr,omitempty" json:"isDown,omitempty" yaml:"isDown,omitempty"`
	HasFlows              bool                   `xml:"hasFlows,attr,omitempty" json:"hasFlows,omitempty" yaml:"hasFlows,omitempty"` // DEPRECATED
	LastIngressFlow       *Time                  `xml:"lastIngressFlow,attr,omitempty" json:"lastIngressFlow,omitempty" yaml:"lastIngressFlow,omitempty"`
	LastEgressFlow        *Time                  `xml:"lastEgressFlow,attr,omitempty" json:"lastEgressFlow,omitempty" yaml:"lastEgressFlow,omitempty"`
	Services              []OnmsMonitoredService `xml:"services,attr,omitempty" json:"services,omitempty" yaml:"services,omitempty"`
	Meta                  []MetaData             `xml:"metaData,attr,omitempty" json:"metaData,omitempty" yaml:"metaData,omitempty"`
}

// ExtractBasic extracts core attributes only
func (obj *OnmsIPInterface) ExtractBasic() *OnmsIPInterface {
	return &OnmsIPInterface{
		ID:          obj.ID,
		IPAddress:   obj.IPAddress,
		IsManaged:   obj.IsManaged,
		SnmpPrimary: obj.SnmpPrimary,
		HostName:    obj.HostName,
	}
}

// Validate verify structure and apply defaults when needed
func (obj *OnmsIPInterface) Validate() error {
	if obj.IPAddress == "" {
		return fmt.Errorf("IP Address cannot be empty")
	}
	ip := net.ParseIP(obj.IPAddress)
	if ip == nil {
		return fmt.Errorf("invalid IP Address: %s", obj.IPAddress)
	}
	if obj.SnmpPrimary == "" {
		obj.SnmpPrimary = "P"
	}
	if obj.IsManaged == "" {
		obj.IsManaged = "M"
	}
	if obj.SNMPInterface != nil {
		if err := obj.SNMPInterface.Validate(); err != nil {
			return err
		}
	}
	for i := range obj.Services {
		svc := &obj.Services[i]
		if err := svc.Validate(); err != nil {
			return err
		}
	}
	for m := range obj.Meta {
		meta := &obj.Meta[m]
		if err := meta.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// OnmsIPInterfaceList a list of nodes
type OnmsIPInterfaceList struct {
	XMLName    xml.Name          `xml:"ipInterfaces" json:"-" yaml:"-"`
	Count      int               `xml:"count,attr" json:"count" yaml:"count"`
	TotalCount int               `xml:"totalCount,attr" json:"totalCount" yaml:"totalCount"`
	Offset     int               `xml:"offset,attr" json:"offset" yaml:"offset"`
	Interfaces []OnmsIPInterface `xml:"ipInterface" json:"ipInterface" yaml:"interfaces"`
}

// OnmsSnmpInterface an entity that represents an OpenNMS SNMP Interface
type OnmsSnmpInterface struct {
	XMLName                 xml.Name `xml:"snmpInterface" json:"-" yaml:"-"`
	ID                      int      `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"`
	IfType                  int      `xml:"ifType,omitempty" json:"ifType,omitempty" yaml:"ifType,omitempty"`
	IfAlias                 string   `xml:"ifAlias,omitempty" json:"ifAlias,omitempty" yaml:"ifAlias,omitempty"`
	IfIndex                 int      `xml:"ifIndex,attr,omitempty" json:"ifIndex,omitempty" yaml:"ifIndex,omitempty"`
	IfDescr                 string   `xml:"ifDescr,omitempty" json:"ifDescr,omitempty" yaml:"ifDescr,omitempty"`
	IfName                  string   `xml:"ifName,omitempty" json:"ifName,omitempty" yaml:"ifName,omitempty"`
	PhysAddress             string   `xml:"physAddress,omitempty" json:"physAddress,omitempty" yaml:"physAddress,omitempty"`
	IfSpeed                 int64    `xml:"ifSpeed,omitempty" json:"ifSpeed,omitempty" yaml:"ifSpeed,omitempty"`
	IfAdminStatus           int      `xml:"ifAdminStatus,omitempty" json:"ifAdminStatus,omitempty" yaml:"ifAdminStatus,omitempty"`
	IfOperStatus            int      `xml:"ifOperStatus,omitempty" json:"ifOperStatus,omitempty" yaml:"ifOperStatus,omitempty"`
	Collect                 bool     `xml:"collect,attr,omitempty" json:"collect,omitempty" yaml:"collect,omitempty"`
	CollectFlag             string   `xml:"collectFlag,attr,omitempty" json:"collectFlag,omitempty" yaml:"collectFlag,omitempty"`
	CollectionUserSpecified bool     `xml:"collectionUserSpecified,omitempty" json:"collectionUserSpecified,omitempty" yaml:"collectionUserSpecified,omitempty"`
	Poll                    bool     `xml:"poll,attr,omitempty" json:"poll,omitempty" yaml:"poll,omitempty"`
	PollFlag                string   `xml:"pollFlag,attr,omitempty" json:"pollFlag,omitempty" yaml:"pollFlag,omitempty"`
	LastPoll                *Time    `xml:"lastCapsdPoll,omitempty" json:"lastCapsdPoll,omitempty" yaml:"lastPoll,omitempty"`
	HasFlows                bool     `xml:"hasFlows,attr,omitempty" json:"hasFlows,omitempty" yaml:"hasFlows,omitempty"` // DEPRECATED
	LastIngressFlow         *Time    `xml:"lastIngressFlow,attr,omitempty" json:"lastIngressFlow,omitempty" yaml:"lastIngressFlow,omitempty"`
	LastEgressFlow          *Time    `xml:"lastEgressFlow,attr,omitempty" json:"lastEgressFlow,omitempty" yaml:"lastEgressFlow,omitempty"`
}

// ExtractBasic extracts core attributes only
func (obj *OnmsSnmpInterface) ExtractBasic() *OnmsSnmpInterface {
	return &OnmsSnmpInterface{
		ID:            obj.ID,
		IfIndex:       obj.IfIndex,
		IfType:        obj.IfType,
		IfAlias:       obj.IfAlias,
		IfName:        obj.IfName,
		IfDescr:       obj.IfDescr,
		IfSpeed:       obj.IfSpeed,
		IfAdminStatus: obj.IfAdminStatus,
		IfOperStatus:  obj.IfOperStatus,
		PhysAddress:   obj.PhysAddress,
		CollectFlag:   obj.CollectFlag,
		PollFlag:      obj.PollFlag,
	}
}

// Validate verify structure and apply defaults when needed
func (obj *OnmsSnmpInterface) Validate() error {
	if obj.IfIndex == 0 {
		return fmt.Errorf("ifIndex cannot be empty")
	}
	if obj.IfName == "" {
		return fmt.Errorf("IfName cannot be empty")
	}
	if obj.IfOperStatus == 0 {
		obj.IfOperStatus = 1
	}
	if obj.IfAdminStatus == 0 {
		obj.IfAdminStatus = 1
	}
	if obj.CollectFlag == "" {
		obj.CollectFlag = "C"
	}
	if obj.PollFlag == "" {
		obj.PollFlag = "N"
	}
	if obj.PhysAddress != "" {
		_, err := net.ParseMAC(obj.PhysAddress)
		if err != nil {
			return fmt.Errorf("invalid Physical Address: %v", err)
		}
	}
	return nil
}

// OnmsSnmpInterfaceList a list of nodes
type OnmsSnmpInterfaceList struct {
	XMLName    xml.Name            `xml:"snmpInterfaces" json:"-" yaml:"-"`
	Count      int                 `xml:"count,attr" json:"count" yaml:"count"`
	TotalCount int                 `xml:"totalCount,attr" json:"totalCount" yaml:"totalCount"`
	Offset     int                 `xml:"offset,attr" json:"offset" yaml:"offset"`
	Interfaces []OnmsSnmpInterface `xml:"snmpInterface" json:"snmpInterface" yaml:"interfaces"`
}

// OnmsNode an entity that represents an OpenNMS node
type OnmsNode struct {
	XMLName         xml.Name            `xml:"node" json:"-" yaml:"-"`
	ID              string              `xml:"id,attr,omitempty" json:"id,omitempty" yaml:"id,omitempty"` // The JSON returns a string instead of an integer
	Type            string              `xml:"type,attr,omitempty" json:"type,omitempty" yaml:"type,omitempty"`
	Label           string              `xml:"label,attr,omitempty" json:"label,omitempty" yaml:"label,omitempty"`
	LabelSource     string              `xml:"labelSource,omitempty" json:"labelSource,omitempty" yaml:"labelSource,omitempty"`
	ForeignSource   string              `xml:"foreignSource,attr,omitempty" json:"foreignSource,omitempty" yaml:"foreignSource,omitempty"`
	ForeignID       string              `xml:"foreignId,attr,omitempty" json:"foreignId,omitempty" yaml:"foreignId,omitempty"`
	Location        string              `xml:"location,attr,omitempty" json:"location,omitempty" yaml:"location,omitempty"`
	SysObjectID     string              `xml:"sysObjectId,omitempty" json:"sysObjectId,omitempty" yaml:"sysObjectId,omitempty"`
	SysName         string              `xml:"sysName,omitempty" json:"sysName,omitempty" yaml:"sysName,omitempty"`
	SysLocation     string              `xml:"sysLocation,omitempty" json:"sysLocation,omitempty" yaml:"sysLocation,omitempty"`
	SysDescription  string              `xml:"sysDescription,omitempty" json:"sysDescription,omitempty" yaml:"sysDescription,omitempty"`
	SysContact      string              `xml:"sysContact,omitempty" json:"sysContact,omitempty" yaml:"sysContact,omitempty"`
	HasFlows        bool                `xml:"hasFlows,attr,omitempty" json:"hasFlows,omitempty" yaml:"hasFlows,omitempty"` // DEPRECATED
	LastIngressFlow *Time               `xml:"lastIngressFlow,attr,omitempty" json:"lastIngressFlow,omitempty" yaml:"lastIngressFlow,omitempty"`
	LastEgressFlow  *Time               `xml:"lastEgressFlow,attr,omitempty" json:"lastEgressFlow,omitempty" yaml:"lastEgressFlow,omitempty"`
	CreateTime      *Time               `xml:"createTime,omitempty" json:"createTime,omitempty" yaml:"createTime,omitempty"`
	LastPoll        *Time               `xml:"lastCapsdPoll,omitempty" json:"lastCapsdPoll,omitempty" yaml:"lastPoll,omitempty"`
	AssetRecord     *OnmsAssetRecord    `xml:"assetRecord,omitempty" json:"assetRecord,omitempty" yaml:"assetRecord,omitempty"`
	Categories      []OnmsCategory      `xml:"categories,omitempty" json:"categories,omitempty" yaml:"categories,omitempty"`
	IPInterfaces    []OnmsIPInterface   `xml:"ipInterfaces,omitempty" json:"ipInterfaces,omitempty" yaml:"ipInterfaces,omitempty"`
	SNMPInterfaces  []OnmsSnmpInterface `xml:"snmpInterfaces,omitempty" json:"snmpInterfaces,omitempty" yaml:"snmpInterfaces,omitempty"`
	Meta            []MetaData          `xml:"metaData,attr,omitempty" json:"metaData,omitempty" yaml:"metaData,omitempty"`
}

// Validate verify structure and apply defaults when needed
func (obj *OnmsNode) Validate() error {
	if obj.Label == "" {
		return fmt.Errorf("label cannot be empty")
	}
	if obj.LabelSource == "" {
		obj.LabelSource = "U"
	}
	if obj.Type == "" {
		obj.Type = "A"
	}
	if obj.ForeignSource == "" && obj.ForeignID != "" {
		return fmt.Errorf("foreign source is required")
	}
	if obj.ForeignSource != "" && obj.ForeignID == "" {
		return fmt.Errorf("foreign ID is required")
	}
	for i := range obj.IPInterfaces {
		intf := &obj.IPInterfaces[i]
		if err := intf.Validate(); err != nil {
			return err
		}
	}
	for i := range obj.SNMPInterfaces {
		intf := &obj.SNMPInterfaces[i]
		if err := intf.Validate(); err != nil {
			return err
		}
	}
	for m := range obj.Meta {
		meta := &obj.Meta[m]
		if err := meta.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ExtractBasic extracts core attributes only
func (obj *OnmsNode) ExtractBasic() *OnmsNode {
	return &OnmsNode{
		ID:             obj.ID,
		Type:           obj.Type,
		Label:          obj.Label,
		LabelSource:    obj.LabelSource,
		Location:       obj.Location,
		SysObjectID:    obj.SysObjectID,
		SysName:        obj.SysName,
		SysLocation:    obj.SysLocation,
		SysDescription: obj.SysDescription,
		SysContact:     obj.SysContact,
	}
}

// GetIPInterface gets a given IP interface by its address (nil if not found)
func (obj *OnmsNode) GetIPInterface(ipAddress string) *OnmsIPInterface {
	for _, ip := range obj.IPInterfaces {
		if ip.IPAddress == ipAddress {
			return &ip
		}
	}
	return nil
}

// GetSnmpInterface gets a given SNMP interface by its ifIndex (nil if not found)
func (obj *OnmsNode) GetSnmpInterface(ifIndex int) *OnmsSnmpInterface {
	for _, snmp := range obj.SNMPInterfaces {
		if snmp.IfIndex == ifIndex {
			return &snmp
		}
	}
	return nil
}

// OnmsNodeList a list of nodes
type OnmsNodeList struct {
	XMLName    xml.Name   `xml:"nodes" json:"-" yaml:"-"`
	Count      int        `xml:"count,attr" json:"count" yaml:"count"`
	TotalCount int        `xml:"totalCount,attr" json:"totalCount" yaml:"totalCount"`
	Offset     int        `xml:"offset,attr" json:"offset" yaml:"offset"`
	Nodes      []OnmsNode `xml:"node" json:"node" yaml:"nodes"`
}
