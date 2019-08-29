package services

import (
	"encoding/json"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

type eventsAPI struct {
	rest api.RestAPI
}

// GetEventsAPI Obtain an implementation of the Foreign Source Definitions API
func GetEventsAPI(rest api.RestAPI) api.EventsAPI {
	return &eventsAPI{rest}
}

func (api eventsAPI) SendEvent(event model.Event) error {
	if err := event.IsValid(); err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(event)
	return api.rest.Post("/rest/events", jsonBytes)
}
