package monitoring

import "github.com/prometheus/client_golang/prometheus"

type PrometheusMonitoring struct {
	httpMonitoringCounter   *prometheus.CounterVec
	httpMonitoringHistogram *prometheus.HistogramVec
}

type IMonitoring interface {
	CountJoiningRoom(endpointName string, statusCode int, errorMsg string, latency float64)
}

func NewPrometheusMonitoring() IMonitoring {
	httpMonitoringCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_endpoints_monitoring",
			Help: "HTTP monitoring for joined method in groupchat service",
		},
		[]string{
			"endpoint_name",
			"status_code",
			"error",
		},
	)

	httpMonitoringHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_endpoints_latency",
			Help: "HTTP latency monitoring",
		},
		[]string{
			"endpoint_name",
			"status_code",
			"error",
		},
	)

	prometheus.MustRegister(httpMonitoringCounter)
	prometheus.MustRegister(httpMonitoringHistogram)

	return &PrometheusMonitoring{
		httpMonitoringCounter:   httpMonitoringCounter,
		httpMonitoringHistogram: httpMonitoringHistogram,
	}
}
