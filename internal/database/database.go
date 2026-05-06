package database

import (
	"context"
	"fmt"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, dsn string, icn string) (*pgxpool.Pool, func() error, error) {
	// Configure the driver to connect to the database
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		err := fmt.Errorf("Connection not established. Config could not be parsed. Error: %w", err)
		return nil, nil, err
	}

	// Create a new dialer with any options
	d, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithIAMAuthN())
	// call cleanup when you're done with the database connection
	cleanup := func() error { return d.Close() }

	if err != nil {
		if d != nil {
			cleanup()
		}
		err := fmt.Errorf("Connection not established. Dialer could not be created. Error: %w", err)
		return nil, nil, err
	}

	// Tell the driver to use the Cloud SQL Go Connector to create connections
	config.ConnConfig.DialFunc = func(ctx context.Context, _ string, instance string) (net.Conn, error) {
		return d.Dial(ctx, icn)
	}

	// Interact with the driver directly as you normally would
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		if pool != nil {
			pool.Close()
		}
		err := fmt.Errorf("Connection not established. Pool could not be created. Error: %w", err)
		cleanup()
		return nil, nil, err
	}

	return pool, cleanup, nil
}
