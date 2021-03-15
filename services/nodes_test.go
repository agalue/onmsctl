package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/OpenNMS/onmsctl/model"

	"gotest.tools/assert"
)

var mockNode = &model.OnmsNode{
	Label:          "test01.example.com",
	ForeignSource:  "Test",
	ForeignID:      "test01",
	SysName:        "test01",
	SysObjectID:    ".1.3.6.1.4.1.8072.3.2.10",
	SysLocation:    "Go Test",
	SysDescription: "Mock Node",
	SNMPInterfaces: []model.OnmsSnmpInterface{
		{
			IfIndex:       1097,
			IfName:        "eth0",
			IfDescr:       "eth0",
			IfSpeed:       1000000000,
			IfType:        6,
			IfOperStatus:  1,
			IfAdminStatus: 1,
		},
		{
			IfIndex:       1098,
			IfName:        "eth1",
			IfDescr:       "eth1",
			IfSpeed:       100000000,
			IfType:        6,
			IfOperStatus:  1,
			IfAdminStatus: 1,
		},
	},
	IPInterfaces: []model.OnmsIPInterface{
		{
			IPAddress:   "10.0.0.1",
			HostName:    "test01.example.com",
			SnmpPrimary: "P",
			IfIndex:     1097,
			Services: []model.OnmsMonitoredService{
				{
					ServiceType: &model.OnmsServiceType{
						Name: "HTTP",
					},
					Meta: []model.MetaData{
						{
							Context: "web",
							Key:     "port",
							Value:   "8080",
						},
					},
				},
			},
			Meta: []model.MetaData{
				{
					Context: "web",
					Key:     "environment",
					Value:   "prod",
				},
			},
		},
		{
			IPAddress:   "10.0.0.2",
			HostName:    "10.0.0.2",
			SnmpPrimary: "N",
			IfIndex:     1098,
		},
	},
	Meta: []model.MetaData{
		{
			Context: "web",
			Key:     "owner",
			Value:   "agalue",
		},
	},
	Categories: []model.OnmsCategory{
		{
			Name: "Production",
		},
	},
}

type mockNodeRest struct {
	t *testing.T
}

func (api mockNodeRest) Get(path string) ([]byte, error) {
	log.Printf("GET PATH: %s", path)
	if strings.Contains(path, "/api/v2/nodes?limit=0&_s=") {
		return nil, nil
	}
	if strings.Contains(path, "/nodes?limit=10") {
		list := createMockNodeList(48, getOffset(path))
		return json.Marshal(list)
	}
	if strings.Contains(path, "/ipinterfaces?limit=10") {
		list := createMockIPList(48, getOffset(path))
		return json.Marshal(list)
	}
	if strings.Contains(path, "/snmpinterfaces?limit=10") {
		list := createMockSNMPList(48, getOffset(path))
		return json.Marshal(list)
	}
	if strings.HasSuffix(path, "/nodes/10/snmpinterfaces/1097") {
		snmp := mockNode.GetSnmpInterface(1097).ExtractBasic()
		snmp.ID = 101 // Fake Primary Key
		return json.Marshal(snmp)
	}
	if strings.HasSuffix(path, "/nodes/10/snmpinterfaces/1098") {
		snmp := mockNode.GetSnmpInterface(1098).ExtractBasic()
		snmp.ID = 102 // Fake Primary Key
		return json.Marshal(snmp)
	}
	if strings.HasSuffix(path, "/nodes/10/ipinterfaces/10.0.0.1") {
		ip := mockNode.GetIPInterface("10.0.0.1")
		ip.ID = "11" // Fake Primary Key
		return json.Marshal(ip)
	}
	return nil, fmt.Errorf("GET should not be called for %s", path)
}

func (api mockNodeRest) Post(path string, jsonBytes []byte) error {
	log.Printf("POST PATH: %s, JSON: %s", path, string(jsonBytes))
	if strings.HasSuffix(path, "/nodes/10/categories") {
		cat := &model.OnmsCategory{}
		if err := json.Unmarshal(jsonBytes, cat); err != nil {
			return err
		}
		if cat.Name != "Production" {
			return fmt.Errorf("Received invalid category %s", string(jsonBytes))
		}
		return nil
	}
	if strings.HasSuffix(path, "/nodes/10/ipinterfaces") {
		intf := &model.OnmsIPInterface{}
		if err := json.Unmarshal(jsonBytes, intf); err != nil {
			return err
		}
		if mockNode.GetIPInterface(intf.IPAddress) == nil {
			return fmt.Errorf("Received IP Interface with invalid IP %s", string(jsonBytes))
		}
		if intf.IfIndex > 0 && (intf.SNMPInterface == nil || intf.SNMPInterface.ID == 0) {
			return fmt.Errorf("Received IP Interface with invalid ifIndex %s", string(jsonBytes))
		}
		return nil
	}
	if strings.HasSuffix(path, "/nodes/10/ipinterfaces/10.0.0.1/services") {
		svc := &model.OnmsMonitoredService{}
		if err := json.Unmarshal(jsonBytes, svc); err != nil {
			return err
		}
		if svc.ServiceType.Name != "HTTP" {
			return fmt.Errorf("Received invalid Service %s", string(jsonBytes))
		}
		return nil
	}
	if strings.HasSuffix(path, "/metadata") {
		meta := &model.MetaData{}
		if err := json.Unmarshal(jsonBytes, meta); err != nil {
			return err
		}
		if meta.Context != "web" {
			return fmt.Errorf("Received invalid metadata %s", string(jsonBytes))
		}
		return nil
	}
	return fmt.Errorf("POST should not be called for %s", path)
}

