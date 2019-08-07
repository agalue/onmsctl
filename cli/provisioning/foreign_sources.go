package provisioning

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// ForeignSourcesCliCommand the CLI command configuration for managing foreign source definitions
var ForeignSourcesCliCommand = cli.Command{
	Name:      "foreign-source",
	ShortName: "fs",
	Usage:     "Manage foreign source definitions",
	Category:  "Foreign Source Definitions",
	Subcommands: []cli.Command{
		{
			Name:      "get",
			Usage:     "Gets a specific foreign source definition by name",
			Action:    showForeignSource,
			ArgsUsage: "<name>",
		},
		{
			Name:      "interval",
			ShortName: "int",
			Usage:     "Sets the scan interval",
			Action:    setScanInterval,
			ArgsUsage: "<interval>",
		},
		{
			Name:   "apply",
			Usage:  "Creates or updates a foreign source definition from a external YAML file",
			Action: applyForeignSource,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name: "format, x",
					Value: &model.EnumValue{
						Enum:    []string{"xml", "json", "yaml"},
						Default: "yaml",
					},
					Usage: "File Format: xml, json, yaml",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<yaml>",
		},
		{
			Name:      "validate",
			ShortName: "v",
			Usage:     "Validates a foreign source definition from a external file",
			Action:    validateForeignSource,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name: "format, x",
					Value: &model.EnumValue{
						Enum:    []string{"xml", "json", "yaml"},
						Default: "xml",
					},
					Usage: "File Format: xml, json, yaml",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<content>",
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes a foreign source definition (restoring defaults)",
			Action:    deleteForeignSource,
			ArgsUsage: "<name>",
		},
	},
}

func showForeignSource(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name required")
	}
	foreignSource := c.Args().First()
	jsonString, err := rest.Instance.Get("/rest/foreignSources/" + foreignSource)
	if err != nil {
		return err
	}
	fsDef := model.ForeignSourceDef{}
	json.Unmarshal(jsonString, &fsDef)
	data, _ := yaml.Marshal(&fsDef)
	fmt.Println(string(data))
	return nil
}

func setScanInterval(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name and scan interval are required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	scanInterval := c.Args().Get(1)
	if scanInterval == "" {
		return fmt.Errorf("Scan interval required")
	}
	if !model.IsValidScanInterval(scanInterval) {
		return fmt.Errorf("Invalid scan interval %s", scanInterval)
	}
	jsonBytes := []byte("scan-interval=" + scanInterval)
	return rest.Instance.Put("/rest/foreignSources/"+foreignSource, jsonBytes, "application/x-www-form-urlencoded")
}

func applyForeignSource(c *cli.Context) error {
	fsDef, err := parseForeignSourceDefinition(c)
	if err != nil {
		return err
	}
	fmt.Printf("Updating foreign source definition %s...\n", fsDef.Name)
	jsonBytes, _ := json.Marshal(fsDef)
	return rest.Instance.Post("/rest/foreignSources", jsonBytes)
}

func validateForeignSource(c *cli.Context) error {
	fsDef, err := parseForeignSourceDefinition(c)
	if err != nil {
		return err
	}
	fmt.Printf("Foreign Source %s is valid!\n", fsDef.Name)
	return nil
}

func deleteForeignSource(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name required")
	}
	foreignSource := c.Args().First()
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	err := rest.Instance.Delete("/rest/foreignSources/deployed/" + foreignSource)
	if err != nil {
		return err
	}
	err = rest.Instance.Delete("/rest/foreignSources/" + foreignSource)
	if err != nil {
		return err
	}
	return nil
}

func isForeignSourceValid(fs model.ForeignSourceDef) error {
	if err := fs.IsValid(); err != nil {
		return err
	}
	if len(fs.Policies) > 0 {
		policiesConfig, err := getPolicies()
		if err != nil {
			return err
		}
		for _, policy := range fs.Policies {
			if err := isPolicyValid(policy, policiesConfig); err != nil {
				return err
			}
		}
	}
	if len(fs.Detectors) > 0 {
		detectorsConfig, err := getDetectors()
		if err != nil {
			return err
		}
		for _, detector := range fs.Detectors {
			if err := isDetectorValid(detector, detectorsConfig); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseForeignSourceDefinition(c *cli.Context) (*model.ForeignSourceDef, error) {
	fsDef := &model.ForeignSourceDef{}
	data, err := common.ReadInput(c, 0)
	if err != nil {
		return fsDef, err
	}
	switch c.String("format") {
	case "xml":
		err = xml.Unmarshal(data, fsDef)
	case "yaml":
		err = yaml.Unmarshal(data, fsDef)
	case "json":
		err = json.Unmarshal(data, fsDef)
	}
	if err != nil {
		return fsDef, err
	}
	return fsDef, isForeignSourceValid(*fsDef)
}
