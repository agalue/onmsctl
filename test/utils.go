package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
	"gotest.tools/assert"
)

var testNode = model.RequisitionNode{
	ForeignID: "n1",
	NodeLabel: "n1",
	Interfaces: []model.RequisitionInterface{
		{IPAddress: "10.0.0.1", SnmpPrimary: "P"},
	},
	Categories: []model.RequisitionCategory{
		{Name: "Server"},
	},
	Assets: []model.RequisitionAsset{
		{Name: "city", Value: "Durham"},
	},
}

var policiesJSON = `
{
	"plugins": [
		{
			"name": "Set Node Category",
			"class": "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy",
			"parameters": [
				{
					"key": "category",
					"required": true,
					"options": []
				},
				{
					"key": "matchBehavior",
					"required": true,
					"options": [
						"ALL_PARAMETERS",
						"ANY_PARAMETER",
						"NO_PARAMETERS"
					]
				},
				{
					"key": "foreignId",
					"required": false,
					"options": []
				},
				{
					"key": "foreignSource",
					"required": false,
					"options": []
				},
				{
					"key": "label",
					"required": false,
					"options": []
				},
				{
					"key": "labelSource",
					"required": false,
					"options": []
				},
				{
					"key": "netBiosDomain",
					"required": false,
					"options": []
				},
				{
					"key": "netBiosName",
					"required": false,
					"options": []
				},
				{
					"key": "operatingSystem",
					"required": false,
					"options": []
				},
				{
					"key": "sysContact",
					"required": false,
					"options": []
				},
				{
					"key": "sysDescription",
					"required": false,
					"options": []
				},
				{
					"key": "sysLocation",
					"required": false,
					"options": []
				},
				{
					"key": "sysName",
					"required": false,
					"options": []
				},
				{
					"key": "sysObjectId",
					"required": false,
					"options": []
				},
				{
					"key": "type",
					"required": false,
					"options": []
				}
			]
		}
	],
	"count": 1,
	"totalCount": 1,
	"offset": 0
}
`

var detectorsJSON = `
{
  "plugins": [
    {
      "name": "ICMP",
      "class": "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
      "parameters": [
        {
          "key": "allowFragmentation",
          "required": false,
          "options": []
        },
        {
          "key": "dscp",
          "required": false,
          "options": []
        },
        {
          "key": "ipMatch",
          "required": false,
          "options": []
        },
        {
          "key": "port",
          "required": false,
          "options": []
        },
        {
          "key": "retries",
          "required": false,
          "options": []
        },
        {
          "key": "serviceName",
          "required": false,
          "options": []
        },
        {
          "key": "timeout",
          "required": false,
          "options": []
        }
      ]
    },
    {
      "name": "SNMP",
      "class": "org.opennms.netmgt.provision.detector.snmp.SnmpDetector",
      "parameters": [
        {
          "key": "forceVersion",
          "required": false,
          "options": []
        },
        {
          "key": "hex",
          "required": false,
          "options": []
        },
        {
          "key": "ipMatch",
          "required": false,
          "options": []
        },
        {
          "key": "isTable",
          "required": false,
          "options": []
        },
        {
          "key": "matchType",
          "required": false,
          "options": []
        },
        {
          "key": "oid",
          "required": false,
          "options": []
        },
        {
          "key": "port",
          "required": false,
          "options": []
        },
        {
          "key": "retries",
          "required": false,
          "options": []
        },
        {
          "key": "serviceName",
          "required": false,
          "options": []
        },
        {
          "key": "timeout",
          "required": false,
          "options": []
        },
        {
          "key": "vbvalue",
          "required": false,
          "options": []
        }
      ]
    }
	],
	"count": 1,
	"totalCount": 1,
	"offset": 0
}
`

// CreateCli creates a CLI Application object
func CreateCli(cmd cli.Command) *cli.App {
	var app = cli.NewApp()
	app.Name = "onmsctl"
	app.Commands = []cli.Command{cmd}
	return app
}

