package provisioning

import (
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/urfave/cli"
)

// AssetsCliCommand the CLI command configuration for managing categories for requisitioned nodes
var AssetsCliCommand = cli.Command{
	Name:     "asset",
	Usage:    "Manage node asset fields",
	Category: "Requisitions",
	Subcommands: []cli.Command{
		{
			Name:      "list",
			Usage:     "List all the assets from a given node",
			ArgsUsage: "<foreignSource> <foreignId>",
			Action:    listAssets,
		},
		{
			Name:      "enumerate",
			ShortName: "enum",
			Usage:     "Enumerate the list of available assets",
			Action:    enumerateAssets,
		},
		{
			Name:      "set",
			Usage:     "Adds or update an asset from a given requisition/node",
			ArgsUsage: "<foreignSource> <foreignId> <assetKey> <assetValue>",
			Action:    setAsset,
		},
		{
			Name:      "delete",
			ShortName: "del",
			Usage:     "Deletes an existing asset from a given node",
			ArgsUsage: "<foreignSource> <foreignId> <assetKey>",
			Action:    deleteAsset,
		},
	},
}

func listAssets(c *cli.Context) error {
	node, err := api.GetNode(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Asset Name\tAsset Value")
	for _, asset := range node.Assets {
		fmt.Fprintf(writer, "%s\t%s\n", asset.Name, asset.Value)
	}
	writer.Flush()
	return nil
}

func enumerateAssets(c *cli.Context) error {
	assets, err := fs.GetAvailableAssets()
	if err != nil {
		return err
	}
	writer := common.NewTableWriter()
	fmt.Fprintln(writer, "Asset Name")
	for _, asset := range assets.Element {
		fmt.Fprintf(writer, "%s\n", asset)
	}
	writer.Flush()
	return nil
}

func setAsset(c *cli.Context) error {
	asset := model.RequisitionAsset{Name: c.Args().Get(2), Value: c.Args().Get(3)}
	return api.SetAsset(c.Args().Get(0), c.Args().Get(1), asset)
}

func deleteAsset(c *cli.Context) error {
	return api.DeleteAsset(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
}
