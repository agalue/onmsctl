package api

import "github.com/OpenNMS/onmsctl/model"

// ProvisioningUtilsAPI The API for common Provisioning operations
type ProvisioningUtilsAPI interface {
	GetRequisitionNames() (*model.RequisitionsList, error)
	RequisitionExists(foreignSource string) bool

	GetAvailableAssets() (*model.ElementList, error)
	GetAvailableDetectors() (*model.PluginList, error)
	GetAvailablePolicies() (*model.PluginList, error)
}
