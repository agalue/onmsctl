package model

import (
	"encoding/xml"
	"fmt"
	"net"
)

// RequisitionMetaData a meta-data entry
type RequisitionMetaData struct {
	XMLName xml.Name `xml:"meta-data" json:"-" yaml:"-"`
	Key     string   `xml:"key,attr" json:"key" yaml:"key"`
	Value   string   `xml:"value,attr" json:"value" yaml:"value"`
}

// RequisitionMonitoredService an IP interface monitored service
type RequisitionMonitoredService struct {
	XMLName  xml.Name              `xml:"monitored-service" json:"-" yaml:"-"`
	Name     string                `xml:"name,attr" json:"service-name" yaml:"name"`
	MetaData []RequisitionMetaData `xml:"meta-data,omitempty" json:"meta-data,omitempty" yaml:"metaData,omitempty"`
}

// IsValid returns an error if the service is invalid
func (s RequisitionMonitoredService) IsValid() error {
	if s.Name == "" {
		return fmt.Errorf("Service name cannot be null")
	}
	return nil
}

// RequisitionAsset a requisition node asset field
type RequisitionAsset struct {
	XMLName xml.Name `xml:"asset" json:"-" yaml:"-"`
	Name    string   `xml:"name,attr" json:"name" yaml:"name"`
	Value   string   `xml:"value,attr" json:"value" yaml:"value"`
}

// IsValid returns an error if asset field is invalid
func (a RequisitionAsset) IsValid() error {
	if a.Name == "" {
		return fmt.Errorf("Asset name cannot be empty")
	}
	if a.Value == "" {
		return fmt.Errorf("Asset value cannot be empty")
	}
	return nil
}

// RequisitionCategory a requisition node category
type RequisitionCategory struct {
	XMLName xml.Name `xml:"category" json:"-" yaml:"-"`
	Name    string   `xml:"name,attr" json:"name" yaml:"name"`
}

// IsValid returns an error if the category is invalid
func (c RequisitionCategory) IsValid() error {
	if c.Name == "" {
		return fmt.Errorf("Category name cannot be null")
	}
	return nil
}

// RequisitionInterface an IP interface of a requisition node
type RequisitionInterface struct {
	XMLName     xml.Name                      `xml:"interface" json:"-" yaml:"-"`
	IPAddress   string                        `xml:"ip-addr,attr" json:"ip-addr" yaml:"ipAddress"`
	Description string                        `xml:"descr,attr,omitempty" json:"descr,omitempty" yaml:"description,omitempty"`
	SnmpPrimary string                        `xml:"snmp-primary,attr,omitempty" json:"snmp-primary" yaml:"snmpPrimary"`
	Status      int                           `xml:"status,attr,omitempty" json:"status" yaml:"status"`
	Services    []RequisitionMonitoredService `xml:"monitored-service,omitempty" json:"monitored-service,omitempty" yaml:"services,omitempty"`
	MetaData    []RequisitionMetaData         `xml:"meta-data,omitempty" json:"meta-data,omitempty" yaml:"metaData,omitempty"`
}

// IsValid returns an error if the interface definition is invalid
func (i *RequisitionInterface) IsValid() error {
	if i.IPAddress == "" {
		return fmt.Errorf("IP Address cannot be empty")
	}
	if i.Status == 0 { // Set a reasonable default when the status is not initialized
		i.Status = 1
	}
	if i.Status != 1 && i.Status != 3 {
		return fmt.Errorf("Invalid Status: %d", i.Status)
	}
	if i.SnmpPrimary == "" { // Set a reasonable default when the primary flag is not initialized
		i.SnmpPrimary = "N"
	}
	if i.SnmpPrimary != "P" && i.SnmpPrimary != "S" && i.SnmpPrimary != "N" {
		return fmt.Errorf("Invalid SnmpPrimary: %s", i.SnmpPrimary)
	}
	ip := net.ParseIP(i.IPAddress)
	if ip == nil {
		addresses, err := net.LookupIP(i.IPAddress)
		if err != nil {
			return fmt.Errorf("Cannot get address from %s (invalid IP or FQDN); %s", i.IPAddress, err)
		}
		fmt.Printf("%s translates to %s, using the first entry.\n", i.IPAddress, addresses)
		i.IPAddress = addresses[0].String()
	}
	serviceMap := make(map[string]int)
	for _, s := range i.Services {
		serviceMap[s.Name]++
		err := s.IsValid()
		if err != nil {
			return err
		}
	}
	for service, count := range serviceMap {
		if count > 1 {
			return fmt.Errorf("Service %s is defined more than once on IP %s", service, i.IPAddress)
		}
	}
	return nil
}

