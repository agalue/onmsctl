package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/OpenNMS/onmsctl/test"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestListRequisitions(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, RequisitionsCliCommand)
	defer server.Close()
	api = services.GetRequisitionsAPI(rest.Instance)

	err = app.Run([]string{app.Name, "req", "list"})
	assert.NilError(t, err)
}

func TestGetRequisition(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, RequisitionsCliCommand)
	defer server.Close()
	api = services.GetRequisitionsAPI(rest.Instance)

	err = app.Run([]string{app.Name, "req", "get"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "req", "get", "Test"})
	assert.NilError(t, err)
}

func TestAddRequisition(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, RequisitionsCliCommand)
	defer server.Close()
	api = services.GetRequisitionsAPI(rest.Instance)

	err = app.Run([]string{app.Name, "req", "add"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "req", "add", "Go"})
	assert.NilError(t, err)
}

func TestDeleteRequisition(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, RequisitionsCliCommand)
	defer server.Close()
	api = services.GetRequisitionsAPI(rest.Instance)

	err = app.Run([]string{app.Name, "req", "delete"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "req", "delete", "Local"})
	assert.NilError(t, err)
}

func TestApplyRequisition(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, RequisitionsCliCommand)
	defer server.Close()
	api = services.GetRequisitionsAPI(rest.Instance)

	err = app.Run([]string{app.Name, "req", "apply"})
	assert.Error(t, err, "Content cannot be empty")

	var testReq = model.Requisition{
		Name: "WebSites",
		Nodes: []model.RequisitionNode{
			{
				ForeignID: "opennms.com",
				Interfaces: []model.RequisitionInterface{
					{IPAddress: "www.opennms.com"},
				},
				Categories: []model.RequisitionCategory{
					{Name: "Server"},
				},
			},
		},
	}

	reqYaml, _ := yaml.Marshal(testReq)
	err = app.Run([]string{app.Name, "req", "apply", string(reqYaml)})
	assert.NilError(t, err)
}
