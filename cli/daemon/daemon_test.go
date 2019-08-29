package daemon

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

	"gotest.tools/assert"
)

var mockData = &model.Event{
	UEI:    "uei.opennms.org/internal/reloadDaemonConfig",
	Source: "onmsctl",
	Parameters: []model.EventParam{
		{Name: "daemonName", Value: "Pollerd"},
	},
}

func TestDaemonMap(t *testing.T) {
	assert.Equal(t, true, isValidDaemon("pollerd"))
	assert.Equal(t, false, isValidDaemon("nonexistent"))
	assert.Equal(t, true, isValidDaemon("correlation:MyEngine"))

	assert.Equal(t, "Collectd", getDaemonName("collectd"))
	assert.Equal(t, "EmailNBI", getDaemonName("nbi:email"))
	assert.Equal(t, "DroolsCorrelationEngine:MyEngine", getDaemonName("correlation:MyEngine"))
}

func TestListDaemon(t *testing.T) {
	app := test.CreateCli(CliCommand)
	err := app.Run([]string{app.Name, "daemon", "list"})
	assert.NilError(t, err)
}

func TestReloadDaemon(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
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
	defer server.Close()

	err = app.Run([]string{app.Name, "daemon", "reload"})
	assert.Error(t, err, "Daemon name required")

	err = app.Run([]string{app.Name, "daemon", "reload", "Weird"})
	assert.Error(t, err, "Invalid daemon name Weird")

	err = app.Run([]string{app.Name, "daemon", "reload", "pollerd"})
	assert.NilError(t, err)
}
