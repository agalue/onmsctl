package nodes

import (
	"fmt"
	"reflect"

	"github.com/OpenNMS/onmsctl/model"
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
	fields := getAvailableFields()
	for _, f := range fields {
		fmt.Printf("%s\n", f)
	}
	return nil
}

func getAssets(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
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
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
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
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
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
		return fmt.Errorf("Invalid field %s", field)
	}
	record, err := services.GetNodesAPI(rest.Instance).GetAssetRecord(criteria)
	if err != nil {
		return err
	}
	reflect.ValueOf(record).Elem().FieldByName(field).SetString(value)
	return services.GetNodesAPI(rest.Instance).SetAssetRecord(criteria, record)
}

func getAvailableFields() []string {
	r := &model.OnmsAssetRecord{}
	s := reflect.ValueOf(r).Elem()
	typeOfT := s.Type()
	fields := make([]string, 0)
	for i := 0; i < s.NumField(); i++ {
		f := typeOfT.Field(i).Name
		if f != "" && f != "ID" && f != "XMLName" {
			fields = append(fields, typeOfT.Field(i).Name)
		}
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
