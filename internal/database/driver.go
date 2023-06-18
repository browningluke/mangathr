package database

import (
	"context"
	"entgo.io/ent/dialect"
	"fmt"
	"github.com/browningluke/mangathr/ent"
	"github.com/browningluke/mangathr/internal/database/postgresql"
	_ "github.com/browningluke/mangathr/internal/database/sqlite3"
	"github.com/browningluke/mangathr/internal/logging"
	"github.com/browningluke/mangathr/internal/utils"
	_ "github.com/lib/pq"
	"strings"
)

type Driver struct {
	client *ent.Client
	ctx    context.Context
}

func GetDriver() (*Driver, error) {
	d := Driver{}

	// Extract driver information
	driverName := ""
	options := ""

	switch config.Driver {
	case "sqlite":
		driverName = dialect.SQLite
		options = fmt.Sprintf("file:%s?cache=shared", utils.ExpandHomePath(config.Sqlite.Path))
	case "postgres":
		driverName = dialect.Postgres
		options = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s %s",
			config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password,
			config.Postgres.DbName, config.Postgres.SSLMode, config.Postgres.Opts)

		// Create database (if needed)
		if config.CreateDatabase {
			logging.Debugln("Creating database `", config.Postgres.DbName, "`")

			psqlConnInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s %s",
				config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password,
				config.Postgres.SSLMode, config.Postgres.Opts)

			if err := postgresql.CreateDatabase("mangathr", psqlConnInfo); err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("%s", "Request driver not implemented")
	}

	// Open driver
	logging.Debugln("Opening database `", driverName, "` with options: ", options)
	client, err := ent.Open(driverName, options)
	if err != nil {
		return nil, err
	}

	// Run the auto migration tool.
	if config.AutoMigrate {
		logging.Debugln("Running auto-migration")
		if err := client.Schema.Create(context.Background()); err != nil {
			return nil, handleConnectionErrors(err)
		}
	}

	// Test connection
	logging.Debugln("Testing connection...")
	_, err = client.Manga.Query().Count(context.Background())
	if err != nil {
		return nil, handleConnectionErrors(err)
	}

	logging.Debugln("Driver connection succeeded")
	d.client = client
	d.ctx = context.Background()
	return &d, nil
}

func (d *Driver) Close() error {
	return d.client.Close()
}

func handleConnectionErrors(err error) error {
	// Postgres: database doesn't exist
	if strings.Contains(err.Error(), fmt.Sprintf("pq: database \"%s\" does not exist", config.Postgres.DbName)) {
		return fmt.Errorf("database doesn't exist, " +
			"ensure database has been created or set 'database.createDatabase' to true")
	}

	// SQLite3: database could not be created
	if strings.Contains(err.Error(),
		"sqlite: check foreign_keys pragma: "+
			"reading schema information unable to open database file: out of memory (14)") {
		return fmt.Errorf("database could not be created, " +
			"ensure database path is valid and user has correct permissions to access path")
	}

	// Postgres: mangas doesn't exist (migration hasn't run)
	if err.Error() == "pq: relation \"mangas\" does not exist" {
		return fmt.Errorf("database hasn't been migrated, " +
			"ensure database is created correctly or set 'database.autoMigrate' to true")
	}

	// Postgres: password auth failed for user
	if strings.Contains(err.Error(), "password authentication failed for user") {
		return fmt.Errorf("password authentication failed, check credentials are valid")
	}
	return err
}
