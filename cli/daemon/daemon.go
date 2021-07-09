package daemon

import (
	"fmt"
	"sort"
	"strings"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// The prefix for correlation engines
const correlatorPrefix = "correlation"

// A map with reloadable daemons
var daemonMap = map[string]string{
	"ackd":                               "Ackd",
	"alarmd":                             "alarmd",
	"bsmd":                               "Bsmd",
	"collectd":                           "Collectd",
	correlatorPrefix:                     "DroolsCorrelationEngine", // Append engine name
	"discoverd":                          "Discovery",
	"enlinkd":                            "Enlinkd",
	"eventd":                             "Eventd",
	"ticketd":                            "Ticketd",
	"syslogd":                            "syslogd",
	"trapd":                              "trapd",
	"telemetryd":                         "telemetryd",
	"nbi-email":                          "EmailNBI",
	"nbi-snmptrap":                       "SnmpTrapNBI",
	"nbi-syslog":                         "SyslogNBI",
	"notifd":                             "Notifd",
	"reportd":                            "Reportd",
	"pollerd":                            "Pollerd",
	"poller-backend":                     "PollerBackEnd",
	"provisiond":                         "Provisiond",
	"provisiond-snmp-asset":              "Provisiond.SnmpAssetProvisioningAdapter",
	"provisiond-snmp-hardware-inventory": "Provisiond.SnmpHardwareInventoryProvisioningAdapter",
	"provisiond-wsman":                   "WsManAssetProvisioningAdapter",
	"scriptd":                            "Scriptd",
	"statsd":                             "Statsd",
	"tl1d":                               "Tl1d",
	"threshd":                            "Threshd",
	"translator":                         "Translator",
	"vacuumd":                            "Vacuumd",
}

// CliCommand the CLI command to manage events
var CliCommand = cli.Command{
	Name:  "daemon",
	Usage: "Manage OpenNMS Daemons",
	Subcommands: []cli.Command{
		{
			Name:         "reload",
			Usage:        "Request reload the configuration of a given OpenNMS daemon",
			ArgsUsage:    "<daemonName>",
			Action:       reloadDaemon,
			BashComplete: reloadBashComplete,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "configFile, f",
					Usage: "Configuration File (used by a few daemons)",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "Show a list of reloadable daemons",
			Action: showReloadableDaemons,
		},
	},
}

func reloadDaemon(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("daemon name required")
	}
	daemonName := c.Args().First()
	if !isValidDaemon(daemonName) {
		return fmt.Errorf("invalid daemon name %s", daemonName)
	}
	event := model.Event{
		UEI:    "uei.opennms.org/internal/reloadDaemonConfig",
		Source: "onmsctl",
	}
	event.AddParameter("daemonName", getDaemonName(daemonName))
	configFile := c.String("configFile")
	if configFile != "" {
		event.AddParameter("configFile", configFile)
	}
	return services.GetEventsAPI(rest.Instance).SendEvent(event)
}

func reloadBashComplete(c *cli.Context) {
	if c.NArg() > 0 {
		return
	}
	for k := range daemonMap {
		fmt.Println(k)
	}
}

func showReloadableDaemons(c *cli.Context) error {
	keys := make([]string, 0, len(daemonMap))
	for k := range daemonMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k)
	}
	return nil
}

func isValidDaemon(daemonName string) bool {
	name := strings.ToLower(daemonName)
	if _, ok := daemonMap[name]; ok {
		return true
	} else if strings.HasPrefix(name, correlatorPrefix) {
		return true
	}
	return false
}

func getDaemonName(id string) string {
	if strings.HasPrefix(id, correlatorPrefix) {
		data := strings.Split(id, ":")
		if len(data) == 2 {
			return daemonMap[correlatorPrefix] + ":" + data[1]
		}
		return daemonMap[correlatorPrefix]
	}
	return daemonMap[id]
}
