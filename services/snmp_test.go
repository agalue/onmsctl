package services

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"gotest.tools/assert"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/google/go-cmp/cmp"
)

var lastPath string
var mutex = &sync.Mutex{}
var mockData = &model.SnmpInfo{
	Version:   "v2c",
	Community: "public",
	Port:      161,
}

type rest struct {
	t *testing.T
}

func (api rest) Get(path string) ([]byte, error) {
	mutex.Lock()
	lastPath = path
	mutex.Unlock()
	assert.Assert(api.t, strings.HasPrefix(path, "/rest/snmpConfig"))
	bytes, _ := json.Marshal(mockData)
	return bytes, nil
}

func (api rest) Post(path string, jsonBytes []byte) error {
	return nil
}

func (api rest) Delete(path string) error {
	return nil
}

func (api rest) Put(path string, jsonBytes []byte, contentType string) error {
	assert.Assert(api.t, strings.HasPrefix(path, "/rest/snmpConfig"))
	snmp := &model.SnmpInfo{}
	json.Unmarshal(jsonBytes, snmp)
	assert.Assert(api.t, cmp.Equal(mockData, snmp))
	return nil
}

func TestGetConfig(t *testing.T) {
	api := GetSnmpAPI(&rest{t})

	snmp, err := api.GetConfig("", "")
	assert.Assert(t, snmp == nil)
	assert.Assert(t, err != nil)

	snmp, err = api.GetConfig("127.0.0.1", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(mockData, snmp))
	assert.Equal(t, "/rest/snmpConfig/127.0.0.1", lastPath)

	snmp, err = api.GetConfig("localhost", "")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(mockData, snmp))
	assert.Equal(t, "/rest/snmpConfig/127.0.0.1", lastPath)

	snmp, err = api.GetConfig("127.0.0.1", "Apex")
	assert.NilError(t, err)
	assert.Assert(t, cmp.Equal(mockData, snmp))
	assert.Equal(t, "/rest/snmpConfig/127.0.0.1?location=Apex", lastPath)
}

func TestSetConfig(t *testing.T) {
	api := GetSnmpAPI(&rest{t})
	err := api.SetConfig("127.0.0.1", *mockData)
	assert.NilError(t, err)
	err = api.SetConfig("localhost", *mockData)
	assert.NilError(t, err)
}
