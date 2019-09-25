package info

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/urfave/cli"

	"gopkg.in/yaml.v2"
)

// CliCommand the CLI command to provide server information
var CliCommand = cli.Command{
	Name:  "info",
	Usage: "Shows version information about the OpenNMS server",
	Action: func(c *cli.Context) error {
		jsonInfo, err := rest.Instance.Get("/rest/info")
		if err != nil {
			return err
		}
		info := model.OnmsInfo{}
		err = json.Unmarshal(jsonInfo, &info)
		if err != nil {
			return err
		}
		data, err := yaml.Marshal(&info)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}
