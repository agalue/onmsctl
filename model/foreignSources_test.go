package model

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

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
	assert.ErrorContains(t, fsDef.Validate(), "invalid scan interval")

	fsDef.ScanInterval = "2w 1d"
	assert.NilError(t, fsDef.Validate())

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
