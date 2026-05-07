package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/baditaflorin/accident-reconstructor/internal/config"
	"github.com/baditaflorin/accident-reconstructor/internal/jobs"
	"github.com/baditaflorin/accident-reconstructor/internal/reconstruction"
)

type App struct {
	Config    config.Config
	Store     *jobs.Store
	Processor reconstruction.Processor
	Metrics   Metrics
	Validate  *validator.Validate
	Logger    *slog.Logger
}

func NewRouter(cfg config.Config, store *jobs.Store, processor reconstruction.Processor) http.Handler {
	reg := prometheus.DefaultRegisterer
	app := App{
		Config:    cfg,
		Store:     store,
		Processor: processor,
		Metrics:   NewMetrics(reg),
		Validate:  validator.New(validator.WithRequiredStructEnabled()),
		Logger:    slog.Default(),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(app.logMiddleware)
	r.Use(app.Metrics.Middleware)
	r.Use(middleware.Timeout(10 * time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			cfg.PagesOrigin,
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	r.Get("/healthz", app.health)
	r.Get("/readyz", app.ready)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/tools", app.tools)
		r.Get("/cases", app.listCases)
		r.Post("/cases", app.createCase)
		r.Get("/cases/{caseID}", app.getCase)
		r.Get("/cases/{caseID}/artifact", app.getArtifact)
		r.Get("/cases/{caseID}/report", app.getReport)
		r.Get("/cases/{caseID}/bundle", app.getBundle)
	})
	return r
}

func (a App) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		a.Logger.Info(
			"http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"trace_id", middleware.GetReqID(r.Context()),
		)
	})
}
