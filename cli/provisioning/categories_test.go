package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

func TestListCategories(t *testing.T) {
	var err error
	app := test.CreateCli(CategoriesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "cat", "list"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "cat", "list", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "cat", "list", "Test", "n1"})
	assert.NilError(t, err)
}

func TestAddCategory(t *testing.T) {
	var err error
	app := test.CreateCli(CategoriesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "cat", "add"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "cat", "add", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "cat", "add", "Test", "n1"})
	assert.Error(t, err, "category name cannot be empty")

	err = app.Run([]string{app.Name, "cat", "add", "Test", "n1", "Production"})
	assert.NilError(t, err)
}

func TestDeleteCategory(t *testing.T) {
	var err error
	app := test.CreateCli(CategoriesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "cat", "delete"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test"})
	assert.Error(t, err, "foreign ID required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test", "n1"})
	assert.Error(t, err, "category name required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test", "n1", "Production"})
	assert.NilError(t, err)
}
