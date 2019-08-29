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
		return nil, fmt.Errorf("Cannot retrieve requisition statistics")
	}
	stats := &model.RequisitionsStats{}
	err = json.Unmarshal(jsonStats, stats)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (api requisitionsAPI) CreateRequisition(foreignSource string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s already exist", foreignSource)
	}
	jsonBytes, _ := json.Marshal(model.Requisition{Name: foreignSource})
	return api.rest.Post("/rest/requisitions", jsonBytes)
}

func (api requisitionsAPI) GetRequisition(foreignSource string) (*model.Requisition, error) {
	if foreignSource == "" {
		return nil, fmt.Errorf("Requisition name required")
	}
	jsonString, err := api.rest.Get("/rest/requisitions/" + foreignSource)
	if err != nil {
		return nil, err
	}
	requisition := &model.Requisition{}
	err = json.Unmarshal(jsonString, requisition)
	if err != nil {
		return nil, err
	}
	return requisition, nil
}

func (api requisitionsAPI) SetRequisition(req model.Requisition) error {
	jsonBytes, _ := json.Marshal(req)
	return api.rest.Post("/rest/requisitions", jsonBytes)
}

func (api requisitionsAPI) DeleteRequisition(foreignSource string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	// Delete all nodes from requisition
	jsonBytes, _ := json.Marshal(model.Requisition{Name: foreignSource})
	err := api.rest.Post("/rest/requisitions", jsonBytes)
	if err != nil {
		return err
	}
	// Import requisition to remove nodes from the database
	err = api.rest.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting=false", nil, "application/json")
	if err != nil {
		return err
	}
	// Delete requisition and foreign-source definitions
	err = api.rest.Delete("/rest/requisitions/deployed/" + foreignSource)
	if err != nil {
		return err
	}
	err = api.rest.Delete("/rest/requisitions/" + foreignSource)
	if err != nil {
		return err
	}
	err = api.rest.Delete("/rest/foreignSources/deployed/" + foreignSource)
	if err != nil {
		return err
	}
	err = api.rest.Delete("/rest/foreignSources/" + foreignSource)
	if err != nil {
		return err
	}
	return nil
}

func (api requisitionsAPI) ImportRequisition(foreignSource string, rescanExisting string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Put("/rest/requisitions/"+foreignSource+"/import?rescanExisting="+rescanExisting, nil, "application/json")
}

func (api requisitionsAPI) GetNode(foreignSource string, foreignID string) (*model.RequisitionNode, error) {
	if foreignSource == "" {
		return nil, fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return nil, fmt.Errorf("Foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return nil, fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	jsonBytes, err := api.rest.Get("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID)
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve node %s from requisition %s", foreignID, foreignSource)
	}
	node := &model.RequisitionNode{}
	json.Unmarshal(jsonBytes, node)
	return node, nil
}

func (api requisitionsAPI) SetNode(foreignSource string, node model.RequisitionNode) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	err := node.IsValid()
	if err != nil {
		return err
	}
	if node.ParentForeignSource != "" && !api.utils.RequisitionExists(node.ParentForeignSource) {
		return fmt.Errorf("Cannot set parent foreign source as requisition %s doesn't exist", node.ParentForeignSource)
	}
	jsonBytes, _ := json.Marshal(node)
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes", jsonBytes)
}

func (api requisitionsAPI) DeleteNode(foreignSource string, foreignID string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID)
}

func (api requisitionsAPI) GetInterface(foreignSource string, foreignID string, ipAddress string) (*model.RequisitionInterface, error) {
	if foreignSource == "" {
		return nil, fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return nil, fmt.Errorf("Foreign ID required")
	}
	if ipAddress == "" {
		return nil, fmt.Errorf("IP Address required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return nil, fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	jsonString, err := api.rest.Get("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress)
	if err != nil {
		return nil, err
	}
	intf := &model.RequisitionInterface{}
	json.Unmarshal(jsonString, intf)
	return intf, nil
}

func (api requisitionsAPI) SetInterface(foreignSource string, foreignID string, intf model.RequisitionInterface) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	err := intf.IsValid()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(intf)
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/interfaces", jsonBytes)
}

func (api requisitionsAPI) DeleteInterface(foreignSource string, foreignID string, ipAddress string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress)

}

func (api requisitionsAPI) SetService(foreignSource string, foreignID string, ipAddress string, svc model.RequisitionMonitoredService) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	err := svc.IsValid()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(svc)
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/interfaces/"+ipAddress+"/services", jsonBytes)
}

func (api requisitionsAPI) DeleteService(foreignSource string, foreignID string, ipAddress string, serviceName string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if ipAddress == "" {
		return fmt.Errorf("IP Address required")
	}
	if serviceName == "" {
		return fmt.Errorf("Service name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/interfaces/" + ipAddress + "/services/" + serviceName)
}

func (api requisitionsAPI) SetCategory(foreignSource string, foreignID string, category model.RequisitionCategory) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	err := category.IsValid()
	if err != nil {
		return err
	}
	jsonBytes, _ := json.Marshal(category)
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/categories", jsonBytes)
}

func (api requisitionsAPI) DeleteCategory(foreignSource string, foreignID string, categoryName string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if categoryName == "" {
		return fmt.Errorf("Category name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/categories/" + categoryName)
}

func (api requisitionsAPI) SetAsset(foreignSource string, foreignID string, asset model.RequisitionAsset) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	err := asset.IsValid()
	if err != nil {
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
		return fmt.Errorf("Invalid Asset Field: %s", asset.Name)
	}
	jsonBytes, _ := json.Marshal(asset)
	return api.rest.Post("/rest/requisitions/"+foreignSource+"/nodes/"+foreignID+"/assets", jsonBytes)
}

func (api requisitionsAPI) DeleteAsset(foreignSource string, foreignID string, assetName string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignID == "" {
		return fmt.Errorf("Foreign ID required")
	}
	if assetName == "" {
		return fmt.Errorf("Asset name required")
	}
	if !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Requisition %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/requisitions/" + foreignSource + "/nodes/" + foreignID + "/assets/" + assetName)
}
