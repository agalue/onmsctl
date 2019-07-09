package provisioning

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
	"gotest.tools/assert"
)

var testNode = Node{
	ForeignID: "n1",
	NodeLabel: "n1",
	Interfaces: []Interface{
		{IPAddress: "10.0.0.1", SnmpPrimary: "P"},
	},
	Categories: []Category{
		{"Server"},
	},
	Assets: []Asset{
		{"city", "Durham"},
	},
}

func sendData(res http.ResponseWriter, data interface{}) {
	bytes, _ := json.Marshal(data)
	res.WriteHeader(http.StatusOK)
	res.Write(bytes)
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
			sendData(res, ElementList{1, []string{"address1", "city", "state", "zip"}})

		case "/rest/requisitionNames":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, RequisitionsList{1, []string{"Test", "Local"}})

		case "/rest/requisitions/deployed/stats":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, RequisitionsStats{1, []RequisitionStats{{"Test", 0, []string{}, 0}}})

		case "/rest/requisitions/Test":
			assert.Equal(t, http.MethodGet, req.Method)
			sendData(res, Requisition{Name: "Test", Nodes: []Node{testNode}})

		case "/rest/requisitions":
			assert.Equal(t, http.MethodPost, req.Method)
			var r Requisition
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			err = json.Unmarshal(bytes, &r)
			assert.NilError(t, err)

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
			var node Node
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &node)
			assert.Equal(t, "n2", node.ForeignID)
			assert.Equal(t, "n2", node.NodeLabel)

		case "/rest/requisitions/Test/nodes/n2":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1/interfaces":
			assert.Equal(t, http.MethodPost, req.Method)
			var intf Interface
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &intf)
			assert.Equal(t, "10.0.0.10", intf.IPAddress)

		case "/rest/requisitions/Test/nodes/n1/interfaces/10.0.0.10":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1/assets":
			assert.Equal(t, http.MethodPost, req.Method)
			var asset Asset
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, &asset)
			assert.Equal(t, "state", asset.Name)
			assert.Equal(t, "NC", asset.Value)

		case "/rest/requisitions/Test/nodes/n1/assets/state":
			assert.Equal(t, http.MethodDelete, req.Method)

		case "/rest/requisitions/Test/nodes/n1/categories":
			assert.Equal(t, http.MethodPost, req.Method)
			var cat Category
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
