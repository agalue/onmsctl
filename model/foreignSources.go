package model

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// ElementList a list of elements/strings
type ElementList struct {
	Count   int      `json:"count" yaml:"count"`
	Element []string `json:"element" yaml:"element"`
}

// Parameter a parameter for a detector or a policy
type Parameter struct {
	XMLName xml.Name `xml:"parameter" json:"-" yaml:"-"`
	Key     string   `xml:"key,attr" json:"key" yaml:"key"`
	Value   string   `xml:"value,attr" json:"value" yaml:"value"`
}

// Detector a provisioning detector
type Detector struct {
	XMLName    xml.Name    `xml:"detector" json:"-" yaml:"-"`
	Name       string      `xml:"name,attr" json:"name" yaml:"name"`
	Class      string      `xml:"class,attr" json:"class" yaml:"class"`
	Parameters []Parameter `xml:"parameter,omitempty" json:"parameter,omitempty" yaml:"parameters,omitempty"`
}

// Validate returns an error if the detector is invalid
func (p *Detector) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("detector name cannot be empty")
	}
	if p.Class == "" {
		return fmt.Errorf("detector class cannot be empty")
	}
	return nil
}

// Policy a provisioning policy
type Policy struct {
	XMLName    xml.Name    `xml:"policy" json:"-" yaml:"-"`
	Name       string      `xml:"name,attr" json:"name" yaml:"name"`
	Class      string      `xml:"class,attr" json:"class" yaml:"class"`
	Parameters []Parameter `xml:"parameter,omitempty" json:"parameter,omitempty" yaml:"parameters,omitempty"`
}

// Validate returns an error if the policy is invalid
func (p *Policy) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("policy name cannot be empty")
	}
	if p.Class == "" {
		return fmt.Errorf("policy class cannot be empty")
	}
	return nil
}

// ForeignSourceDef a foreign source definition
type ForeignSourceDef struct {
	XMLName      xml.Name   `xml:"foreign-source" json:"-" yaml:"-"`
	Name         string     `xml:"name,attr" json:"name" yaml:"name"`
	DateStamp    *Time      `xml:"date-stamp,attr,omitempty" json:"date-stamp,omitempty" yaml:"dateStamp,omitempty"`
	ScanInterval string     `xml:"scan-interval" json:"scan-interval" yaml:"scanInterval"`
	Detectors    []Detector `xml:"detectors>detector" json:"detectors,omitempty" yaml:"detectors,omitempty"`
	Policies     []Policy   `xml:"policies>policy" json:"policies,omitempty" yaml:"policies,omitempty"`
}

// Validate returns an error if the node definition is invalid
func (fs *ForeignSourceDef) Validate() error {
	if fs.Name == "" {
		return fmt.Errorf("the name of a Foreign Source definition cannot be empty")
	}
	if matched, _ := regexp.MatchString(`[/\\?:&*'"]`, fs.Name); matched {
		return fmt.Errorf("invalid characters on Foreign Source name %s:, /, \\, ?, &, *, ', \"", fs.Name)
	}
	if fs.ScanInterval == "" {
		return fmt.Errorf("the scan interval of a Foreign Source definition cannot be empty")
	}
	for !IsValidScanInterval(fs.ScanInterval) {
		return fmt.Errorf("invalid scan interval %s", fs.ScanInterval)
	}
	for _, d := range fs.Detectors {
		err := d.Validate()
		if err != nil {
			return err
		}
	}
	for _, p := range fs.Policies {
		err := p.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetDetector gets a detector by its name or class
func (fs ForeignSourceDef) GetDetector(detectorID string) (*Detector, error) {
	if detectorID == "" {
		return nil, fmt.Errorf("detector name or class required")
	}
	for _, detector := range fs.Detectors {
		if detector.Class == detectorID || detector.Name == detectorID {
			return &detector, nil
		}
	}
	return nil, fmt.Errorf("cannot find detector for %s", detectorID)
}

// GetPolicy gets a policy by its name or class
func (fs ForeignSourceDef) GetPolicy(policyID string) (*Policy, error) {
	if policyID == "" {
		return nil, fmt.Errorf("policy name or class required")
	}
	for _, policy := range fs.Policies {
		if policy.Class == policyID || policy.Name == policyID {
			return &policy, nil
		}
	}
	return nil, fmt.Errorf("cannot find policy for %s", policyID)
}

// Plugin a definition class for a detector or a policy
type Plugin struct {
	Name       string        `json:"name" yaml:"name"`
	Class      string        `json:"class" yaml:"class"`
	Parameters []PluginParam `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// FindParameter finds a parameter on the plugin
func (p Plugin) FindParameter(paramName string) *PluginParam {
	for _, param := range p.Parameters {
		if param.Key == paramName {
			return &param
		}
	}
	return nil
}

// VerifyParameters verify detector/policy parameters
func (p Plugin) VerifyParameters(parameters []Parameter) error {
	for _, param := range parameters {
		config := p.FindParameter(param.Key)
		if config == nil {
			return fmt.Errorf("invalid parameter %s for %s", param.Key, p.Class)
		}
	}
	for _, param := range p.Parameters {
		if param.Required {
			pa := FindParameter(parameters, param.Key)
			if pa == nil {
				return fmt.Errorf("missing required parameter %s on %s", param.Key, p.Class)
			}
			if len(param.Options) > 0 {
				found := false
				for _, opt := range param.Options {
					if opt == pa.Value {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("invalid parameter value %s on %s. Valid values are: %s", pa.Key, p.Class, param.Options)
				}
			}
		}
	}
	return nil
}

// PluginParam a parameter of a given plugin
type PluginParam struct {
	Key      string   `json:"key" yaml:"key"`
	Required bool     `json:"required" yaml:"required"`
	Options  []string `json:"options,omitempty" yaml:"options,omitempty"`
}

// PluginList a list of plugins
type PluginList struct {
	Count   int      `json:"count" yaml:"count"`
	Plugins []Plugin `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

// FindPlugin finds a plugin by class name
func (list PluginList) FindPlugin(cls string) *Plugin {
	for _, p := range list.Plugins {
		if p.Class == cls {
			return &p
		}
	}
	return nil
}

// FindParameter finds a parameter by name on a slice of parameters
func FindParameter(parameters []Parameter, paramName string) *Parameter {
	for _, param := range parameters {
		if param.Key == paramName {
			return &param
		}
	}
	return nil
}

// IsValidScanInterval checks if a given scan-interval is valid
func IsValidScanInterval(scanInterval string) bool {
	if scanInterval == "" {
		return false
	}
	re, _ := regexp.Compile(`^[0-9]+(w|d|h|m|s|ms)$`)
	for _, el := range strings.Split(scanInterval, " ") {
		if !re.MatchString(el) {
			return false
		}
	}
	return true
}
