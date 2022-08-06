package sqlite3

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"modernc.org/sqlite"
)

// SQLite wrapper to bridge modernc/sqlite and ent provided by tux21b here:
// https://github.com/ent/ent/discussions/1667#discussioncomment-1132296

type sqliteDriver struct {
	*sqlite.Driver
}

func (d sqliteDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		conn.Close()
		return nil, fmt.Errorf("%s%s", err, " failed to enable enable foreign keys")
	}
	return conn, nil
}

func init() {
	sql.Register("sqlite3", sqliteDriver{Driver: &sqlite.Driver{}})
}
