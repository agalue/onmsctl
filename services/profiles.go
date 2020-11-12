package services

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/OpenNMS/onmsctl/api"
	"github.com/OpenNMS/onmsctl/model"
	"github.com/OpenNMS/onmsctl/rest"
	"gopkg.in/yaml.v2"
)

type profilesAPI struct {
	rest api.RestAPI
}

// GetProfilesAPI Obtain an implementation of the Profiles API
func GetProfilesAPI(rest api.RestAPI) api.ProfilesAPI {
	return &profilesAPI{rest}
}

func (api profilesAPI) GetProfilesConfig() (*model.ProfilesConfig, error) {
	configFile := getConfigFile()
	cfg := &model.ProfilesConfig{}
	if fileExists(configFile) {
		data, _ := ioutil.ReadFile(configFile)
		if err := yaml.Unmarshal(data, cfg); err == nil {
			if !cfg.IsEmpty() {
				updateRestInstance(cfg.GetDefaultProfile())
			}
		} else {
			fmt.Println("Please use the config subcommand to configure your server profiles.")
		}
	}
	return cfg, nil
}

func (api profilesAPI) SetDefault(profileName string) error {
	cfg, err := api.GetProfilesConfig()
	if err != nil {
		return err
	}
	if idx := findProfileIndex(profileName, cfg); idx > -1 {
		cfg.Default = profileName
		return saveConfig(cfg)
	}
	return fmt.Errorf("Cannot find profile %s", profileName)
}

func (api profilesAPI) SetProfile(profile model.Profile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	cfg, err := api.GetProfilesConfig()
	if err != nil {
		return err
	}
	if idx := findProfileIndex(profile.Name, cfg); idx > -1 {
		cfg.Profiles[idx] = profile
	} else {
		cfg.Profiles = append(cfg.Profiles, profile)
	}
	if len(cfg.Profiles) == 1 && cfg.Default == "" {
		cfg.Default = profile.Name
	}
	return saveConfig(cfg)
}

func (api profilesAPI) DeleteProfile(profileName string) error {
	cfg, err := api.GetProfilesConfig()
	if err != nil {
		return err
	}
	if idx := findProfileIndex(profileName, cfg); idx > -1 {
		cfg.Profiles = append(cfg.Profiles[:idx], cfg.Profiles[idx+1:]...)
		return saveConfig(cfg)
	}
	return fmt.Errorf("Cannot find profile %s", profileName)
}

func findProfileIndex(profileName string, cfg *model.ProfilesConfig) int {
	found := -1
	for idx, p := range cfg.Profiles {
		if p.Name == profileName {
			found = idx
			break
		}
	}
	return found
}

func updateRestInstance(profile *model.Profile) {
	if profile != nil {
		rest.Instance.URL = profile.URL
		rest.Instance.Timeout = profile.Timeout
		rest.Instance.Username = profile.Username
		rest.Instance.Password = profile.Password
		rest.Instance.Insecure = profile.Insecure
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getConfigFile() string {
	homeDir, _ := os.UserHomeDir()
	configFile := homeDir + string(os.PathSeparator) + ".onms" + string(os.PathSeparator) + "config.yaml"
	return getEnv("ONMSCONFIG", configFile)
}

func saveConfig(cfg *model.ProfilesConfig) error {
	sort.SliceStable(cfg.Profiles, func(i, j int) bool {
		return cfg.Profiles[i].Name < cfg.Profiles[j].Name
	})
	if data, err := yaml.Marshal(cfg); err == nil {
		ioutil.WriteFile(getConfigFile(), data, 0644)
	}
	return nil
}
