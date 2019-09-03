package api

import "github.com/OpenNMS/onmsctl/model"

// SnmpAPI the API to manipulate SNMP Configuration
type SnmpAPI interface {
	GetConfig(ipAddress string, location string) (*model.SnmpInfo, error)
	SetConfig(ipAddress string, config model.SnmpInfo) error
}
