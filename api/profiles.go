package api

import "github.com/OpenNMS/onmsctl/model"

// ProfilesAPI the API to manipulate configuration profiles
type ProfilesAPI interface {
	GetProfilesConfig() (*model.ProfilesConfig, error)
	SetProfile(profile model.Profile) error
	SetDefault(profileName string) error
	DeleteProfile(profileName string) error
}
