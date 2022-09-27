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

// WithMetricNameTransform creates an Option for transforming Prometheus metric names
// Especially useful since Prometheus has a very limited range of characters allowed for
// metric names.
func WithMetricNameTransform(fun types.MetricNameTransformer) RecursiveMetricsOption {
	return &recursiveMetricNameTransform{
		fun: fun,
	}
}

type recursiveMetricsLabelNameTransform struct {
	fun types.LabelNameTransformer
}

type recursiveMetricNameTransform struct {
	fun types.MetricNameTransformer
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

	labelNameTransforms := []types.LabelNameTransformer{}
	metricNameTransforms := []types.MetricNameTransformer{}
	for _, opt := range opts {
		switch trans := opt.(type) {
		case *recursiveMetricsLabelNameTransform:
			labelNameTransforms = append(labelNameTransforms, trans.fun)
		case *recursiveMetricNameTransform:
			metricNameTransforms = append(metricNameTransforms, trans.fun)
		default:
			panic(fmt.Sprintf("Unrecognized option %#v", opt))
		}
	}

	var labelNameTransform types.LabelNameTransformer
	switch len(labelNameTransforms) {
	case 0:
		labelNameTransform = func(value string) (labelName string) { return value }
	case 1:
		labelNameTransform = labelNameTransforms[0]
	default:
		labelNameTransform = func(value string) (labelName string) {
			for _, trans := range labelNameTransforms {
				value = trans(value)
			}
			return value
		}
	}

	var metricNameTransform types.MetricNameTransformer
	switch len(metricNameTransforms) {
	case 0:
		metricNameTransform = func(value string) (metricName string) { return value }
	case 1:
		metricNameTransform = metricNameTransforms[0]
	default:
		metricNameTransform = func(value string) (metricName string) {
			for _, trans := range metricNameTransforms {
				value = trans(value)
			}
			return value
		}
	}

	rlr := label.RecursiveReflector{}
	fillLabels(t, &rlr, "", types.LabelNames{}, labelNameTransform)

	cs := &collector.Collectors{}
	cs.Fill(t, &rlr, "", metricNameTransform)
	u := &updater{
		c:                   cs,
		labelNameTransform:  labelNameTransform,
		metricNameTransform: metricNameTransform,
	}
	return cs, u
}
