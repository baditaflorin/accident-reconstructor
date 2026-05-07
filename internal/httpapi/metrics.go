package httpapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Metrics groups Prometheus collectors used by the API.
type Metrics struct {
	Requests       *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
	Jobs           *prometheus.CounterVec
	UploadedBytes  prometheus.Counter
	PipelineTime   prometheus.Histogram
}

// NewMetrics registers API and reconstruction metrics.
func NewMetrics(reg prometheus.Registerer) Metrics {
	m := Metrics{
		Requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "accident_http_requests_total",
			Help: "Total HTTP requests.",
		}, []string{"method", "path", "status"}),
		RequestLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "accident_http_request_duration_seconds",
			Help:    "HTTP request duration.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method", "path"}),
		Jobs: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "accident_reconstruction_jobs_total",
			Help: "Reconstruction jobs by result.",
		}, []string{"result"}),
		UploadedBytes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "accident_uploaded_bytes_total",
			Help: "Total uploaded bytes.",
		}),
		PipelineTime: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "accident_pipeline_duration_seconds",
			Help:    "Reconstruction pipeline duration.",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60, 180},
		}),
	}
	reg.MustRegister(m.Requests, m.RequestLatency, m.Jobs, m.UploadedBytes, m.PipelineTime)
	return m
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Middleware instruments HTTP requests.
func (m Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		route := r.URL.Path
		status := strconv.Itoa(recorder.status)
		m.Requests.WithLabelValues(r.Method, route, status).Inc()
		m.RequestLatency.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}
