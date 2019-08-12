package snmp

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/test"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestGetSnmp(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "snmp", "get"})
	assert.Error(t, err, "IP Address or FQDN required")

	err = app.Run([]string{app.Name, "snmp", "get", "10.0.0.500"})
	assert.ErrorContains(t, err, "Cannot parse address from 10.0.0.500")

	err = app.Run([]string{app.Name, "snmp", "get", "10.0.0.1"})
	assert.NilError(t, err)
}

func TestSetSnmp(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "snmp", "set"})
	assert.Error(t, err, "IP Address or FQDN required")

	err = app.Run([]string{app.Name, "snmp", "set", "-c", "private", "-v", "v1", "10.0.0.1"})
	assert.NilError(t, err)
}

func TestApplySnmp(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	info := &model.SnmpInfo{
		Version:   "v1",
		Community: "private",
	}
	yamlBytes, _ := yaml.Marshal(info)
	err = app.Run([]string{app.Name, "snmp", "apply", "10.0.0.1", string(yamlBytes)})
	assert.NilError(t, err)
}
