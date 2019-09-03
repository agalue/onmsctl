package events

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/test"
	"github.com/google/go-cmp/cmp"

	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

var mockData = &model.Event{
	UEI:       "uei.opennms.org/test",
	NodeID:    10,
	Interface: "10.0.0.1",
	Service:   "SNMP",
	Source:    "onmsctl",
	Parameters: []model.EventParam{
		{Name: "owner", Value: "agalue"},
	},
}

func createMockServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Assert(t, strings.HasPrefix(req.URL.Path, "/rest/events"))
		assert.Equal(t, http.MethodPost, req.Method)
		event := &model.Event{}
		bytes, err := ioutil.ReadAll(req.Body)
		assert.NilError(t, err)
		json.Unmarshal(bytes, event)
		assert.Assert(t, cmp.Equal(mockData, event))
		res.WriteHeader(http.StatusOK)
	}))
	rest.Instance.URL = server.URL
	return server
}

func TestSendEvent(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	server := createMockServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "events", "send"})
	assert.Error(t, err, "UEI required")

	err = app.Run([]string{app.Name, "events", "send", "-n", "10", "-i", "10.0.0.1", "-s", "SNMP", "-p", "owner=agalue", "uei.opennms.org/test"})
	assert.NilError(t, err)
}

func TestApplyEvent(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	server := createMockServer(t)
	defer server.Close()

	yamlBytes, _ := yaml.Marshal(mockData)
	err = app.Run([]string{app.Name, "events", "apply", string(yamlBytes)})
	assert.NilError(t, err)
}