func (api mockNodeRest) PostRaw(path string, dataBytes []byte, contentType string) (*http.Response, error) {
	log.Printf("POST PATH (raw): %s, JSON: %s", path, string(dataBytes))
	if strings.HasSuffix(path, "/nodes") {
		response := &http.Response{
			Status:     http.StatusText(http.StatusCreated),
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Location": []string{"/nodes/10"},
			},
		}
		return response, nil
	}
	if strings.HasSuffix(path, "/nodes/10/snmpinterfaces") {
		intf := &model.OnmsSnmpInterface{}
		if err := xml.Unmarshal(dataBytes, intf); err != nil {
			return &http.Response{}, err
		}
		return &http.Response{
			Status:     http.StatusText(http.StatusCreated),
			StatusCode: http.StatusCreated,
			Header: http.Header{
				"Location": []string{fmt.Sprintf("/nodes/10/snmpinterfaces/%d", intf.IfIndex)},
			},
		}, nil
	}
	return nil, fmt.Errorf("POST-Raw should not be called for %s", path)
}

func (api mockNodeRest) Delete(path string) error {
	return fmt.Errorf("DELETE should not be called for %s", path)
}

func (api mockNodeRest) Put(path string, dataBytes []byte, contentType string) error {
	return fmt.Errorf("should not be called")
}

func (api mockNodeRest) IsValid(r *http.Response) error {
	return nil
}

func getOffset(path string) int {
	offset := 0
	re := regexp.MustCompile(`offset=(\d+)`)
	match := re.FindStringSubmatch(path)
	if match != nil {
		offset, _ = strconv.Atoi(match[1])
	}
	return offset
}

func createMockNodeList(total int, offset int) *model.OnmsNodeList {
	list := &model.OnmsNodeList{
		TotalCount: total,
		Offset:     offset,
	}
	for i := offset; i < offset+defaultLimit; i++ {
		if i < total {
			list.Nodes = append(list.Nodes, model.OnmsNode{Label: fmt.Sprintf("node%d", i)})
		}
	}
	list.Count = len(list.Nodes)
	return list
}

func createMockIPList(total int, offset int) *model.OnmsIPInterfaceList {
	list := &model.OnmsIPInterfaceList{
		TotalCount: total,
		Offset:     offset,
	}
	for i := offset; i < offset+defaultLimit; i++ {
		if i < total {
			list.Interfaces = append(list.Interfaces, model.OnmsIPInterface{IPAddress: fmt.Sprintf("10.0.%d.1", i)})
		}
	}
	list.Count = len(list.Interfaces)
	return list
}

func createMockSNMPList(total int, offset int) *model.OnmsSnmpInterfaceList {
	list := &model.OnmsSnmpInterfaceList{
		TotalCount: total,
		Offset:     offset,
	}
	for i := offset; i < offset+defaultLimit; i++ {
		if i < total {
			list.Interfaces = append(list.Interfaces, model.OnmsSnmpInterface{IfIndex: 1000 + i, IfName: fmt.Sprintf("eth%d", i)})
		}
	}
	list.Count = len(list.Interfaces)
	return list
}

func TestCriteria(t *testing.T) {
	var err error
	api := nodesAPI{}
	assert.NilError(t, api.isCriteriaValid("42"))
	assert.NilError(t, api.isCriteriaValid("routers:r01"))
	err = api.isCriteriaValid("routers:r01:b")
	assert.Assert(t, err != nil)
	err = api.isCriteriaValid("routers")
	assert.Assert(t, err != nil)
}

func TestAddNode(t *testing.T) {
	api := GetNodesAPI(&mockNodeRest{t})
	err := api.AddNode(mockNode)
	assert.NilError(t, err)
}

func TestGetNodes(t *testing.T) {
	api := GetNodesAPI(&mockNodeRest{t})
	list, err := api.GetNodes()
	assert.NilError(t, err)
	for _, n := range list.Nodes {
		log.Println(n.Label)
	}
	assert.Equal(t, 48, len(list.Nodes))
}

func TestGetIPInterfaces(t *testing.T) {
	api := GetNodesAPI(&mockNodeRest{t})
	list, err := api.GetIPInterfaces("1")
	assert.NilError(t, err)
	for _, n := range list.Interfaces {
		log.Println(n.IPAddress)
	}
	assert.Equal(t, 48, len(list.Interfaces))
}

func TestGetSnmpInterfaces(t *testing.T) {
	api := GetNodesAPI(&mockNodeRest{t})
	list, err := api.GetSnmpInterfaces("1")
	assert.NilError(t, err)
	for _, n := range list.Interfaces {
		log.Println(n.IfName)
	}
	assert.Equal(t, 48, len(list.Interfaces))
}
