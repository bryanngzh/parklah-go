package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/bryanngzh/parklah-go/internal/config"
	"github.com/bryanngzh/parklah-go/internal/db"
	"github.com/bryanngzh/parklah-go/internal/handlers"
	"github.com/bryanngzh/parklah-go/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(context.Background(), cfg.DSN())
	if err != nil {
		log.Fatalf("[api] DB connect: %v", err)
	}
	defer pool.Close()
	log.Println("[api] Connected to database")

	// Fetch Singapore public holidays for current and next year at startup
	phDates := make(map[string]bool)
	now := time.Now()
	for _, year := range []int{now.Year(), now.Year() + 1} {
		ph, err := util.FetchSGPublicHolidays(context.Background(), year)
		if err != nil {
			log.Printf("[api] Warning: failed to fetch SG public holidays for %d: %v", year, err)
			continue
		}
		for k := range ph {
			phDates[k] = true
		}
	}
	log.Printf("[api] Loaded %d SG public holiday dates", len(phDates))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1/carparks", func(r chi.Router) {
		r.Get("/nearby", handlers.GetNearby(pool))
		r.Post("/batch", handlers.GetBatch(pool))
		r.Get("/{code}", handlers.GetCarpark(pool))
		r.Get("/{code}/availability", handlers.GetAvailability(pool))
		r.Get("/{code}/rates", handlers.GetRates(pool, phDates))
	})

	addr := ":" + cfg.APIPort
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[api] Listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[api] Server error: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("[api] Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[api] Shutdown error: %v", err)
	}
	log.Println("[api] Stopped")
}
