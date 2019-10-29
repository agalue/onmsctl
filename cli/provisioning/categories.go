package provisioning

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/urfave/cli"
)

// CategoriesCliCommand the CLI command configuration for managing categories for requisitioned nodes
var CategoriesCliCommand = cli.Command{
	Name:      "category",
	ShortName: "cat",
	Usage:     "Manage Surveillance Categories",
	Category:  "Requisitions",
	Subcommands: []cli.Command{
		{
			Name:         "list",
			Usage:        "List all monitored services from a given node",
			ArgsUsage:    "<foreignSource> <foreignId>",
			Action:       listCategories,
			BashComplete: foreignIDBashComplete,
		},
		{
			Name:      "add",
			Usage:     "Adds a new category to a given node",
			ArgsUsage: "<foreignSource> <foreignId> <categoryName>",
			Action:    addCategory,
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes a category from a given node",
			ArgsUsage: "<foreignSource> <foreignId> <categoryName>",
			Action:    deleteCategory,
		},
	},
}

func listCategories(c *cli.Context) error {
	node, err := getReqAPI().GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	if len(node.Categories) == 0 {
		fmt.Println("There are no categories on the chosen node")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Category Name")
	for _, cat := range node.Categories {
		fmt.Fprintf(writer, "%s\n", cat.Name)
	}
	writer.Flush()
	return nil
}

func addCategory(c *cli.Context) error {
	cat := model.RequisitionCategory{Name: c.Args().Get(2)}
	return getReqAPI().SetCategory(c.Args().Get(0), c.Args().Get(1), cat)
}

func deleteCategory(c *cli.Context) error {
	return getReqAPI().DeleteCategory(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
}
