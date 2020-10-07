package nodes

import (
	"fmt"
	"strconv"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// SnmpInterfacesCliCommand the CLI command to manage nodes
var SnmpInterfacesCliCommand = cli.Command{
	Name:  "snmpInterfaces",
	Usage: "Manage SNMP Interfaces",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "Lists all the SNMP interfaces for a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID>",
			Action:    getSnmpInterfaces,
		},
		{
			Name:      "add",
			Usage:     "Add a new SNMP interface to a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID>",
			Action:    addSnmpInterface,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "ifIndex, i",
					Usage: "The IF-MIB::ifIndex",
				},
				cli.IntFlag{
					Name:  "ifOper, o",
					Usage: "The IF-MIB::ifOperStatus",
				},
				cli.IntFlag{
					Name:  "ifAdmin, A",
					Usage: "The IF-MIB::ifAdminStatus",
				},
				cli.Int64Flag{
					Name:  "ifSpeed, s",
					Usage: "The IF-MIB::ifSpeed expressed in bits per second",
				},
				cli.IntFlag{
					Name:  "ifType, t",
					Usage: "The IF-MIB::ifType",
				},
				cli.StringFlag{
					Name:  "ifName, n",
					Usage: "The IF-MIB::ifName",
				},
				cli.StringFlag{
					Name:  "ifDescr, d",
					Usage: "The IF-MIB::ifDescr",
				},
				cli.StringFlag{
					Name:  "ifAlias, a",
					Usage: "The IF-MIB::ifAlias",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing SNMP interface from a given node (cannot be undone)",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <ifIndex>",
			Action:    deleteSnmpInterface,
		},
	},
}

func getSnmpInterfaces(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	list, err := services.GetNodesAPI(rest.Instance).GetSnmpInterfaces(criteria)
	if err != nil {
		return err
	}
	if len(list.Interfaces) == 0 {
		fmt.Println("There are no SNMP interfaces")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "ID\tifIndex\tifDescr\tifName\tifAlias\tifType\tifOperStatus\tifAdminStatus")
	for _, n := range list.Interfaces {
		fmt.Fprintf(writer, "%d\t%d\t%s\t%s\t%s\t%d\t%d\t%d\n", n.ID, n.IfIndex, n.IfDescr, n.IfName, n.IfAlias, n.IfType, n.IfOperStatus, n.IfAdminStatus)
	}
	writer.Flush()
	return nil
}

func deleteSnmpInterface(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	idx := c.Args().Get(1)
	if idx == "" {
		return fmt.Errorf("ifIndex is required")
	}
	ifIndex, err := strconv.Atoi(idx)
	if err != nil {
		return fmt.Errorf("Cannot parse ifIndex: %s", idx)
	}
	return services.GetNodesAPI(rest.Instance).DeleteSnmpInterface(criteria, ifIndex)
}

func addSnmpInterface(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	snmp := &model.OnmsSnmpInterface{
		IfIndex:       c.Int("ifIndex"),
		IfOperStatus:  c.Int("ifOper"),
		IfAdminStatus: c.Int("ifAdmin"),
		IfType:        c.Int("ifType"),
		IfSpeed:       c.Int64("ifSpeed"),
		IfName:        c.String("ifName"),
		IfDescr:       c.String("ifDescr"),
		IfAlias:       c.String("ifAlias"),
	}
	if err := snmp.Validate(); err != nil {
		return err
	}
	return services.GetNodesAPI(rest.Instance).SetSnmpInterface(criteria, snmp)
}