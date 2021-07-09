package services

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

type requisitionsAPI struct {
	rest  api.RestAPI
	utils api.ProvisioningUtilsAPI
}

// GetRequisitionsAPI Obtain an implementation of the Requisitions API
func GetRequisitionsAPI(rest api.RestAPI) api.RequisitionsAPI {
	return &requisitionsAPI{rest, GetProvisioningUtilsAPI(rest)}
}

func (api requisitionsAPI) GetRequisitionsStats() (*model.RequisitionsStats, error) {
	jsonStats, err := api.rest.Get("/rest/requisitions/deployed/stats")
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve requisition statistics")
	}
	stats := &model.RequisitionsStats{}
	if err = json.Unmarshal(jsonStats, stats); err != nil {
		return nil, err
	}
	return stats, nil
}

func (api requisitionsAPI) CreateRequisition(foreignSource string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s already exist", foreignSource)
	}
	jsonBytes, err := json.Marshal(model.Requisition{Name: foreignSource})
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/requisitions", jsonBytes)
}

func (api requisitionsAPI) GetRequisition(foreignSource string) (*model.Requisition, error) {
	if foreignSource == "" {
		return nil, fmt.Errorf("requisition name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return nil, fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	jsonString, err := api.rest.Get("/rest/requisitions/" + foreignSource)
	if err != nil {
		return nil, err
	}
	requisition := &model.Requisition{}
	if err = json.Unmarshal(jsonString, requisition); err != nil {
		return nil, err
	}
	return requisition, nil
}

func (api requisitionsAPI) SetRequisition(req model.Requisition) error {
	if req.Name == "default" {
		return fmt.Errorf("requisition cannot be named 'default'")
	}
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/requisitions", jsonBytes)
}

func (api requisitionsAPI) DeleteRequisition(foreignSource string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	// Delete all nodes from requisition
	jsonBytes, err := json.Marshal(model.Requisition{Name: foreignSource})
	if err != nil {
		return err
	}
	if err = api.rest.Post("/rest/requisitions", jsonBytes); err != nil {
		return err
	}
	// Import requisition to remove nodes from the database
	if err = api.rest.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting=false", nil, "application/json"); err != nil {
		return err
	}
	// Delete requisition and foreign-source definitions
	if err = api.rest.Delete("/rest/requisitions/deployed/" + foreignSource); err != nil {
		return err
	}
	if err = api.rest.Delete("/rest/requisitions/" + foreignSource); err != nil {
		return err
	}
	if err = api.rest.Delete("/rest/foreignSources/deployed/" + foreignSource); err != nil {
		return err
	}
	if err = api.rest.Delete("/rest/foreignSources/" + foreignSource); err != nil {
		return err
	}
	return nil
}

func (api requisitionsAPI) ImportRequisition(foreignSource string, rescanExisting string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting="+rescanExisting, nil, "application/json")
}

func (api requisitionsAPI) GetNode(foreignSource string, foreignID string) (*model.RequisitionNode, error) {
	if foreignSource == "" {
		return nil, fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return nil, fmt.Errorf("foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return nil, fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	jsonBytes, err := api.rest.Get("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve node %s from requisition %s", foreignID, foreignSource)
	}
	node := &model.RequisitionNode{}
	if err := json.Unmarshal(jsonBytes, node); err != nil {
		return nil, err
	}
	return node, nil
}

func (api requisitionsAPI) SetNode(foreignSource string, node model.RequisitionNode) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	if err := node.Validate(); err != nil {
		return err
	}
	if node.ParentForeignSource != "" && !api.utils.RequisitionExists(node.ParentForeignSource) {
		return fmt.Errorf("cannot set parent foreign source as requisition %s doesn't exist", node.ParentForeignSource)
	}
	jsonBytes, err := json.Marshal(node)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes", jsonBytes)
}

func (api requisitionsAPI) DeleteNode(foreignSource string, foreignID string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID)
}

func (api requisitionsAPI) GetInterface(foreignSource string, foreignID string, ipAddress string) (*model.RequisitionInterface, error) {
	if foreignSource == "" {
		return nil, fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return nil, fmt.Errorf("foreign ID required")
	}
	if ipAddress == "" {
		return nil, fmt.Errorf("IP Address required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return nil, fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	jsonString, err := api.rest.Get("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress)
	if err != nil {
		return nil, err
	}
	intf := &model.RequisitionInterface{}
	if err := json.Unmarshal(jsonString, intf); err != nil {
		return nil, err
	}
	return intf, nil
}

func (api requisitionsAPI) SetInterface(foreignSource string, foreignID string, intf model.RequisitionInterface) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	if err := intf.Validate(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(intf)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/interfaces", jsonBytes)
}

func (api requisitionsAPI) DeleteInterface(foreignSource string, foreignID string, ipAddress string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress)

}

func (api requisitionsAPI) SetService(foreignSource string, foreignID string, ipAddress string, svc model.RequisitionMonitoredService) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	if err := svc.Validate(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(svc)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/interfaces/"+ipAddress+"/services", jsonBytes)
}

func (api requisitionsAPI) DeleteService(foreignSource string, foreignID string, ipAddress string, serviceName string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	if serviceName == "" {
		return fmt.Errorf("service name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress + "/services/" + serviceName)
}

func (api requisitionsAPI) SetCategory(foreignSource string, foreignID string, category model.RequisitionCategory) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	if err := category.Validate(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(category)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/categories", jsonBytes)
}

func (api requisitionsAPI) DeleteCategory(foreignSource string, foreignID string, categoryName string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if categoryName == "" {
		return fmt.Errorf("category name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/categories/" + categoryName)
}

func (api requisitionsAPI) SetAsset(foreignSource string, foreignID string, asset model.RequisitionAsset) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	if err := asset.Validate(); err != nil {
		return err
	}
	assets, err := api.utils.GetAvailableAssets()
	if err != nil {
		return err
	}
	var found = false
	for _, a := range assets.Element {
		if a == asset.Name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid Asset Field: %s", asset.Name)
	}
	jsonBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/assets", jsonBytes)
}

func (api requisitionsAPI) DeleteAsset(foreignSource string, foreignID string, assetName string) error {
	if foreignSource == "" {
		return fmt.Errorf("requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("foreign ID required")
	}
	if assetName == "" {
		return fmt.Errorf("asset name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/assets/" + assetName)
}