// RequisitionNode a requisitioned node
type RequisitionNode struct {
	XMLName             xml.Name               `xml:"node" json:"-" yaml:"-"`
	NodeLabel           string                 `xml:"node-label,attr" json:"node-label" yaml:"nodeLabel"`
	ForeignID           string                 `xml:"foreign-id,attr" json:"foreign-id" yaml:"foreignID"`
	Location            string                 `xml:"location,attr,omitempty" json:"location,omitempty" yaml:"location,omitempty"`
	City                string                 `xml:"city,attr,omitempty" json:"city,omitempty" yaml:"city,omitempty"`
	Building            string                 `xml:"building,attr,omitempty" json:"building,omitempty" yaml:"building,omitempty"`
	ParentForeignSource string                 `xml:"parent-foreign-source,attr,omitempty" json:"parent-foreign-source,omitempty" yaml:"parentForeignSource,omitempty"`
	ParentForeignID     string                 `xml:"parent-foreign-id,attr,omitempty" json:"parent-foreign-id,omitempty" yaml:"parentForeignID,omitempty"`
	ParentNodeLabel     string                 `xml:"parent-node-label,omitempty" json:"parent-node-label,omitempty" yaml:"parentNodeLabel,omitempty"`
	Interfaces          []RequisitionInterface `xml:"interface,omitempty" json:"interface,omitempty" yaml:"interfaces,omitempty"`
	Categories          []RequisitionCategory  `xml:"category,omitempty" json:"category,omitempty" yaml:"categories,omitempty"`
	Assets              []RequisitionAsset     `xml:"asset,omitempty" json:"asset,omitempty" yaml:"assets,omitempty"`
	MetaData            []RequisitionMetaData  `xml:"meta-data,omitempty" json:"meta-data,omitempty" yaml:"metaData,omitempty"`
}

// IsValid returns an error if the node definition is invalid
func (n *RequisitionNode) IsValid() error {
	if n.ForeignID == "" {
		return fmt.Errorf("Foreign ID cannot be empty")
	}
	if n.NodeLabel == "" { // Set a reasonable default when the label is not initialized
		n.NodeLabel = n.ForeignID
	}
	if n.ParentForeignID != "" && n.ParentNodeLabel != "" {
		return fmt.Errorf("Cannot set both parent foreign ID and parent node label, choose one")
	}
	if n.ParentNodeLabel == n.NodeLabel {
		return fmt.Errorf("The parent node cannot be the node itself. The parent-nodel-label has to be different than the node-label")
	}
	if n.ParentForeignID == n.ForeignID {
		return fmt.Errorf("The parent node cannot be the node itself. The parent-foreign-id has to be different than the foreign-id")
	}
	primaryCount := 0
	intfMap := make(map[string]int)
	for i := range n.Interfaces {
		intf := &n.Interfaces[i]
		intfMap[intf.IPAddress]++
		if intf.SnmpPrimary == "P" {
			primaryCount++
		}
		err := intf.IsValid()
		if err != nil {
			return err
		}
	}
	for ipAddr, count := range intfMap {
		if count > 1 {
			return fmt.Errorf("IP Address %s is defined more than once on node %s", ipAddr, n.NodeLabel)
		}
	}
	if primaryCount > 1 {
		return fmt.Errorf("Node %s cannot have more than one primary interface", n.NodeLabel)
	}
	for _, c := range n.Categories {
		err := c.IsValid()
		if err != nil {
			return err
		}
	}
	for _, a := range n.Assets {
		err := a.IsValid()
		if err != nil {
			return err
		}
	}
	return nil
}

