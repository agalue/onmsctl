package provisioning

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// NodesCliCommand the CLI command configuration for managing requisitioned nodes
var NodesCliCommand = cli.Command{
	Name:      "node",
	ShortName: "n",
	Usage:     "Manage Nodes",
	Category:  "Requisitions",
	Subcommands: []cli.Command{
		{
			Name:         "list",
			Usage:        "List all nodes from a given requisition",
			ArgsUsage:    "<foreignSource>",
			BashComplete: requisitionNameBashComplete,
			Action:       listNodes,
		},
		{
			Name:         "get",
			Usage:        "Gets a specific node from a given requisition",
			ArgsUsage:    "<foreignSource> <foreignId>",
			Action:       showNode,
			BashComplete: foreignIDBashComplete,
		},
		{
			Name:         "set",
			ShortName:    "add",
			Usage:        "Adds or updates a node from a given requisition",
			ArgsUsage:    "<foreignSource> <foreignId>",
			BashComplete: foreignIDBashComplete,
			Action:       setNode,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "label, l",
					Usage: "Node Label",
				},
				cli.StringFlag{
					Name:  "location, L",
					Usage: "Node Location (when using Minion)",
				},
				cli.StringFlag{
					Name:  "city, c",
					Usage: "City",
				},
				cli.StringFlag{
					Name:  "building, b",
					Usage: "Building",
				},
				cli.StringFlag{
					Name:  "parentForeignSource, pfs",
					Usage: "Parent Foreign Source",
				},
				cli.StringFlag{
					Name:  "parentForeignID, pfid",
					Usage: "Parent Foreign ID",
				},
				cli.StringFlag{
					Name:  "parentNodeLabel, pnl",
					Usage: "Parent Node Label",
				},
				cli.StringSliceFlag{
					Name:  "metadata, m",
					Usage: "A metadata entry (e.x. --metadata 'foo=bar')",
				},
			},
		},
		{
			Name:  "apply",
			Usage: "Creates or updates a node on a given requisition from a external YAML file, overriding any existing content",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage:    "<foreignSource> <yaml>",
			Action:       applyNode,
			BashComplete: requisitionNameBashComplete,
		},
		{
			Name:         "delete",
			ShortName:    "del",
			Usage:        "Deletes a node from a given requisition",
			ArgsUsage:    "<foreignSource> <foreignId>",
			Action:       deleteNode,
			BashComplete: foreignIDBashComplete,
		},
		{
			Name:      "meta",
			ShortName: "m",
			Usage:     "Manage metadata",
			Subcommands: []cli.Command{
				{
					Name:         "list",
					Usage:        "Gets all metadata for a given node",
					ArgsUsage:    "<foreignSource> <foreignId>",
					Action:       nodeListMetaData,
					BashComplete: foreignIDBashComplete,
				},
				{
					Name:         "set",
					Usage:        "Adds or updates a metadata entry for a given node",
					ArgsUsage:    "<foreignSource> <foreignId> <metadata-key> <metadata-value>",
					Action:       nodeSetMetaData,
					BashComplete: foreignIDBashComplete,
				},
				{
					Name:         "delete",
					Usage:        "Deletes a metadata entry from a given node",
					ArgsUsage:    "<foreignSource> <foreignId> <metadata-key>",
					Action:       nodeDeleteMetaData,
					BashComplete: foreignIDBashComplete,
				},
			},
		},
	},
}

func listNodes(c *cli.Context) error {
	requisition, err := getReqAPI().GetRequisition(c.Args().Get(0))
	if err != nil {
		return err
	}
	if len(requisition.Nodes) == 0 {
		fmt.Println("There are no nodes on the chosen requisition")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Foreign ID\tLabel\tLocation\tInterfaces\tAssets\tCategories")
	for _, node := range requisition.Nodes {
		location := node.Location
		if location == "" {
			location = "Default"
		}
		fmt.Fprintf(writer, "%s\t%s\t%s\t%d\t%d\t%d\n", node.ForeignID, node.NodeLabel, location, len(node.Interfaces), len(node.Assets), len(node.Categories))
	}
	writer.Flush()
	return nil
}

func showNode(c *cli.Context) error {
	node, err := getReqAPI().GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(&node)
	fmt.Println(string(data))
	return nil
}

func setNode(c *cli.Context) error {
	node := model.RequisitionNode{
		ForeignID:           c.Args().Get(1),
		NodeLabel:           c.String("label"),
		Location:            c.String("location"),
		City:                c.String("city"),
		Building:            c.String("building"),
		ParentForeignSource: c.String("parentForeignSource"),
		ParentForeignID:     c.String("parentForeignID"),
		ParentNodeLabel:     c.String("parentNodeLabel"),
	}

	api := getReqAPI()
	current, err := api.GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		mergeNodeMetaData(c, &node)
		return api.SetNode(c.Args().Get(0), node)
	}
	err = current.Merge(node)
	if err != nil {
		return err
	}
	mergeNodeMetaData(c, current)
	return api.SetNode(c.Args().Get(0), *current)
}

func applyNode(c *cli.Context) error {
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	node := model.RequisitionNode{}
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		return err
	}
	return getReqAPI().SetNode(c.Args().Get(0), node)
}

func deleteNode(c *cli.Context) error {
	return getReqAPI().DeleteNode(c.Args().Get(0), c.Args().Get(1))
}

func nodeListMetaData(c *cli.Context) error {
	node, err := getReqAPI().GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	if len(node.MetaData) == 0 {
		fmt.Println("There is no meta-data for the chosen node")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Context\tKey\tValue")
	for _, m := range node.MetaData {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", m.Context, m.Key, m.Value)
	}
	writer.Flush()
	return nil
}

func nodeSetMetaData(c *cli.Context) error {
	node, err := getReqAPI().GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	node.SetMetaData(c.Args().Get(2), c.Args().Get(3))
	if err := node.Validate(); err != nil {
		return err
	}
	return getReqAPI().SetNode(c.Args().Get(0), *node)
}

func nodeDeleteMetaData(c *cli.Context) error {
	node, err := getReqAPI().GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	node.DeleteMetaData(c.Args().Get(2))
	if err := node.Validate(); err != nil {
		return err
	}
	return getReqAPI().SetNode(c.Args().Get(0), *node)
}

func mergeNodeMetaData(c *cli.Context, target *model.RequisitionNode) {
	metaData := c.StringSlice("metaData")
	for _, p := range metaData {
		data := strings.Split(p, "=")
		target.AddMetaData(data[0], data[1])
	}
}
