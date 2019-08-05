package provisioning

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
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
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name required")
	}
	foreignSource := c.Args().First()
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	jsonBytes, err := rest.Instance.Get("/rest/requisitions/" + foreignSource)
	if err != nil {
		return fmt.Errorf("Cannot retrieve nodes from requisition %s", foreignSource)
	}
	requisition := model.Requisition{}
	json.Unmarshal(jsonBytes, &requisition)
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
	node, err := GetNode(c)
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(&node)
	fmt.Println(string(data))
	return nil
}

func setNode(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name and foreign ID required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	node := model.RequisitionNode{
		ForeignID:           foreignID,
		NodeLabel:           c.String("label"),
		Location:            c.String("location"),
		City:                c.String("city"),
		Building:            c.String("building"),
		ParentForeignSource: c.String("parentForeignSource"),
		ParentForeignID:     c.String("parentForeignID"),
		ParentNodeLabel:     c.String("parentNodeLabel"),
	}
	err := node.IsValid()
	if err != nil {
		return err
	}
	if node.ParentForeignSource != "" && !RequisitionExists(node.ParentForeignSource) {
		return fmt.Errorf("Cannot set parent foreign source as requisition %s doesn't exist", node.ParentForeignSource)
	}
	jsonBytes, _ := json.Marshal(node)
	return rest.Instance.Post("/rest/requisitions/"+foreignSource+"/nodes", jsonBytes)
}

func applyNode(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	node := &model.RequisitionNode{}
	yaml.Unmarshal(data, node)
	err = node.IsValid()
	if err != nil {
		return err
	}
	fmt.Printf("Adding node %s to requisition %s...\n", node.ForeignID, foreignSource)
	jsonBytes, _ := json.Marshal(node)
	return rest.Instance.Post("/rest/requisitions/"+foreignSource+"/nodes", jsonBytes)
}

func deleteNode(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name and foreign ID required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	return rest.Instance.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID)
}
