package provisioning

import (
	"testing"

	"gotest.tools/assert"
)

func TestListCategories(t *testing.T) {
	var err error
	app := CreateCli(CategoriesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "cat", "list"})
	assert.Error(t, err, "Requisition name and foreign ID required")

	err = app.Run([]string{app.Name, "cat", "list", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "cat", "list", "Test", "n1"})
	assert.NilError(t, err)
}

func TestAddCategory(t *testing.T) {
	var err error
	app := CreateCli(CategoriesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "cat", "add"})
	assert.Error(t, err, "Requisition name, foreign ID, category name required")

	err = app.Run([]string{app.Name, "cat", "add", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "cat", "add", "Test", "n1"})
	assert.Error(t, err, "Category name required")

	err = app.Run([]string{app.Name, "cat", "add", "Test", "n1", "Production"})
	assert.NilError(t, err)
}

func TestDeleteCategory(t *testing.T) {
	var err error
	app := CreateCli(CategoriesCliCommand)
	testServer := CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "cat", "delete"})
	assert.Error(t, err, "Requisition name, foreign ID and category name required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test", "n1"})
	assert.Error(t, err, "Category name required")

	err = app.Run([]string{app.Name, "cat", "delete", "Test", "n1", "Production"})
	assert.NilError(t, err)
}
