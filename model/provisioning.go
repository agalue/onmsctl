package model

import (
	"encoding/xml"
	"fmt"
	"net"
	"regexp"

	"github.com/imdario/mergo"
)

// AllowFqdnOnRequisitionedInterfaces when this is true, if the content of an IP Address is a FQDN it will be translated into a valid IPv4
var AllowFqdnOnRequisitionedInterfaces = true

// RequisitionMetaData a meta-data entry
type RequisitionMetaData struct {
	XMLName xml.Name `xml:"meta-data" json:"-" yaml:"-"`
	Key     string   `xml:"key,attr" json:"key" yaml:"key"`
	Value   string   `xml:"value,attr" json:"value" yaml:"value"`
	Context string   `xml:"context,attr,omitempty" json:"context,omitempty" yaml:"context,omitempty"`
}

// Validate returns an error if asset field is invalid
func (m *RequisitionMetaData) Validate() error {
	if m.Context == "" {
		m.Context = "requisition"
	}
	if m.Key == "" {
		return fmt.Errorf("Meta-data key cannot be empty")
	}
	if m.Value == "" {
		return fmt.Errorf("Meta-data value for key %s cannot be empty", m.Key)
	}
	return nil
}

// RequisitionMonitoredService an IP interface monitored service
type RequisitionMonitoredService struct {
	XMLName  xml.Name              `xml:"monitored-service" json:"-" yaml:"-"`
	Name     string                `xml:"service-name,attr" json:"service-name" yaml:"name"`
	MetaData []RequisitionMetaData `xml:"meta-data,omitempty" json:"meta-data,omitempty" yaml:"metaData,omitempty"`
}

// AddMetaData adds a meta-data entry to the node
func (s *RequisitionMonitoredService) AddMetaData(key string, value string) {
	s.MetaData = append(s.MetaData, RequisitionMetaData{Context: "requisition", Key: key, Value: value})
}

// SetMetaData adds or updates an existing meta-data entry on the service
func (s *RequisitionMonitoredService) SetMetaData(key string, value string) {
	var found *RequisitionMetaData
	for i := range s.MetaData {
		m := &s.MetaData[i]
		if m.Key == key {
			found = m
		}
	}
	if found == nil {
		s.AddMetaData(key, value)
	} else {
		found.Value = value
	}
}

// DeleteMetaData deletes an existing meta-data entry from the service
func (s *RequisitionMonitoredService) DeleteMetaData(key string) {
	var found int
	for i, m := range s.MetaData {
		if m.Key == key {
			found = i
		}
	}
	if len(s.MetaData) > 0 {
		s.MetaData = append(s.MetaData[:found], s.MetaData[found+1:]...)
	}
}

// Validate returns an error if the service is invalid
func (s RequisitionMonitoredService) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("Service name cannot be empty")
	}
	if matched, _ := regexp.MatchString(`[/\\?:&*'"]`, s.Name); matched {
		return fmt.Errorf("Invalid characters on service name %s:, /, \\, ?, &, *, ', \"", s.Name)
	}
	for i := range s.MetaData {
		m := &s.MetaData[i]
		err := m.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// RequisitionAsset a requisition node asset field
type RequisitionAsset struct {
	XMLName xml.Name `xml:"asset" json:"-" yaml:"-"`
	Name    string   `xml:"name,attr" json:"name" yaml:"name"`
	Value   string   `xml:"value,attr" json:"value" yaml:"value"`
}

// Validate returns an error if asset field is invalid
func (a RequisitionAsset) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("Asset name cannot be empty")
	}
	if matched, _ := regexp.MatchString(`[/\\?:&*'"]`, a.Name); matched {
		return fmt.Errorf("Invalid characters on asset name %s:, /, \\, ?, &, *, ', \"", a.Name)
	}
	if a.Value == "" {
		return fmt.Errorf("Asset value for %s cannot be empty", a.Name)
	}
	return nil
}

// RequisitionCategory a requisition node category
type RequisitionCategory struct {
	XMLName xml.Name `xml:"category" json:"-" yaml:"-"`
	Name    string   `xml:"name,attr" json:"name" yaml:"name"`
}

