package gen

import "github.com/prometheus/client_golang/prometheus"

// UpdatingCollector is a Collector which can be updated.
type UpdatingCollector interface {
	prometheus.Collector

	Update(v interface{})
}
