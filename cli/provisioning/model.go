package provisioning

import (
	"fmt"
	"net"
)

// Meta a meta-data entry
type Meta struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// Service an IP interface monitored service
type Service struct {
	Name     string `json:"service-name" yaml:"name"`
	MetaData []Meta `json:"meta-data,omitempty" yaml:"metaData,omitempty"`
}

// IsValid returns an error if the service is invalid
func (s Service) IsValid() error {
	if s.Name == "" {
		return fmt.Errorf("Service name cannot be null")
	}
	return nil
}

// Asset a requisition node asset field
type Asset struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// IsValid returns an error if asset field is invalid
func (a Asset) IsValid() error {
	if a.Name == "" {
		return fmt.Errorf("Asset name cannot be empty")
	}
	if a.Value == "" {
		return fmt.Errorf("Asset value cannot be empty")
	}
	return nil
}

// Category a requisition node category
type Category struct {
	Name string `json:"name" yaml:"name"`
}

// IsValid returns an error if the category is invalid
func (c Category) IsValid() error {
	if c.Name == "" {
		return fmt.Errorf("Category name cannot be null")
	}
	return nil
}

// Interface an IP interface of a requisition node
type Interface struct {
	IPAddress   string    `json:"ip-addr" yaml:"ipAddress"`
	Description string    `json:"descr,omitempty" yaml:"description,omitempty"`
	SnmpPrimary string    `json:"snmp-primary" yaml:"snmpPrimary"`
	Status      int       `json:"status" yaml:"status"`
	Services    []Service `json:"monitored-service,omitempty" yaml:"services,omitempty"`
	MetaData    []Meta    `json:"meta-data,omitempty" yaml:"metaData,omitempty"`
}

// IsValid returns an error if the interface definition is invalid
func (i *Interface) IsValid() error {
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
		fmt.Printf("%s translates to %s, using the first entry.\n", i.IPAddress, addresses)
		if err != nil {
			return fmt.Errorf("Cannot get address from %s: %s", i.IPAddress, err)
		}
		i.IPAddress = addresses[0].String()
	}
	for _, s := range i.Services {
		err := s.IsValid()
		if err != nil {
			return err
		}
	}
	return nil
}

// Node a requisition node
type Node struct {
	NodeLabel           string      `json:"node-label" yaml:"nodeLabel"`
	ForeignID           string      `json:"foreign-id" yaml:"foreignID"`
	Location            string      `json:"location,omitempty" yaml:"location,omitempty"`
	City                string      `json:"city,omitempty" yaml:"city,omitempty"`
	Building            string      `json:"building,omitempty" yaml:"building,omitempty"`
	ParentForeignSource string      `json:"parent-foreign-source,omitempty" yaml:"parentForeignSource,omitempty"`
	ParentForeignID     string      `json:"parent-foreign-id,omitempty" yaml:"parentForeignID,omitempty"`
	ParentNodeLabel     string      `json:"parent-node-label,omitempty" yaml:"parentNodeLabel,omitempty"`
	Interfaces          []Interface `json:"interface,omitempty" yaml:"interfaces,omitempty"`
	Categories          []Category  `json:"category,omitempty" yaml:"categories,omitempty"`
	Assets              []Asset     `json:"asset,omitempty" yaml:"assets,omitempty"`
	MetaData            []Meta      `json:"meta-data,omitempty" yaml:"metaData,omitempty"`
}

// IsValid returns an error if the node definition is invalid
func (n *Node) IsValid() error {
	if n.ForeignID == "" {
		return fmt.Errorf("Foreign ID cannot be empty")
	}
	if n.NodeLabel == "" { // Set a reasonable default when the label is not initialized
		n.NodeLabel = n.ForeignID
	}
	if n.ParentForeignID != "" && n.ParentNodeLabel != "" {
		return fmt.Errorf("Cannot set both parent foreign ID and parent node label, choose one")
	}
	primaryCount := 0
	for i := range n.Interfaces {
		intf := &n.Interfaces[i]
		if intf.SnmpPrimary == "P" {
			primaryCount++
		}
		err := intf.IsValid()
		if err != nil {
			return err
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
	DateStamp  int64  `json:"date-stamp" yaml:"dateStamp"`
	LastImport int64  `json:"last-import" yaml:"lastImport"`
	Name       string `json:"foreign-source" yaml:"name"`
	Nodes      []Node `json:"node,omitempty" yaml:"nodes,omitempty"`
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
			return err
		}
	}
	for k, v := range foreignIDs {
		if v > 1 {
			return fmt.Errorf("Duplicate Foreign ID %s", k)
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
	LastImport int64    `json:"last-imported" yaml:"lastImport"`
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
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// Policy a provisioning policy
type Policy struct {
	Name       string      `json:"name" yaml:"name"`
	Class      string      `json:"class" yaml:"class"`
	Parameters []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// ForeignSourceDef a foreign source definition
type ForeignSourceDef struct {
	Name         string     `json:"name" yaml:"name"`
	DateStamp    int64      `json:"date-stamp" yaml:"dateStamp"`
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
