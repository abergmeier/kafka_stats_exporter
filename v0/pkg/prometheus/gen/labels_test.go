package gen

import (
	"reflect"
	"sort"
	"testing"

	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	expectedLabels = prometheus.Labels{
		"name":      "MyName",
		"client_id": "MyClientId",
		"type":      "MyType",
	}
	expectedLabelNames = LabelNames{
		"client_id", "name", "type",
	}
)

func TestMakeLabelGenerator(t *testing.T) {
	tpe := reflect.TypeOf(typed.Stats{})
	lg, lns := MakeLabelReflector(tpe, "")
	ls := lg(typed.Stats{
		Name:     "MyName",
		ClientId: "MyClientId",
		Type:     "MyType",
	})
	if !reflect.DeepEqual(ls, expectedLabels) {
		t.Fatal("Invalid labels generated:", ls)
	}
	sort.Strings(lns)
	if !reflect.DeepEqual(lns, expectedLabelNames) {
		t.Fatal("Invalid label names generated:", lns)
	}
}
