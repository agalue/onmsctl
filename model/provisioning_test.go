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

var testNode = RequisitionNode{
	ForeignID: "opennms.com",
	Interfaces: []RequisitionInterface{
		{
			IPAddress:   "www.opennms.com",
			SnmpPrimary: "P",
			Services: []RequisitionMonitoredService{
				{
					Name: "HTTPS",
					MetaData: []RequisitionMetaData{
						{
							Key:   "url",
							Value: "/index.html",
						},
					},
				},
			},
			MetaData: []RequisitionMetaData{
				{
					Key:   "mpls",
					Value: "false",
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
	MetaData: []RequisitionMetaData{
		{
			Key:   "owner",
			Value: "agalue",
		},
	},
}

func TestRequisitionObject(t *testing.T) {
	req := &Requisition{
		Name:      "Test",
		DateStamp: &Time{time.Now()},
		Nodes:     []RequisitionNode{testNode},
	}
	var err error

	assert.NilError(t, req.Validate())
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

	assert.NilError(t, req.Validate())
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
	err = req.Validate()
	assert.ErrorContains(t, err, "not a valid IPv4")
}

func TestMetaData(t *testing.T) {
	node := &RequisitionNode{
		ForeignID: "n1",
		NodeLabel: "n1",
	}

	node.AddMetaData("k1", "v1")
	assert.Equal(t, 1, len(node.MetaData))
	assert.Equal(t, "v1", node.MetaData[0].Value)

	node.SetMetaData("k2", "v2")
	assert.Equal(t, 2, len(node.MetaData))
	assert.Equal(t, "v2", node.MetaData[1].Value)
	node.SetMetaData("k2", "v20")
	assert.Equal(t, "v20", node.MetaData[1].Value)

	assert.NilError(t, node.Validate())

	intf := &RequisitionInterface{
		IPAddress: "10.0.0.1",
	}

	intf.AddMetaData("k3", "v3")
	assert.Equal(t, 1, len(intf.MetaData))
	assert.Equal(t, "v3", intf.MetaData[0].Value)

	intf.SetMetaData("k4", "v4")
	assert.Equal(t, 2, len(intf.MetaData))
	assert.Equal(t, "v4", intf.MetaData[1].Value)
	intf.SetMetaData("k4", "v40")
	assert.Equal(t, "v40", intf.MetaData[1].Value)

	assert.NilError(t, intf.Validate())

	node.AddInterface(intf)

	intfP := node.GetInterface("10.0.0.1")
	assert.Assert(t, intfP != nil)

	svc := &RequisitionMonitoredService{
		Name: "HTTP",
	}

	svc.AddMetaData("k5", "v5")
	assert.Equal(t, 1, len(svc.MetaData))
	assert.Equal(t, "v5", svc.MetaData[0].Value)

	svc.SetMetaData("k6", "v6")
	assert.Equal(t, 2, len(svc.MetaData))
	assert.Equal(t, "v6", svc.MetaData[1].Value)
	svc.SetMetaData("k6", "v60")
	assert.Equal(t, "v60", svc.MetaData[1].Value)

	assert.NilError(t, svc.Validate())
}

func TestMergeNodes(t *testing.T) {
	source := RequisitionNode{
		ParentForeignSource: "Switch",
		ParentForeignID:     "SW01",
		MetaData: []RequisitionMetaData{
			{
				Key:   "important",
				Value: "true",
			},
		},
	}

	assert.Equal(t, "", testNode.ParentForeignSource)
	assert.Equal(t, "", testNode.ParentForeignID)
	assert.Equal(t, "owner", testNode.MetaData[0].Key)

	testNode.Merge(source)

	assert.Equal(t, "opennms.com", testNode.ForeignID)
	assert.Equal(t, "Switch", testNode.ParentForeignSource)
	assert.Equal(t, "SW01", testNode.ParentForeignID)
	assert.Equal(t, "important", testNode.MetaData[0].Key)
}
