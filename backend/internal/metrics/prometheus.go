// internal/metrics/prometheus.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// EventMetrics holds the Prometheus metrics for our event system
type EventMetrics struct {
	EventsPublished   *prometheus.CounterVec
	EventsReceived    *prometheus.CounterVec
	EventProcessTime  *prometheus.HistogramVec
	EventsErrorCount  *prometheus.CounterVec
	ActiveSubscribers *prometheus.GaugeVec
}

// NewEventMetrics creates and registers Prometheus metrics for the event system
func NewEventMetrics() *EventMetrics {
	return &EventMetrics{
		EventsPublished: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "coffee_commerce_events_published_total",
				Help: "Total number of events published",
			},
			[]string{"topic", "service"},
		),
		EventsReceived: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "coffee_commerce_events_received_total",
				Help: "Total number of events received",
			},
			[]string{"topic", "service"},
		),
		EventProcessTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "coffee_commerce_event_process_time_seconds",
				Help:    "Time taken to process events",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"topic", "service"},
		),
		EventsErrorCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "coffee_commerce_events_error_total",
				Help: "Total number of errors processing events",
			},
			[]string{"topic", "service", "error_type"},
		),
		ActiveSubscribers: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "coffee_commerce_event_subscribers_active",
				Help: "Number of active subscribers per topic",
			},
			[]string{"topic"},
		),
	}
}