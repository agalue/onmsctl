package services

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

type foreignSourcesAPI struct {
	rest  api.RestAPI
	utils api.ProvisioningUtilsAPI
}

// GetForeignSourcesAPI Obtain an implementation of the Foreign Source Definitions API
func GetForeignSourcesAPI(rest api.RestAPI) api.ForeignSourcesAPI {
	return &foreignSourcesAPI{rest, GetProvisioningUtilsAPI(rest)}
}

func (api foreignSourcesAPI) GetForeignSourceDef(foreignSource string) (*model.ForeignSourceDef, error) {
	if foreignSource == "" {
		return nil, fmt.Errorf("Requisition name required")
	}
	if foreignSource != "default" && !api.utils.RequisitionExists(foreignSource) {
		return nil, fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	fsDef := &model.ForeignSourceDef{}
	jsonBytes, err := api.rest.Get("/rest/foreignSources/" + foreignSource)
	if err != nil {
		return nil, fmt.Errorf("Cannot retrieve foreign source definition %s", foreignSource)
	}
	if err := json.Unmarshal(jsonBytes, fsDef); err != nil {
		return nil, err
	}
	return fsDef, nil
}

func (api foreignSourcesAPI) SetForeignSourceDef(fs model.ForeignSourceDef) error {
	if err := fs.Validate(); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(fs)
	if err != nil {
		return err
	}
	if fs.Name != "default" && !api.utils.RequisitionExists(fs.Name) {
		return fmt.Errorf("A requisition called '%s' must exist before creating the FS definition", fs.Name)
	}
	return api.rest.Post("/rest/foreignSources", jsonBytes)
}

func (api foreignSourcesAPI) SetScanInterval(foreignSource string, scanInterval string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if scanInterval == "" {
		return fmt.Errorf("Scan interval required")
	}
	if !model.IsValidScanInterval(scanInterval) {
		return fmt.Errorf("Invalid scan interval %s", scanInterval)
	}
	bytes := []byte("scanInterval=" + scanInterval)
	return api.rest.Put("/rest/foreignSources/"+foreignSource, bytes, "application/x-www-form-urlencoded")
}

func (api foreignSourcesAPI) DeleteForeignSourceDef(foreignSource string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignSource != "default" && !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	err := api.rest.Delete("/rest/foreignSources/deployed/" + foreignSource)
	if err != nil {
		return err
	}
	err = api.rest.Delete("/rest/foreignSources/" + foreignSource)
	if err != nil {
		return err
	}
	return nil
}

func (api foreignSourcesAPI) IsPolicyValid(policy model.Policy) error {
	config, err := api.utils.GetAvailablePolicies()
	if err != nil {
		return nil
	}
	return api.isPolicyValid(*config, policy)
}

func (api foreignSourcesAPI) isPolicyValid(config model.PluginList, policy model.Policy) error {
	if err := policy.Validate(); err != nil {
		return err
	}
	plugin := config.FindPlugin(policy.Class)
	if plugin == nil {
		return fmt.Errorf("Cannot find policy with class %s", policy.Class)
	}
	if err := plugin.VerifyParameters(policy.Parameters); err != nil {
		return err
	}
	return nil
}

func (api foreignSourcesAPI) IsDetectorValid(detector model.Detector) error {
	config, err := api.utils.GetAvailableDetectors()
	if err != nil {
		return nil
	}
	return api.isDetectorValid(*config, detector)
}

func (api foreignSourcesAPI) isDetectorValid(config model.PluginList, detector model.Detector) error {
	if err := detector.Validate(); err != nil {
		return err
	}
	plugin := config.FindPlugin(detector.Class)
	if plugin == nil {
		return fmt.Errorf("Cannot find detector with class %s", detector.Class)
	}
	if err := plugin.VerifyParameters(detector.Parameters); err != nil {
		return err
	}
	return nil
}

func (api foreignSourcesAPI) IsForeignSourceValid(fsDef model.ForeignSourceDef) error {
	if err := fsDef.Validate(); err != nil {
		return err
	}
	if len(fsDef.Policies) > 0 {
		policiesConfig, err := api.utils.GetAvailablePolicies()
		if err != nil {
			return err
		}
		for _, policy := range fsDef.Policies {
			if err := api.isPolicyValid(*policiesConfig, policy); err != nil {
				return fmt.Errorf("Problem with foreign source %s: %s", fsDef.Name, err.Error())
			}
		}
	}
	if len(fsDef.Detectors) > 0 {
		detectorsConfig, err := api.utils.GetAvailableDetectors()
		if err != nil {
			return err
		}
		for _, detector := range fsDef.Detectors {
			if err := api.isDetectorValid(*detectorsConfig, detector); err != nil {
				return fmt.Errorf("Problem with foreign source %s: %s", fsDef.Name, err.Error())
			}
		}
	}
	return nil
}

func (api foreignSourcesAPI) GetDetectorConfig(detectorID string) (*model.Plugin, error) {
	if detectorID == "" {
		return nil, fmt.Errorf("Detector name or class required")
	}
	detectors, err := api.utils.GetAvailableDetectors()
	if err != nil {
		return nil, err
	}
	for _, plugin := range detectors.Plugins {
		if plugin.Class == detectorID || plugin.Name == detectorID {
			return &plugin, nil
		}
	}
	return nil, fmt.Errorf("Cannot find detector for %s", detectorID)
}

func (api foreignSourcesAPI) GetDetector(foreignSource string, detectorID string) (*model.Detector, error) {
	fsDef, err := api.GetForeignSourceDef(foreignSource)
	if err != nil {
		return nil, err
	}
	detector, err := fsDef.GetDetector(detectorID)
	if err != nil {
		return nil, err
	}
	return detector, nil
}

func (api foreignSourcesAPI) SetDetector(foreignSource string, detector model.Detector) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignSource != "default" && !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	if err := api.IsDetectorValid(detector); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(detector)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/foreignSources/"+foreignSource+"/detectors", jsonBytes)
}

