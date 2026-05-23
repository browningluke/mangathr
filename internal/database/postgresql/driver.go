package postgresql

import (
	"database/sql"
	"github.com/browningluke/mangathr/v2/internal/logging"
	"github.com/lib/pq"
)

func CreateDatabase(dbName, connInfo string) error {
	logging.Debugln("Opening postgres connection")
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check whether the database already exists.
	// db.Exec on a SELECT returns 0 RowsAffected regardless of result count
	// with the pq driver, so we use QueryRow with EXISTS instead.
	var exists bool
	logging.Debugln("Checking if database exists")
	err = db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM pg_catalog.pg_database WHERE lower(datname) = lower($1))",
		dbName,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		logging.Debugln("Skipping creation. Database already exists")
		return nil
	}

	logging.Debugln("Executing db create SQL command")
	// pq.QuoteIdentifier safely quotes the identifier to prevent injection.
	_, err = db.Exec("CREATE DATABASE " + pq.QuoteIdentifier(dbName))
	if err != nil {
		// 42P04 = duplicate_database; harmless race where another process
		// created the DB between our existence check and this CREATE.
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42P04" {
			logging.Debugln("Database created concurrently, ignoring duplicate error")
			return nil
		}
		return err
	}

	return nil
}
