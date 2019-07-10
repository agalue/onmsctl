package provisioning

import (
	"testing"

	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestListRequisitions(t *testing.T) {
	var err error
	app := CreateCli(RequisitionsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "req", "list"})
	assert.NilError(t, err)
}

func TestGetRequisition(t *testing.T) {
	var err error
	app := CreateCli(RequisitionsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "req", "get"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "req", "get", "Test"})
	assert.NilError(t, err)
}

func TestAddRequisition(t *testing.T) {
	var err error
	app := CreateCli(RequisitionsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "req", "add"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "req", "add", "Go"})
	assert.NilError(t, err)
}

func TestDeleteRequisition(t *testing.T) {
	var err error
	app := CreateCli(RequisitionsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "req", "delete"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "req", "delete", "Local"})
	assert.NilError(t, err)
}

func TestApplyRequisition(t *testing.T) {
	var err error
	app := CreateCli(RequisitionsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "req", "apply"})
	assert.Error(t, err, "YAML content cannot be empty")

	var testReq = Requisition{
		Name: "WebSites",
		Nodes: []Node{
			{
				ForeignID: "opennms.com",
				Interfaces: []Interface{
					{IPAddress: "www.opennms.com"},
				},
				Categories: []Category{
					{"Server"},
				},
			},
		},
	}

	reqYaml, _ := yaml.Marshal(testReq)
	err = app.Run([]string{app.Name, "req", "apply", string(reqYaml)})
	assert.NilError(t, err)
}
