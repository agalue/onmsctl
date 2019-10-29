package provisioning

import (
	"fmt"
	"strings"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
	"github.com/urfave/cli"
)

// Formats the available file formats for requisitions and foreign source definitions
var Formats = []string{"xml", "json", "yaml"}

func getReqAPI() api.RequisitionsAPI {
	return services.GetRequisitionsAPI(rest.Instance)
}

func getFsAPI() api.ForeignSourcesAPI {
	return services.GetForeignSourcesAPI(rest.Instance)
}

func getUtilsAPI() api.ProvisioningUtilsAPI {
	return services.GetProvisioningUtilsAPI(rest.Instance)
}

func requisitionNameBashComplete(c *cli.Context) {
	if c.NArg() > 0 {
		return
	}
	list, err := getUtilsAPI().GetRequisitionNames()
	if err != nil {
		return
	}
	for _, fs := range list.ForeignSources {
		fmt.Println(fs)
	}
}

func foreignIDBashComplete(c *cli.Context) {
	requisitionNameBashComplete(c)
	if c.NArg() == 1 {
		req, err := getReqAPI().GetRequisition(c.Args().First())
		if err != nil {
			return
		}
		for _, n := range req.Nodes {
			fmt.Println(zshNormalize(n.ForeignID))
		}
	}
}

func ipAddressBashComplete(c *cli.Context) {
	foreignIDBashComplete(c)
	if c.NArg() == 2 {
		node, err := getReqAPI().GetNode(c.Args().Get(0), c.Args().Get(1))
		if err != nil {
			return
		}
		for _, intf := range node.Interfaces {
			fmt.Println(zshNormalize(intf.IPAddress))
		}
	}
}

func servicesBashComplete(c *cli.Context) {
	ipAddressBashComplete(c)
	if c.NArg() == 2 {
		intf, err := getReqAPI().GetInterface(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
		if err != nil {
			return
		}
		for _, svc := range intf.Services {
			fmt.Println(zshNormalize(svc.Name))
		}
	}
}

// TODO Bash or ZSH don't like spaces when doing auto-complete
//      The characters '.' and ':' are escaped as they have a special meaning
func zshNormalize(src string) string {
	dst := strings.ReplaceAll(src, ":", "\\:")
	return strings.ReplaceAll(dst, ".", "\\.")
}
