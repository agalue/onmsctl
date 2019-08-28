package api

import (
	"github.com/OpenNMS/onmsctl/model"
)

// RequisitionsAPI the API to manipulate Requisitions
type RequisitionsAPI interface {
	GetRequisitionNames() (*model.RequisitionsList, error)
	GetRequisitionsStats() (*model.RequisitionsStats, error)
	RequisitionExists(foreignSource string) bool

	CreateRequisition(foreignSource string) error
	GetRequisition(foreignSource string) (*model.Requisition, error)
	SetRequisition(req model.Requisition) error
	DeleteRequisition(foreignSource string) error
	ImportRequisition(foreignSource string, rescanExisting string) error

	GetNode(foreignSource string, foreignID string) (*model.RequisitionNode, error)
	SetNode(foreignSource string, node model.RequisitionNode) error
	DeleteNode(foreignSource string, foreignID string) error

	GetInterface(foreignSource string, foreignID string, ipAddress string) (*model.RequisitionInterface, error)
	SetInterface(foreignSource string, foreignID string, intf model.RequisitionInterface) error
	DeleteInterface(foreignSource string, foreignID string, ipAddress string) error

	SetService(foreignSource string, foreignID string, ipAddress string, svc model.RequisitionMonitoredService) error
	DeleteService(foreignSource string, foreignID string, ipAddress string, serviceName string) error

	SetCategory(foreignSource string, foreignID string, category model.RequisitionCategory) error
	DeleteCategory(foreignSource string, foreignID string, categoryName string) error

	SetAsset(foreignSource string, foreignID string, asset model.RequisitionAsset) error
	DeleteAsset(foreignSource string, foreignID string, assetName string) error
}
