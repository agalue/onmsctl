package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

func TestListCategories(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, CategoriesCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "cat", "list"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "cat", "list", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "cat", "list", "Test", "n1"})
	assert.NilError(t, err)
}

func TestAddCategory(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, CategoriesCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "cat", "add"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "cat", "add", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "cat", "add", "Test", "n1"})
	assert.Error(t, err, "Category name cannot be empty")

	err = app.Run([]string{app.Name, "cat", "add", "Test", "n1", "Production"})
	assert.NilError(t, err)
}

func TestDeleteCategory(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, CategoriesCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "cat", "delete"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test", "n1"})
	assert.Error(t, err, "Category name required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test", "n1", "Production"})
	assert.NilError(t, err)
}
