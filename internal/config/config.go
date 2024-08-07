package config

import (
	"errors"
	"github.com/browningluke/mangathr/v2/internal/database"
	"github.com/browningluke/mangathr/v2/internal/downloader"
	"github.com/browningluke/mangathr/v2/internal/sources/cubari"
	"github.com/browningluke/mangathr/v2/internal/sources/mangadex"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"gopkg.in/yaml.v3"
	"k8s.io/helm/pkg/strvals"
	"os"
)

type Config struct {
	Database   database.Config
	Downloader downloader.Config
	Sources    struct {
		Mangadex mangadex.Config
		Cubari   cubari.Config
	}
	LogLevel string `yaml:"logLevel"`
}

func (c *Config) Propagate() {
	downloader.SetConfig(c.Downloader)
	database.SetConfig(c.Database)

	// Sources
	mangadex.SetConfig(c.Sources.Mangadex)
	cubari.SetConfig(c.Sources.Cubari)
}

func (c *Config) Load(path string, inContainer bool) (exists bool, err error) {
	c.useDefaults(inContainer)

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return false, err
	}
	if err = c.validate(); err != nil {
		return true, err
	}

	return true, nil
}

// Merge combines sVals `key1=value1` pairs with current config (overwriting existing values)
func (c *Config) Merge(sVals *[]string) error {
	// Convert current config to map[string]interface{}
	currentConfig, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	currentData := map[string]interface{}{}
	if err = yaml.Unmarshal(currentConfig, &currentData); err != nil {
		return err
	}

	// Parse strvals into map[string]interface{}
	strValData := map[string]interface{}{}
	for _, val := range *sVals {
		err := strvals.ParseInto(val, strValData)
		if err != nil {
			return err
		}
	}

	currentData = mergeMaps(currentData, strValData)

	// Convert map[string]interface{} back to config object
	d, err := yaml.Marshal(&currentData)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(d, &c); err != nil {
		return err
	}

	return nil
}

// mergeMaps used with "k8s.io/helm/pkg/strvals". Sourced from:
// https://github.com/helm/helm/blob/3bb50bbbdd9c946ba9989fbe4fb4104766302a64/pkg/cli/values/options.go#L108
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}

	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func (c *Config) useDefaults(inContainer bool) {
	downloadConf := downloader.Config{}
	downloadConf.Default(inContainer)
	c.Downloader = downloadConf

	databaseConf := database.Config{}
	databaseConf.Default(inContainer)
	c.Database = databaseConf

	mangadexConf := mangadex.Config{}
	mangadexConf.Default()
	c.Sources.Mangadex = mangadexConf

	c.LogLevel = ""

	// Overwrite defaults if we are in a container
	if inContainer {
	}
}

func (c *Config) validate() error {
	// Validate database
	if err := c.Database.Validate(); err != nil {
		return err
	}
	if !validateMetadataAgent(c.Downloader.Metadata.Agent) {
		return errors.New("InvalidMetadataAgentError: " + c.Downloader.Metadata.Agent + " is not a valid agent.")
	}
	if !validateMetadataLocation(c.Downloader.Metadata.Location) {
		return errors.New("InvalidMetadataLocationError: " + c.Downloader.Metadata.Location + " is not a valid location.")
	}
	return nil
}

func validateMetadataAgent(agent string) bool {
	_, exists := utils.FindInSliceFold([]string{"comicinfo", "json"}, agent)
	return exists
}

func validateMetadataLocation(location string) bool {
	_, exists := utils.FindInSliceFold([]string{"internal", "external", "both"}, location)
	return exists
}
