package provisioning

import (
	"github.com/urfave/cli"
)

// CliCommand the CLI command configuration for managing inventory
var CliCommand = cli.Command{
	Name:      "provision",
	ShortName: "inv",
	Usage:     "Manage Provisioning / Inventory",
	Subcommands: []cli.Command{
		RequisitionsCliCommand,
		NodesCliCommand,
		InterfacesCliCommand,
		ServicesCliCommand,
		CategoriesCliCommand,
		AssetsCliCommand,
		ForeignSourcesCliCommand,
		DetectorsCliCommand,
		PoliciesCliCommand,
	},
}
