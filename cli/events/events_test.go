package events

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/test"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestSendEvent(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "events", "send"})
	assert.Error(t, err, "UEI required")

	err = app.Run([]string{app.Name, "events", "send", "-p", "owner=agalue", "uei.opennms.org/test"})
	assert.NilError(t, err)
}

func TestApplyEvent(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	event := &model.Event{
		UEI:       "uei.opennms.org/test",
		NodeID:    10,
		Interface: "10.0.0.1",
		Service:   "SNMP",
		Parameters: []model.EventParam{
			{Name: "owner", Value: "agalue"},
		},
	}
	yamlBytes, _ := yaml.Marshal(event)
	err = app.Run([]string{app.Name, "events", "apply", string(yamlBytes)})
	assert.NilError(t, err)
}
