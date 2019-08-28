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
	Name:     "node",
	Usage:    "Manage Nodes",
	Category: "Requisitions",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "List all nodes from a given requisition",
			ArgsUsage: "<foreignSource>",
			Action:    listNodes,
		},
		{
			Name:      "get",
			Usage:     "Gets a specific node from a given requisition",
			ArgsUsage: "<foreignSource> <foreignId>",
			Action:    showNode,
		},
		{
			Name:      "set",
			ShortName: "add",
			Usage:     "Adds or updates a node from a given requisition",
			ArgsUsage: "<foreignSource> <foreignId>",
			Action:    setNode,
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
					Name:  "metaData, m",
					Usage: "A meta-data entry (e.x. --metaData 'foo=bar')",
				},
			},
		},
		{
			Name:  "apply",
			Usage: "Creates or updates a node on a given requisition from a external YAML file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<foreignSource> <yaml>",
			Action:    applyNode,
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes a node from a given requisition",
			ArgsUsage: "<foreignSource> <foreignId>",
			Action:    deleteNode,
		},
	},
}

func listNodes(c *cli.Context) error {
	requisition, err := api.GetRequisition(c.Args().Get(0))
	if err != nil {
		return err
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
	node, err := api.GetNode(c.Args().Get(0), c.Args().Get(1))
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
	metaData := c.StringSlice("metaData")
	for _, p := range metaData {
		data := strings.Split(p, "=")
		node.AddMetaData(data[0], data[1])
	}
	return api.SetNode(c.Args().Get(0), node)
}

func applyNode(c *cli.Context) error {
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	node := model.RequisitionNode{}
	yaml.Unmarshal(data, &node)
	return api.SetNode(c.Args().Get(0), node)
}

func deleteNode(c *cli.Context) error {
	return api.DeleteNode(c.Args().Get(0), c.Args().Get(1))
}
