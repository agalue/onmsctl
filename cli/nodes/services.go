package nodes

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// ServicesCliCommand the CLI command to manage nodes
var ServicesCliCommand = cli.Command{
	Name:  "services",
	Usage: "Manage Monitored Services",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "Lists all the monitoring services for a given IP interface on a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress>",
			Action:    getServices,
		},
		{
			Name:      "add",
			Usage:     "Add a new monitored service a given IP interface on a given node",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress> <serviceName>",
			Action:    addService,
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing monitored service from a given IP interface on a given node (cannot be undone)",
			ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress> <serviceName>",
			Action:    deleteService,
		},
		{
			Name:  "metadata",
			Usage: "Manage service-level metadata",
			Subcommands: []cli.Command{
				{
					Name:      "list",
					Usage:     "Lists all the service-level metadata",
					Action:    listServiceMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress> <serviceName>",
				},
				{
					Name:      "set",
					Usage:     "Adds or updates a service-level metadata entry",
					Action:    setServiceMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress> <serviceName>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "context, c",
							Usage:    "Metadata Context",
							Required: true,
						},
						cli.StringFlag{
							Name:     "key, k",
							Usage:    "Metadata Key",
							Required: true,
						},
						cli.StringFlag{
							Name:     "value, v",
							Usage:    "Metadata Value",
							Required: true,
						},
					},
				},
				{
					Name:      "delete",
					Usage:     "Deletes an existing service-level metadata entry",
					Action:    deleteServiceMetadata,
					ArgsUsage: "<nodeId|foreignSource:foreignID> <ipAddress> <serviceName>",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "context, c",
							Usage:    "Metadata Context",
							Required: true,
						},
						cli.StringFlag{
							Name:     "key, k",
							Usage:    "Metadata Key",
							Required: true,
						},
					},
				},
			},
		},
	},
}

func getServices(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if criteria == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	list, err := services.GetNodesAPI(rest.Instance).GetMonitoredServices(criteria, ipaddr)
	if err != nil {
		return err
	}
	if len(list.Services) == 0 {
		fmt.Println("There are no monitored services")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "ID\tService Name\tIs Down\tLast Good\tLast Fail")
	for _, n := range list.Services {
		fmt.Fprintf(writer, "%d\t%s\t%v\t%s\t%s\n", n.ID, n.ServiceType.Name, n.IsDown, n.LastGood, n.LastFail)
	}
	writer.Flush()
	return nil
}

func deleteService(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if criteria == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	svc := c.Args().Get(2)
	if criteria == "" {
		return fmt.Errorf("Monitored Service Name is required")
	}
	return services.GetNodesAPI(rest.Instance).DeleteMonitoredService(criteria, ipaddr, svc)
}

func addService(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if criteria == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	svc := c.Args().Get(2)
	if criteria == "" {
		return fmt.Errorf("Monitored Service Name is required")
	}
	service := &model.OnmsMonitoredService{
		ServiceType: &model.OnmsServiceType{
			Name: svc,
		},
	}
	if err := service.Validate(); err != nil {
		return err
	}
	return services.GetNodesAPI(rest.Instance).SetMonitoredService(criteria, ipaddr, service)
}

func listServiceMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if ipaddr == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	svc := c.Args().Get(2)
	if svc == "" {
		return fmt.Errorf("Monitored Service Name is required")
	}
	meta, err := services.GetNodesAPI(rest.Instance).GetMonitoredServiceMetadata(criteria, ipaddr, svc)
	if err != nil {
		return err
	}
	if len(meta) == 0 {
		fmt.Println("There is no metadata")
		return nil
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Context\tKey\tValue")
	for _, m := range meta {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", m.Context, m.Key, m.Value)
	}
	writer.Flush()
	return nil
}

func setServiceMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if ipaddr == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	svc := c.Args().Get(2)
	if svc == "" {
		return fmt.Errorf("Monitored Service Name is required")
	}
	meta := model.MetaData{
		Context: c.String("context"),
		Key:     c.String("key"),
		Value:   c.String("value"),
	}
	if err := meta.Validate(); err != nil {
		return err
	}
	return services.GetNodesAPI(rest.Instance).SetMonitoredServiceMetadata(criteria, ipaddr, svc, meta)
}

func deleteServiceMetadata(c *cli.Context) error {
	criteria := c.Args().Get(0)
	if criteria == "" {
		return fmt.Errorf("Either the nodeID or the foreignSource:foreignID combination is required")
	}
	ipaddr := c.Args().Get(1)
	if ipaddr == "" {
		return fmt.Errorf("Interface IP Address is required")
	}
	svc := c.Args().Get(2)
	if svc == "" {
		return fmt.Errorf("Monitored Service Name is required")
	}
	ctx := c.String("context")
	if ctx == "" {
		return fmt.Errorf("Context is required")
	}
	key := c.String("key")
	if key == "" {
		return fmt.Errorf("Key is required")
	}
	return services.GetNodesAPI(rest.Instance).DeleteMonitoredServiceMetadata(criteria, ipaddr, svc, ctx, key)
}
