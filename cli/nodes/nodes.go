package nodes

import (
	"fmt"
	"log"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// CliCommand the CLI command to manage nodes
var CliCommand = cli.Command{
	Name:        "nodes",
	Usage:       "Manage OpenNMS Nodes / Inventory",
	Description: "Manage OpenNMS Nodes / Inventory\n   The recommended way to populate the inventory is via Provisioning Requisitions or Auto-Discover.\n   When neither of those methods are possible, you can manually build the inventory using this sub-command.\n   It is expected a YAML file with 'nodes' as root tag and an array of elements inside of it.",
	Subcommands: []cli.Command{
		{
			Name:   "list",
			Usage:  "Lists all the nodes",
			Action: getNodes,
		},
		{
			Name:   "add",
			Usage:  "Add a new node",
			Action: addNode,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "label, l",
					Usage:    "Node Label",
					Required: true,
				},
				cli.StringFlag{
					Name:  "location, L",
					Usage: "Node Minion Location",
				},
				cli.StringFlag{
					Name:  "sysOID, so",
					Usage: "SNMP System Object ID",
				},
				cli.StringFlag{
					Name:  "sysName, sn",
					Usage: "SNMP System Name",
				},
				cli.StringFlag{
					Name:  "sysDescr, sd",
					Usage: "SNMP System Description",
				},
				cli.StringFlag{
					Name:  "sysLocation, sl",
					Usage: "SNMP System Location",
				},
				cli.StringFlag{
					Name:  "sysContact, sc",
					Usage: "SNMP System Contact",
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "Adds a new set of nodes from a external file",
			Action: addNodes,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<content>",
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing node (cannot be undone)",
			ArgsUsage: "<nodeId|foreignSource:foreignID>",
			Action:    deleteNode,
		},
		{
			Name:  "metadata",
			Usage: "Manage node-level metadata",
			Subcommands: []cli.Command{
				{
					Name:      "list",
					Usage:     "Lists all the node-level metadata",
					Action:    listNodeMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID>",
				},
				{
					Name:      "set",
					Usage:     "Adds or updates a node-level metadata entry",
					Action:    setNodeMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "context, c",
							Usage:    "Metadata Context",
							Required: true,
						},
						cli.StringFlag{
							Name:     "key, k",
							Usage:    "Metadata Key",
							Required: true,
						},
						cli.StringFlag{
							Name:     "value, v",
							Usage:    "Metadata Value",
							Required: true,
						},
					},
				},
				{
					Name:      "delete",
					Usage:     "Deletes an existing node-level metadata entry",
					Action:    deleteNodeMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "context, c",
							Usage:    "Metadata Context",
							Required: true,
						},
						cli.StringFlag{
							Name:     "key, k",
							Usage:    "Metadata Key",
							Required: true,
						},
					},
				},
			},
		},
		IPInterfacesCliCommand,
		SnmpInterfacesCliCommand,
		CategoriesCliCommand,
		AssetsCliCommand,
	},
}

func getNodes(c *cli.Context) error {
	list, err := services.GetNodesAPI(rest.Instance).GetNodes()
	if err != nil {
		return err
	}
	if len(list.Nodes) == 0 {
		fmt.Println("There are no nodes")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Node ID\tNode Label\tForeign Source\tForeign ID\tSNMP sysObjectID")
	for _, n := range list.Nodes {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", n.ID, n.Label, n.ForeignSource, n.ForeignID, n.SysObjectID)
	}
	writer.Flush()
	return nil
}

func deleteNode(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	return services.GetNodesAPI(rest.Instance).DeleteNode(criteria)
}

func addNode(c *cli.Context) error {
	n := &model.OnmsNode{
		Label:          c.String("label"),
		Location:       c.String("location"),
		SysObjectID:    c.String("sysOID"),
		SysName:        c.String("sysName"),
		SysDescription: c.String("sysDescr"),
		SysContact:     c.String("sysContact"),
		SysLocation:    c.String("sysLocation"),
	}
	if err := n.Validate(); err != nil {
		return err
	}
	return services.GetNodesAPI(rest.Instance).AddNode(n)
}

func addNodes(c *cli.Context) error {
	list := &model.OnmsNodeList{}
	data, err := common.ReadInput(c, 0)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, list)
	if err != nil {
		return err
	}
	api := services.GetNodesAPI(rest.Instance)
	for _, n := range list.Nodes {
		err := api.AddNode(&n)
		if err != nil {
			log.Printf("[ERROR] %v\n", err)
		}
	}
	return nil
}

func listNodeMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	meta, err := services.GetNodesAPI(rest.Instance).GetNodeMetadata(criteria)
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

func setNodeMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	meta := model.MetaData{
		Context: c.String("context"),
		Key:     c.String("key"),
		Value:   c.String("value"),
	}
	if err := meta.Validate(); err != nil {
		return err
	}
	return services.GetNodesAPI(rest.Instance).SetNodeMetadata(criteria, meta)
}

func deleteNodeMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ctx := c.String("context")
	if ctx == "" {
		return fmt.Errorf("Context is required")
	}
	key := c.String("key")
	if key == "" {
		return fmt.Errorf("Key is required")
	}
	return services.GetNodesAPI(rest.Instance).DeleteNodeMetadata(criteria, ctx, key)
}
