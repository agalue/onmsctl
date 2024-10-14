package model

import "fmt"

// Profile provides information about accessing a given OpenNMS server
type Profile struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Insecure bool   `yaml:"insecure"`
	Timeout  int    `yaml:"timeout"`
}

// Validate verify required fields
func (p *Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("Profile name cannot be empty")
	}
	if p.URL == "" {
		return fmt.Errorf("OpenNMS URL cannot be empty")
	}
	if p.Username == "" {
		return fmt.Errorf("OpenNMS username cannot be empty")
	}
	if p.Password == "" {
		return fmt.Errorf("OpenNMS user password cannot be empty")
	}
	return nil
}

// ProfilesConfig provides information about the configured OpenNMS servers
type ProfilesConfig struct {
	Default  string    `yaml:"defaultProfile"`
	Profiles []Profile `yaml:"profiles,omitempty"`
}

// IsEmpty checks if configuration is empty
func (cfg ProfilesConfig) IsEmpty() bool {
	return cfg.Default == "" && len(cfg.Profiles) == 0
}

// GetDefaultProfile gets the default profile
func (cfg ProfilesConfig) GetDefaultProfile() *Profile {
	for _, p := range cfg.Profiles {
		if p.Name == cfg.Default {
			return &p
		}
	}
	return nil
}
