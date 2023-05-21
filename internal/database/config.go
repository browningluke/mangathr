package database

import (
	"fmt"
	"github.com/browningluke/mangathr/internal/config/defaults"
	"github.com/browningluke/mangathr/internal/utils"
	"path/filepath"
	"strings"
)

var config Config

type Config struct {
	Driver         string
	CreateDatabase bool `yaml:"createDatabase"`
	AutoMigrate    bool `yaml:"autoMigrate"`
	Sqlite         struct {
		Path string
	}
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		DbName   string `yaml:"dbName"`
		SSLMode  string `yaml:"sslMode"`
		Opts     string
	}
}

func SetConfig(cfg Config) {
	config = cfg
}

// Validate ensures all values are valid, and cleans up values (folding, etc.)
func (c *Config) Validate() error {
	// Valid options
	validDrivers := []string{"sqlite", "postgres"}
	validPostgresSSLMode := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}

	// Driver
	if _, exists := utils.FindInSliceFold(validDrivers, c.Driver); !exists {
		return fmt.Errorf("database.driver not in [%s]", strings.Join(validDrivers, ", "))
	}

	// Postgres.SSLMode
	if _, exists := utils.FindInSliceFold(validPostgresSSLMode, c.Postgres.SSLMode); !exists {
		return fmt.Errorf("database.postgres.sslMode not in [%s]",
			strings.Join(validPostgresSSLMode, ", "))
	}

	// Postgres.Password
	if len(c.Postgres.Password) == 0 {
		return fmt.Errorf("database.postgres.password is empty")
	}

	return nil
}

func (c *Config) Default(inContainer bool) {
	c.Driver = "sqlite"
	c.CreateDatabase = true
	c.AutoMigrate = true

	// Sqlite
	c.Sqlite.Path = filepath.Join(defaults.ConfigDir(), "db.sqlite")

	// Postgres
	c.Postgres.Host = "127.0.0.1"
	c.Postgres.Port = "5432"
	c.Postgres.User = "mangathr"
	c.Postgres.DbName = "mangathr"
	c.Postgres.SSLMode = "disable"
	c.Postgres.Opts = ""

	if inContainer {
		c.Sqlite.Path = filepath.Join(defaults.ConfigDirDocker(), "db.sqlite")
	}
}
