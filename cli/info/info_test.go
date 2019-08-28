package info

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

func TestSendEvent(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, CliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "info"})
	assert.NilError(t, err)
}
