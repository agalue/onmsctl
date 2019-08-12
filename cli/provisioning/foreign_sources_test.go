package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/test"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestGetForeignSource(t *testing.T) {
	var err error
	app := test.CreateCli(ForeignSourcesCliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "fs", "get"})
	assert.Error(t, err, "Foreign source name required")

	err = app.Run([]string{app.Name, "fs", "get", "Test"})
	assert.NilError(t, err)
}

func TestSetScanInterval(t *testing.T) {
	var err error
	app := test.CreateCli(ForeignSourcesCliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "fs", "int"})
	assert.Error(t, err, "Foreign source name and scan interval are required")

	err = app.Run([]string{app.Name, "fs", "int", "Test"})
	assert.Error(t, err, "Scan interval required")

	err = app.Run([]string{app.Name, "fs", "int", "Test", "5YEARS"})
	assert.Error(t, err, "Invalid scan interval 5YEARS")

	err = app.Run([]string{app.Name, "fs", "int", "Test", "1w"})
	assert.NilError(t, err)
}

func TestDeleteForeignSource(t *testing.T) {
	var err error
	app := test.CreateCli(ForeignSourcesCliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "fs", "del"})
	assert.Error(t, err, "Foreign source name required")

	err = app.Run([]string{app.Name, "fs", "del", "Local"})
	assert.NilError(t, err)
}

func TestApplyForeignSource(t *testing.T) {
	var err error
	app := test.CreateCli(ForeignSourcesCliCommand)
	testServer := test.CreateTestServer(t)
	defer testServer.Close()

	err = app.Run([]string{app.Name, "fs", "apply"})
	assert.Error(t, err, "Content cannot be empty")

	fsDef := &model.ForeignSourceDef{
		Name:         "Local",
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
	fsYaml, _ := yaml.Marshal(fsDef)
	err = app.Run([]string{app.Name, "fs", "apply", string(fsYaml)})
	assert.NilError(t, err)
}
