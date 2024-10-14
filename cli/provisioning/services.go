package provisioning

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
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
			Name:         "list",
			Usage:        "List all monitored services from a given IP interface",
			ArgsUsage:    "<foreignSource> <foreignId> <ipAddress>",
			Action:       listServices,
			BashComplete: ipAddressBashComplete,
		},
		{
			Name:         "set",
			ShortName:    "add",
			Usage:        "Adds or update a monitored service to a given IP interface, overriding any existing content",
			ArgsUsage:    "<foreignSource> <foreignId> <ipAddress> <serviceName>",
			Action:       setService,
			BashComplete: servicesBashComplete,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "metadata, m",
					Usage: "A metadata entry (e.x. --metadata 'foo=bar')",
				},
			},
		},
		{
			Name:         "delete",
			ShortName:    "del",
			Usage:        "Deletes a monitored service from a given IP interface",
			ArgsUsage:    "<foreignSource> <foreignId> <ipAddress> <serviceName>",
			Action:       deleteService,
			BashComplete: servicesBashComplete,
		},
		{
			Name:      "meta",
			ShortName: "m",
			Usage:     "Manage metadata",
			Subcommands: []cli.Command{
				{
					Name:         "list",
					Usage:        "Gets all metadata for a given service",
					ArgsUsage:    "<foreignSource> <foreignId> <ipAddress> <serviceName>",
					Action:       svcListMetaData,
					BashComplete: servicesBashComplete,
				},
				{
					Name:         "set",
					Usage:        "Adds or updates a metadata entry for a given service",
					ArgsUsage:    "<foreignSource> <foreignId> <ipAddress> <serviceName> <metadata-key> <metadata-value>",
					Action:       svcSetMetaData,
					BashComplete: servicesBashComplete,
				},
				{
					Name:         "delete",
					Usage:        "Deletes a metadata entry from a given service",
					ArgsUsage:    "<foreignSource> <foreignId> <ipAddress> <serviceName> <metadata-key>",
					Action:       svcDeleteMetaData,
					BashComplete: servicesBashComplete,
				},
			},
		},
	},
}

func listServices(c *cli.Context) error {
	intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		return err
	}
	if len(intf.Services) == 0 {
		fmt.Println("There are no monitored services on the chosen IP interface")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Service Name")
	for _, svc := range intf.Services {
		fmt.Fprintf(writer, "%s\n", svc.Name)
	}
	writer.Flush()
	return nil
}

func setService(c *cli.Context) error {
	svc := model.RequisitionMonitoredService{Name: c.Args().Get(3)}
	metaData := c.StringSlice("metaData")
	for _, p := range metaData {
		data := strings.Split(p, "=")
		svc.AddMetaData(data[0], data[1])
	}
	return getReqAPI().SetService(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2), svc)
}

func deleteService(c *cli.Context) error {
	return getReqAPI().DeleteService(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2), c.Args().Get(3))
}

func svcListMetaData(c *cli.Context) error {
	service, err := getMonitoredService(c)
	if err == nil {
		return err
	}
	if len(service.MetaData) == 0 {
		fmt.Println("There is no meta-data for the chosen service")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Context\tKey\tValue")
	for _, m := range service.MetaData {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", m.Context, m.Key, m.Value)
	}
	writer.Flush()
	return nil
}

func svcSetMetaData(c *cli.Context) error {
	service, err := getMonitoredService(c)
	if err == nil {
		return err
	}
	service.SetMetaData(c.Args().Get(4), c.Args().Get(5))
	if err := service.Validate(); err != nil {
		return err
	}
	return getReqAPI().SetService(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2), *service)
}

func svcDeleteMetaData(c *cli.Context) error {
	service, err := getMonitoredService(c)
	if err == nil {
		return err
	}
	service.DeleteMetaData(c.Args().Get(4))
	if err := service.Validate(); err != nil {
		return err
	}
	return getReqAPI().SetService(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2), *service)
}

func getMonitoredService(c *cli.Context) (*model.RequisitionMonitoredService, error) {
	intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		return nil, err
	}
	service := intf.GetService(c.Args().Get(3))
	if service == nil {
		return nil, fmt.Errorf("the service doesn't exist on the interface")
	}
	return service, nil
}
