package events

import (
	"fmt"
	"net"
	"time"

	"github.com/OpenNMS/onmsctl/common"
)

// SNMP an event SNMP object
type SNMP struct {
	ID        string       `json:"id" yaml:"id"`
	Version   string       `json:"version,omitempty" yaml:"version,omitempty"`
	Specific  int          `json:"specific" yaml:"specific,omitempty"`
	Generic   int          `json:"generic" yaml:"generic,omitempty"`
	Community string       `json:"community,omitempty" yaml:"community,omitempty"`
	Timestamp *common.Time `json:"time-stamp,omitempty" yaml:"timeStamp,omitempty"`
}

// Parameter an event parameter object
type Parameter struct {
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

// IsValid returns an error if the log message is invalid
func (lm *LogMsg) IsValid() error {
	if lm.Message == "" {
		return fmt.Errorf("Message cannot be null")
	}
	if lm.Destination == "" {
		lm.Destination = "logndisplay"
	}
	return nil
}

// Event an event object
// Time uses a string format. Example: "Saturday, July 13, 2019 2:13:43 PM GMT"
type Event struct {
	SnmpMask      *Mask       `json:"mask,omitempty" json:"mask,omitempty"`
	Snmp          *SNMP       `json:"snmp,omitempty" json:"snmp,omitempty"`
	LogMessage    *LogMsg     `json:"logmsg,omitempty" json:"logmsg,omitempty"`
	UEI           string      `json:"uei" yaml:"uei"`
	Source        string      `json:"source" yaml:"source"`
	Time          string      `json:"time,omitempty" yaml:"time,omitempty"`
	Host          string      `json:"host,omitempty" yaml:"host,omitempty"`
	MasterStation string      `json:"master-station,omitempty" yaml:"masterStation,omitempty"`
	NodeID        int64       `json:"nodeid,omitempty" yaml:"nodeID,omitempty"`
	Interface     string      `json:"interface,omitempty" yaml:"interface,omitempty"`
	Service       string      `json:"service,omitempty" yaml:"service,omitempty"`
	IfIndex       int         `json:"ifindex,omitempty" yaml:"ifIndex,omitempty"`
	SnmpHost      string      `json:"snmphost,omitempty" yaml:"snmpHost,omitempty"`
	Parameters    []Parameter `json:"parms,omitempty" yaml:"parameters,omitempty"`
	Description   string      `json:"descr,omitempty" yaml:"description,omitempty"`
	Severity      string      `json:"severity,omitempty" yaml:"severity,omitempty"`
	PathOutage    string      `json:"pathoutage,omitempty" yaml:"pathOutage,omitempty"`
	OperInstruct  string      `json:"operinstruct,omitempty" yaml:"operInstruct,omitempty"`
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

// IsValid returns an error if the event object is invalid
func (e *Event) IsValid() error {
	if e.UEI == "" {
		return fmt.Errorf("UEI cannot be null")
	}
	if e.LogMessage != nil {
		err := e.LogMessage.IsValid()
		if err != nil {
			return err
		}
	}
	if e.Interface != "" {
		ip := net.ParseIP(e.Interface)
		if ip == nil {
			return fmt.Errorf("Invalid Interface: %s", e.Interface)
		}
	}
	if e.Severity != "" {
		severities := common.EnumValue{
			Enum: []string{"Indeterminate", "Normal", "Warning", "Minor", "Major", "Critical"},
		}
		err := severities.Set(e.Severity)
		if err != nil {
			return err
		}
	}
	return nil
}
