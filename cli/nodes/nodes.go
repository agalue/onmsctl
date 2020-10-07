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
			Name:   "apply",
			Usage:  "Creates a set of nodes from a external file",
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
