package events

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

var severities = &model.EnumValue{
	Enum: model.Severities.Enum,
}

var destinations = &model.EnumValue{
	Enum:    model.Destinations.Enum,
	Default: model.Destinations.Enum[0],
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
					Name:  "host, H",
					Usage: "IP address or FQDN of the host that sends the event",
				},
				cli.Int64Flag{
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
				cli.IntFlag{
					Name:  "ifindex, f",
					Usage: "ifIndex of the interface",
				},
				cli.StringFlag{
					Name:  "descr, d",
					Usage: "A description for the event",
				},
				//              Unfortunately, it is not possible to override the logMsg field.
				//				cli.StringFlag{
				//					Name:  "logmsg, m",
				//					Usage: "The log message for the event",
				//				},
				cli.GenericFlag{
					Name:  "severity, x",
					Value: severities,
					Usage: "The severity of the event: " + severities.EnumAsString(),
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
		Host:        c.String("host"),
		Source:      "onmsctl",
	}
	params := c.StringSlice("parm")
	for _, p := range params {
		data := strings.Split(p, "=")
		if len(data) == 2 {
			event.AddParameter(data[0], data[1])
		}
	}
	logmsg := c.String("logmsg")
	if logmsg != "" {
		event.LogMessage = &model.LogMsg{
			Message: logmsg,
		}
	}
	return getAPI().SendEvent(event)
}

func applyEvent(c *cli.Context) error {
	data, err := common.ReadInput(c, 0)
	if err != nil {
		return err
	}
	event := model.Event{}
	if err := yaml.Unmarshal(data, &event); err != nil {
		return err
	}
	if err := event.Validate(); err != nil {
		return err
	}
	return getAPI().SendEvent(event)
}

func getAPI() api.EventsAPI {
	return services.GetEventsAPI(rest.Instance)
}
