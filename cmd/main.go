package main

import (
	"context"
	"log/slog"
	"os"
	"uptime-monitor/internal/database"
)

func main() {
	ctx := context.Background()

	portWithoutColin, portExists := os.LookupEnv("PORT")
	port := ":" + portWithoutColin
	if !portExists {
		port = ":8080"
	}
	databaseServiceAccount, databaseServiceAccountExists := os.LookupEnv("DATABASE_SERVICE_ACCOUNT")
	if !databaseServiceAccountExists {
		databaseServiceAccount = ""
	}
	icnString, icnStringExists := os.LookupEnv("ICN_STRING")
	if !icnStringExists {
		icnString = ""
	}

	cfg := config{
		addr: port,
		db: dbConfig{
			dsn: databaseServiceAccount,
			icn: icnString,
		},
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	pool, cleanup, err := database.Connect(ctx, cfg.db.dsn, cfg.db.icn)
	if err != nil {
		slog.Error("Database has failed to connect", "error:", err)
		os.Exit(1)
	}

	api := application{
		config: cfg,
		pool: pool,
	}

	defer cleanup()
	defer pool.Close()

	logger.Info("Connected to database")

	if err := api.run(api.mount()); err != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}

}
