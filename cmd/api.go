package main

import (
	"log"
	"net/http"
	"time"
	repo "uptime-monitor/internal/adapters/postgresql/sqlc"
	"uptime-monitor/internal/sites"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All good."))
	})

	siteService := sites.NewService(repo.New(app.pool))
	siteHandler := sites.NewHandler(siteService)
	r.Get("/sites", siteHandler.ListSites)

	r.Post("/sites/add", siteHandler.AddSite)

	r.Post("/sites/remove", siteHandler.RemoveSite)

	r.Post("/sites/poll/enqueue", siteHandler.EnqueuePollSites)

	r.Post("/sites/poll/worker", siteHandler.PollSite)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

type application struct {
	config config
	pool   *pgxpool.Pool
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
	icn string
}
