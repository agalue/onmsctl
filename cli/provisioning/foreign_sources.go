package provisioning

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
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
			ArgsUsage: "<name> <interval>",
		},
		{
			Name:   "apply",
			Usage:  "Creates or updates a foreign source definition from a external YAML file",
			Action: applyForeignSource,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name: "format, x",
					Value: &model.EnumValue{
						Enum:    Formats,
						Default: "yaml",
					},
					Usage: "File Format: " + strings.Join(Formats, ", "),
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<content>",
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
						Enum:    Formats,
						Default: "xml",
					},
					Usage: "File Format: " + strings.Join(Formats, ", "),
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
	fsDef, err := getFsAPI().GetForeignSourceDef(c.Args().Get(0))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(fsDef)
	fmt.Println(string(data))
	return nil
}

func setScanInterval(c *cli.Context) error {
	return getFsAPI().SetScanInterval(c.Args().Get(0), c.Args().Get(1))
}

func applyForeignSource(c *cli.Context) error {
	fsDef, err := parseForeignSourceDefinition(c)
	if err != nil {
		return err
	}
	return getFsAPI().SetForeignSourceDef(*fsDef)
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
	return getFsAPI().DeleteForeignSourceDef(c.Args().Get(0))
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
	return fsDef, getFsAPI().IsForeignSourceValid(*fsDef)
}
