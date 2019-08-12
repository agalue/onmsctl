package info

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

func TestSendEvent(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "info"})
	assert.NilError(t, err)
}
