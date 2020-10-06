package services

import (
	"encoding/json"
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
			IfIndex:       99,
			IfName:        "eth0",
			IfDescr:       "eth0",
			IfSpeed:       1000000000,
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
			IfIndex:     99,
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
	if strings.Contains(path, "/nodes?limit=10") {
		offset := 0
		re := regexp.MustCompile(`offset=(\d+)`)
		match := re.FindStringSubmatch(path)
		if match != nil {
			offset, _ = strconv.Atoi(match[1])
		}
		list := createMockNodeList(48, offset)
		return json.Marshal(list)
	}
	if strings.HasSuffix(path, "/nodes/10/snmpinterfaces/99") {
		snmp := mockNode.GetSnmpInterface(99).ExtractBasic()
		snmp.ID = 11 // Fake Primary Key
		return json.Marshal(snmp)
	}
	return nil, fmt.Errorf("GET should not be called for %s", path)
}

func (api mockNodeRest) Post(path string, jsonBytes []byte) error {
	log.Printf("PATH: %s, JSON: %s", path, string(jsonBytes))
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
			return fmt.Errorf("Received invalid IP Interface %s", string(jsonBytes))
		}
		if intf.IfIndex > 0 && (intf.SNMPInterface == nil || intf.SNMPInterface.ID != 11) {
			return fmt.Errorf("Received invalid IP Interface %s", string(jsonBytes))
		}
		return nil
	}
	if strings.HasSuffix(path, "/nodes/10/snmpinterfaces") {
		intf := &model.OnmsSnmpInterface{}
		if err := json.Unmarshal(jsonBytes, intf); err != nil {
			return err
		}
		if mockNode.GetSnmpInterface(intf.IfIndex) == nil {
			return fmt.Errorf("Received invalid SNMP Interface %s", string(jsonBytes))
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

func (api mockNodeRest) PostRaw(path string, jsonBytes []byte) (*http.Response, error) {
	log.Printf("PATH: %s, JSON: %s", path, string(jsonBytes))
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
	return nil, fmt.Errorf("POST-Raw should not be called for %s", path)
}

func (api mockNodeRest) Delete(path string) error {
	return fmt.Errorf("DELETE should not be called for %s", path)
}

func (api mockNodeRest) Put(path string, jsonBytes []byte, contentType string) error {
	return fmt.Errorf("should not be called")
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
