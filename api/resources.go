package api

import "github.com/OpenNMS/onmsctl/model"

// ResourcesAPI the API to manipulate Resources
type ResourcesAPI interface {
	GetResourceForNode(nodeCriteria string) (*model.Resource, error)
	GetResources() (*model.ResourceList, error)
	GetResource(resourceID string) (*model.Resource, error)
	DeleteResource(resourceID string) error
}
