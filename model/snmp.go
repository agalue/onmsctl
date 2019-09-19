package model

import (
	"fmt"
)

// SNMPVersions the SNMP version enumeration
var SNMPVersions = &EnumValue{
	Enum:    []string{"v1", "v2c", "v3"},
	Default: "v2c",
}

// SNMPPrivProtocols the Private Protocols enumeration
var SNMPPrivProtocols = &EnumValue{
	Enum: []string{"DES", "AES", "AES192", "AES256"},
}

// SNMPAuthProtocols the Authentication Protocols enumeration
var SNMPAuthProtocols = &EnumValue{
	Enum: []string{"MD5", "SHA"},
}

// SnmpInfo SNMP Configuration for a give IP Interface;
// it provides partial information compared with what's available on snmp-config.xml
type SnmpInfo struct {
	Version         string `json:"version,omitempty" yaml:"version,omitempty"`
	Location        string `json:"location,omitempty" yaml:"location,omitempty"`
	Port            int    `json:"port,omitempty" yaml:"port,omitempty"`
	Retries         int    `json:"retries,omitempty" yaml:"retries,omitempty"`
	Timeout         int    `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Community       string `json:"community,omitempty" yaml:"community,omitempty"`
	ContextName     string `json:"contextName,omitempty" yaml:"contextName,omitempty"`
	SecurityLevel   int    `json:"securityLevel,omitempty" yaml:"securityLevel,omitempty"`
	SecurityName    string `json:"securityName,omitempty" yaml:"securityName,omitempty"`
	PrivProtocol    string `json:"privProtocol,omitempty" yaml:"privProtocol,omitempty"`
	PrivPassPhrase  string `json:"privPassPhrase,omitempty" yaml:"privPassPhrase,omitempty"`
	AuthProtocol    string `json:"authProtocol,omitempty" yaml:"authProtocol,omitempty"`
	AuthPassPhrase  string `json:"authPassPhrase,omitempty" yaml:"authPassPhrase,omitempty"`
	EngineID        string `json:"engineID,omitempty" yaml:"engineID,omitempty"`
	ContextEngineID string `json:"ContextEngineID,omitempty" yaml:"ContextEngineID,omitempty"`
	EnterpriseID    string `json:"enterpriseID,omitempty" yaml:"enterpriseID,omitempty"`
	MaxRequestSize  int    `json:"maxRequestSize,omitempty" yaml:"maxRequestSize,omitempty"`
	MaxRepetitions  int    `json:"maxRepetitions,omitempty" yaml:"maxRepetitions,omitempty"`
	MaxVarsPerPdu   int    `json:"maxVarsPerPdu,omitempty" yaml:"maxVarsPerPdu,omitempty"`
	ProxyHost       string `json:"proxyHost,omitempty" yaml:"proxyHost,omitempty"`
	TTL             int    `json:"ttl,omitempty" yaml:"ttl,omitempty"`
}

// IsValid returns an error if the service is invalid
func (s *SnmpInfo) IsValid() error {
	if s.Community == "" {
		return fmt.Errorf("SNMP Community String cannot be null")
	}
	if s.SecurityLevel != 0 {
		if s.SecurityLevel < 0 || s.SecurityLevel > 3 {
			return fmt.Errorf("Invalid Security Level. Allowed values: 1, 2, or 3")
		}
		if s.Version != "v3" {
			s.SecurityLevel = 0
		}
	}
	if s.Version != "" {
		if err := SNMPVersions.Set(s.Version); err != nil {
			return fmt.Errorf("Invalid SNMP Version. Allowed values: %s", SNMPVersions.EnumAsString())
		}
	}
	if s.PrivProtocol != "" {
		if err := SNMPPrivProtocols.Set(s.PrivProtocol); err != nil {
			return fmt.Errorf("Invalid Priv Protocol. Allowed values: %s", SNMPPrivProtocols.EnumAsString())
		}
	}
	if s.AuthProtocol != "" {
		if err := SNMPAuthProtocols.Set(s.AuthProtocol); err != nil {
			return fmt.Errorf("Invalid Auth Protocol. Allowed values: %s", SNMPAuthProtocols.EnumAsString())
		}
	}
	return nil
}