// Validate returns an error if the category is invalid
func (c RequisitionCategory) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("Category name cannot be empty")
	}
	if matched, _ := regexp.MatchString(`[/\\?:&*'"]`, c.Name); matched {
		return fmt.Errorf("Invalid characters on category name %s:, /, \\, ?, &, *, ', \"", c.Name)
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

// AddMetaData adds a meta-data entry to the interface
func (intf *RequisitionInterface) AddMetaData(key string, value string) {
	intf.MetaData = append(intf.MetaData, RequisitionMetaData{Context: "requisition", Key: key, Value: value})
}

// SetMetaData adds or updates an existing meta-data entry on the IP interface
func (intf *RequisitionInterface) SetMetaData(key string, value string) {
	var found *RequisitionMetaData
	for i := range intf.MetaData {
		m := &intf.MetaData[i]
		if m.Key == key {
			found = m
		}
	}
	if found == nil {
		intf.AddMetaData(key, value)
	} else {
		found.Value = value
	}
}

// DeleteMetaData deletes an existing meta-data entry from the IP interface
func (intf *RequisitionInterface) DeleteMetaData(key string) {
	var found int
	for i, m := range intf.MetaData {
		if m.Key == key {
			found = i
		}
	}
	if len(intf.MetaData) > 0 {
		intf.MetaData = append(intf.MetaData[:found], intf.MetaData[found+1:]...)
	}
}

// GetService gets a given monitored service from the IP interface if exists
func (intf *RequisitionInterface) GetService(serviceName string) *RequisitionMonitoredService {
	var service *RequisitionMonitoredService
	for _, svc := range intf.Services {
		if svc.Name == serviceName {
			service = &svc
			break
		}
	}
	return service
}

// AddService add an service to the IP interface
func (intf *RequisitionInterface) AddService(svc *RequisitionMonitoredService) {
	intf.Services = append(intf.Services, *svc)
}

// Validate returns an error if the interface definition is invalid
func (intf *RequisitionInterface) Validate() error {
	if intf.IPAddress == "" {
		return fmt.Errorf("IP Address cannot be empty")
	}
	if intf.Status == 0 { // Set a reasonable default when the status is not initialized
		intf.Status = 1
	}
	if intf.Status != 1 && intf.Status != 3 {
		return fmt.Errorf("Invalid status for interface %s: %d", intf.IPAddress, intf.Status)
	}
	if intf.SnmpPrimary == "" { // Set a reasonable default when the primary flag is not initialized
		intf.SnmpPrimary = "N"
	}
	if intf.SnmpPrimary != "P" && intf.SnmpPrimary != "S" && intf.SnmpPrimary != "N" {
		return fmt.Errorf("Invalid snmp-primary for interface %s: %s", intf.IPAddress, intf.SnmpPrimary)
	}
	if err := intf.validateIP(); err != nil {
		return err
	}
	if err := intf.validateServices(); err != nil {
		return err
	}
	for i := range intf.MetaData {
		m := intf.MetaData[i]
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Merge merges the fields from the provided source
func (intf *RequisitionInterface) Merge(source RequisitionInterface) error {
	return mergo.Merge(intf, source, mergo.WithOverride)
}

func (intf *RequisitionInterface) validateIP() error {
	ip := net.ParseIP(intf.IPAddress)
	if ip == nil {
		if AllowFqdnOnRequisitionedInterfaces {
			addresses, err := net.LookupIP(intf.IPAddress)
			if err != nil || len(addresses) == 0 {
				return fmt.Errorf("Cannot get address from %s (invalid IP or FQDN); %s", intf.IPAddress, err)
			}
			fmt.Printf("%s translates to %s.\n", intf.IPAddress, addresses[0].String())
			intf.IPAddress = addresses[0].String()
		} else {
			return fmt.Errorf("%s is not a valid IPv4 or IPv6 address", intf.IPAddress)
		}
	}
	return nil
}

func (intf *RequisitionInterface) validateServices() error {
	serviceMap := make(map[string]int)
	for i := range intf.Services {
		s := &intf.Services[i]
		if err := s.Validate(); err != nil {
			return err
		}
		serviceMap[s.Name]++
	}
	for service, count := range serviceMap {
		if count > 1 {
			return fmt.Errorf("Service %s is defined more than once on interface %s", service, intf.IPAddress)
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

// AddMetaData adds a meta-data entry to the node
func (n *RequisitionNode) AddMetaData(key string, value string) {
	n.MetaData = append(n.MetaData, RequisitionMetaData{Context: "requisition", Key: key, Value: value})
}

// SetMetaData adds or updates an existing meta-data entry on the node
func (n *RequisitionNode) SetMetaData(key string, value string) {
	var found *RequisitionMetaData
	for i := range n.MetaData {
		m := &n.MetaData[i]
		if m.Key == key {
			found = m
		}
	}
	if found == nil {
		n.AddMetaData(key, value)
	} else {
		found.Value = value
	}
}

// DeleteMetaData deletes an existing meta-data entry from the node
func (n *RequisitionNode) DeleteMetaData(key string) {
	var found int
	for i, m := range n.MetaData {
		if m.Key == key {
			found = i
		}
	}
	if len(n.MetaData) > 0 {
		n.MetaData = append(n.MetaData[:found], n.MetaData[found+1:]...)
	}
}

// AddInterface add an  IP interface to the node
func (n *RequisitionNode) AddInterface(intf *RequisitionInterface) {
	n.Interfaces = append(n.Interfaces, *intf)
}

// GetInterface gets an existing IP interface from the node
func (n *RequisitionNode) GetInterface(ipAddress string) *RequisitionInterface {
	for _, intf := range n.Interfaces {
		if intf.IPAddress == ipAddress {
			return &intf
		}
	}
	return nil
}

// Merge merges the fields from the provided source
func (n *RequisitionNode) Merge(source RequisitionNode) error {
	return mergo.Merge(n, source, mergo.WithOverride)
}

// Validate returns an error if the node definition is invalid
func (n *RequisitionNode) Validate() error {
	if n.ForeignID == "" {
		return fmt.Errorf("Foreign ID cannot be empty")
	}
	if matched, _ := regexp.MatchString(`[/\\?:&*'"]`, n.ForeignID); matched {
		return fmt.Errorf("Invalid characters on Foreign ID %s:, /, \\, ?, &, *, ', \"", n.ForeignID)
	}
	if n.NodeLabel == "" { // Set a reasonable default when the label is not initialized
		n.NodeLabel = n.ForeignID
	}
	if n.ParentForeignID != "" && n.ParentNodeLabel != "" {
		return fmt.Errorf("Cannot set both parent foreign ID and parent node label on node %s, choose one", n.NodeLabel)
	}
	if n.ParentNodeLabel == n.NodeLabel {
		return fmt.Errorf("The parent node cannot be the node itself. The parent-nodel-label has to be different than the node-label")
	}
	if n.ParentForeignID == n.ForeignID {
		return fmt.Errorf("The parent node cannot be the node itself. The parent-foreign-id has to be different than the foreign-id")
	}
	if err := n.validateInterfaces(); err != nil {
		return err
	}
	for i := range n.Categories {
		c := &n.Categories[i]
		if err := c.Validate(); err != nil {
			return err
		}
	}
	for i := range n.Assets {
		a := &n.Assets[i]
		if err := a.Validate(); err != nil {
			return err
		}
	}
	for i := range n.MetaData {
		m := &n.MetaData[i]
		if err := m.Validate(); err != nil {
			return err
		}
	}
	for i := range n.Interfaces {
		intf := &n.Interfaces[i]
		if err := intf.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (n *RequisitionNode) validateInterfaces() error {
	primaryCount := 0
	intfMap := make(map[string]int)
	for i := range n.Interfaces {
		intf := &n.Interfaces[i]
		intfMap[intf.IPAddress]++
		if intf.SnmpPrimary == "P" {
			primaryCount++
		}
		if err := intf.Validate(); err != nil {
			return err
		}
	}
	if primaryCount > 1 {
		return fmt.Errorf("Node %s cannot have more than one primary interface", n.NodeLabel)
	}
	for ipAddr, count := range intfMap {
		if count > 1 {
			return fmt.Errorf("IP Address %s is defined more than once on node %s", ipAddr, n.NodeLabel)
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

// AddNode add a node to the requisition
func (r *Requisition) AddNode(node *RequisitionNode) {
	r.Nodes = append(r.Nodes, *node)
}

// Validate returns an error if the requisition definition is invalid
func (r *Requisition) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("Requisition name cannot be empty")
	}
	if matched, _ := regexp.MatchString(`[/\\?:&*'"]`, r.Name); matched {
		return fmt.Errorf("Invalid characters on requisition name %s:, /, \\, ?, &, *, ', \"", r.Name)
	}
	foreignIDs := make(map[string]int)
	for i := range r.Nodes {
		n := &r.Nodes[i]
		foreignIDs[n.ForeignID]++
		err := n.Validate()
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

// GetRequisitionStats gets the stats of a given requisition
func (stats RequisitionsStats) GetRequisitionStats(foreignSource string) RequisitionStats {
	for _, req := range stats.ForeignSources {
		if req.Name == foreignSource {
			return req
		}
	}
	return RequisitionStats{}
}
