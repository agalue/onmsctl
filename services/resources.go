package services

import (
	"encoding/json"
	"fmt"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
)

type resourcesAPI struct {
	rest api.RestAPI
}

// GetResourcesAPI Obtain an implementation of the Resources API
func GetResourcesAPI(rest api.RestAPI) api.ResourcesAPI {
	return &resourcesAPI{rest}
}

func (api resourcesAPI) GetResourceForNode(nodeCriteria string) (*model.Resource, error) {
	if nodeCriteria == "" {
		return nil, fmt.Errorf("node ID or foreignSource:foreignID combination required")
	}
	jsonInfo, err := api.rest.Get("/rest/resources/fornode/" + nodeCriteria)
	if err != nil {
		return nil, err
	}
	resource := &model.Resource{}
	if err := json.Unmarshal(jsonInfo, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func (api resourcesAPI) GetResources() (*model.ResourceList, error) {
	jsonInfo, err := api.rest.Get("/rest/resources")
	if err != nil {
		return nil, err
	}
	resourceList := &model.ResourceList{}
	if err := json.Unmarshal(jsonInfo, resourceList); err != nil {
		return nil, err
	}
	return resourceList, nil
}

func (api resourcesAPI) GetResource(resourceID string) (*model.Resource, error) {
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID required")
	}
	jsonInfo, err := api.rest.Get("/rest/resources/" + resourceID)
	if err != nil {
		return nil, err
	}
	resource := &model.Resource{}
	if err := json.Unmarshal(jsonInfo, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func (api resourcesAPI) DeleteResource(resourceID string) error {
	if resourceID == "" {
		return fmt.Errorf("resource ID required")
	}
	err := api.rest.Delete("/rest/resources/" + resourceID)
	if err != nil {
		return err
	}
	return nil
}
