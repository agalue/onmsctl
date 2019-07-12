package provisioning

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// RequisitionsCliCommand the CLI command configuration for managing requisitions
var RequisitionsCliCommand = cli.Command{
	Name:      "requisition",
	ShortName: "req",
	Usage:     "Manage Requisitions",
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
				cli.StringFlag{
					Name:  "file, f",
					Usage: "External YAML file (use '-' for STDIN Pipe)",
				},
			},
			ArgsUsage: "<yaml>",
		},
		{
			Name:      "import",
			ShortName: "sync",
			Usage:     "Imports or synchronize a requisition",
			Action:    importRequisition,
			Flags: []cli.Flag{
				cli.GenericFlag{
					Name: "rescanExisting, r",
					Value: &common.EnumValue{
						Enum:    []string{"true", "false", "dbonly"},
						Default: "true",
					},
					Usage: "Rescan Existing: true, false, dbonly",
				},
			},
			ArgsUsage: "<name>",
		},
		{
			Name:      "delete",
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
	stats := RequisitionsStats{}
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
	requisition := Requisition{}
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
	jsonBytes, _ := json.Marshal(Requisition{Name: foreignSource})
	return rest.Instance.Post("/rest/requisitions", jsonBytes)
}

func applyRequisition(c *cli.Context) error {
	data, err := common.ReadInput(c, 0)
	if err != nil {
		return err
	}
	requisition := &Requisition{}
	yaml.Unmarshal(data, requisition)
	err = requisition.IsValid()
	if err != nil {
		return err
	}
	fmt.Printf("Adding requisition %s...\n", requisition.Name)
	jsonBytes, _ := json.Marshal(requisition)
	return rest.Instance.Post("/rest/requisitions", jsonBytes)
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
	return rest.Instance.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting="+rescanExisting, nil)
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
	jsonBytes, _ := json.Marshal(Requisition{Name: foreignSource})
	err := rest.Instance.Post("/rest/requisitions", jsonBytes)
	if err != nil {
		return err
	}
	// Import requisition to remove nodes from the database
	err = rest.Instance.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting=false", nil)
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

func getStats(stats RequisitionsStats, foreignSource string) RequisitionStats {
	for _, req := range stats.ForeignSources {
		if req.Name == foreignSource {
			return req
		}
	}
	return RequisitionStats{}
}

func getDisplayTime(lastImport common.Time) string {
	if lastImport.IsZero() {
		return "Never"
	}
	return lastImport.In(time.Local).String()
}
