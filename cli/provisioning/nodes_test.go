package provisioning

import (
	"testing"

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
