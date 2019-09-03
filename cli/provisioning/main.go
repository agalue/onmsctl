package provisioning

import (
	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/rest"
	"github.com/OpenNMS/onmsctl/services"
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
