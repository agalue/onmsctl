package snmp

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
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
					Name:  "version, v",
					Value: model.SNMPVersions,
					Usage: "SNMP Version: " + model.SNMPVersions.EnumAsString(),
				},
				cli.StringFlag{
					Name:  "location, l",
					Usage: "Minion Location",
				},
				cli.IntFlag{
					Name:  "port, p",
					Value: 161,
					Usage: "The UDP Port of the SNMP agent",
				},
				cli.IntFlag{
					Name:  "retry, r",
					Value: 2,
					Usage: "The number of retries before giving up",
				},
				cli.IntFlag{
					Name:  "timeout, t",
					Value: 1800,
					Usage: "Timeout in milliseconds",
				},
				cli.IntFlag{
					Name:  "maxRepetitions, mr",
					Value: 2,
					Usage: "Maximum repetitions",
				},
				cli.IntFlag{
					Name:  "maxVarsPerPdu, mvpp",
					Value: 10,
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
				cli.IntFlag{
					Name:  "securityLevel, sl",
					Value: 1,
					Usage: "SNMPv3 Security Level: 1 noAuthNoPriv, 2: authNoPriv, 3: authPriv",
				},
				cli.GenericFlag{
					Name:  "privProtocol, pp",
					Value: model.SNMPPrivProtocols,
					Usage: "SNMPv3 Privacy Protocol: " + model.SNMPPrivProtocols.EnumAsString(),
				},
				cli.StringFlag{
					Name:  "privPassPhrase, ppp",
					Usage: "SNMPv3 Password Phrase for Privacy Protocol",
				},
				cli.GenericFlag{
					Name:  "authProtocol, ap",
					Value: model.SNMPAuthProtocols,
					Usage: "SNMPv3 Authentication Protocol: " + model.SNMPAuthProtocols.EnumAsString(),
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
	snmp, err := getAPI().GetConfig(c.Args().Get(0), c.String("location"))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(&snmp)
	fmt.Println(string(data))
	return nil
}

func setSnmpConfig(c *cli.Context) error {
	snmp := model.SnmpInfo{
		Version:         c.String("version"),
		Location:        c.String("location"),
		Port:            c.Int("port"),
		Retries:         c.Int("retry"),
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
	if err := checkLocation(snmp); err != nil {
		return err
	}
	return getAPI().SetConfig(c.Args().Get(0), snmp)
}

func applySnmpConfig(c *cli.Context) error {
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	snmp := model.SnmpInfo{}
	err = yaml.Unmarshal(data, &snmp)
	if err != nil {
		return err
	}
	if err = checkLocation(snmp); err != nil {
		return err
	}
	return getAPI().SetConfig(c.Args().Get(0), snmp)
}

func checkLocation(snmp model.SnmpInfo) error {
	if snmp.Location != "" {
		ok, err := getMonitoringLocationsAPI().LocationExists(snmp.Location)
		if err == nil && !ok {
			return fmt.Errorf("Location %s doesn't exists", snmp.Location)
		}
	}
	return nil
}

func getAPI() api.SnmpAPI {
	return services.GetSnmpAPI(rest.Instance)
}

func getMonitoringLocationsAPI() api.MonitoringLocationsAPI {
	return services.GetMonitoringLocationsAPI(rest.Instance)
}
