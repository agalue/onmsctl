package resources

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
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
	resourceList, err := getAPI().GetResources()
	if err != nil {
		return err
	}
	resourceList.Enumerate("")
	return nil
}

func showResource(c *cli.Context) error {
	resource, err := getAPI().GetResource(c.Args().Get(0))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(resource)
	fmt.Println(string(data))
	return nil
}

func showNode(c *cli.Context) error {
	resource, err := getAPI().GetResourceForNode(c.Args().Get(0))
	if err != nil {
		return err
	}
	data, _ := yaml.Marshal(resource)
	fmt.Println(string(data))
	return nil
}

func deleteResource(c *cli.Context) error {
	getAPI().DeleteResource(c.Args().Get(0))
	fmt.Println("Resource has been deleted.")
	return nil
}

func getAPI() api.ResourcesAPI {
	return services.GetResourcesAPI(rest.Instance)
}
