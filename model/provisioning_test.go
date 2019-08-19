package model

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestRequisitionObject(t *testing.T) {
	req := &Requisition{
		Name:      "Test",
		DateStamp: &Time{time.Now()},
		Nodes: []RequisitionNode{
			{
				ForeignID: "opennms.com",
				Interfaces: []RequisitionInterface{
					{
						IPAddress:   "www.opennms.com",
						SnmpPrimary: "P",
						Services: []RequisitionMonitoredService{
							{
								Name: "HTTPS",
							},
						},
					},
				},
				Assets: []RequisitionAsset{
					{
						Name:  "city",
						Value: "Apex",
					},
				},
				Categories: []RequisitionCategory{
					{
						Name: "Production",
					},
				},
			},
		},
	}
	var err error

	err = req.IsValid()
	assert.NilError(t, err)
	assert.Equal(t, req.Nodes[0].NodeLabel, req.Nodes[0].ForeignID)
	assert.Equal(t, req.Nodes[0].Interfaces[0].IPAddress, "34.194.50.139")

	bytes, err := json.MarshalIndent(req, "", "  ")
	assert.NilError(t, err)
	fmt.Println(string(bytes))
	assert.NilError(t, json.Unmarshal(bytes, &Requisition{}))

	bytes, err = xml.MarshalIndent(req, "", "  ")
	assert.NilError(t, err)
	fmt.Println(string(bytes))
	assert.NilError(t, xml.Unmarshal(bytes, &Requisition{}))

	bytes, err = yaml.Marshal(req)
	assert.NilError(t, err)
	fmt.Println(string(bytes))
	assert.NilError(t, yaml.Unmarshal(bytes, &Requisition{}))
}

func TestForeignSourceObject(t *testing.T) {
	fsDef := &ForeignSourceDef{
		Name: "Test",
		Detectors: []Detector{
			{
				Name:  "ICMP",
				Class: "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
			},
			{
				Name:  "SNMP",
				Class: "org.opennms.netmgt.provision.detector.snmp.SnmpDetector",
			},
		},
		Policies: []Policy{
			{
				Name:  "Production",
				Class: "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy",
				Parameters: []Parameter{
					{
						Key:   "category",
						Value: "Production",
					},
					{
						Key:   "matchBehavior",
						Value: "NO_PARAMETERS",
					},
				},
			},
		},
	}

	var err error

	fsDef.ScanInterval = "2YEARS" // This is wrong on purpose
	err = fsDef.IsValid()
	assert.ErrorContains(t, err, "Invalid scan interval")

	fsDef.ScanInterval = "2w 1d"
	err = fsDef.IsValid()
	assert.NilError(t, err)

	bytes, err := json.MarshalIndent(fsDef, "", "  ")
	assert.NilError(t, err)
	fmt.Println(string(bytes))
	assert.NilError(t, json.Unmarshal(bytes, &ForeignSourceDef{}))

	bytes, err = xml.MarshalIndent(fsDef, "", "  ")
	assert.NilError(t, err)
	fmt.Println(string(bytes))
	assert.NilError(t, xml.Unmarshal(bytes, &ForeignSourceDef{}))

	bytes, err = yaml.Marshal(fsDef)
	assert.NilError(t, err)
	fmt.Println(string(bytes))
	assert.NilError(t, yaml.Unmarshal(bytes, &ForeignSourceDef{}))
}

func TestRequisitionXML(t *testing.T) {
	reqXML := `
	<model-import xmlns="http://xmlns.opennms.org/xsd/config/model-import" date-stamp="2018-10-25T04:10:15.355-05:00" foreign-source="Cassandra" last-import="2018-10-25T04:10:21.944-05:00">
		<node foreign-id="cass01" node-label="cass01">
		 <interface descr="bond0" ip-addr="10.0.0.10" status="1" snmp-primary="P">
				<monitored-service service-name="ICMP"/>
				<monitored-service service-name="SNMP"/>
				<monitored-service service-name="JMX-Cassandra-Newts-I1"/>
				<monitored-service service-name="JMX-Cassandra-I1"/>
		 </interface>
		 <interface descr="bond0:1" ip-addr="10.0.0.11" status="1" snmp-primary="N">
				<monitored-service service-name="JMX-Cassandra-Newts-I2"/>
				<monitored-service service-name="JMX-Cassandra-I2"/>
		 </interface>
		 <interface descr="bond0:2" ip-addr="10.0.0.12" status="1" snmp-primary="N">
				<monitored-service service-name="JMX-Cassandra-Newts-I3"/>
				<monitored-service service-name="JMX-Cassandra-I3"/>
		 </interface>
		 <category name="Servers"/>
		 <category name="Cassandra"/>
		</node>
	</model-import>
	`
	req := &Requisition{}
	err := xml.Unmarshal([]byte(reqXML), req)
	assert.NilError(t, err)

	bytes, err := yaml.Marshal(req)
	assert.NilError(t, err)
	fmt.Println(string(bytes))

	assert.Equal(t, 1, len(req.Nodes))
	n := req.Nodes[0]
	assert.Equal(t, 3, len(n.Interfaces))
	i1 := n.Interfaces[0]
	assert.Equal(t, 4, len(i1.Services))

	err = req.IsValid()
	assert.NilError(t, err)
}

func TestInvalidRequisitionXML(t *testing.T) {
	reqXML := `
	<model-import xmlns="http://xmlns.opennms.org/xsd/config/model-import" date-stamp="2018-10-25T04:10:15.355-05:00" foreign-source="Test" last-import="2018-10-25T04:10:21.944-05:00">
		<node foreign-id="www.opennms.org" node-label="www.opennms.org">
		 <interface descr="eth0" ip-addr="www.opennms.org" status="1" snmp-primary="N"/>
		</node>
	</model-import>
	`
	req := &Requisition{}
	err := xml.Unmarshal([]byte(reqXML), req)
	assert.NilError(t, err)
	AllowFqdnOnRequisitionedInterfaces = false
	err = req.IsValid()
	assert.ErrorContains(t, err, "not a valid IPv4")
}
