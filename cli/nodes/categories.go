package nodes

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// CategoriesCliCommand the CLI command to manage nodes
var CategoriesCliCommand = cli.Command{
	Name:  "categories",
	Usage: "Manage Surveillance or Node Categories",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "Lists all the categories for a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID>",
			Action:    getCategories,
		},
		{
			Name:      "add",
			Usage:     "Add a new category to a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <categoryName>",
			Action:    addCategory,
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing category from a given node (cannot be undone)",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <category>",
			Action:    deleteCategory,
		},
	},
}

func getCategories(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("either the nodeID or the foreignSource:foreignID combination is required")
	}
	list, err := services.GetNodesAPI(rest.Instance).GetCategories(criteria)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		fmt.Println("There are no categories")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "ID\tName")
	for _, cat := range list {
		fmt.Fprintf(writer, "%d\t%s\n", cat.ID, cat.Name)
	}
	writer.Flush()
	return nil
}

func deleteCategory(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("either the nodeID or the foreignSource:foreignID combination is required")
	}
	name := c.Args().Get(1)
	if name == "" {
		return fmt.Errorf("category name is required")
	}
	return services.GetNodesAPI(rest.Instance).DeleteCategory(criteria, name)
}

func addCategory(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("either the nodeID or the foreignSource:foreignID combination is required")
	}
	name := c.Args().Get(1)
	if name == "" {
		return fmt.Errorf("category name is required")
	}
	category := &model.OnmsCategory{Name: name}
	return services.GetNodesAPI(rest.Instance).AddCategory(criteria, category)
}
