package model

import (
	"fmt"
	"net"
	"time"
)

var (
	// Severities list of valid event severities
	Severities = EnumValue{
		Enum: []string{"Indeterminate", "Normal", "Warning", "Minor", "Major", "Critical"},
	}

	// Destinations for logMsg
	Destinations = EnumValue{
		Enum: []string{"logndisplay", "displayonly", "logonly", "suppress", "donotpersist"},
	}
)

// SNMP an event SNMP object
type SNMP struct {
	ID        string `json:"id" yaml:"id"`
	Version   string `json:"version,omitempty" yaml:"version,omitempty"`
	Specific  int    `json:"specific" yaml:"specific,omitempty"`
	Generic   int    `json:"generic" yaml:"generic,omitempty"`
	Community string `json:"community,omitempty" yaml:"community,omitempty"`
	Timestamp *Time  `json:"time-stamp,omitempty" yaml:"timeStamp,omitempty"`
}

// EventParam an event parameter object
type EventParam struct {
	Name  string `json:"parmName" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// MaskElement an event mask element object
type MaskElement struct {
	Name   string   `json:"mename" yaml:"mename"`
	Values []string `json:"mevalue" yaml:"mevalue"`
}

// Mask an event mask object
type Mask struct {
	Elements []MaskElement `json:"maskelement,omitempty" yaml:"maskElement,omitempty"`
}

// LogMsg the event log message
type LogMsg struct {
	Message     string `json:"value" yaml:"message"`
	Notify      bool   `json:"notify" yaml:"notify"`
	Destination string `json:"dest" yaml:"destination"`
}

// Validate returns an error if the log message is invalid
func (lm *LogMsg) Validate() error {
	if lm.Message == "" {
		return fmt.Errorf("message cannot be null")
	}
	if lm.Destination == "" {
		lm.Destination = "logndisplay"
	}
	return nil
}

// Event an event object
// Time uses a string format. Example: "Saturday, July 13, 2019 2:13:43 PM GMT"
type Event struct {
	SnmpMask      *Mask        `json:"mask,omitempty" yaml:"mask,omitempty"`
	Snmp          *SNMP        `json:"snmp,omitempty" yaml:"snmp,omitempty"`
	LogMessage    *LogMsg      `json:"logmsg,omitempty" yaml:"logmsg,omitempty"`
	UEI           string       `json:"uei" yaml:"uei"`
	Source        string       `json:"source" yaml:"source"`
	Time          string       `json:"time,omitempty" yaml:"time,omitempty"`
	Host          string       `json:"host,omitempty" yaml:"host,omitempty"`
	MasterStation string       `json:"master-station,omitempty" yaml:"masterStation,omitempty"`
	NodeID        int64        `json:"nodeid,omitempty" yaml:"nodeID,omitempty"`
	Interface     string       `json:"interface,omitempty" yaml:"interface,omitempty"`
	Service       string       `json:"service,omitempty" yaml:"service,omitempty"`
	IfIndex       int          `json:"ifIndex,omitempty" yaml:"ifIndex,omitempty"`
	SnmpHost      string       `json:"snmphost,omitempty" yaml:"snmpHost,omitempty"`
	Parameters    []EventParam `json:"parms,omitempty" yaml:"parameters,omitempty"`
	Description   string       `json:"descr,omitempty" yaml:"description,omitempty"`
	Severity      string       `json:"severity,omitempty" yaml:"severity,omitempty"`
	PathOutage    string       `json:"pathoutage,omitempty" yaml:"pathOutage,omitempty"`
	OperInstruct  string       `json:"operinstruct,omitempty" yaml:"operInstruct,omitempty"`
}

// AddParameter adds a new parameter to the event
func (e *Event) AddParameter(key string, value string) {
	e.Parameters = append(e.Parameters, EventParam{Name: key, Value: value})
}

// SetTime sets the string date based on a Time object
func (e *Event) SetTime(date time.Time) {
	d := date.UTC()
	txt := "AM"
	hour := d.Hour()
	if d.Hour() > 11 {
		txt = "PM"
	}
	if d.Hour() > 12 {
		hour -= 12
	}
	e.Time = fmt.Sprintf("%s, %s %d, %d %d:%02d:%02d %s GMT", d.Weekday(), d.Month(), d.Day(), d.Year(), hour, d.Minute(), d.Second(), txt)
}

// Validate returns an error if the event object is invalid
func (e Event) Validate() error {
	if e.UEI == "" {
		return fmt.Errorf("UEI cannot be null")
	}
	if e.LogMessage != nil {
		err := e.LogMessage.Validate()
		if err != nil {
			return err
		}
	}
	if e.Interface != "" {
		ip := net.ParseIP(e.Interface)
		if ip == nil {
			return fmt.Errorf("invalid Interface: %s", e.Interface)
		}
	}
	if e.Severity != "" {
		if err := Severities.Set(e.Severity); err != nil {
			return err
		}
	}
	return nil
}

// OnmsEventParam parameters of an OnmsEvent entity
type OnmsEventParam struct {
	Name  string
	Value string
	Type  string
}

// OnmsEvent OpenNMS event entity
type OnmsEvent struct {
	ID                   int              `json:"id" yaml:"id"`
	UEI                  string           `json:"uei" yaml:"uei"`
	EventTime            *Time            `json:"time,omitempty" yaml:"time,omitempty"`
	EventHost            string           `json:"host,omitempty" yaml:"host,omitempty"`
	EventSource          string           `json:"source,omitempty" yaml:"source,omitempty"`
	CreateTime           *Time            `json:"createTime,omitempty" yaml:"createTime,omitempty"`
	SnmpHost             string           `json:"snmpHost,omitempty" yaml:"snmpHost,omitempty"`
	Snmp                 string           `json:"snmp,omitempty" yaml:"snmp,omitempty"`
	NodeID               int              `json:"nodeId,omitempty" yaml:"nodeId,omitempty"`
	NodeLabel            string           `json:"nodeLabel,omitempty" yaml:"nodeLabel,omitempty"`
	IPAddress            string           `json:"ipAddress,omitempty" yaml:"ipAddress,omitempty"`
	ServiceType          OnmsServiceType  `json:"serviceType,omitempty" yaml:"serviceType,omitempty"`
	IfIndex              int              `json:"ifIndex,omitempty" yaml:"ifIndex,omitempty"`
	Severity             string           `json:"severity,omitempty" yaml:"severity,omitempty"`
	Log                  string           `json:"log,omitempty" yaml:"log,omitempty"`
	LogGroup             string           `json:"logGroup,omitempty" yaml:"logGroup,omitempty"`
	LogMessage           string           `json:"logMessage,omitempty" yaml:"logMessage,omitempty"`
	Display              string           `json:"display,omitempty" yaml:"display,omitempty"`
	Description          string           `json:"description,omitempty" yaml:"description,omitempty"`
	PathOutage           string           `json:"pathOutage,omitempty" yaml:"pathOutage,omitempty"`
	Correlation          string           `json:"correlation,omitempty" yaml:"correlation,omitempty"`
	SuppressedCount      int              `json:"suppressedCount,omitempty" yaml:"suppressedCount,omitempty"`
	OperatorInstructions string           `json:"operatorInstructions,omitempty" yaml:"operatorInstructions,omitempty"`
	OperatorAction       string           `json:"operatorAction,omitempty" yaml:"operatorAction,omitempty"`
	AutoAction           string           `json:"autoAction,omitempty" yaml:"autoAction,omitempty"`
	Parameters           []OnmsEventParam `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// OnmsEventList a list of events
type OnmsEventList struct {
	Count      int         `json:"count" yaml:"count"`
	TotalCount int         `json:"totalCount" yaml:"totalCount"`
	Offset     int         `json:"offset" yaml:"offset"`
	Events     []OnmsEvent `json:"event" yaml:"events"`
}
