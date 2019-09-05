package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestResources(t *testing.T) {
	jsonBytes := []byte(`
  {
    "id": "node[NetworkEquipment:1506541767935].interfaceSnmp[vlan_12-a80000000000]",
    "label": "vlan.12 (wireless_network, 10.0.0.1)",
    "name": "vlan_12-a80000000000",
    "link": "element/interface.jsp?ipinterfaceid=366",
    "typeLabel": "SNMP Interface Data",
    "parentId": "node[NetworkEquipment:1506541767935]",
    "stringPropertyAttributes": {
      "ifHighSpeed": "0",
      "ifName": "vlan.12"
    },
    "externalValueAttributes": {
      "ifIndex": "514",
      "ifSpeed": "0",
      "hasFlows": "false",
      "ifSpeedFriendly": "0 bps",
      "nodeId": "4"
    },
    "rrdGraphAttributes": {
      "ifHCOutUcastPkts": {
        "name": "ifHCOutUcastPkts",
        "relativePath": "snmp/4/vlan_12-a80000000000",
        "rrdFile": "ifHCOutUcastPkts.jrb"
      },
      "ifInDiscards": {
        "name": "ifInDiscards",
        "relativePath": "snmp/4/vlan_12-a80000000000",
        "rrdFile": "ifInDiscards.jrb"
      }
    }
  }  
  `)
	resource := &Resource{}
	err := json.Unmarshal(jsonBytes, resource)
	assert.NilError(t, err)
	assert.Equal(t, "vlan_12-a80000000000", resource.Name)
	assert.Equal(t, 2, len(resource.StringAttributes))
	assert.Equal(t, "vlan.12", resource.StringAttributes["ifName"])
	assert.Equal(t, 5, len(resource.ExternalAttributes))
	assert.Equal(t, "514", resource.ExternalAttributes["ifIndex"])
	assert.Equal(t, 2, len(resource.NumericAttributes))
	assert.Equal(t, "ifInDiscards.jrb", resource.NumericAttributes["ifInDiscards"].RrdFile)

	yamlBytes, err := yaml.Marshal(resource)
	assert.NilError(t, err)
	fmt.Println(string(yamlBytes))
}
