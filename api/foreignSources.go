package api

import (
	"github.com/OpenNMS/onmsctl/model"
)

// ForeignSourcesAPI the API to manipulate Foreign Source definitions
type ForeignSourcesAPI interface {
	GetForeignSourceDef(foreignSource string) (*model.ForeignSourceDef, error)
	SetForeignSourceDef(fs model.ForeignSourceDef) error
	SetScanInterval(foreignSource string, scanInterval string) error
	DeleteForeignSourceDef(foreignSource string) error

	GetAvailableAssets() (*model.ElementList, error)
	GetAvailableDetectors() (*model.PluginList, error)
	GetAvailablePolicies() (*model.PluginList, error)

	IsForeignSourceValid(fsDef model.ForeignSourceDef) error
	IsPolicyValid(policy model.Policy) error
	IsDetectorValid(detector model.Detector) error

	GetDetectorConfig(detectorID string) (*model.Plugin, error)
	GetDetector(foreignSource string, detectorID string) (*model.Detector, error)
	SetDetector(foreignSource string, detector model.Detector) error
	DeleteDetector(foreignSource string, detectorName string) error

	GetPolicyConfig(policyID string) (*model.Plugin, error)
	GetPolicy(foreignSource string, policyID string) (*model.Policy, error)
	SetPolicy(foreignSource string, policy model.Policy) error
	DeletePolicy(foreignSource string, policyName string) error
}
