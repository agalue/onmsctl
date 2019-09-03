package info

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/test"

	"gotest.tools/assert"
)

var mockData = &model.OnmsInfo{
	DisplayVersion:     "24.1.2",
	Version:            "24.1.2",
	PackageName:        "opennms",
	PackageDescription: "OpenNMS",
}

func TestSendEvent(t *testing.T) {
	app := test.CreateCli(CliCommand)
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Assert(t, strings.HasPrefix(req.URL.Path, "/rest/info"))
		assert.Equal(t, http.MethodGet, req.Method)
		bytes, _ := json.Marshal(mockData)
		res.WriteHeader(http.StatusOK)
		res.Write(bytes)
	}))
	rest.Instance.URL = server.URL
	defer server.Close()

	err := app.Run([]string{app.Name, "info"})
	assert.NilError(t, err)
}
