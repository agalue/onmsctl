package resources

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// CliCommand the CLI command to manage events
var CliCommand = cli.Command{
	Name:  "resources",
	Usage: "Manage Collected Resources",
	Subcommands: []cli.Command{
		{
			Name:   "list",
			Usage:  "Shows the ID of each resource and its children",
			Action: showResources,
		},
		{
			Name:      "show",
			Usage:     "Shows all details of a given resource",
			Action:    showResource,
			ArgsUsage: "<resourceId>",
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes a given resource",
			Action:    deleteResource,
			ArgsUsage: "<resourceId>",
		},
		{
			Name:      "node",
			Usage:     "Shows all the resources for a given node",
			Action:    showNode,
			ArgsUsage: "<nodeId|FS:FID>",
		},
	},
}

func showResources(c *cli.Context) error {
	jsonInfo, err := rest.Instance.Get("/rest/resources")
	if err != nil {
		return err
	}
	resourceList := model.ResourceList{}
	json.Unmarshal(jsonInfo, &resourceList)
	resourceList.Enumerate("")
	return nil
}

func showResource(c *cli.Context) error {
	resourceID := c.Args().Get(0)
	if resourceID == "" {
		return fmt.Errorf("Resource ID required")
	}
	jsonInfo, err := rest.Instance.Get("/rest/resources/" + resourceID)
	if err != nil {
		return err
	}
	resource := model.Resource{}
	json.Unmarshal(jsonInfo, &resource)
	data, _ := yaml.Marshal(&resource)
	fmt.Println(string(data))
	return nil
}

func showNode(c *cli.Context) error {
	nodeCriteria := c.Args().Get(0)
	if nodeCriteria == "" {
		return fmt.Errorf("Node ID or Foreign-Source:Foreign-ID combination required")
	}
	jsonInfo, err := rest.Instance.Get("/rest/resources/fornode/" + nodeCriteria)
	if err != nil {
		return err
	}
	resource := model.Resource{}
	json.Unmarshal(jsonInfo, &resource)
	data, _ := yaml.Marshal(&resource)
	fmt.Println(string(data))
	return nil
}

func deleteResource(c *cli.Context) error {
	resourceID := c.Args().Get(0)
	if resourceID == "" {
		return fmt.Errorf("Resource ID required")
	}
	err := rest.Instance.Delete("/rest/resources/" + resourceID)
	if err != nil {
		return err
	}
	fmt.Println("Resource has been deleted.")
	return nil
}
