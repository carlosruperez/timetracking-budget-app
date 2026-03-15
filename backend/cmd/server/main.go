package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rupi/timetracking/internal/auth"
	"github.com/rupi/timetracking/internal/budget"
	"github.com/rupi/timetracking/internal/category"
	"github.com/rupi/timetracking/internal/config"
	"github.com/rupi/timetracking/internal/db"
	"github.com/rupi/timetracking/internal/report"
	"github.com/rupi/timetracking/internal/timer"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	log.Info().Msg("migrations applied")

	// Wire up dependencies
	jwtSvc := auth.NewJWTService(cfg.JWTSecret)

	authSvc := auth.NewService(database, jwtSvc)
	authHandler := auth.NewHandler(authSvc)
	authMiddleware := auth.Middleware(jwtSvc)

	catRepo := category.NewRepository(database)
	catSvc := category.NewService(catRepo)
	catHandler := category.NewHandler(catSvc)

	timerRepo := timer.NewRepository(database)
	timerSvc := timer.NewService(timerRepo)
	timerHandler := timer.NewHandler(timerSvc)

	budgetRepo := budget.NewRepository(database)
	budgetSvc := budget.NewService(budgetRepo)
	budgetHandler := budget.NewHandler(budgetSvc)

	reportSvc := report.NewService(database)
	reportHandler := report.NewHandler(reportSvc)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://*.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Auth (public)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/auth/me", authHandler.Me)

			// Categories
			r.Get("/categories", catHandler.List)
			r.Post("/categories", catHandler.Create)
			r.Get("/categories/{id}", catHandler.Get)
			r.Put("/categories/{id}", catHandler.Update)
			r.Delete("/categories/{id}", catHandler.Delete)

			// Timer
			r.Get("/timer/active", timerHandler.GetActive)
			r.Get("/timer/stream", timerHandler.Stream)
			r.Post("/timer/start", timerHandler.Start)
			r.Post("/timer/pause", timerHandler.Pause)
			r.Post("/timer/resume", timerHandler.Resume)
			r.Post("/timer/stop", timerHandler.Stop)
			r.Get("/timer/entries", timerHandler.ListEntries)
			r.Get("/timer/entries/{id}", timerHandler.GetEntry)
			r.Put("/timer/entries/{id}", timerHandler.UpdateEntry)
			r.Delete("/timer/entries/{id}", timerHandler.DeleteEntry)

			// Budgets
			r.Get("/budgets", budgetHandler.List)
			r.Post("/budgets", budgetHandler.Create)
			r.Put("/budgets/{id}", budgetHandler.Update)
			r.Delete("/budgets/{id}", budgetHandler.Delete)
			r.Get("/budgets/status", budgetHandler.GetStatus)

			// Reports
			r.Get("/reports/summary", reportHandler.Summary)
			r.Get("/reports/daily", reportHandler.Daily)
			r.Get("/reports/weekly", reportHandler.Weekly)
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Info().Str("port", cfg.Port).Msg("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
	}
}
