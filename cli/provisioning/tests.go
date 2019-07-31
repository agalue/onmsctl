package provisioning

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
		{"Server"},
	},
	Assets: []model.RequisitionAsset{
		{"city", "Durham"},
	},
}

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

		case "/rest/foreignSourcesConfig/assets":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.ElementList{1, []string{"address1", "city", "state", "zip"}})

		case "/rest/requisitionNames":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.RequisitionsList{1, []string{"Test", "Local"}})

		case "/rest/requisitions/deployed/stats":
			assert.Equal(t, http.MethodGet, req.Method)
			now := model.Time{Time: time.Now()}
			sendData(res, model.RequisitionsStats{1, []model.RequisitionStats{{"Test", 0, nil, &now}}})

		case "/rest/requisitions/Test":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, model.Requisition{Name: "Test", Nodes: []model.RequisitionNode{testNode}})

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
