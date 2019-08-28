package provisioning

import (
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
)

// Formats the available file formats for requisitions and foreign source definitions
var Formats = []string{"xml", "json", "yaml"}

var api = services.GetRequisitionsAPI(rest.Instance)
var fs = services.GetForeignSourcesAPI(rest.Instance, api)
