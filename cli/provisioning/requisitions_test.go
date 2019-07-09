package provisioning

import (
	"testing"

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
