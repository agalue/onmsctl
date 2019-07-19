package provisioning

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

// ServicesCliCommand the CLI command configuration for managing services on IP interfaces for requisitioned nodes
var ServicesCliCommand = cli.Command{
	Name:      "service",
	ShortName: "svc",
	Usage:     "Manage Monitored Services",
	Category:  "Requisitions",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "List all monitored services from a given IP interface",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress>",
			Action:    listServices,
		},
		{
			Name:      "add",
			Usage:     "Adds a new monitored service to a given IP interface",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress> <serviceName>",
			Action:    addService,
		},
		{
			Name:      "delete",
			Usage:     "Deletes a monitored service from a given IP interface",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress> <serviceName>",
			Action:    deleteService,
		},
	},
}

func listServices(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID and IP address required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("ForeignID is required")
	}
	ipAddress := c.Args().Get(2)
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	jsonBytes, err := rest.Instance.Get("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress)
	if err != nil {
		return fmt.Errorf("Cannot retrieve interfaces")
	}
	intf := model.Interface{}
	json.Unmarshal(jsonBytes, &intf)
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Service Name")
	for _, svc := range intf.Services {
		fmt.Fprintf(writer, "%s\n", svc.Name)
	}
	writer.Flush()
	return nil
}

func addService(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, IP address and service name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	ipAddress := c.Args().Get(2)
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	service := c.Args().Get(3)
	if service == "" {
		return fmt.Errorf("Service name required")
	}
	svc := model.Service{Name: service}
	jsonBytes, _ := json.Marshal(svc)
	return rest.Instance.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/interfaces/"+ipAddress+"/services", jsonBytes)
}

func deleteService(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, IP address and service name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	ipAddress := c.Args().Get(2)
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	service := c.Args().Get(3)
	if service == "" {
		return fmt.Errorf("Service name required")
	}
	return rest.Instance.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress + "/services/" + service)
}
