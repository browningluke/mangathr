package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/browningluke/mangathr/v2/internal/logging"
	_ "github.com/lib/pq"
)

func CreateDatabase(dbName, connInfo string) error {
	logging.Debugln("Opening postgres connection")
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return err
	}

	logging.Debugln("Executing db search SQL command")
	existsResult, err := db.Exec(
		fmt.Sprintf("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s');", dbName))
	if err != nil {
		return err
	}

	rowsAffected, err := existsResult.RowsAffected()
	if err != nil {
		return err
	}

	// If database doesn't exist, create it
	if rowsAffected == 0 {
		logging.Debugln("Executing db create SQL command")
		_, err = db.Exec(fmt.Sprintf("create database %s", dbName))
		if err != nil {
			return err
		}
	} else {
		logging.Debugln("Skipping creation. Database already exists")
	}

	return nil
}
