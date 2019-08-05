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

	bytes, err = xml.MarshalIndent(req, "", "  ")
	assert.NilError(t, err)
	fmt.Println(string(bytes))

	bytes, err = yaml.Marshal(req)
	assert.NilError(t, err)
	fmt.Println(string(bytes))
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

	bytes, err = xml.MarshalIndent(fsDef, "", "  ")
	assert.NilError(t, err)
	fmt.Println(string(bytes))

	bytes, err = yaml.Marshal(fsDef)
	assert.NilError(t, err)
	fmt.Println(string(bytes))
}
