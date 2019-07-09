package provisioning

import (
	"testing"

	"gotest.tools/assert"
)

func TestListAssets(t *testing.T) {
	var err error
	app := CreateCli(AssetsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "assets", "list"})
	assert.Error(t, err, "Requisition name and foreign ID required")

	err = app.Run([]string{app.Name, "assets", "list", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "assets", "list", "Test", "n1"})
	assert.NilError(t, err)
}

func TestAddAsset(t *testing.T) {
	var err error
	app := CreateCli(AssetsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "assets", "set"})
	assert.Error(t, err, "Requisition name, foreign ID, asset name and value required")

	err = app.Run([]string{app.Name, "assets", "set", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "assets", "set", "Test", "n1"})
	assert.Error(t, err, "Asset name required")

	err = app.Run([]string{app.Name, "assets", "set", "Test", "n1", "state"})
	assert.Error(t, err, "Asset value required")

	err = app.Run([]string{app.Name, "assets", "set", "Test", "n1", "state", "NC"})
	assert.NilError(t, err)
}

func TestDeleteAsset(t *testing.T) {
	var err error
	app := CreateCli(AssetsCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "assets", "delete"})
	assert.Error(t, err, "Requisition name, foreign ID, asset name required")

	err = app.Run([]string{app.Name, "assets", "delete", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "assets", "delete", "Test", "n1"})
	assert.Error(t, err, "Asset name required")

	err = app.Run([]string{app.Name, "assets", "delete", "Test", "n1", "state"})
	assert.NilError(t, err)
}
