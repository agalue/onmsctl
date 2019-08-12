package snmp

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// SNMPVersions the SNMP version enumeration
var SNMPVersions = &model.EnumValue{
	Enum:    []string{"v1", "v2c", "v3"},
	Default: "v2c",
}

// SNMPPrivProtocols the Private Protocols enumeration
var SNMPPrivProtocols = &model.EnumValue{
	Enum: []string{"DES", "AES", "AES192", "AES256"},
}

// SNMPAuthProtocols the Authentication Protocols enumeration
var SNMPAuthProtocols = &model.EnumValue{
	Enum: []string{"MD5", "SHA"},
}

// CliCommand the CLI command to provide server information
var CliCommand = cli.Command{
	Name:  "snmp",
	Usage: "Manage SNMP configuration",
	Subcommands: []cli.Command{
		{
			Name:      "get",
			Usage:     "Gets the SNMP configuration for a given IP address",
			ArgsUsage: "<ipAddress|fqdn>",
			Action:    showSnmpConfig,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "location, l",
					Usage: "Minion Location",
				},
			},
		},
		{
			Name:      "set",
			Usage:     "Sets the SNMP Configuration for a given IP address",
			ArgsUsage: "<ipAddress|fqdn>",
			Action:    setSnmpConfig,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name:  "version, v",
					Value: SNMPVersions,
					Usage: "SNMP Version: " + SNMPVersions.EnumAsString(),
				},
				cli.StringFlag{
					Name:  "location, l",
					Usage: "Minion Location",
				},
				cli.StringFlag{
					Name:  "port, p",
					Usage: "The UDP Port of the SNMP agent",
				},
				cli.StringFlag{
					Name:  "retry, r",
					Usage: "The number of retries before giving up",
				},
				cli.StringFlag{
					Name:  "timeout, t",
					Usage: "Timeout in milliseconds",
				},
				cli.StringFlag{
					Name:  "maxRepetitions, mr",
					Usage: "Maximum repetitions",
				},
				cli.StringFlag{
					Name:  "maxVarsPerPdu, mvpp",
					Usage: "Maximum variables per PDU",
				},
				cli.StringFlag{
					Name:  "community, c",
					Usage: "Community String for SNMPv1 or SNMPv2c",
				},
				cli.StringFlag{
					Name:  "securityName, sn",
					Usage: "SNMPv3 Security Name",
				},
				cli.StringFlag{
					Name:  "securityLevel, sl",
					Value: "1",
					Usage: "SNMPv3 Security Level: 1 noAuthNoPriv, 2: authNoPriv, 3: authPriv",
				},
				cli.GenericFlag{
					Name:  "privProtocol, pp",
					Value: SNMPPrivProtocols,
					Usage: "SNMPv3 Privacy Protocol: " + SNMPPrivProtocols.EnumAsString(),
				},
				cli.StringFlag{
					Name:  "privPassPhrase, ppp",
					Usage: "SNMPv3 Password Phrase for Privacy Protocol",
				},
				cli.GenericFlag{
					Name:  "authProtocol, ap",
					Value: SNMPAuthProtocols,
					Usage: "SNMPv3 Authentication Protocol: " + SNMPAuthProtocols.EnumAsString(),
				},
				cli.StringFlag{
					Name:  "authPassPhrase, app",
					Usage: "SNMPv3 Password Phrase for Authentication Protocol",
				},
				cli.StringFlag{
					Name:  "engineID, eid",
					Usage: "SNMPv3 Unique Engine ID of the SNMP agent",
				},
				cli.StringFlag{
					Name:  "contextEngineID, ceid",
					Usage: "SNMPv3 Context Engine ID",
				},
				cli.StringFlag{
					Name:  "enterpriseID, entid",
					Usage: "SNMPv3 Enterprise ID",
				},
				cli.StringFlag{
					Name:  "contextName, ctx",
					Usage: "SNMPv3 Context Name",
				},
			},
		},
		{
			Name:      "apply",
			Usage:     "Creates or updates the SNMP configuration for a given IP address",
			Action:    applySnmpConfig,
			ArgsUsage: "<ipAddress|fqdn> <yaml>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
		},
	},
}

func showSnmpConfig(c *cli.Context) error {
	ipAddress, err := getIPAddress(c)
	if err != nil {
		return err
	}
	url := "/rest/snmpConfig/" + ipAddress
	location := c.String("location")
	if location != "" {
		url += "?location=" + location
	}
	jsonString, err := rest.Instance.Get(url)
	if err != nil {
		return err
	}
	snmp := model.SnmpInfo{}
	json.Unmarshal(jsonString, &snmp)
	data, _ := yaml.Marshal(&snmp)
	fmt.Println(string(data))
	return nil
}

func setSnmpConfig(c *cli.Context) error {
	ipAddress, err := getIPAddress(c)
	if err != nil {
		return err
	}
	snmp := model.SnmpInfo{
		Version:         c.String("version"),
		Location:        c.String("location"),
		Port:            c.Int("port"),
		Retries:         c.Int("retries"),
		Timeout:         c.Int("timeout"),
		Community:       c.String("community"),
		ContextName:     c.String("contextName"),
		SecurityLevel:   c.Int("securityLevel"),
		SecurityName:    c.String("securityName"),
		PrivProtocol:    c.String("privProtocol"),
		PrivPassPhrase:  c.String("privPassPhrase"),
		AuthProtocol:    c.String("authProtocol"),
		AuthPassPhrase:  c.String("authPassPhrase"),
		EngineID:        c.String("engineID"),
		ContextEngineID: c.String("contextEngineID"),
		EnterpriseID:    c.String("enterpriseID"),
		MaxRequestSize:  c.Int("maxRequestSize"),
		MaxRepetitions:  c.Int("maxRepetitions"),
		MaxVarsPerPdu:   c.Int("maxVarsPerPdu"),
	}
	err = snmp.IsValid()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(snmp)
	return rest.Instance.Put("/rest/snmpConfig/"+ipAddress, jsonBytes, "application/json")
}

func applySnmpConfig(c *cli.Context) error {
	ipAddress, err := getIPAddress(c)
	if err != nil {
		return err
	}
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	snmp := model.SnmpInfo{}
	yaml.Unmarshal(data, &snmp)
	err = snmp.IsValid()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(snmp)
	return rest.Instance.Put("/rest/snmpConfig/"+ipAddress, jsonBytes, "application/json")
}

func getIPAddress(c *cli.Context) (string, error) {
	ipAddress := c.Args().Get(0)
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
