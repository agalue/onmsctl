package services

import (
	"os"
	"testing"

	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
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

func TestProfiles(t *testing.T) {
	fileName := "/tmp/_onms_config.yaml"
	os.Setenv("ONMSCONFIG", fileName)
	api := GetProfilesAPI(rest.Instance)

	// With an unexisting file, the configuration should be empty
	cfg, err := api.GetProfilesConfig()
	assert.NilError(t, err)
	assert.Equal(t, true, cfg.IsEmpty())
	assert.Equal(t, "http://localhost:8980/opennms", rest.Instance.URL)
	assert.Equal(t, false, fileExists(fileName))

	// Save a profile
	p1 := model.Profile{
		Name:     "Demo",
		URL:      "https://demo.opennms.org/opennms",
		Username: "demo",
		Password: "demo",
	}
	err = api.SetProfile(p1)
	assert.NilError(t, err)
	assert.Equal(t, true, fileExists(fileName))

	// File must exist and the only entry is marked as default
	cfg, err = api.GetProfilesConfig()
	assert.NilError(t, err)
	assert.Equal(t, false, cfg.IsEmpty())
	assert.Equal(t, "https://demo.opennms.org/opennms", rest.Instance.URL)
	assert.Equal(t, true, fileExists(fileName))

	// Save a second profile
	p2 := model.Profile{
		Name:     "Prod",
		URL:      "https://onms.agalue.net/opennms",
		Username: "operator",
		Password: "0p3r@t0r",
	}
	err = api.SetProfile(p2)
	assert.NilError(t, err)

	// File must exist with 2 entries
	cfg, err = api.GetProfilesConfig()
	assert.NilError(t, err)
	assert.Equal(t, false, cfg.IsEmpty())
	assert.Equal(t, "https://demo.opennms.org/opennms", rest.Instance.URL)
	assert.Equal(t, 2, len(cfg.Profiles))

	// Change Default
	err = api.SetDefault(p2.Name)
	assert.NilError(t, err)

	// File must exist with 2 entries and a new default
	cfg, err = api.GetProfilesConfig()
	assert.NilError(t, err)
	assert.Equal(t, false, cfg.IsEmpty())
	assert.Equal(t, "https://onms.agalue.net/opennms", rest.Instance.URL)
	assert.Equal(t, 2, len(cfg.Profiles))

	// Remove temporary file
	os.Remove(fileName)
}
