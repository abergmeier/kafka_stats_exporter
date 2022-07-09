package gen

import "github.com/prometheus/client_golang/prometheus"

//Updater allows for updating prometheus Metrics from Value
type Updater interface {
	Update(v interface{}, labels prometheus.Labels)
}

// UpdatingCollector is a Collector which can be updated.
type UpdatingCollector interface {
	prometheus.Collector
	Updater
}
