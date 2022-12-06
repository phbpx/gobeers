// Package metrics constructs the metrics the application will track.
package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// This holds the single instance of the metrics value needed for
// collecting metrics. The expvar package is already based on a singleton
// for the different metrics that are registered with the package so there
// isn't much choice here.
var m *metrics

// =============================================================================

// metrics represents the set of metrics we gather. These fields are
// safe to be accessed concurrently thanks to expvar. No extra abstraction is required.
type metrics struct {
	requests prometheus.Counter
	errors   prometheus.Counter
	panics   prometheus.Counter
}

// init constructs the metrics value that will be used to capture metrics.
// The metrics value is stored in a package level variable since everything
// inside of expvar is registered as a singleton. The use of once will make
// sure this initialization only happens once.
func init() {
	m = &metrics{
		requests: promauto.NewCounter(prometheus.CounterOpts{
			Name: "requests",
			Help: "Number of HTTP requests.",
		}),
		errors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "errors",
			Help: "Number of errors.",
		}),
		panics: promauto.NewCounter(prometheus.CounterOpts{
			Name: "panics",
			Help: "Number of panics.",
		}),
	}
}

// =============================================================================

// Metrics will be supported through the context.

// ctxKeyMetric represents the type of value for the context key.
type ctxKey int

// key is how metric values are stored/retrieved.
const key ctxKey = 1

// =============================================================================

// Set sets the metrics data into the context.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// Add more of these functions when a metric needs to be collected in
// different parts of the codebase. This will keep this package the
// central authority for metrics and metrics won't get lost.

// AddRequests increments the request metric by 1.
func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Inc()
	}
}

// AddErrors increments the errors metric by 1.
func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Inc()
	}
}

// AddPanics increments the panics metric by 1.
func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Inc()
	}
}
