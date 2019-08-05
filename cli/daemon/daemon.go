package daemon

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

// DaemonMap a map with reloadable daemons
var DaemonMap = map[string]string{
	"ackd":                               "Ackd",
	"alarmd":                             "alarmd",
	"bsmd":                               "Bsmd",
	"collectd":                           "Collectd",
	"correlation:<engine-name>":          "DroolsCorrelationEngine", // Append engine name
	"discoverd":                          "Discovery",
	"enlinkd":                            "Enlinkd",
	"eventd":                             "Eventd",
	"ticketd":                            "Ticketd",
	"syslogd":                            "syslogd",
	"trapd":                              "trapd",
	"telemetryd":                         "telemetryd",
	"nbi:email":                          "EmailNBI",
	"nbi:snmptrap":                       "SnmpTrapNBI",
	"nbi:syslog":                         "SyslogNBI",
	"notifd":                             "Notifd",
	"reportd":                            "Reportd",
	"pollerd":                            "Pollerd",
	"poller-backend":                     "PollerBackEnd",
	"provisiond":                         "Provisiond",
	"scriptd":                            "Scriptd",
	"statsd":                             "Statsd",
	"tl1d":                               "Tl1d",
	"threshd":                            "Threshd",
	"translator":                         "Translator",
	"provisiond:snmp-asset":              "Provisiond.SnmpAssetProvisioningAdapter",
	"provisiond:snmp-hardware-inventory": "Provisiond.SnmpHardwareInventoryProvisioningAdapter",
	"provisiond:wsman":                   "WsManAssetProvisioningAdapter",
	"vacuumd":                            "Vacuumd",
}

// CliCommand the CLI command to manage events
var CliCommand = cli.Command{
	Name:  "daemon",
	Usage: "Manage OpenNMS daemons",
	Subcommands: []cli.Command{
		{
			Name:      "reload",
			Usage:     "Request reload the configuration of a given OpenNMS daemon",
			ArgsUsage: "<daemonName>",
			Action:    reloadDaemon,
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
		return fmt.Errorf("Daemon name required")
	}
	daemonName := c.Args().First()
	if !isValidDaemon(daemonName) {
		return fmt.Errorf("Invalid daemon name %s", daemonName)
	}
	event := model.Event{
		UEI:    "uei.opennms.org/internal/reloadDaemonConfig",
		Source: "onmsctl",
	}
	event.AddParameter("daemonName", DaemonMap[daemonName])
	configFile := c.String("configFile")
	if configFile != "" {
		event.AddParameter("configFile", configFile)
	}
	jsonBytes, _ := json.Marshal(event)
	return rest.Instance.Post("/rest/events/", jsonBytes)
}

func showReloadableDaemons(c *cli.Context) error {
	keys := make([]string, 0, len(DaemonMap))
	for k := range DaemonMap {
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
	if _, ok := DaemonMap[name]; ok {
		return true
	} else if strings.HasPrefix(name, "correlation:") {
		return true
	}
	return false
}
