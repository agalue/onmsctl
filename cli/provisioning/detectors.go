package provisioning

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// DetectorsCliCommand the CLI command configuration for managing foreign source detectors
var DetectorsCliCommand = cli.Command{
	Name:      "detector",
	ShortName: "d",
	Usage:     "Manage foreign source detectors",
	Category:  "Foreign Source Definitions",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "List all the detectors from a given foreign source definition",
			ArgsUsage: "<foreignSource>",
			Action:    listDetectors,
		},
		{
			Name:      "enumerate",
			ShortName: "enum",
			Usage:     "Enumerate the list of available detector classes",
			Action:    enumerateDetectorClasses,
		},
		{
			Name:      "describe",
			ShortName: "desc",
			Usage:     "Describe a given detector class",
			ArgsUsage: "<detectorName|ClassName>",
			Action:    describeDetectorClass,
		},
		{
			Name:      "get",
			Usage:     "Gets a detector from a given foreign source definition",
			ArgsUsage: "<foreignSource> <detectorName|className>",
			Action:    getDetector,
		},
		{
			Name:      "set",
			Usage:     "Adds or update a detector for a given foreign source definition, overriding any existing content",
			ArgsUsage: "<foreignSource> <detectorName> <className>",
			Action:    setDetector,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "parameter, p",
					Usage: "A detector parameter (e.x. -p 'retry=2')",
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "Creates or updates a detector from a external YAML file, overriding any existing content",
			Action: applyDetector,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<foreignSource> <yaml>",
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes an existing detector from a given foreign source definition",
			ArgsUsage: "<foreignSource> <detectorName>",
			Action:    deleteDetector,
		},
	},
}

func listDetectors(c *cli.Context) error {
	fsDef, err := getFsAPI().GetForeignSourceDef(c.Args().Get(0))
	if err != nil {
		return err
	}
	if len(fsDef.Detectors) == 0 {
		fmt.Println("There are no detectors on the chosen foreign source definition")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Detector Name\tDetector Class")
	for _, detector := range fsDef.Detectors {
		fmt.Fprintf(writer, "%s\t%s\n", detector.Name, detector.Class)
	}
	writer.Flush()
	return nil
}

func enumerateDetectorClasses(c *cli.Context) error {
	detectors, err := getUtilsAPI().GetAvailableDetectors()
	if err != nil {
		return err
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Detector Name\tDetector Class")
	for _, plugin := range detectors.Plugins {
		fmt.Fprintf(writer, "%s\t%s\n", plugin.Name, plugin.Class)
	}
	writer.Flush()
	return nil
}

func describeDetectorClass(c *cli.Context) error {
	plugin, err := getFsAPI().GetDetectorConfig(c.Args().Get(0))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(plugin)
	fmt.Println(string(data))
	return nil
}

func getDetector(c *cli.Context) error {
	detector, err := getFsAPI().GetDetector(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(detector)
	fmt.Println(string(data))
	return nil
}

func setDetector(c *cli.Context) error {
	detector := model.Detector{Name: c.Args().Get(1), Class: c.Args().Get(2)}
	params := c.StringSlice("parameter")
	for _, p := range params {
		data := strings.Split(p, "=")
		param := model.Parameter{Key: data[0], Value: data[1]}
		detector.Parameters = append(detector.Parameters, param)
	}
	return getFsAPI().SetDetector(c.Args().Get(0), detector)
}

func applyDetector(c *cli.Context) error {
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	detector := model.Detector{}
	yaml.Unmarshal(data, &detector)
	return getFsAPI().SetDetector(c.Args().Get(0), detector)
}

func deleteDetector(c *cli.Context) error {
	return getFsAPI().DeleteDetector(c.Args().Get(0), c.Args().Get(1))
}
