package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/test"
	"gotest.tools/assert"
)

var mockForeignSource = &model.ForeignSourceDef{
	Name:         "default",
	ScanInterval: "1w",
	Detectors: []model.Detector{
		{
			Name:  "ICMP",
			Class: "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
		},
		{
			Name:  "SNMP",
			Class: "org.opennms.netmgt.provision.detector.snmp.SnmpDetector",
		},
	},
	Policies: []model.Policy{
		{
			Name:  "Production",
			Class: "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy",
			Parameters: []model.Parameter{
				{
					Key:   "category",
					Value: "Production",
				},
				{
					Key:   "matchBehavior",
					Value: "NO_PARAMETERS",
				},
			},
		},
	},
}

type mockForeignSourcesRest struct {
	t *testing.T
}

func (api mockForeignSourcesRest) Get(path string) ([]byte, error) {
	switch path {
	case "/rest/requisitionNames":
		return []byte(`{"count":2,"foreign-source":["Test1","Test2"]}`), nil
	case "/rest/foreignSourcesConfig/assets":
		return []byte(`{"count":3,"element":["address1","city","zip"]}`), nil
	case "/rest/foreignSourcesConfig/policies":
		return []byte(test.PoliciesJSON), nil
	case "/rest/foreignSourcesConfig/detectors":
		return []byte(test.DetectorsJSON), nil
	case "/rest/foreignSources/default":
		return json.Marshal(mockForeignSource)
	}
	return nil, fmt.Errorf("GET: should not be called with path %s", path)
}

func (api mockForeignSourcesRest) PostRaw(path string, dataBytes []byte, contentType string) (*http.Response, error) {
	return &http.Response{}, fmt.Errorf("should not be called")
}

func (api mockForeignSourcesRest) Post(path string, jsonBytes []byte) error {
	switch path {
	case "/rest/foreignSources":
		fs := &model.ForeignSourceDef{}
		if err := json.Unmarshal(jsonBytes, fs); err != nil {
			return err
		}
		assert.Equal(api.t, "default", fs.Name)
		return nil
	case "/rest/foreignSources/default/detectors":
		detector := &model.Detector{}
		if err := json.Unmarshal(jsonBytes, detector); err != nil {
			return err
		}
		assert.Equal(api.t, "ICMP", detector.Name)
		return nil
	case "/rest/foreignSources/default/policies":
		policy := &model.Policy{}
		if err := json.Unmarshal(jsonBytes, policy); err != nil {
			return err
		}
		assert.Equal(api.t, "Production", policy.Name)
		return nil
	}
	return fmt.Errorf("POST: should not be called with %s", path)
}

func (api mockForeignSourcesRest) Delete(path string) error {
	switch path {
	case "/rest/foreignSources/Test1":
		return nil
	case "/rest/foreignSources/deployed/Test1":
		return nil
	case "/rest/foreignSources/default/detectors/ICMP":
		return nil
	case "/rest/foreignSources/default/policies/Production":
		return nil
	}
	return fmt.Errorf("DELETE: should not be called with %s", path)
}

func (api mockForeignSourcesRest) Put(path string, dataBytes []byte, contentType string) error {
	switch path {
	case "/rest/foreignSources/Example":
		assert.Equal(api.t, "scanInterval=3d", string(dataBytes))
		return nil
	}
	return fmt.Errorf("PUT: should not be called with %s", path)
}

func (api mockForeignSourcesRest) IsValid(r *http.Response) error {
	return nil
}

func TestGetForeignSourceDef(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	fs, err := api.GetForeignSourceDef("default")
	assert.NilError(t, err)
	assert.Equal(t, "default", fs.Name)
}

func TestSetForeignSourceDef(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.SetForeignSourceDef(*mockForeignSource)
	assert.NilError(t, err)
}

func TestSetScanInterval(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.SetScanInterval("Example", "3d")
	assert.NilError(t, err)
}

func TestDeleteForeignSourceDef(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.DeleteForeignSourceDef("Test1")
	assert.NilError(t, err)
}

