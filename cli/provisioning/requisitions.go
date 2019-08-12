package provisioning

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
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
			Name:      "get",
			Usage:     "Gets a specific requisition by name",
			Action:    showRequisition,
			ArgsUsage: "<name>",
		},
		{
			Name:      "add",
			Usage:     "Adds a new requisition",
			Action:    addRequisition,
			ArgsUsage: "<name>",
		},
		{
			Name:   "apply",
			Usage:  "Creates or updates a requisition from a external YAML file",
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
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<yaml>",
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
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<content>",
		},
		{
			Name:      "import",
			ShortName: "sync",
			Usage:     "Import or synchronize a requisition",
			Action:    importRequisition,
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
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes a requisition",
			Action:    deleteRequisition,
			ArgsUsage: "<name>",
		},
	},
}

func listRequisitions(c *cli.Context) error {
	requisitions, err := GetRequisitionNames()
	if err != nil {
		return err
	}
	jsonStats, err := rest.Instance.Get("/rest/requisitions/deployed/stats")
	if err != nil {
		return fmt.Errorf("Cannot retrieve requisition statistics")
	}
	stats := model.RequisitionsStats{}
	json.Unmarshal(jsonStats, &stats)
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Requisition\tNodes in DB\tLast Import")
	for _, req := range requisitions.ForeignSources {
		stats := getStats(stats, req)
		fmt.Fprintf(writer, "%s\t%d\t%s\n", req, len(stats.ForeignIDs), getDisplayTime(stats.LastImport))
	}
	writer.Flush()
	return nil
}

func showRequisition(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name required")
	}
	foreignSource := c.Args().First()
	jsonString, err := rest.Instance.Get("/rest/requisitions/" + foreignSource)
	if err != nil {
		return err
	}
	requisition := model.Requisition{}
	json.Unmarshal(jsonString, &requisition)
	data, _ := yaml.Marshal(&requisition)
	fmt.Println(string(data))
	return nil
}

func addRequisition(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name required")
	}
	foreignSource := c.Args().First()
	if RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s already exist", foreignSource)
	}
	jsonBytes, _ := json.Marshal(model.Requisition{Name: foreignSource})
	return rest.Instance.Post("/rest/requisitions", jsonBytes)
}

func applyRequisition(c *cli.Context) error {
	requisition, err := parseRequisition(c)
	if err != nil {
		return err
	}
	fmt.Printf("Updating requisition %s...\n", requisition.Name)
	jsonBytes, _ := json.Marshal(requisition)
	return rest.Instance.Post("/rest/requisitions", jsonBytes)
}

func validateRequisition(c *cli.Context) error {
	requisition, err := parseRequisition(c)
	if err != nil {
		return err
	}
	fmt.Printf("Requisition %s is valid!\n", requisition.Name)
	return nil
}

func importRequisition(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name required")
	}
	foreignSource := c.Args().First()
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	rescanExisting := c.String("rescanExisting")
	fmt.Printf("Importing requisition %s (rescanExisting? %s)...\n", foreignSource, rescanExisting)
	return rest.Instance.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting="+rescanExisting, nil, "application/json")
}

func deleteRequisition(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name required")
	}
	foreignSource := c.Args().First()
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	// Delete all nodes from requisition
	jsonBytes, _ := json.Marshal(model.Requisition{Name: foreignSource})
	err := rest.Instance.Post("/rest/requisitions", jsonBytes)
	if err != nil {
		return err
	}
	// Import requisition to remove nodes from the database
	err = rest.Instance.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting=false", nil, "application/json")
	if err != nil {
		return err
	}
	// Delete requisition and foreign-source definitions
	err = rest.Instance.Delete("/rest/requisitions/deployed/" + foreignSource)
	if err != nil {
		return err
	}
	err = rest.Instance.Delete("/rest/requisitions/" + foreignSource)
	if err != nil {
		return err
	}
	err = rest.Instance.Delete("/rest/foreignSources/deployed/" + foreignSource)
	if err != nil {
		return err
	}
	err = rest.Instance.Delete("/rest/foreignSources/" + foreignSource)
	if err != nil {
		return err
	}
	return nil
}

func getStats(stats model.RequisitionsStats, foreignSource string) model.RequisitionStats {
	for _, req := range stats.ForeignSources {
		if req.Name == foreignSource {
			return req
		}
	}
	return model.RequisitionStats{}
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
	switch c.String("format") {
	case "xml":
		err = xml.Unmarshal(data, requisition)
	case "yaml":
		err = yaml.Unmarshal(data, requisition)
	case "json":
		err = json.Unmarshal(data, requisition)
	}
	if err != nil {
		return requisition, err
	}
	return requisition, requisition.IsValid()
}
