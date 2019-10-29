package provisioning

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// InterfacesCliCommand the CLI command configuration for managing IP interfaces on requisitioned nodes
var InterfacesCliCommand = cli.Command{
	Name:      "interface",
	ShortName: "intf",
	Usage:     "Manage IP Interfaces",
	Category:  "Requisitions",
	Subcommands: []cli.Command{
		{
			Name:         "list",
			Usage:        "List all IP interfaces from a given node",
			ArgsUsage:    "<foreignSource> <foreignId>",
			Action:       listInterfaces,
			BashComplete: foreignIDBashComplete,
		},
		{
			Name:         "get",
			Usage:        "Gets a specific IP interface from a given node",
			ArgsUsage:    "<foreignSource> <foreignId> <ipAddress>",
			Action:       showInterface,
			BashComplete: ipAddressBashComplete,
		},
		{
			Name:      "set",
			ShortName: "add",
			Usage:     "Adds or update an IP interface from a given node",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress|fqdn>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "description, d",
					Usage: "IP Interface Description",
				},
				cli.GenericFlag{
					Name: "snmpPrimary, p",
					Value: &model.EnumValue{
						Enum:    []string{"P", "N", "S"},
						Default: "N",
					},
					Usage: "Primary Interface Flag: P (primary), S (secondary), N (Not Elegible)",
				},
				cli.IntFlag{
					Name:  "status, s",
					Value: 1,
					Usage: "Interface Status: 1 for managed, 3 for unmanaged (yes, I know)",
				},
				cli.StringSliceFlag{
					Name:  "metaData, m",
					Usage: "A meta-data entry (e.x. --metaData 'foo=bar')",
				},
			},
			Action:       setInterface,
			BashComplete: foreignIDBashComplete,
		},
		{
			Name:         "apply",
			Usage:        "Creates or updates an IP interface on a given node from a external YAML file, overriding any existing content",
			ArgsUsage:    "<foreignSource> <foreignId> <yaml>",
			Action:       applyInterface,
			BashComplete: foreignIDBashComplete,
		},
		{
			Name:         "delete",
			ShortName:    "del",
			Usage:        "Deletes an IP interface from a given node",
			ArgsUsage:    "<foreignSource> <foreignId> <ipAddress>",
			Action:       deleteInterface,
			BashComplete: ipAddressBashComplete,
		},
		{
			Name:      "meta",
			ShortName: "m",
			Usage:     "Manage meta-data",
			Subcommands: []cli.Command{
				{
					Name:         "list",
					Usage:        "Gets all meta-data for a given IP Interface",
					ArgsUsage:    "<foreignSource> <foreignId> <ipAddress>",
					Action:       intfListMetaData,
					BashComplete: ipAddressBashComplete,
				},
				{
					Name:         "set",
					Usage:        "Adds or updates a meta-data entry for a given IP Interface",
					ArgsUsage:    "<foreignSource> <foreignId> <ipAddress> <metaData-key> <metaData-value>",
					Action:       intfSetMetaData,
					BashComplete: ipAddressBashComplete,
				},
				{
					Name:         "delete",
					Usage:        "Deletes a meta-data entry from a given IP Interface",
					ArgsUsage:    "<foreignSource> <foreignId> <ipAddress> <metaData-key>",
					Action:       intfDeleteMetaData,
					BashComplete: ipAddressBashComplete,
				},
			},
		},
	},
}

func listInterfaces(c *cli.Context) error {
	node, err := getReqAPI().GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	if len(node.Interfaces) == 0 {
		fmt.Println("There are no IP interfaces on the chosen node")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "IP Address\tDescription\tSNMP Primary\tServices")
	for _, intf := range node.Interfaces {
		desc := intf.Description
		if desc == "" {
			desc = "N/A"
		}
		fmt.Fprintf(writer, "%s\t%s\t%s\t%d\n", intf.IPAddress, desc, intf.SnmpPrimary, len(intf.Services))
	}
	writer.Flush()
	return nil
}

func showInterface(c *cli.Context) error {
	intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(&intf)
	fmt.Println(string(data))
	return nil
}

func setInterface(c *cli.Context) error {
	intf := model.RequisitionInterface{
		IPAddress:   c.Args().Get(2),
		Description: c.String("description"),
		SnmpPrimary: c.String("snmpPrimary"),
		Status:      c.Int("status"),
	}

	api := getReqAPI()
	current, err := api.GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		mergeInterfaceMetaData(c, &intf)
		return api.SetInterface(c.Args().Get(0), c.Args().Get(1), intf)
	}
	err = current.Merge(intf)
	if err != nil {
		return err
	}
	mergeInterfaceMetaData(c, current)
	return api.SetInterface(c.Args().Get(0), c.Args().Get(1), *current)
}

func applyInterface(c *cli.Context) error {
	data, err := common.ReadInput(c, 2)
	if err != nil {
		return err
	}
	intf := model.RequisitionInterface{}
	err = yaml.Unmarshal(data, &intf)
	if err != nil {
		return err
	}
	return getReqAPI().SetInterface(c.Args().Get(0), c.Args().Get(1), intf)
}

func deleteInterface(c *cli.Context) error {
	return getReqAPI().DeleteInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
}

func intfListMetaData(c *cli.Context) error {
	intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		return err
	}
	if len(intf.MetaData) == 0 {
		fmt.Println("There is no meta-data for the chosen IP interface")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Context\tKey\tValue")
	for _, m := range intf.MetaData {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", m.Context, m.Key, m.Value)
	}
	writer.Flush()
	return nil
}

func intfSetMetaData(c *cli.Context) error {
	intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		return err
	}
	intf.SetMetaData(c.Args().Get(3), c.Args().Get(4))
	if err := intf.Validate(); err != nil {
		return err
	}
	return getReqAPI().SetInterface(c.Args().Get(0), c.Args().Get(1), *intf)
}

func intfDeleteMetaData(c *cli.Context) error {
	intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		return err
	}
	intf.DeleteMetaData(c.Args().Get(3))
	if err := intf.Validate(); err != nil {
		return err
	}
	return getReqAPI().SetInterface(c.Args().Get(0), c.Args().Get(1), *intf)
}

func mergeInterfaceMetaData(c *cli.Context, target *model.RequisitionInterface) {
	metaData := c.StringSlice("metaData")
	for _, p := range metaData {
		data := strings.Split(p, "=")
		target.AddMetaData(data[0], data[1])
	}
}
