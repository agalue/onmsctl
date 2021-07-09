package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

func TestListInterfaces(t *testing.T) {
	var err error
	app := test.CreateCli(InterfacesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "intf", "list"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "intf", "list", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "intf", "list", "Test", "n1"})
	assert.NilError(t, err)
}

func TestGetInterface(t *testing.T) {
	var err error
	app := test.CreateCli(InterfacesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "intf", "get"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "intf", "get", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "intf", "get", "Test", "n1"})
	assert.Error(t, err, "IP Address required")

	err = app.Run([]string{app.Name, "intf", "get", "Test", "n1", "10.0.0.1"})
	assert.NilError(t, err)
}

func TestAddInterface(t *testing.T) {
	var err error
	app := test.CreateCli(InterfacesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "intf", "add"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "intf", "add", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "intf", "add", "Test", "n1"})
	assert.Error(t, err, "IP Address cannot be empty")

	err = app.Run([]string{app.Name, "intf", "add", "Test", "n1", "10.0.0.10"})
	assert.NilError(t, err)
}

func TestDeleteInterface(t *testing.T) {
	var err error
	app := test.CreateCli(InterfacesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "intf", "delete"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "intf", "delete", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "intf", "delete", "Test", "n1"})
	assert.Error(t, err, "IP Address required")

	err = app.Run([]string{app.Name, "intf", "delete", "Test", "n1", "10.0.0.10"})
	assert.NilError(t, err)
}

func TestInterfaceMetaData(t *testing.T) {
	var err error
	app := test.CreateCli(InterfacesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "intf", "meta", "set", "Test", "n1", "10.0.0.1", "active", "true"})
	assert.NilError(t, err)

	err = app.Run([]string{app.Name, "intf", "meta", "set", "Test", "n1", "10.0.0.1", "mpls", "false"})
	assert.NilError(t, err)
}
