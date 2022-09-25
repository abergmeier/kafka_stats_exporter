package gen

import (
	"fmt"
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/internal/collector"
	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/types"
	"github.com/prometheus/client_golang/prometheus"
)

// RecursiveMetricsOption represents an opaque option implementation
// for creating RecursiveMetrics
type RecursiveMetricsOption interface {
}

// WithLabelNameTransform creates an Option for transforming Prometheus label names
// Especially useful since Prometheus has a very limited range of characters allowed for
// label names.
// Does not transform passed in Labels!
func WithLabelNameTransform(fun types.LabelNameTransformer) RecursiveMetricsOption {
	return &recursiveMetricsLabelNameTransform{
		fun: fun,
	}
}

type recursiveMetricsLabelNameTransform struct {
	fun types.LabelNameTransformer
}

// NewRecursiveMetricsFromTags builds Metrics recursively for the type of `tagged`.
// Uses tags to build Metrics. Does not expose metrics directly.
// Returns `Collector` for reading created Metrics and `Updater` for
// updating the Metrics with a Type of `tagged`.
func NewRecursiveMetricsFromTags(tagged interface{}, opts ...RecursiveMetricsOption) (prometheus.Collector, Updater) {
	t := reflect.TypeOf(tagged)
	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
	}

	transforms := []types.LabelNameTransformer{}
	for _, opt := range opts {
		switch trans := opt.(type) {
		case *recursiveMetricsLabelNameTransform:
			transforms = append(transforms, trans.fun)
		default:
			panic(fmt.Sprintf("Unrecognized option %#v", opt))
		}
	}

	var transform types.LabelNameTransformer
	switch len(transforms) {
	case 0:
		transform = func(value string) (labelName string) { return value }
	case 1:
		transform = transforms[0]
	default:
		transform = func(value string) (labelName string) {
			for _, trans := range transforms {
				value = trans(value)
			}
			return value
		}
	}

	rlr := label.RecursiveReflector{}
	fillLabels(t, &rlr, "", types.LabelNames{}, transform)

	cs := &collector.Collectors{}
	cs.Fill(t, &rlr, "")
	u := &updater{
		c:                  cs,
		labelNameTransform: transform,
	}
	return cs, u
}
