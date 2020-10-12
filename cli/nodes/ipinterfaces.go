package nodes

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// IPInterfacesCliCommand the CLI command to manage nodes
var IPInterfacesCliCommand = cli.Command{
	Name:  "ipInterfaces",
	Usage: "Manage IP Interfaces",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "Lists all the IP interfaces for a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID>",
			Action:    getIPInterfaces,
		},
		{
			Name:      "add",
			Usage:     "Add a new IP interface to a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID>",
			Action:    addIPInterface,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ipAddr, i",
					Usage: "IP Address",
				},
				cli.StringFlag{
					Name:  "hostname, n",
					Usage: "Hostname or FQDN",
				},
				cli.StringFlag{
					Name:  "isManaged, m",
					Usage: "Managed Flag",
				},
				cli.StringFlag{
					Name:  "snmpPrimary, p",
					Usage: "SNMP Primary Flag",
				},
				cli.IntFlag{
					Name:  "ifIndex, I",
					Usage: "ifIndex of the existing SNMP interface to associte with the IP interface",
				},
			},
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing IP interface from a given node (cannot be undone)",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress>",
			Action:    deleteIPInterface,
		},
		{
			Name:  "metadata",
			Usage: "Manage interface-level metadata",
			Subcommands: []cli.Command{
				{
					Name:      "list",
					Usage:     "Lists all the interface-level metadata",
					Action:    listInterfaceMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress>",
				},
				{
					Name:      "set",
					Usage:     "Adds or updates a interface-level metadata entry",
					Action:    setInterfaceMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "context, c",
							Usage: "Metadata Context",
						},
						cli.StringFlag{
							Name:  "key, k",
							Usage: "Metadata Key",
						},
						cli.StringFlag{
							Name:  "value, v",
							Usage: "Metadata Value",
						},
					},
				},
				{
					Name:      "delete",
					Usage:     "Deletes an existing interface-level metadata entry",
					Action:    deleteInterfaceMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "context, c",
							Usage: "Metadata Context",
						},
						cli.StringFlag{
							Name:  "key, k",
							Usage: "Metadata Key",
						},
					},
				},
			},
		},
		ServicesCliCommand,
	},
}

func getIPInterfaces(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	list, err := services.GetNodesAPI(rest.Instance).GetIPInterfaces(criteria)
	if err != nil {
		return err
	}
	if len(list.Interfaces) == 0 {
		fmt.Println("There are no IP interfaces")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "ID\tIP Address\tifIndex\tIs Managed\tSNMP Primary\tIs Down")
	for _, n := range list.Interfaces {
		fmt.Fprintf(writer, "%s\t%s\t%d\t%s\t%s\t%v\n", n.ID, n.IPAddress, n.IfIndex, n.IsManaged, n.SnmpPrimary, n.IsDown)
	}
	writer.Flush()
	return nil
}

func deleteIPInterface(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if ipaddr == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	return services.GetNodesAPI(rest.Instance).DeleteIPInterface(criteria, ipaddr)
}

func addIPInterface(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ip := &model.OnmsIPInterface{
		IPAddress:   c.String("ipAddr"),
		HostName:    c.String("hostname"),
		IsManaged:   c.String("isManaged"),
		SnmpPrimary: c.String("snmpInterface"),
		IfIndex:     c.Int("ifIndex"),
	}
	if err := ip.Validate(); err != nil {
		return err
	}
	return services.GetNodesAPI(rest.Instance).SetIPInterface(criteria, ip)
}

func listInterfaceMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if criteria == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	meta, err := services.GetNodesAPI(rest.Instance).GetIPInterfaceMetadata(criteria, ipaddr)
	if err != nil {
		return err
	}
	if len(meta) == 0 {
		fmt.Println("There is no metadata")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Context\tKey\tValue")
	for _, m := range meta {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", m.Context, m.Key, m.Value)
	}
	writer.Flush()
	return nil
}

func setInterfaceMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if ipaddr == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	meta := model.MetaData{
		Context: c.String("context"),
		Key:     c.String("key"),
		Value:   c.String("value"),
	}
	if err := meta.Validate(); err != nil {
		return err
	}
	return services.GetNodesAPI(rest.Instance).SetIPInterfaceMetadata(criteria, ipaddr, meta)
}

func deleteInterfaceMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if ipaddr == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	ctx := c.String("context")
	if ctx == "" {
		return fmt.Errorf("Context is required")
	}
	key := c.String("key")
	if key == "" {
		return fmt.Errorf("Key is required")
	}
	return services.GetNodesAPI(rest.Instance).DeleteIPInterfaceMetadata(criteria, ipaddr, ctx, key)
}
