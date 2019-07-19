package provisioning

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/common"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"
)

// AssetsCliCommand the CLI command configuration for managing categories for requisitioned nodes
var AssetsCliCommand = cli.Command{
	Name:     "assets",
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
			Usage:     "Adds a new asset or update an existing based on its name, for a given requisition/node",
			ArgsUsage: "<foreignSource> <foreignId> <assetKey> <assetValue>",
			Action:    setAsset,
		},
		{
			Name:      "delete",
			Usage:     "Deletes an existing asset from a given node",
			ArgsUsage: "<foreignSource> <foreignId> <assetKey>",
			Action:    deleteAsset,
		},
	},
}

func listAssets(c *cli.Context) error {
	node, err := GetNode(c)
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
	assets, err := getAssets(c)
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
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, asset name and value required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	assetKey := c.Args().Get(2)
	if assetKey == "" {
		return fmt.Errorf("Asset name required")
	}
	assetValue := c.Args().Get(3)
	if assetValue == "" {
		return fmt.Errorf("Asset value required")
	}
	assets, err := getAssets(c)
	if err != nil {
		return err
	}
	var found = false
	for _, asset := range assets.Element {
		if asset == assetKey {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Invalid Asset Field: %s", assetKey)
	}
	asset := model.RequisitionAsset{Name: assetKey, Value: assetValue}
	jsonBytes, _ := json.Marshal(asset)
	return rest.Instance.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/assets", jsonBytes)
}

func deleteAsset(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Requisition name, foreign ID, asset name required")
	}
	foreignSource := c.Args().Get(0)
	if !RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	foreignID := c.Args().Get(1)
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	asset := c.Args().Get(2)
	if asset == "" {
		return fmt.Errorf("Asset name required")
	}
	return rest.Instance.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/assets/" + asset)
}

func getAssets(c *cli.Context) (model.ElementList, error) {
	assets := model.ElementList{}
	jsonAssets, err := rest.Instance.Get("/rest/foreignSourcesConfig/assets")
	if err != nil {
		return assets, fmt.Errorf("Cannot retrieve asset names list")
	}
	json.Unmarshal(jsonAssets, &assets)
	return assets, nil
}
