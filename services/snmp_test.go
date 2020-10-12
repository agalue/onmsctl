package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"gotest.tools/assert"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/google/go-cmp/cmp"
)

var mockSnmpInfo = &model.SnmpInfo{
	Version:   "v2c",
	Community: "public",
	Port:      161,
}

type mockSnmpInfoRest struct {
	test     *testing.T
	lastPath string
}

func (api *mockSnmpInfoRest) Get(path string) ([]byte, error) {
	api.lastPath = path
	assert.Assert(api.test, strings.HasPrefix(path, "/rest/snmpConfig"))
	bytes, _ := json.Marshal(mockSnmpInfo)
	return bytes, nil
}

func (api mockSnmpInfoRest) Post(path string, jsonBytes []byte) error {
	return fmt.Errorf("should not be called")
}

func (api mockSnmpInfoRest) PostRaw(path string, dataBytes []byte, contentType string) (*http.Response, error) {
	return nil, fmt.Errorf("should not be called")
}

func (api mockSnmpInfoRest) Delete(path string) error {
	return fmt.Errorf("should not be called")
}

func (api mockSnmpInfoRest) Put(path string, dataBytes []byte, contentType string) error {
	assert.Assert(api.test, strings.HasPrefix(path, "/rest/snmpConfig"))
	snmp := &model.SnmpInfo{}
	json.Unmarshal(dataBytes, snmp)
	assert.Assert(api.test, cmp.Equal(mockSnmpInfo, snmp))
	return nil
}

func (api mockSnmpInfoRest) IsValid(r *http.Response) error {
	return nil
}

func TestGetConfig(t *testing.T) {
	rest := &mockSnmpInfoRest{test: t}
	api := GetSnmpAPI(rest)

	snmp, err := api.GetConfig("", "")
	assert.Assert(t, snmp == nil)
	assert.Assert(t, err != nil)

	snmp, err = api.GetConfig("127.0.0.1", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(mockSnmpInfo, snmp))
	assert.Equal(t, "/rest/snmpConfig/127.0.0.1", rest.lastPath)

	snmp, err = api.GetConfig("localhost", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(mockSnmpInfo, snmp))
	assert.Assert(t, rest.lastPath == "/rest/snmpConfig/127.0.0.1" || rest.lastPath == "/rest/snmpConfig/::1")

	snmp, err = api.GetConfig("127.0.0.1", "Apex")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(mockSnmpInfo, snmp))
	assert.Equal(t, "/rest/snmpConfig/127.0.0.1?location=Apex", rest.lastPath)
}

func TestSetConfig(t *testing.T) {
	rest := &mockSnmpInfoRest{test: t}
	api := GetSnmpAPI(rest)
	err := api.SetConfig("127.0.0.1", *mockSnmpInfo)
	assert.NilError(t, err)
	err = api.SetConfig("localhost", *mockSnmpInfo)
	assert.NilError(t, err)
}
