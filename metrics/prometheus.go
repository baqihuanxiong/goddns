package metrics

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"goddns/config"
	"net/http"
)

// Prometheus metrics service implementation
// Reference: https://prometheus.io/docs/guides/go-application/
type PrometheusService struct {
	conf *config.Metrics

	httpRequestSummary *prometheus.SummaryVec
	ddnsCounter        prometheus.Counter
}

func NewPrometheusService(config *config.Metrics) (*PrometheusService, error) {
	svc := &PrometheusService{conf: config}

	svc.httpRequestSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "goddns",
		Name:       "http_request_duration_seconds",
		Help:       "The latency of HTTP requests.",
		Objectives: map[float64]float64{0.5: 0.05, 0.8: 0.05, 0.95: 0.05},
	}, []string{"handler", "method", "code"})

	svc.ddnsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "goddns",
		Name:      "ddns_total",
		Help:      "Total times a DDNS occur.",
	})

	return svc, nil
}

func (s *PrometheusService) CollectHttp(handler, method, statusCode string, duration float64) {
	s.httpRequestSummary.WithLabelValues(handler, method, statusCode).Observe(duration)
}

func (s *PrometheusService) IncreaseDDNSCounter() {
	s.ddnsCounter.Inc()
}

func (s *PrometheusService) Serve() error {
	router := mux.NewRouter()
	router.Handle(s.conf.MetricsPath, promhttp.Handler())
	router.HandleFunc(s.conf.HealthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	svr := &http.Server{
		Addr:    s.conf.Listen,
		Handler: router,
	}
	defer func() { _ = svr.Shutdown(context.TODO()) }()

	return svr.ListenAndServe()
}
