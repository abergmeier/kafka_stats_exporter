package gen

import (
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/internal/collector"
	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/prometheus/client_golang/prometheus"
)

func NewRecursiveUpdaterFromTags(tagged interface{}) (prometheus.Collector, Updater) {
	t := reflect.TypeOf(tagged)
	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
	}

	rlr := label.RecursiveLabelReflector{}
	fillLabels(t, &rlr, "", nil)

	cs := &collector.Collectors{}
	cs.Fill(t, &rlr, "", LabelNames{})
	u := &updater{
		c: cs,
	}
	return cs, u
}
