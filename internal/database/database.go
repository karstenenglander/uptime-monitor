package database

import (
	"context"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, dsn string, icn string) (*pgxpool.Pool, func() error, error){
	// Configure the driver to connect to the database
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		/* handle error */
	}

	// Create a new dialer with any options
	d, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithIAMAuthN())
	if err != nil {
		/* handle error */
	}

	// Tell the driver to use the Cloud SQL Go Connector to create connections
	config.ConnConfig.DialFunc = func(ctx context.Context, _ string, instance string) (net.Conn, error) {
		return d.Dial(ctx, icn)
	}

	// Interact with the driver directly as you normally would
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		/* handle error */
	}

	// call cleanup when you're done with the database connection
	cleanup := func() error { return d.Close() }
	// ... etc

	return pool, cleanup, nil
}