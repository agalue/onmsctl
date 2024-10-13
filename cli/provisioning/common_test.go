package provisioning

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

var testNode = model.RequisitionNode{
	ForeignID: "n1",
	NodeLabel: "n1",
	Interfaces: []model.RequisitionInterface{
		{
			IPAddress:   "10.0.0.1",
			SnmpPrimary: "P",
			Services: []model.RequisitionMonitoredService{
				{
					Name: "HTTP",
					MetaData: []model.RequisitionMetaData{
						{Key: "url", Value: "/index.html"},
					},
				},
			},
			MetaData: []model.RequisitionMetaData{
				{Key: "mpls", Value: "false"},
			},
		},
	},
	Categories: []model.RequisitionCategory{
		{Name: "Server"},
	},
	Assets: []model.RequisitionAsset{
		{Name: "city", Value: "Durham"},
	},
	MetaData: []model.RequisitionMetaData{
		{Key: "owner", Value: "agalue"},
	},
}

var testForeignSource = model.ForeignSourceDef{
	Name: "Test",
	Detectors: []model.Detector{
		{
			Name:  "ICMP",
			Class: "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
		},
	},
	Policies: []model.Policy{
		{
			Name:  "Production",
			Class: "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy",
			Parameters: []model.Parameter{
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

func createTestServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("Received %s request from %s\n", req.Method, req.URL.Path)

		switch req.URL.Path {

		case "/rest/foreignSourcesConfig/policies":
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(test.PoliciesJSON))

		case "/rest/foreignSourcesConfig/detectors":
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(test.DetectorsJSON))

		case "/rest/foreignSourcesConfig/assets":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.ElementList{
				Count:   1,
				Element: []string{"address1", "city", "state", "zip"},
			})

		case "/rest/requisitionNames":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.RequisitionsList{
				Count:          2,
				ForeignSources: []string{"Test", "Local"},
			})

		case "/rest/requisitions/deployed/stats":
			assert.Equal(t, http.MethodGet, req.Method)
			now := model.Time{Time: time.Now()}
			sendData(res, model.RequisitionsStats{
				Count: 1,
				ForeignSources: []model.RequisitionStats{
					{Name: "Test", Count: 0, LastImport: &now},
				},
			})

		case "/rest/requisitions/Test":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.Requisition{
				Name:  "Test",
				Nodes: []model.RequisitionNode{testNode},
			})

		case "/rest/requisitions":
			assert.Equal(t, http.MethodPost, req.Method)
			var r model.Requisition
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			err = json.Unmarshal(bytes, &r)
			assert.NilError(t, err)
			if r.Name == "WebSites" {
				node := r.Nodes[0]
				assert.Equal(t, "opennms.com", node.ForeignID)
				assert.Equal(t, "opennms.com", node.NodeLabel)
				assert.Assert(t, net.ParseIP(node.Interfaces[0].IPAddress) != nil)
			}

		case "/rest/requisitions/Local/import":
			assert.Equal(t, http.MethodPut, req.Method)

		case "/rest/requisitions/Local":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/deployed/Local":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/foreignSources/Local":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/foreignSources/deployed/Local":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Go":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, testNode)

		case "/rest/requisitions/Test/nodes":
			assert.Equal(t, http.MethodPost, req.Method)
			var node model.RequisitionNode
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &node)
			if node.ForeignID == "n2" {
				assert.Equal(t, "n2", node.NodeLabel)
			}
			if node.ForeignID == "opennms.com" {
				assert.Equal(t, "opennms.com", node.NodeLabel)
				assert.Assert(t, net.ParseIP(node.Interfaces[0].IPAddress) != nil)
			}

		case "/rest/requisitions/Test/nodes/n2":
			assert.Assert(t, http.MethodDelete == req.Method || http.MethodGet == req.Method)

		case "/rest/requisitions/Test/nodes/n1/interfaces":
			assert.Equal(t, http.MethodPost, req.Method)
			var intf model.RequisitionInterface
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &intf)
			assert.Assert(t, strings.HasPrefix(intf.IPAddress, "10.0.0.1"))
			if intf.IPAddress == "10.0.0.1" {
				switch len(intf.MetaData) {
				case 1:
					assert.Equal(t, "false", intf.MetaData[0].Value)
				case 2:
					assert.Equal(t, "false", intf.MetaData[0].Value)
					assert.Equal(t, "active", intf.MetaData[1].Key)
				}
			}

		case "/rest/requisitions/Test/nodes/n1/interfaces/10.0.0.1":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, testNode.Interfaces[0])

		case "/rest/requisitions/Test/nodes/n1/interfaces/10.0.0.10":
			assert.Assert(t, http.MethodDelete == req.Method || http.MethodGet == req.Method)

		case "/rest/requisitions/Test/nodes/n1/assets":
			assert.Equal(t, http.MethodPost, req.Method)
			var asset model.RequisitionAsset
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &asset)
			assert.Equal(t, "state", asset.Name)
			assert.Equal(t, "NC", asset.Value)

		case "/rest/requisitions/Test/nodes/n1/assets/state":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1/categories":
			assert.Equal(t, http.MethodPost, req.Method)
			var cat model.RequisitionCategory
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &cat)
			assert.Equal(t, "Production", cat.Name)

		case "/rest/requisitions/Test/nodes/n1/categories/Production":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/foreignSources/Test":
			if req.Method == http.MethodPut {
				interval := req.FormValue("scanInterval")
				if interval == "" {
					res.WriteHeader(http.StatusBadRequest)
				}
				return
			}
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, testForeignSource)

		case "/rest/foreignSources":
			assert.Equal(t, http.MethodPost, req.Method)
			var fs model.ForeignSourceDef
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &fs)
			assert.Equal(t, "Local", fs.Name)

		case "/rest/foreignSources/Test/detectors":
			assert.Equal(t, http.MethodPost, req.Method)
			var detector model.Detector
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &detector)
			assert.Equal(t, "ICMP", detector.Name)

		case "/rest/foreignSources/Test/detectors/HTTP":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/foreignSources/Test/policies":
			assert.Equal(t, http.MethodPost, req.Method)
			var policy model.Policy
			bytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &policy)
			assert.Equal(t, "Switches", policy.Name)
			assert.Equal(t, 2, len(policy.Parameters))

		case "/rest/foreignSources/Test/policies/Production":
			assert.Equal(t, http.MethodDelete, req.Method)

		default:
			res.WriteHeader(http.StatusForbidden)
		}
	}))

	rest.Instance.URL = server.URL
	return server
}

func sendData(res http.ResponseWriter, data interface{}) {
	bytes, _ := json.Marshal(data)
	res.WriteHeader(http.StatusOK)
	res.Write(bytes)
}