func TestIsForeignSourceValid(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.IsForeignSourceValid(*mockForeignSource)
	assert.NilError(t, err)

	err = api.IsForeignSourceValid(model.ForeignSourceDef{
		Name:         "Example",
		ScanInterval: "1d",
		Policies: []model.Policy{
			{
				Name:  "Wrong",
				Class: "org.opennms.netmgt.provision.persist.policies.WrongPolicy",
			},
		},
	})
	assert.ErrorContains(t, err, "cannot find policy with class")

	err = api.IsForeignSourceValid(model.ForeignSourceDef{
		Name:         "Example",
		ScanInterval: "1d",
		Detectors: []model.Detector{
			{
				Name:  "Wrong",
				Class: "org.opennms.netmgt.provision.detector.mock.WrongDetector",
			},
		},
	})
	assert.ErrorContains(t, err, "cannot find detector with class")
}

func TestIsPolicyValid(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.IsPolicyValid(mockForeignSource.Policies[0])
	assert.NilError(t, err)

	err = api.IsPolicyValid(model.Policy{
		Name:  "Wrong",
		Class: "org.opennms.netmgt.provision.persist.policies.WrongPolicy",
	})
	assert.ErrorContains(t, err, "cannot find policy with class")
}

func TestIsDetectorValid(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.IsDetectorValid(mockForeignSource.Detectors[0])
	assert.NilError(t, err)

	err = api.IsDetectorValid(model.Detector{
		Name:  "Wrong",
		Class: "org.opennms.netmgt.provision.detector.mock.WrongDetector",
	})
	assert.ErrorContains(t, err, "cannot find detector with class")
}

func TestGetDetectorConfig(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})

	plugin, err := api.GetDetectorConfig("ICMP")
	assert.NilError(t, err)
	assert.Equal(t, "ICMP", plugin.Name)
	assert.Equal(t, 7, len(plugin.Parameters))

	plugin, err = api.GetDetectorConfig("org.opennms.netmgt.provision.detector.icmp.IcmpDetector")
	assert.NilError(t, err)
	assert.Equal(t, "ICMP", plugin.Name)
	assert.Equal(t, 7, len(plugin.Parameters))

	_, err = api.GetDetectorConfig("Wrong")
	assert.ErrorContains(t, err, "cannot find detector")
}

func TestPolicyConfig(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})

	plugin, err := api.GetPolicyConfig("Set Node Category")
	assert.NilError(t, err)
	assert.Equal(t, "Set Node Category", plugin.Name)
	assert.Equal(t, 15, len(plugin.Parameters))

	plugin, err = api.GetPolicyConfig("org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy")
	assert.NilError(t, err)
	assert.Equal(t, "Set Node Category", plugin.Name)
	assert.Equal(t, 15, len(plugin.Parameters))

	_, err = api.GetPolicyConfig("Wrong")
	assert.ErrorContains(t, err, "cannot find policy")
}

func TestGetDetector(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	detector, err := api.GetDetector("default", "ICMP")
	assert.NilError(t, err)
	assert.Equal(t, "ICMP", detector.Name)
	_, err = api.GetDetector("default", "Wrong")
	assert.ErrorContains(t, err, "cannot find detector")
}

func TestGetPolicy(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	policy, err := api.GetPolicy("default", "Production")
	assert.NilError(t, err)
	assert.Equal(t, "Production", policy.Name)
	_, err = api.GetPolicy("default", "Wrong")
	assert.ErrorContains(t, err, "cannot find policy")
}

func TestSetDetector(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.SetDetector("default", mockForeignSource.Detectors[0])
	assert.NilError(t, err)
}

func TestSetPolicy(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.SetPolicy("default", mockForeignSource.Policies[0])
	assert.NilError(t, err)
}

func TestDeleteDetector(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.DeleteDetector("default", "ICMP")
	assert.NilError(t, err)
}

func TestDeletePolicy(t *testing.T) {
	api := GetForeignSourcesAPI(&mockForeignSourcesRest{t})
	err := api.DeletePolicy("default", "Production")
	assert.NilError(t, err)
}
