package gen

import (
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/internal/collector"
	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/types"
	"github.com/prometheus/client_golang/prometheus"
)

// NewRecursiveMetricsFromTags builds Metrics recursively for the type of `tagged`.
// Uses tags to build Metrics. Does not expose metrics directly.
// Returns `Collector` for reading created Metrics and `Updater` for
// updating the Metrics with a Type of `tagged`.
func NewRecursiveMetricsFromTags(tagged interface{}) (prometheus.Collector, Updater) {
	t := reflect.TypeOf(tagged)
	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
	}

	rlr := label.RecursiveReflector{}
	fillLabels(t, &rlr, "", types.LabelNames{})

	cs := &collector.Collectors{}
	cs.Fill(t, &rlr, "")
	u := &updater{
		c: cs,
	}
	return cs, u
}