func (api foreignSourcesAPI) DeleteDetector(foreignSource string, detectorName string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if detectorName == "" {
		return fmt.Errorf("Detector name required")
	}
	if foreignSource != "default" && !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/foreignSources/" + foreignSource + "/detectors/" + detectorName)
}

func (api foreignSourcesAPI) GetPolicyConfig(policyID string) (*model.Plugin, error) {
	if policyID == "" {
		return nil, fmt.Errorf("Policy name or class required")
	}
	policies, err := api.utils.GetAvailablePolicies()
	if err != nil {
		return nil, err
	}
	for _, plugin := range policies.Plugins {
		if plugin.Class == policyID || plugin.Name == policyID {
			return &plugin, nil
		}
	}
	return nil, fmt.Errorf("Cannot find policy for %s", policyID)
}

func (api foreignSourcesAPI) GetPolicy(foreignSource string, policyID string) (*model.Policy, error) {
	fsDef, err := api.GetForeignSourceDef(foreignSource)
	if err != nil {
		return nil, err
	}
	policy, err := fsDef.GetPolicy(policyID)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (api foreignSourcesAPI) SetPolicy(foreignSource string, policy model.Policy) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if foreignSource != "default" && !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	if err := api.IsPolicyValid(policy); err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	return api.rest.Post("/rest/foreignSources/"+foreignSource+"/policies", jsonBytes)
}

func (api foreignSourcesAPI) DeletePolicy(foreignSource string, policyName string) error {
	if foreignSource == "" {
		return fmt.Errorf("Requisition name required")
	}
	if policyName == "" {
		return fmt.Errorf("Policy name required")
	}
	if foreignSource != "default" && !api.utils.RequisitionExists(foreignSource) {
		return fmt.Errorf("Foreign source %s doesn't exist", foreignSource)
	}
	return api.rest.Delete("/rest/foreignSources/" + foreignSource + "/policies/" + policyName)
}
