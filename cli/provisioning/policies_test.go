package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/test"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestEnumPolicies(t *testing.T) {
	var err error
	app := test.CreateCli(PoliciesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "policy", "enum"})
	assert.NilError(t, err)
}

func TestDescribePolicy(t *testing.T) {
	var err error
	app := test.CreateCli(PoliciesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "policy", "desc"})
	assert.Error(t, err, "policy name or class required")

	err = app.Run([]string{app.Name, "policy", "desc", "Set Node Category"})
	assert.NilError(t, err)
}

func TestListPolicy(t *testing.T) {
	var err error
	app := test.CreateCli(PoliciesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "policy", "list"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "policy", "list", "Test"})
	assert.NilError(t, err)
}

func TestGetPolicy(t *testing.T) {
	var err error
	app := test.CreateCli(PoliciesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "policy", "get"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "policy", "get", "Test"})
	assert.Error(t, err, "policy name or class required")

	err = app.Run([]string{app.Name, "policy", "get", "Test", "Production"})
	assert.NilError(t, err)
}

func TestDeletePolicy(t *testing.T) {
	var err error
	app := test.CreateCli(PoliciesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "policy", "del"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "policy", "del", "Test"})
	assert.Error(t, err, "policy name required")

	err = app.Run([]string{app.Name, "policy", "del", "Test", "Production"})
	assert.NilError(t, err)
}

func TestApplyPolicy(t *testing.T) {
	var err error
	app := test.CreateCli(PoliciesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "policy", "apply"})
	assert.Error(t, err, "content cannot be empty")

	err = app.Run([]string{app.Name, "policy", "apply", "Test"})
	assert.Error(t, err, "content cannot be empty")

	var testPolicy = model.Policy{
		Name:  "Avoid discover IP interfaces",
		Class: "org.opennms.netmgt.provision.persist.policies.MatchingIpInterfacePolicy",
		Parameters: []model.Parameter{
			{
				Key:   "action",
				Value: "DO_NOT_PERSIST",
			},
			{
				Key:   "matchBehavior",
				Value: "NO_PARAMETERS",
			},
		},
	}
	policyYaml, _ := yaml.Marshal(testPolicy)
	err = app.Run([]string{app.Name, "policy", "apply", "Test", string(policyYaml)})
	assert.Error(t, err, "cannot find policy with class org.opennms.netmgt.provision.persist.policies.MatchingIpInterfacePolicy")

	testPolicy = model.Policy{
		Name:  "Switches",
		Class: "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy",
		Parameters: []model.Parameter{
			{
				Key:   "category",
				Value: "Switches",
			},
		},
	}
	policyYaml, _ = yaml.Marshal(testPolicy)
	err = app.Run([]string{app.Name, "policy", "apply", "Test", string(policyYaml)})
	assert.Error(t, err, "missing required parameter matchBehavior on org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy")

	testPolicy.Parameters = append(testPolicy.Parameters, model.Parameter{
		Key:   "matchBehavior",
		Value: "NoIdea",
	})
	policyYaml, _ = yaml.Marshal(testPolicy)
	err = app.Run([]string{app.Name, "policy", "apply", "Test", string(policyYaml)})
	assert.Error(t, err, "invalid parameter value matchBehavior on org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy. Valid values are: [ALL_PARAMETERS ANY_PARAMETER NO_PARAMETERS]")

	testPolicy.Parameters[1].Value = "NO_PARAMETERS"
	policyYaml, _ = yaml.Marshal(testPolicy)
	err = app.Run([]string{app.Name, "policy", "apply", "Test", string(policyYaml)})
	assert.NilError(t, err)
}

func TestSetPolicy(t *testing.T) {
	var err error
	app := test.CreateCli(PoliciesCliCommand)
	server := createTestServer(t)
	defer server.Close()

	err = app.Run([]string{app.Name, "policy", "set"})
	assert.Error(t, err, "requisition name required")

	err = app.Run([]string{app.Name, "policy", "set", "Test"})
	assert.Error(t, err, "policy name cannot be empty")

	err = app.Run([]string{app.Name, "policy", "set", "Test", "Switches"})
	assert.Error(t, err, "policy class cannot be empty")

	err = app.Run([]string{app.Name, "policy", "set", "Test", "Switches", "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy"})
	assert.Error(t, err, "missing required parameter category on org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy")

	err = app.Run([]string{app.Name, "policy", "set", "-p", "category=Switches", "Test", "Switches", "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy"})
	assert.Error(t, err, "missing required parameter matchBehavior on org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy")

	err = app.Run([]string{app.Name, "policy", "set", "-p", "category=Switches", "-p", "matchBehavior=NO_PARAMETERS", "Test", "Switches", "org.opennms.netmgt.provision.persist.policies.NodeCategorySettingPolicy"})
	assert.NilError(t, err)
}
