package main

import (
	"context"
	"log/slog"
	"os"
	"uptime-monitor/internal/database"
)

func main() {
	ctx := context.Background()
	
	cfg := config{
		addr: ":8080",
		db:   dbConfig{
			dsn: os.Getenv("DATABASE_SERVICE_ACCOUNT"), 
			icn: os.Getenv("ICN_STRING"),
		},
	}

	api := application{
		config: cfg,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	pool, cleanup, err := database.Connect(ctx, cfg.db.dsn, cfg.db.icn)
	if err != nil {
		slog.Error("Database has failed to connect", "error", err)
		os.Exit(1)
	}

	defer cleanup()
	defer pool.Close()

	logger.Info("Connected to database")


	if err := api.run(api.mount()); err != nil {
		slog.Error("Server has failed to start", "error", err)
		os.Exit(1)
	}

}
