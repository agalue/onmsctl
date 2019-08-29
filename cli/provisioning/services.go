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
			Name:      "list",
			Usage:     "List all monitored services from a given IP interface",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress>",
			Action:    listServices,
		},
		{
			Name:      "set",
			ShortName: "add",
			Usage:     "Adds or update a monitored service to a given IP interface",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress> <serviceName>",
			Action:    setService,
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "metaData, m",
					Usage: "A meta-data entry (e.x. --metaData 'foo=bar')",
				},
			},
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes a monitored service from a given IP interface",
			ArgsUsage: "<foreignSource> <foreignId> <ipAddress> <serviceName>",
			Action:    deleteService,
		},
	},
}

func listServices(c *cli.Context) error {
	intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
	if err != nil {
		return err
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
