package services

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

type provisioningUtilsAPI struct {
	rest api.RestAPI
}

// GetProvisioningUtilsAPI Obtain an implementation of the Provisioning Utils API
func GetProvisioningUtilsAPI(rest api.RestAPI) api.ProvisioningUtilsAPI {
	return &provisioningUtilsAPI{rest}
}

func (api provisioningUtilsAPI) GetRequisitionNames() (*model.RequisitionsList, error) {
	jsonRequisitions, err := api.rest.Get("/rest/requisitionNames")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve requisition names: %s", err)
	}
	requisitions := &model.RequisitionsList{}
	if err := json.Unmarshal(jsonRequisitions, requisitions); err != nil {
		return nil, err
	}
	return requisitions, nil
}

func (api provisioningUtilsAPI) RequisitionExists(foreignSource string) bool {
	requisitions, err := api.GetRequisitionNames()
	if err != nil {
		return false
	}
	var found = false
	for _, fs := range requisitions.ForeignSources {
		if fs == foreignSource {
			found = true
			break
		}
	}
	return found
}

func (api provisioningUtilsAPI) GetAvailableAssets() (*model.ElementList, error) {
	assets := &model.ElementList{}
	jsonAssets, err := api.rest.Get("/rest/foreignSourcesConfig/assets")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve asset names list")
	}
	if err := json.Unmarshal(jsonAssets, assets); err != nil {
		return nil, err
	}
	return assets, nil
}

func (api provisioningUtilsAPI) GetAvailableDetectors() (*model.PluginList, error) {
	detectors := &model.PluginList{}
	jsonData, err := api.rest.Get("/rest/foreignSourcesConfig/detectors")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve detector list")
	}
	if err := json.Unmarshal(jsonData, detectors); err != nil {
		return nil, err
	}
	return detectors, nil
}

func (api provisioningUtilsAPI) GetAvailablePolicies() (*model.PluginList, error) {
	policies := &model.PluginList{}
	jsonData, err := api.rest.Get("/rest/foreignSourcesConfig/policies")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve policy list")
	}
	if err := json.Unmarshal(jsonData, policies); err != nil {
		return nil, err
	}
	return policies, nil
}
