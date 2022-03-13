package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mangathrV2/internal/sources/connections/mangadex"
	"strings"
)

type Config struct {
	Database struct {
		Driver   string
		Uri      string
	}
	Metadata struct {
		Agent             string
		Location          string
	}
	Downloader struct {
		SimultaneousPages int `yaml:"simultaneousPages"`
		PageRetries       int `yaml:"pageRetries"`
		Delay struct {
			Page          int
			Chapter       int
			UpdateChapter int `yaml:"updateChapter"`
		}
		Output struct {
			Path          string
			UpdatePath    string `yaml:"updatePath"`
			Zip           bool
		}
	}
	Sources struct {
		Mangadex mangadex.Config
	}

}

func (c* Config) Load(path string) error {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return err
	}
	if err = c.validate(); err != nil {
		return err
	}

	return nil
}

func (c* Config) validate() error {
	if !validateDatabaseDriver(c.Database.Driver) {
		return errors.New("InvalidDatabaseError: " + c.Database.Driver + " is not a valid database.")
	}
	if !validateMetadataAgent(c.Metadata.Agent) {
		return errors.New("InvalidMetadataAgentError: " + c.Metadata.Agent + " is not a valid agent.")
	}
	if !validateMetadataLocation(c.Metadata.Location) {
		return errors.New("InvalidMetadataLocationError: " + c.Metadata.Location + " is not a valid location.")
	}
	return nil
}


func validateDatabaseDriver(driver string) bool {
	return isInSlice(driver, []string{"sqlite"})
}

func validateMetadataAgent(agent string) bool {
	return isInSlice(agent, []string{"comicinfo", "json"})
}

func validateMetadataLocation(location string) bool {
	return isInSlice(location, []string{"internal", "external", "both"})
}


func isInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == strings.ToLower(s) {
			return true
		}
	}
	return false
}