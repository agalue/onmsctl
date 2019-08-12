package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"github.com/urfave/cli"
	"gotest.tools/assert"
)

func TestUtils(t *testing.T) {
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	names, err := GetRequisitionNames()
	assert.NilError(t, err)
	assert.Equal(t, names.Count, 1)
	assert.Equal(t, names.ForeignSources[0], "Test")

	assert.Equal(t, true, RequisitionExists("Test"))
	assert.Equal(t, false, RequisitionExists("Unexisting"))

	app := test.CreateCli(cli.Command{
		Name: "node",
		Action: func(c *cli.Context) error {
			node, err := GetNode(c)
			if err != nil {
				return err
			}
			assert.NilError(t, err)
			assert.Equal(t, "n1", node.ForeignID)
			return nil
		},
	})

	err = app.Run([]string{app.Name, "node"})
	assert.Error(t, err, "Requisition name and foreign ID required")

	err = app.Run([]string{app.Name, "node", "Test"})
	assert.Error(t, err, "Foreign ID required")

	err = app.Run([]string{app.Name, "node", "Test", "n1"})
	assert.NilError(t, err)
}
