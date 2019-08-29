package events

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

var api = services.GetEventsAPI(rest.Instance)

// Severities list of valid event severities
var Severities = &model.EnumValue{
	Enum: []string{"Indeterminate", "Normal", "Warning", "Minor", "Major", "Critical"},
}

// CliCommand the CLI command to manage events
var CliCommand = cli.Command{
	Name:  "events",
	Usage: "Manage events",
	Subcommands: []cli.Command{
		{
			Name:      "send",
			Usage:     "Sends an event to OpenNMS",
			ArgsUsage: "<uei>",
			Action:    sendEvent,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "nodeid, n",
					Usage: "The numeric node identifier",
				},
				cli.StringFlag{
					Name:  "interface, i",
					Usage: "IP address of the interface",
				},
				cli.StringFlag{
					Name:  "service, s",
					Usage: "Service name",
				},
				cli.StringFlag{
					Name:  "ifindex, f",
					Usage: "ifIndex of the interface",
				},
				cli.StringFlag{
					Name:  "descr, d",
					Usage: "A description for the event browser",
				},
				cli.GenericFlag{
					Name:  "severity, x",
					Value: Severities,
					Usage: "The severity of the event: " + Severities.EnumAsString(),
				},
				cli.StringSliceFlag{
					Name:  "parm, p",
					Usage: "An event parameter (e.x. --parm 'url=http://www.google.com/')",
				},
			},
		},
		{
			Name:      "apply",
			Usage:     "Sends an event to OpenNMS in YAML format",
			Action:    applyEvent,
			ArgsUsage: "<yaml>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
		},
	},
}

func sendEvent(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("UEI required")
	}
	uei := c.Args().First()
	event := model.Event{
		UEI:         uei,
		NodeID:      c.Int64("nodeid"),
		Interface:   c.String("interface"),
		Service:     c.String("service"),
		IfIndex:     c.Int("ifindex"),
		Description: c.String("descr"),
		Severity:    c.String("severity"),
		Source:      "onmsctl",
	}
	params := c.StringSlice("parm")
	for _, p := range params {
		data := strings.Split(p, "=")
		event.AddParameter(data[0], data[1])
	}
	return api.SendEvent(event)
}

func applyEvent(c *cli.Context) error {
	data, err := common.ReadInput(c, 0)
	if err != nil {
		return err
	}
	event := model.Event{}
	yaml.Unmarshal(data, &event)
	err = event.IsValid()
	if err != nil {
		return err
	}
	return api.SendEvent(event)
}
