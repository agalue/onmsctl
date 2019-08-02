package model

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// EnumValue a enumaration array of strings
type EnumValue struct {
	Enum     []string
	Default  string
	selected string
}

// Set sets a value of the enum
func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.selected = value
			return nil
		}
	}
	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

// String gets the value of the enum as string
func (e EnumValue) String() string {
	if e.selected == "" {
		return e.Default
	}
	return e.selected
}

// Time an object to seamlessly manage times in multiple formats
type Time struct {
	time.Time
}

// MarshalJSON converts time object into timestamp in milliseconds
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(t.UnixNano() / int64(time.Millisecond))
}

// UnmarshalJSON converts timestamp in milliseconds into time object
func (t *Time) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	sec := i / 1000
	nsec := (i % 1000) * 1000
	t.Time = time.Unix(sec, nsec)
	return nil
}

// MarshalYAML converts time object into time as string
func (t Time) MarshalYAML() ([]byte, error) {
	if t.IsZero() {
		return json.Marshal("")
	}
	return yaml.Marshal(t.String())
}

// UnmarshalYAML converts time string into time object
func (t *Time) UnmarshalYAML(data []byte) error {
	var s string
	var err error
	if err = yaml.Unmarshal(data, &s); err != nil {
		return err
	}
	t.Time, err = time.Parse(s, s)
	if err != nil {
		return err
	}
	return nil
}

// MarshalXML converts time object into time as string
func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if t.IsZero() {
		return e.EncodeElement("", start)
	}
	return e.EncodeElement(t.String(), start)
}

// UnmarshalXML converts time string into time object
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	var err error
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	t.Time, err = time.Parse(s, s)
	if err != nil {
		return err
	}
	return nil
}
