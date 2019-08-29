package snmp

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/OpenNMS/onmsctl/test"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestGetSnmp(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, CliCommand)
	defer server.Close()
	api = services.GetSnmpAPI(rest.Instance)

	err = app.Run([]string{app.Name, "snmp", "get"})
	assert.Error(t, err, "IP Address or FQDN required")

	err = app.Run([]string{app.Name, "snmp", "get", "10.0.0.500"})
	assert.ErrorContains(t, err, "Cannot parse address from 10.0.0.500")

	err = app.Run([]string{app.Name, "snmp", "get", "10.0.0.1"})
	assert.NilError(t, err)
}

func TestSetSnmp(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, CliCommand)
	defer server.Close()
	api = services.GetSnmpAPI(rest.Instance)

	err = app.Run([]string{app.Name, "snmp", "set"})
	assert.Error(t, err, "IP Address or FQDN required")

	err = app.Run([]string{app.Name, "snmp", "set", "-c", "private", "-v", "v1", "10.0.0.1"})
	assert.NilError(t, err)
}

func TestApplySnmp(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, CliCommand)
	defer server.Close()
	api = services.GetSnmpAPI(rest.Instance)

	info := &model.SnmpInfo{
		Version:   "v1",
		Community: "private",
	}
	yamlBytes, _ := yaml.Marshal(info)
	err = app.Run([]string{app.Name, "snmp", "apply", "10.0.0.1", string(yamlBytes)})
	assert.NilError(t, err)
}
