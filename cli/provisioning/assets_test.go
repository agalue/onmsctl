package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

func TestListAssets(t *testing.T) {
	var err error
	app := test.CreateCli(AssetsCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "asset", "list"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "asset", "list", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "asset", "list", "Test", "n1"})
	assert.NilError(t, err)
}

func TestAddAsset(t *testing.T) {
	var err error
	app := test.CreateCli(AssetsCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "asset", "set"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "asset", "set", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "asset", "set", "Test", "n1"})
	assert.Error(t, err, "asset name cannot be empty")

	err = app.Run([]string{app.Name, "asset", "set", "Test", "n1", "state"})
	assert.Error(t, err, "asset value for state cannot be empty")

	err = app.Run([]string{app.Name, "asset", "set", "Test", "n1", "state", "NC"})
	assert.NilError(t, err)
}

func TestDeleteAsset(t *testing.T) {
	var err error
	app := test.CreateCli(AssetsCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "asset", "delete"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "asset", "delete", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "asset", "delete", "Test", "n1"})
	assert.Error(t, err, "asset name required")

	err = app.Run([]string{app.Name, "asset", "delete", "Test", "n1", "state"})
	assert.NilError(t, err)
}
