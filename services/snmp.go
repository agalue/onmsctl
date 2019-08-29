package services

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

type snmpAPI struct {
	rest api.RestAPI
}

// GetSnmpAPI Obtain an implementation of the SNMP API
func GetSnmpAPI(rest api.RestAPI) api.SnmpAPI {
	return &snmpAPI{rest}
}

func (api snmpAPI) GetConfig(ipAddress string, location string) (*model.SnmpInfo, error) {
	ipAddress, err := api.validateIPAddress(ipAddress)
	if err != nil {
		return nil, err
	}
	url := "/rest/snmpConfig/" + ipAddress
	if location != "" {
		url += "?location=" + location
	}
	jsonString, err := api.rest.Get(url)
	if err != nil {
		return nil, err
	}
	snmp := &model.SnmpInfo{}
	if err := json.Unmarshal(jsonString, snmp); err != nil {
		return nil, err
	}
	return snmp, nil
}

func (api snmpAPI) SetConfig(ipAddress string, config model.SnmpInfo) error {
	ipAddress, err := api.validateIPAddress(ipAddress)
	if err != nil {
		return err
	}
	err = config.IsValid()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(config)
	return api.rest.Put("/rest/snmpConfig/"+ipAddress, jsonBytes, "application/json")
}

func (api snmpAPI) validateIPAddress(ipAddress string) (string, error) {
	if ipAddress == "" {
		return "", fmt.Errorf("IP Address or FQDN required")
	}
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		addresses, err := net.LookupIP(ipAddress)
		if err != nil || len(addresses) == 0 {
			return "", fmt.Errorf("Cannot parse address from %s (invalid IP or FQDN); %s", ipAddress, err)
		}
		fmt.Printf("%s translates to %s\n", ipAddress, addresses[0].String())
		ipAddress = addresses[0].String()
	}
	return ipAddress, nil
}
