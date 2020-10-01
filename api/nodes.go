package api

import "github.com/OpenNMS/onmsctl/model"

// NodesAPI the API to manipulate nodes
type NodesAPI interface {
	GetNodes() (*model.OnmsNodeList, error)

	GetNode(nodeCriteria string) (*model.OnmsNode, error)
	AddNode(node *model.OnmsNode) error
	DeleteNode(nodeCriteria string) error
	GetNodeMetadata(nodeCriteria string) ([]model.MetaData, error)
	SetNodeMetadata(nodeCriteria string, meta model.MetaData) error
	DeleteNodeMetadata(nodeCriteria string, context string, key string) error

	GetIPInterfaces(nodeCriteria string) (*model.OnmsIPInterfaceList, error)

	GetIPInterface(nodeCriteria string, ipAddress string) (*model.OnmsIPInterface, error)
	SetIPInterface(nodeCriteria string, intf *model.OnmsIPInterface) error
	DeleteIPInterface(nodeCriteria string, ipAddress string) error
	GetIPInterfaceMetadata(nodeCriteria string, ipAddress string) ([]model.MetaData, error)
	SetIPInterfaceMetadata(nodeCriteria string, ipAddress string, meta model.MetaData) error
	DeleteIPInterfaceMetadata(nodeCriteria string, ipAddress string, context string, key string) error

	GetSnmpInterfaces(nodeCriteria string) (*model.OnmsSnmpInterfaceList, error)

	GetSnmpInterface(nodeCriteria string, ifIndex int) (*model.OnmsSnmpInterface, error)
	SetSnmpInterface(nodeCriteria string, intf *model.OnmsSnmpInterface) error
	DeleteSnmpInterface(nodeCriteria string, ifIndex int) error

	LinkInterfaces(nodeCriteria string, ifIndex int, ipAddress string) error

	GetMonitoredServices(nodeCriteria string, ipAddress string) (*model.OnmsMonitoredServiceList, error)

	GetMonitoredService(nodeCriteria string, ipAddress string, service string) (*model.OnmsMonitoredService, error)
	SetMonitoredService(nodeCriteria string, ipAddress string, svc *model.OnmsMonitoredService) error
	DeleteMonitoredService(nodeCriteria string, ipAddress string, service string) error
	GetMonitoredServiceMetadata(nodeCriteria string, ipAddress string, service string) ([]model.MetaData, error)
	SetMonitoredServiceMetadata(nodeCriteria string, ipAddress string, service string, meta model.MetaData) error
	DeleteMonitoredServiceMetadata(nodeCriteria string, ipAddress string, service string, context string, key string) error

	GetCategories(nodeCriteria string) ([]model.OnmsCategory, error)
	AddCategory(nodeCriteria string, category *model.OnmsCategory) error
	DeleteCategory(nodeCriteria string, category string) error

	GetAssetRecord(nodeCriteria string) (*model.OnmsAssetRecord, error)
	SetAssetRecord(nodeCriteria string, record *model.OnmsAssetRecord) error
}
