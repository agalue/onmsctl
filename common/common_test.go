package common

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestGetConfigFile(t *testing.T) {
	var file string

	os.Setenv("HOME", "/home/agalue")
	file = getConfigFile()
	assert.Equal(t, "/home/agalue/.onms/config.yaml", file)

	os.Setenv("ONMSCONFIG", "/opt/opennms/etc/onmsctl.yaml")
	file = getConfigFile()
	assert.Equal(t, "/opt/opennms/etc/onmsctl.yaml", file)

	assert.Equal(t, true, fileExists("/etc/hosts"))
	assert.Equal(t, false, fileExists("/_unknown"))
}