// CreateTestServer creates a test HTTP server
func CreateTestServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("Received %s request from %s\n", req.Method, req.URL.Path)

		switch req.URL.Path {

		case "/rest/foreignSourcesConfig/policies":
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(policiesJSON))

		case "/rest/foreignSourcesConfig/detectors":
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(detectorsJSON))

		case "/rest/foreignSourcesConfig/assets":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.ElementList{
				Count:   1,
				Element: []string{"address1", "city", "state", "zip"},
			})

		case "/rest/requisitionNames":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.RequisitionsList{
				Count:          1,
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
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			err = json.Unmarshal(bytes, &r)
			assert.NilError(t, err)
			if r.Name == "WebSites" {
				node := r.Nodes[0]
				assert.Equal(t, "opennms.com", node.ForeignID)
				assert.Equal(t, "opennms.com", node.NodeLabel)
				assert.Equal(t, "34.194.50.139", node.Interfaces[0].IPAddress)
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
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &node)
			if node.ForeignID == "n2" {
				assert.Equal(t, "n2", node.NodeLabel)
			}
			if node.ForeignID == "opennms.com" {
				assert.Equal(t, "opennms.com", node.NodeLabel)
				assert.Equal(t, "34.194.50.139", node.Interfaces[0].IPAddress)
			}

		case "/rest/requisitions/Test/nodes/n2":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1/interfaces":
			assert.Equal(t, http.MethodPost, req.Method)
			var intf model.RequisitionInterface
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &intf)
			assert.Equal(t, "10.0.0.10", intf.IPAddress)

		case "/rest/requisitions/Test/nodes/n1/interfaces/10.0.0.1":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, testNode.Interfaces[0])

		case "/rest/requisitions/Test/nodes/n1/interfaces/10.0.0.10":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1/assets":
			assert.Equal(t, http.MethodPost, req.Method)
			var asset model.RequisitionAsset
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &asset)
			assert.Equal(t, "state", asset.Name)
			assert.Equal(t, "NC", asset.Value)

		case "/rest/requisitions/Test/nodes/n1/assets/state":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1/categories":
			assert.Equal(t, http.MethodPost, req.Method)
			var cat model.RequisitionCategory
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &cat)
			assert.Equal(t, "Production", cat.Name)

		case "/rest/requisitions/Test/nodes/n1/categories/Production":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/foreignSources/Test":
			if req.Method == http.MethodPut {
				interval := req.FormValue("scan-interval")
				if interval == "" {
					res.WriteHeader(http.StatusBadRequest)
				}
				return
			}
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.ForeignSourceDef{
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
			})

		case "/rest/foreignSources":
			assert.Equal(t, http.MethodPost, req.Method)
			var fs model.ForeignSourceDef
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &fs)
			assert.Equal(t, "Local", fs.Name)

		case "/rest/foreignSources/Test/detectors":
			assert.Equal(t, http.MethodPost, req.Method)
			var detector model.Detector
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &detector)
			assert.Equal(t, "ICMP", detector.Name)

		case "/rest/foreignSources/Test/detectors/HTTP":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/foreignSources/Test/policies":
			assert.Equal(t, http.MethodPost, req.Method)
			var policy model.Policy
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &policy)
			assert.Equal(t, "Switches", policy.Name)
			assert.Equal(t, 2, len(policy.Parameters))

		case "/rest/foreignSources/Test/policies/Production":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/info":
			data := &model.OnmsInfo{
				DisplayVersion:     "24.1.2",
				Version:            "24.1.2",
				PackageName:        "opennms",
				PackageDescription: "OpenNMS",
			}
			sendData(res, data)

		case "/rest/events":
			assert.Equal(t, http.MethodPost, req.Method)
			var event model.Event
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &event)
			assert.Assert(t, event.UEI != "")

		case "/rest/snmpConfig/10.0.0.1":
			if http.MethodGet == req.Method {
				data := &model.SnmpInfo{
					Version:   "v2c",
					Community: "public",
					Port:      161,
				}
				sendData(res, data)
				return
			}
			if http.MethodPut == req.Method {
				var snmp model.SnmpInfo
				bytes, err := ioutil.ReadAll(req.Body)
				assert.NilError(t, err)
				json.Unmarshal(bytes, &snmp)
				assert.Equal(t, "v1", snmp.Version)
				assert.Equal(t, "private", snmp.Community)

				return
			}
			res.WriteHeader(http.StatusForbidden)

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