// Requisition a requisition or set of nodes
type Requisition struct {
	XMLName    xml.Name          `xml:"model-import" json:"-" yaml:"-"`
	DateStamp  *Time             `xml:"date-stamp,attr,omitempty" json:"date-stamp,omitempty" yaml:"dateStamp,omitempty"`
	LastImport *Time             `xml:"last-import,attr,omitempty" json:"last-import,omitempty" yaml:"lastImport,omitempty"`
	Name       string            `xml:"foreign-source,attr" json:"foreign-source" yaml:"name"`
	Nodes      []RequisitionNode `xml:"node,omitempty" json:"node,omitempty" yaml:"nodes,omitempty"`
}

// IsValid returns an error if the requisition definition is invalid
func (r Requisition) IsValid() error {
	if r.Name == "" {
		return fmt.Errorf("Requisition name cannot be null")
	}
	foreignIDs := make(map[string]int)
	for i := range r.Nodes {
		n := &r.Nodes[i]
		foreignIDs[n.ForeignID]++
		err := n.IsValid()
		if err != nil {
			return fmt.Errorf("Problem on node %s on requisition %s: %s", n.NodeLabel, r.Name, err.Error())
		}
	}
	for id, count := range foreignIDs {
		if count > 1 {
			return fmt.Errorf("Duplicate Foreign ID %s on requisition %s", id, r.Name)
		}
	}
	return nil
}

// RequisitionsList a list of requisitions names
type RequisitionsList struct {
	Count          int      `json:"count" yaml:"count"`
	ForeignSources []string `json:"foreign-source" yaml:"foreignSources"`
}

// RequisitionStats statistics about the requisition
type RequisitionStats struct {
	Name       string   `json:"name" yaml:"name"`
	Count      int      `json:"count" yaml:"count"`
	ForeignIDs []string `json:"foreign-id" yaml:"foreignID"`
	LastImport *Time    `json:"last-imported,omitempty" yaml:"lastImport,omitempty"`
}

// RequisitionsStats statistics about all the requisitions
type RequisitionsStats struct {
	Count          int                `json:"count"`
	ForeignSources []RequisitionStats `json:"foreign-source"`
}

// ElementList a list of elements/strings
type ElementList struct {
	Count   int      `json:"count" yaml:"count"`
	Element []string `json:"element" yaml:"element"`
}

// Parameter a parameter for a detector or a policy
type Parameter struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// Detector a provisioning detector
type Detector struct {
	Name       string      `json:"name" yaml:"name"`
	Class      string      `json:"class" yaml:"class"`
	Parameters []Parameter `json:"parameter,omitempty" yaml:"parameters,omitempty"`
}

// Policy a provisioning policy
type Policy struct {
	Name       string      `json:"name" yaml:"name"`
	Class      string      `json:"class" yaml:"class"`
	Parameters []Parameter `json:"parameter,omitempty" yaml:"parameters,omitempty"`
}

// ForeignSourceDef a foreign source definition
type ForeignSourceDef struct {
	Name         string     `json:"name" yaml:"name"`
	DateStamp    *Time      `json:"date-stamp,omitempty" yaml:"dateStamp,omitempty"`
	ScanInterval string     `json:"scan-interval" yaml:"scanInterval"`
	Detectors    []Detector `json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Policies     []Policy   `json:"policies,omitempty" yaml:"policies,omitempty"`
}

// Plugin a definiton class for a detector or a policy
type Plugin struct {
	Name       string        `json:"name" yaml:"name"`
	Class      string        `json:"class" yaml:"class"`
	Parameters []PluginParam `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// PluginParam a parameter of a given plugin
type PluginParam struct {
	Key      string   `json:"key" yaml:"key"`
	Required bool     `json:"required" yaml:"required"`
	Options  []string `json:"options,omitempty" yaml:"options,omitempty"`
}

// PluginList a list of plugins
type PluginList struct {
	Count   int      `json:"count" yaml:"count"`
	Plugins []Plugin `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}
