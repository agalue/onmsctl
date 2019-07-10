package provisioning

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// InterfacesCliCommand the CLI command configuration for managing IP interfaces on requisitioned nodes
var InterfacesCliCommand = cli.Command{
	Name:      "interface",
	ShortName: "intf",
	Usage:     "Manage IP Interfaces",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "List all IP interfaces from a given node",
			ArgsUsage: "<foreignSource> <foreignId>",
			Action:    listInterfaces,
		},
		{
			Name:      "get",
			Usage:     "Gets a specific interface from a given node",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress>",
			Action:    showInterface,
		},
		{
			Name:      "set",
			ShortName: "add",
			Usage:     "Adds a new IP interface, or update an existing one based on IP address, on a given requisition/node",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress|fqdn>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "description, d",
					Usage: "IP Interface Description",
				},
				cli.GenericFlag{
					Name: "snmpPrimary, p",
					Value: &EnumValue{
						Enum:    []string{"P", "N", "S"},
						Default: "N",
					},
					Usage: "Primary Interface Flag: P (primary), S (secondary), N (Not Elegible)",
				},
				cli.StringFlag{
					Name:  "status, s",
					Value: "1",
					Usage: "Interface Status: 1 for managed, 3 for unmanaged (yes, I know)",
				},
			},
			Action: setInterface,
		},
		{
			Name:      "apply",
			Usage:     "Creates or updates an interface on a given requisition/node from a external YAML file",
			ArgsUsage: "<foreignSource> <foreignId> <yaml>",
			Action:    applyInterface,
		},
		{
			Name:      "delete",
			Usage:     "Deletes an interface from a given requisition/node",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress>",
			Action:    deleteInterface,
		},
	},
}

func listInterfaces(c *cli.Context) error {
	node, err := GetNode(c)
	if err != nil {
		return err
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "IP Address\tDescription\tSNMP Primary\tServices")
	for _, intf := range node.Interfaces {
		desc := intf.Description
		if desc == "" {
			desc = "N/A"
		}
		fmt.Fprintf(writer, "%s\t%s\t%s\t%d\n", intf.IPAddress, desc, intf.SnmpPrimary, len(intf.Services))
	}
	writer.Flush()
	return nil
}

func showInterface(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, IP address required")
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
	jsonString, err := rest.Instance.Get("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress)
	if err != nil {
		return err
	}
	intf := Interface{}
	json.Unmarshal(jsonString, &intf)
	data, _ := yaml.Marshal(&intf)
	fmt.Println(string(data))
	return nil
}

func setInterface(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, IP address required")
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
	intf := Interface{
		IPAddress:   ipAddress,
		Description: c.String("description"),
		SnmpPrimary: c.String("snmpPrimary"),
		Status:      c.Int("status"),
	}
	err := intf.IsValid()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(intf)
	return rest.Instance.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/interfaces", jsonBytes)
}

func applyInterface(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	data, err := common.ReadInput(c, 2)
	if err != nil {
		return err
	}
	intf := Interface{}
	yaml.Unmarshal(data, &intf)
	err = intf.IsValid()
	if err != nil {
		return err
	}
	fmt.Printf("Adding interface %s to requisition %s on node %s...\n", intf.IPAddress, foreignSource, foreignID)
	jsonBytes, _ := json.Marshal(intf)
	return rest.Instance.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/interfaces", jsonBytes)
}

func deleteInterface(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, IP address required")
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
		return fmt.Errorf("IP address required")
	}
	return rest.Instance.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress)
}
