package snmp

import (
	"encoding/json"
	"fmt"
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

var mockData = &model.SnmpInfo{
	Version:        "v2c",
	Community:      "public",
	Port:           161,
	MaxVarsPerPdu:  10,
	MaxRepetitions: 2,
	Retries:        2,
	Timeout:        1800,
	TTL:            20000,
}

func createMockServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		assert.Assert(t, strings.HasPrefix(req.URL.Path, "/rest/snmpConfig"))
		switch req.Method {
		case http.MethodGet:
			res.WriteHeader(http.StatusOK)
			bytes, _ := json.Marshal(mockData)
			res.Write(bytes)
		case http.MethodPut:
			snmp := &model.SnmpInfo{}
			bytes, err := ioutil.ReadAll(req.Body)
			assert.NilError(t, err)
			json.Unmarshal(bytes, snmp)
			assert.Assert(t, cmp.Equal(mockData, snmp))
			res.WriteHeader(http.StatusOK)
		default:
			res.WriteHeader(http.StatusForbidden)
		}
	}))
	rest.Instance.URL = server.URL
	return server
}

func TestGetSnmp(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	server := createMockServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "snmp", "get"})
	assert.Error(t, err, "IP Address or FQDN required")

	err = app.Run([]string{app.Name, "snmp", "get", "10.0.0.500"})
	assert.ErrorContains(t, err, "Cannot parse address from 10.0.0.500")

	err = app.Run([]string{app.Name, "snmp", "get", "10.0.0.1"})
	assert.NilError(t, err)

	err = app.Run([]string{app.Name, "snmp", "get", "localhost"})
	assert.NilError(t, err)
}

func TestSetSnmp(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	server := createMockServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "snmp", "set"})
	assert.Error(t, err, "IP Address or FQDN required")

	err = app.Run([]string{app.Name, "snmp", "set", "-c", mockData.Community, "-v", mockData.Version, "10.0.0.1", "-ttl", fmt.Sprintf("%d", mockData.TTL)})
	assert.NilError(t, err)
}

func TestApplySnmp(t *testing.T) {
	var err error
	app := test.CreateCli(CliCommand)
	server := createMockServer(t)
	defer server.Close()

	yamlBytes, _ := yaml.Marshal(mockData)
	err = app.Run([]string{app.Name, "snmp", "apply", "10.0.0.1", string(yamlBytes)})
	assert.NilError(t, err)
}
