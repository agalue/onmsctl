package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/google/go-cmp/cmp"

	"gotest.tools/assert"
)

var mockEvent = &model.Event{
	UEI:       "uei.opennms.org/test",
	NodeID:    10,
	Interface: "10.0.0.1",
	Service:   "SNMP",
	Source:    "onmsctl",
	Parameters: []model.EventParam{
		{Name: "owner", Value: "agalue"},
	},
}

type mockEventRest struct {
	t *testing.T
}

func (api mockEventRest) Get(path string) ([]byte, error) {
	return nil, fmt.Errorf("sould not be called")
}

func (api mockEventRest) Post(path string, jsonBytes []byte) error {
	assert.Assert(api.t, strings.HasPrefix(path, "/rest/events"))
	event := &model.Event{}
	json.Unmarshal(jsonBytes, event)
	assert.Assert(api.t, cmp.Equal(mockEvent, event))
	return nil
}

func (api mockEventRest) Delete(path string) error {
	return fmt.Errorf("sould not be called")
}

func (api mockEventRest) Put(path string, jsonBytes []byte, contentType string) error {
	return fmt.Errorf("sould not be called")
}

func TestSendEvent(t *testing.T) {
	api := GetEventsAPI(&mockEventRest{t})
	err := api.SendEvent(*mockEvent)
	assert.NilError(t, err)
}
