package provisioning

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// RequisitionsCliCommand the CLI command configuration for managing requisitions
var RequisitionsCliCommand = cli.Command{
	Name:      "requisition",
	ShortName: "req",
	Usage:     "Manage Requisitions",
	Category:  "Requisitions",
	Subcommands: []cli.Command{
		{
			Name:   "list",
			Usage:  "List all requisitions",
			Action: listRequisitions,
		},
		{
			Name:         "get",
			Usage:        "Gets a specific requisition by name",
			Action:       showRequisition,
			BashComplete: requisitionNameBashComplete,
			ArgsUsage:    "<name>",
		},
		{
			Name:      "add",
			Usage:     "Adds a new requisition",
			Action:    addRequisition,
			ArgsUsage: "<name>",
		},
		{
			Name:   "apply",
			Usage:  "Creates or updates a requisition from a external file, overriding any existing content",
			Action: applyRequisition,
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
			Usage:     "Validates a requisition from a external file",
			Action:    validateRequisition,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name: "format, x",
					Value: &model.EnumValue{
						Enum:    Formats,
						Default: "xml",
					},
					Usage: "File Format: " + strings.Join(Formats, ", "),
				},
				cli.BoolFlag{
					Name:  "forceParseFQDN, F",
					Usage: "Force parsing FQDN (for XML and JSON)",
				},
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
				cli.BoolFlag{
					Name:  "yaml, y",
					Usage: "To generate the YAML representation on success",
				},
			},
			ArgsUsage: "<content>",
		},
		{
			Name:         "import",
			ShortName:    "sync",
			Usage:        "Import or synchronize a requisition",
			Action:       importRequisition,
			BashComplete: requisitionNameBashComplete,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name: "rescanExisting, r",
					Value: &model.EnumValue{
						Enum:    []string{"true", "false", "dbonly"},
						Default: "true",
					},
					Usage: `
	true, to update the database and execute the scan phase
	false, to add/delete nodes on the DB skipping the scan phase
	dbonly, to add/detete/update nodes on the DB skipping the scan phase
	`,
				},
			},
			ArgsUsage: "<name>",
		},
		{
			Name:         "delete",
			ShortName:    "del",
			Usage:        "Deletes a requisition",
			Action:       deleteRequisition,
			BashComplete: requisitionNameBashComplete,
			ArgsUsage:    "<name>",
		},
	},
}

func listRequisitions(c *cli.Context) error {
	requisitions, err := getUtilsAPI().GetRequisitionNames()
	if err != nil {
		return err
	}
	if len(requisitions.ForeignSources) == 0 {
		fmt.Println("There are no requisitions")
		return nil
	}
	statistics, err := getReqAPI().GetRequisitionsStats()
	if err != nil {
		return err
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Requisition\tNodes in DB\tLast Import")
	for _, req := range requisitions.ForeignSources {
		stats := statistics.GetRequisitionStats(req)
		fmt.Fprintf(writer, "%s\t%d\t%s\n", req, len(stats.ForeignIDs), getDisplayTime(stats.LastImport))
	}
	writer.Flush()
	return nil
}

func showRequisition(c *cli.Context) error {
	requisition, err := getReqAPI().GetRequisition(c.Args().First())
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(requisition)
	fmt.Println(string(data))
	return nil
}

func addRequisition(c *cli.Context) error {
	return getReqAPI().CreateRequisition(c.Args().First())
}

func applyRequisition(c *cli.Context) error {
	requisition, err := parseRequisition(c)
	if err != nil {
		return err
	}
	return getReqAPI().SetRequisition(*requisition)
}

func validateRequisition(c *cli.Context) error {
	requisition, err := parseRequisition(c)
	if err != nil {
		return err
	}
	if c.Bool("yaml") {
		data, err := yaml.Marshal(requisition)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	fmt.Printf("Requisition %s is valid!\n", requisition.Name)
	return nil
}

func importRequisition(c *cli.Context) error {
	return getReqAPI().ImportRequisition(c.Args().First(), c.String("rescanExisting"))
}

func deleteRequisition(c *cli.Context) error {
	return getReqAPI().DeleteRequisition(c.Args().First())
}

func getDisplayTime(lastImport *model.Time) string {
	if lastImport == nil || lastImport.IsZero() {
		return "Never"
	}
	return lastImport.In(time.Local).String()
}

func parseRequisition(c *cli.Context) (*model.Requisition, error) {
	requisition := &model.Requisition{}
	data, err := common.ReadInput(c, 0)
	if err != nil {
		return requisition, err
	}
	forceParse := c.Bool("forceParseFQDN")
	switch c.String("format") {
	case "xml":
		err = xml.Unmarshal(data, requisition)
		model.AllowFqdnOnRequisitionedInterfaces = forceParse
	case "json":
		err = json.Unmarshal(data, requisition)
		model.AllowFqdnOnRequisitionedInterfaces = forceParse
	case "yaml":
		err = yaml.Unmarshal(data, requisition)
	}
	if err != nil {
		return requisition, err
	}
	return requisition, requisition.Validate()
}
