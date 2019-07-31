package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestListNodes(t *testing.T) {
	var err error
	app := CreateCli(NodesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "node", "list"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "node", "list", "Test"})
	assert.NilError(t, err)
}

func TestGetNode(t *testing.T) {
	var err error
	app := CreateCli(NodesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "node", "get"})
	assert.Error(t, err, "Requisition name and foreign ID required")

	err = app.Run([]string{app.Name, "node", "get", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "node", "get", "Test", "n1"})
	assert.NilError(t, err)
}

func TestAddNode(t *testing.T) {
	var err error
	app := CreateCli(NodesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "node", "add"})
	assert.Error(t, err, "Requisition name and foreign ID required")

	err = app.Run([]string{app.Name, "node", "add", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "node", "add", "Test", "n2"})
	assert.NilError(t, err)
}

func TestDeleteNode(t *testing.T) {
	var err error
	app := CreateCli(NodesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "node", "delete"})
	assert.Error(t, err, "Requisition name and foreign ID required")

	err = app.Run([]string{app.Name, "node", "delete", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "node", "delete", "Test", "n2"})
	assert.NilError(t, err)
}

func TestApplyNode(t *testing.T) {
	var err error
	app := CreateCli(NodesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "node", "apply"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "node", "apply", "Test"})
	assert.Error(t, err, "YAML content cannot be empty")

	var testNode = model.RequisitionNode{
		ForeignID: "opennms.com",
		Interfaces: []model.RequisitionInterface{
			{IPAddress: "www.opennms.com"},
		},
		Categories: []model.RequisitionCategory{
			{"Server"},
		},
	}
	nodeYaml, _ := yaml.Marshal(testNode)
	err = app.Run([]string{app.Name, "node", "apply", "Test", string(nodeYaml)})
	assert.NilError(t, err)
}
