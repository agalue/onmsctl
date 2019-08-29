package model

import (
	"fmt"
)

// SnmpInfo SNMP Configuration for a give IP Interface
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
		if s.Version != "v1" && s.Version != "v2c" && s.Version != "v3" {
			return fmt.Errorf("Invalid SNMP Version. Allowed values: v1, v2c, or v3")
		}
	}
	if s.PrivProtocol != "" {
		if s.PrivProtocol != "DES" && s.PrivProtocol != "AES" && s.PrivProtocol != "AES192" && s.PrivProtocol != "AES256" {
			return fmt.Errorf("Invalid Priv Protocol. Allowed values: DES, AES, AES192, AES256")
		}
	}
	if s.AuthProtocol != "" {
		if s.AuthProtocol != "MD5" && s.AuthProtocol != "SHA" {
			return fmt.Errorf("Invalid Auth Protocol. Allowed values: MD5, SHA")
		}
	}
	return nil
}
