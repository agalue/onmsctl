package model

import (
	"fmt"
	"strconv"
	"time"
)

// ScheduledTypes list of valid scheduled outage types
var ScheduledTypes = EnumValue{
	Enum: []string{"specific", "daily", "weekly", "monthly"},
}

// WeekDays list of valid week days
var WeekDays = EnumValue{
	Enum: []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"},
}

// ScheduledNode represents a node by its Node ID under a scheduled outage
type ScheduledNode struct {
	ID string `json:"id" yaml:"id"`
}

// ScheduledInterface represents an IP interface by its Address under a scheduled outage
type ScheduledInterface struct {
	Address string `json:"address" yaml:"address"`
}

// ScheduledTime the time of a given scheduled outage (contents depend on the type)
type ScheduledTime struct {
	Day    string `json:"day,omitempty" yaml:"day,omitempty"`
	Begins string `json:"begins" yaml:"begins"`
	Ends   string `json:"ends" yaml:"ends"`
}

// IsValid verifies if the schedule time is correct for a given type
func (sc *ScheduledTime) IsValid(scheduleType string) error {
	switch scheduleType {
	case "specific": // { "begins": "01-Jun-2017 00:00:00", "ends": "30-Jun-2017 23:59:59" }
		if _, err := time.Parse("02-Jan-2006 15:04:05", sc.Begins); err != nil {
			return fmt.Errorf("invalid specific begin date: %s", err.Error())
		}
		if _, err := time.Parse("02-Jan-2006 15:04:05", sc.Ends); err != nil {
			return fmt.Errorf("invalid specific end date: %s", err.Error())
		}
		if sc.Day != "" {
			return fmt.Errorf("specific schedule only requires begins and ends dates")
		}
		return nil
	case "daily": // { "begins": "17:00:00", "ends": "20:00:00" }
		if err := sc.hasValidHourlyRange(); err != nil {
			return fmt.Errorf("daily schedule error: %s", err.Error())
		}
		if sc.Day != "" {
			return fmt.Errorf("daily schedule only requires begins and ends hours")
		}
		return nil
	case "weekly": // { "begins": "00:00:00", "ends": "23:59:59", "day": "saturday" }
		if err := WeekDays.Set(sc.Day); err != nil {
			return fmt.Errorf("invalid day for weekly schedule. Allowed values: %s", WeekDays.EnumAsString())
		}
		if err := sc.hasValidHourlyRange(); err != nil {
			return fmt.Errorf("weekly schedule error: %s", err.Error())
		}
		return nil
	case "monthly": // { "begins": "00:00:00", "ends": "23:59:59", "day": "1" }
		if _, err := strconv.Atoi(sc.Day); err != nil {
			return fmt.Errorf("invalid monthly day: %s", err.Error())
		}
		if err := sc.hasValidHourlyRange(); err != nil {
			return fmt.Errorf("monthly schedule error: %s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("invalid type %s", scheduleType)
}

func (sc *ScheduledTime) hasValidHourlyRange() error {
	if sc.Begins == "" {
		return fmt.Errorf("begin hour cannot be empty")
	}
	if sc.Ends == "" {
		return fmt.Errorf("end hour cannot be empty")
	}
	if _, err := time.Parse("15:04:05", sc.Begins); err != nil {
		return fmt.Errorf("invalid begin hour: %s", sc.Begins)
	}
	if _, err := time.Parse("15:04:05", sc.Ends); err != nil {
		return fmt.Errorf("invalid end hour: %s", sc.Ends)
	}
	return nil
}

// ScheduledOutage a scheduled outage
type ScheduledOutage struct {
	Name       string               `json:"name" yaml:"name"`
	Type       string               `json:"type" yaml:"type"`
	Nodes      []ScheduledNode      `json:"node,omitempty" yaml:"nodes,omitempty"`
	Interfaces []ScheduledInterface `json:"interface,omitempty" yaml:"interfaces,omitempty"`
	Times      []ScheduledTime      `json:"time,omitempty" yaml:"times,omitempty"`
}

// IsValid verifies if the scheduled outage is correct
func (o *ScheduledOutage) IsValid() error {
	if o.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if o.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}
	if err := ScheduledTypes.Set(o.Type); err != nil {
		return fmt.Errorf("invalid scheduled type. Allowed values: %s", ScheduledTypes.EnumAsString())
	}
	for _, time := range o.Times {
		if err := time.IsValid(o.Type); err != nil {
			return err
		}
	}
	return nil
}

// ScheduledOutageList list of scheduled outages
type ScheduledOutageList struct {
	Outages []ScheduledOutage
}
