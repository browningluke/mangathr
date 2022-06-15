package database

import (
	"context"
	"fmt"
	"mangathrV2/ent"
	_ "mangathrV2/internal/database/sqlite3"
)

const (
	SQLITE = iota
	// PSQL
)

type Driver struct {
	client *ent.Client
	ctx    context.Context
}

func GetDriver(driver int) (*Driver, error) {
	d := Driver{}

	// Extract driver information
	driverName := ""
	options := ""

	switch driver {
	case SQLITE:
		driverName = "sqlite3"
		options = "file:examples/db.sqlite?cache=shared&_fk=1"
		break
	default:
		return nil, fmt.Errorf("%s", "Request driver not implemented")
	}

	// Open driver
	client, err := ent.Open(driverName, options)
	if err != nil {
		return nil, err
	}

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	d.client = client
	d.ctx = context.Background()
	return &d, nil
}

func (d *Driver) Close() error {
	return d.client.Close()
}
