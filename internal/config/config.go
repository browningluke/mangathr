package config

import (
	"errors"
	"github.com/browningluke/mangathr/v2/internal/database"
	"github.com/browningluke/mangathr/v2/internal/downloader"
	"github.com/browningluke/mangathr/v2/internal/hooks"
	"github.com/browningluke/mangathr/v2/internal/sources/cubari"
	"github.com/browningluke/mangathr/v2/internal/sources/mangadex"
	"github.com/browningluke/mangathr/v2/internal/sources/mangaplus"
	"github.com/browningluke/mangathr/v2/internal/utils"
	"gopkg.in/yaml.v3"
	"k8s.io/helm/pkg/strvals"
	"os"
)

type Config struct {
	Database   database.Config
	Downloader downloader.Config
	Sources    struct {
		Mangadex  mangadex.Config
		Cubari    cubari.Config
		MangaPlus mangaplus.Config
	}
	Hooks    hooks.HooksConfig `yaml:"hooks"`
	LogLevel string            `yaml:"logLevel"`
}

func (c *Config) Propagate() {
	downloader.SetConfig(c.Downloader)
	database.SetConfig(c.Database)

	// Sources
	mangadex.SetConfig(c.Sources.Mangadex)
	cubari.SetConfig(c.Sources.Cubari)
	mangaplus.SetConfig(c.Sources.MangaPlus)

	hooks.SetConfig(c.Hooks)
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

	cubariConf := cubari.Config{}
	cubariConf.Default()
	c.Sources.Cubari = cubariConf

	mangaplusConf := mangaplus.Config{}
	mangaplusConf.Default()
	c.Sources.MangaPlus = mangaplusConf

	c.LogLevel = ""
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
	if err := c.Sources.MangaPlus.Validate(); err != nil {
		return err
	}
	if err := validateHooks(c.Hooks); err != nil {
		return err
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

var validHookEvents = []string{
	hooks.EventDownloadChapter,
	hooks.EventUpdateSuccess,
	hooks.EventUpdateError,
}

var validWebhookMethods = []string{"GET", "POST", "PUT", "PATCH"}

func validateHooks(cfg hooks.HooksConfig) error {
	for _, d := range cfg.Discord {
		if d.Name == "" {
			return errors.New("hooks.discord: each hook must have a name")
		}
		if d.WebhookURL == "" {
			return errors.New("hooks.discord: hook " + d.Name + " is missing webhookURL")
		}
		if err := validateHookEvents(d.On, "hooks.discord."+d.Name); err != nil {
			return err
		}
	}
	for _, w := range cfg.Webhook {
		if w.Name == "" {
			return errors.New("hooks.webhook: each hook must have a name")
		}
		if w.WebhookURL == "" {
			return errors.New("hooks.webhook: hook " + w.Name + " is missing webhookURL")
		}
		if w.RequestType != "" {
			if _, ok := utils.FindInSliceFold(validWebhookMethods, w.RequestType); !ok {
				return errors.New("hooks.webhook: hook " + w.Name + " has invalid requestType: " + w.RequestType)
			}
		}
		if err := validateHookEvents(w.On, "hooks.webhook."+w.Name); err != nil {
			return err
		}
	}
	for _, s := range cfg.Subcommand {
		if s.Name == "" {
			return errors.New("hooks.subcommand: each hook must have a name")
		}
		if s.Command == "" {
			return errors.New("hooks.subcommand: hook " + s.Name + " is missing command")
		}
		if err := validateHookEvents(s.On, "hooks.subcommand."+s.Name); err != nil {
			return err
		}
	}
	return nil
}

func validateHookEvents(on []string, hookPath string) error {
	if len(on) == 0 {
		return errors.New(hookPath + ": on must have at least one event")
	}
	for _, e := range on {
		if _, ok := utils.FindInSliceFold(validHookEvents, e); !ok {
			return errors.New(hookPath + ": unknown event: " + e +
				" (valid: download.chapter, update.success, update.error)")
		}
	}
	return nil
}
