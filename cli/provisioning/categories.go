package provisioning

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
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
			Name:      "list",
			Usage:     "List all monitored services from a given node",
			ArgsUsage: "<foreignSource> <foreignId>",
			Action:    listCategories,
		},
		{
			Name:      "add",
			Usage:     "Adds a new category to a given node",
			ArgsUsage: "<foreignSource> <foreignId> <categoryName>",
			Action:    addCategory,
		},
		{
			Name:      "delete",
			Usage:     "Deletes a category from a given node",
			ArgsUsage: "<foreignSource> <foreignId> <categoryName>",
			Action:    deleteCategory,
		},
	},
}

func listCategories(c *cli.Context) error {
	node, err := GetNode(c)
	if err != nil {
		return err
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
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, category name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	category := c.Args().Get(2)
	if category == "" {
		return fmt.Errorf("Category name required")
	}
	cat := model.RequisitionCategory{Name: category}
	jsonBytes, _ := json.Marshal(cat)
	return rest.Instance.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/categories", jsonBytes)
}

func deleteCategory(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID and category name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	category := c.Args().Get(2)
	if category == "" {
		return fmt.Errorf("Category name required")
	}
	return rest.Instance.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/categories/" + category)
}
