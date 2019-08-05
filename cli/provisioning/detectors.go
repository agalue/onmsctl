package provisioning

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// DetectorsCliCommand the CLI command configuration for managing foreign source detectors
var DetectorsCliCommand = cli.Command{
	Name:     "detector",
	Usage:    "Manage foreign source detectors",
	Category: "Foreign Source Definitions",
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
			Usage:     "Adds or update a detector for a given foreign source definition",
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
			Usage:  "Creates or updates a detector from a external YAML file",
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
	fsDef, err := GetForeignSourceDef(c)
	if err != nil {
		return err
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
	detectors, err := getDetectors()
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
	if !c.Args().Present() {
		return fmt.Errorf("Detector name or class required")
	}
	src := c.Args().Get(0)
	detectors, err := getDetectors()
	if err != nil {
		return err
	}
	for _, plugin := range detectors.Plugins {
		if plugin.Class == src || plugin.Name == src {
			data, _ := yaml.Marshal(&plugin)
			fmt.Println(string(data))
			return nil
		}
	}
	return fmt.Errorf("Cannot find detector for %s", src)
}

func getDetector(c *cli.Context) error {
	fsDef, err := GetForeignSourceDef(c)
	if err != nil {
		return err
	}
	src := c.Args().Get(1)
	if src == "" {
		return fmt.Errorf("Detector name or class required")
	}
	for _, detector := range fsDef.Detectors {
		if detector.Class == src || detector.Name == src {
			data, _ := yaml.Marshal(&detector)
			fmt.Println(string(data))
			return nil
		}
	}
	return fmt.Errorf("Cannot find detector for %s", src)
}

func setDetector(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name, detector name and class required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	detectorName := c.Args().Get(1)
	if detectorName == "" {
		return fmt.Errorf("Detector name required")
	}
	detectorClass := c.Args().Get(2)
	if detectorClass == "" {
		return fmt.Errorf("Detector class required")
	}
	detector := model.Detector{Name: detectorName, Class: detectorClass}
	params := c.StringSlice("parameter")
	for _, p := range params {
		data := strings.Split(p, "=")
		param := model.Parameter{Key: data[0], Value: data[1]}
		detector.Parameters = append(detector.Parameters, param)
	}
	detectors, err := getDetectors()
	if err != nil {
		return err
	}
	err = isDetectorValid(detector, detectors)
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(detector)
	return rest.Instance.Post("/rest/foreignSources/"+foreignSource+"/detectors", jsonBytes)
}

func applyDetector(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	data, err := common.ReadInput(c, 1)
	if err != nil {
		return err
	}
	detector := &model.Detector{}
	yaml.Unmarshal(data, detector)
	detectors, err := getDetectors()
	if err != nil {
		return err
	}
	err = isDetectorValid(*detector, detectors)
	if err != nil {
		return err
	}
	fmt.Printf("Updating detector %s...\n", detector.Name)
	jsonBytes, _ := json.Marshal(detector)
	return rest.Instance.Post("/rest/foreignSources/"+foreignSource+"/detectors", jsonBytes)
}

func deleteDetector(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Foreign source name and detector name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	detector := c.Args().Get(1)
	if detector == "" {
		return fmt.Errorf("Detector name required")
	}
	return rest.Instance.Delete("/rest/foreignSources/" + foreignSource + "/detectors/" + detector)
}

func getDetectors() (model.PluginList, error) {
	detectors := model.PluginList{}
	jsonData, err := rest.Instance.Get("/rest/foreignSourcesConfig/detectors")
	if err != nil {
		return detectors, fmt.Errorf("Cannot retrieve detector list")
	}
	json.Unmarshal(jsonData, &detectors)
	return detectors, nil
}

func isDetectorValid(detector model.Detector, config model.PluginList) error {
	if err := detector.IsValid(); err != nil {
		return err
	}
	plugin := config.FindPlugin(detector.Class)
	if plugin == nil {
		return fmt.Errorf("Cannot find detector with class %s", detector.Class)
	}
	if err := plugin.VerifyParameters(detector.Parameters); err != nil {
		return err
	}
	return nil
}
