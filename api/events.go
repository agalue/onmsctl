package api

import "github.com/OpenNMS/onmsctl/model"

// EventsAPI the API to manipulate Events
type EventsAPI interface {
	SendEvent(event model.Event) error
}
