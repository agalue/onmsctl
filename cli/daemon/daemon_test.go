package daemon

import (
	"testing"

	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

func TestDaemonMap(t *testing.T) {
	assert.Equal(t, true, isValidDaemon("pollerd"))
	assert.Equal(t, false, isValidDaemon("nonexistent"))
	assert.Equal(t, true, isValidDaemon("correlation:MyEngine"))

	assert.Equal(t, "Collectd", getDaemonName("collectd"))
	assert.Equal(t, "EmailNBI", getDaemonName("nbi:email"))
	assert.Equal(t, "DroolsCorrelationEngine:MyEngine", getDaemonName("correlation:MyEngine"))
}

func TestListDaemon(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "daemon", "list"})
	assert.NilError(t, err)
}

func TestReloadDaemon(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "daemon", "reload"})
	assert.Error(t, err, "Daemon name required")

	err = app.Run([]string{app.Name, "daemon", "reload", "SomethignWeird"})
	assert.Error(t, err, "Invalid daemon name SomethignWeird")

	err = app.Run([]string{app.Name, "daemon", "reload", "pollerd"})
	assert.NilError(t, err)
}
