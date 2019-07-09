package provisioning

import (
	"testing"

	"gotest.tools/assert"
)

func TestListInterfaces(t *testing.T) {
	var err error
	app := CreateCli(InterfacesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "intf", "list"})
	assert.Error(t, err, "Requisition name and foreign ID required")

	err = app.Run([]string{app.Name, "intf", "list", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "intf", "list", "Test", "n1"})
	assert.NilError(t, err)
}

func TestAddInterface(t *testing.T) {
	var err error
	app := CreateCli(InterfacesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "intf", "add"})
	assert.Error(t, err, "Requisition name, foreign ID, IP address required")

	err = app.Run([]string{app.Name, "intf", "add", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "intf", "add", "Test", "n1"})
	assert.Error(t, err, "IP Address required")

	err = app.Run([]string{app.Name, "intf", "add", "Test", "n1", "10.0.0.10"})
	assert.NilError(t, err)
}

func TestDeleteInterface(t *testing.T) {
	var err error
	app := CreateCli(InterfacesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "intf", "delete"})
	assert.Error(t, err, "Requisition name, foreign ID, IP address required")

	err = app.Run([]string{app.Name, "intf", "delete", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "intf", "delete", "Test", "n1"})
	assert.Error(t, err, "IP address required")

	err = app.Run([]string{app.Name, "intf", "delete", "Test", "n1", "10.0.0.10"})
	assert.NilError(t, err)
}
