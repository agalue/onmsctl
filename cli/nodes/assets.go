package nodes

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// AssetsCliCommand the CLI command to manage nodes
var AssetsCliCommand = cli.Command{
	Name:  "assets",
	Usage: "Manage Assets Record",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "Lists all the assets for a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID>",
			Action:    getAssets,
		},
		{
			Name:   "enum",
			Usage:  "Enumerates the available asset fields",
			Action: getFields,
		},
		{
			Name:      "set",
			Usage:     "Adds or updates an asset field for a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <fieldName> <fieldValue>",
			Action:    addAssetField,
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing asset field from a given node (cannot be undone)",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <fieldName>",
			Action:    deleteAssetField,
		},
	},
}

func getFields(c *cli.Context) error {
	assets, err := services.GetProvisioningUtilsAPI(rest.Instance).GetAvailableAssets()
	if err != nil {
		return err
	}
	for _, asset := range assets.Element {
		fmt.Printf("%s\n", asset)
	}
	return nil
}

func getAssets(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("either the nodeID or the foreignSource:foreignID combination is required")
	}
	record, err := services.GetNodesAPI(rest.Instance).GetAssetRecord(criteria)
	if err != nil {
		return err
	}
	text, _ := yaml.Marshal(record)
	fmt.Println(string(text))
	return nil
}

func deleteAssetField(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("either the nodeID or the foreignSource:foreignID combination is required")
	}
	field := c.Args().Get(1)
	if field == "" {
		return fmt.Errorf("field name is required")
	}
	return setAssetField(criteria, field, "")
}

func addAssetField(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("either the nodeID or the foreignSource:foreignID combination is required")
	}
	field := c.Args().Get(1)
	if field == "" {
		return fmt.Errorf("field name is required")
	}
	value := c.Args().Get(2)
	if value == "" {
		return fmt.Errorf("field value is required")
	}
	return setAssetField(criteria, field, value)
}

func setAssetField(criteria string, field string, value string) error {
	if !isValidField(field) {
		return fmt.Errorf("invalid field %s", field)
	}
	return services.GetNodesAPI(rest.Instance).SetAssetField(criteria, field, value)
}

func getAvailableFields() []string {
	assets, err := services.GetProvisioningUtilsAPI(rest.Instance).GetAvailableAssets()
	fields := make([]string, 0)
	if err != nil {
		return fields
	}
	for _, asset := range assets.Element {
		fields = append(fields, asset)
	}
	return fields
}

func isValidField(field string) bool {
	for _, f := range getAvailableFields() {
		if f == field {
			return true
		}
	}
	return false
}
