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
					Name: "version, v",
					Value: &model.EnumValue{
						Enum:    []string{"v1", "v2c", "v3"},
						Default: "v2c",
					},
					Usage: "SNMP Version: v1, v2c, v3",
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
					Name: "privProtocol, pp",
					Value: &model.EnumValue{
						Enum: []string{"DES", "AES", "AES192", "AES256"},
					},
					Usage: "SNMPv3 Privacy Protocol: DES, AES, AES192, AES256",
				},
				cli.StringFlag{
					Name:  "privPassPhrase, ppp",
					Usage: "SNMPv3 Password Phrase for Privacy Protocol",
				},
				cli.GenericFlag{
					Name: "authProtocol, ap",
					Value: &model.EnumValue{
						Enum: []string{"MD5", "SHA"},
					},
					Usage: "SNMPv3 Authentication Protocol: MD5, SHA",
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
	if !c.Args().Present() {
		return fmt.Errorf("IP Address required")
	}
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
	if !c.Args().Present() {
		return fmt.Errorf("IP Address required")
	}
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
	if !c.Args().Present() {
		return fmt.Errorf("IP Address required")
	}
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
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		addresses, err := net.LookupIP(ipAddress)
		fmt.Printf("%s translates to %s, using the first entry.\n", ipAddress, addresses)
		if err != nil {
			return "", fmt.Errorf("Cannot get address from %s: %s", ipAddress, err)
		}
		ipAddress = addresses[0].String()
	}
	return ipAddress, nil
}
