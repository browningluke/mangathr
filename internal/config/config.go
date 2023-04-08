package config

import (
	"errors"
	"github.com/browningluke/mangathr/internal/config/defaults"
	"github.com/browningluke/mangathr/internal/downloader"
	"github.com/browningluke/mangathr/internal/sources/mangadex"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Config struct {
	Database struct {
		Driver string
		Uri    string
	}
	Downloader downloader.Config
	Sources    struct {
		Mangadex mangadex.Config
		Cubari   cubari.Config
	}
	LogLevel string `yaml:"logLevel"`
}

func (c *Config) Propagate() {
	downloader.SetConfig(c.Downloader)

	// Sources
	mangadex.SetConfig(c.Sources.Mangadex)
	cubari.SetConfig(c.Sources.Cubari)
}

func (c *Config) Load(path string, inContainer bool) error {
	c.useDefaults(inContainer)

	yamlFile, err := os.ReadFile(path)
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

func (c *Config) useDefaults(inContainer bool) {
	c.Database.Driver = defaults.DatabaseDriver()
	c.Database.Uri = defaults.DatabaseUri()

	downloadConf := downloader.Config{}
	downloadConf.Default(inContainer)
	c.Downloader = downloadConf

	mangadexConf := mangadex.Config{}
	mangadexConf.Default()
	c.Sources.Mangadex = mangadexConf

	c.LogLevel = ""

	// Overwrite defaults if we are in a container
	if inContainer {
		c.Database.Uri = defaults.DatabaseUriDocker()
	}
}

func (c *Config) validate() error {
	if !validateDatabaseDriver(c.Database.Driver) {
		return errors.New("InvalidDatabaseError: " + c.Database.Driver + " is not a valid database.")
	}
	if !validateMetadataAgent(c.Downloader.Metadata.Agent) {
		return errors.New("InvalidMetadataAgentError: " + c.Downloader.Metadata.Agent + " is not a valid agent.")
	}
	if !validateMetadataLocation(c.Downloader.Metadata.Location) {
		return errors.New("InvalidMetadataLocationError: " + c.Downloader.Metadata.Location + " is not a valid location.")
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
