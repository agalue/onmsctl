package provisioning

import (
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/test"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
)

func TestEnumDetectors(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, DetectorsCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "detector", "enum"})
	assert.NilError(t, err)
}

func TestDescribeDetector(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, DetectorsCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "detector", "desc"})
	assert.Error(t, err, "Detector name or class required")

	err = app.Run([]string{app.Name, "detector", "desc", "ICMP"})
	assert.NilError(t, err)
}

func TestListDetectors(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, DetectorsCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "detector", "list"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "detector", "list", "Test"})
	assert.NilError(t, err)
}

func TestGetDetector(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, DetectorsCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "detector", "get"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "detector", "get", "Test"})
	assert.Error(t, err, "Detector name or class required")

	err = app.Run([]string{app.Name, "detector", "get", "Test", "ICMP"})
	assert.NilError(t, err)
}

func TestDeleteDetector(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, DetectorsCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "detector", "del"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "detector", "del", "Test"})
	assert.Error(t, err, "Detector name required")

	err = app.Run([]string{app.Name, "detector", "del", "Test", "HTTP"})
	assert.NilError(t, err)
}

func TestApplyDetector(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, DetectorsCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "detector", "apply"})
	assert.Error(t, err, "Content cannot be empty")

	err = app.Run([]string{app.Name, "detector", "apply", "Test"})
	assert.Error(t, err, "Content cannot be empty")

	var testDetector = model.Detector{
		Name:  "HTTP",
		Class: "org.opennms.netmgt.provision.detector.http.HttpDetector",
	}
	nodeYaml, _ := yaml.Marshal(testDetector)
	err = app.Run([]string{app.Name, "detector", "apply", "Test", string(nodeYaml)})
	assert.Error(t, err, "Cannot find detector with class org.opennms.netmgt.provision.detector.http.HttpDetector")

	testDetector = model.Detector{
		Name:  "ICMP",
		Class: "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
	}
	nodeYaml, _ = yaml.Marshal(testDetector)
	err = app.Run([]string{app.Name, "detector", "apply", "Test", string(nodeYaml)})
	assert.NilError(t, err)
}

func TestSetDetector(t *testing.T) {
	var err error
	app, server := test.InitializeMocks(t, DetectorsCliCommand)
	defer server.Close()

	err = app.Run([]string{app.Name, "detector", "set"})
	assert.Error(t, err, "Requisition name required")

	err = app.Run([]string{app.Name, "detector", "set", "Test"})
	assert.Error(t, err, "Detector name cannot be empty")

	err = app.Run([]string{app.Name, "detector", "set", "Test", "ICMP"})
	assert.Error(t, err, "Detector class cannot be empty")

	err = app.Run([]string{app.Name, "detector", "set", "Test", "ICMP", "org.opennms.netmgt.provision.detector.icmp.IcmpDetector"})
	assert.NilError(t, err)
}
