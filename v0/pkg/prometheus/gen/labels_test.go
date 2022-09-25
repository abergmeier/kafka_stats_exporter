package gen

import (
	"reflect"
	"testing"

	"github.com/abergmeier/kafka_stats_exporter/internal"
	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/types"
	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	expectedLabels = prometheus.Labels{
		"name":      "MyName",
		"client_id": "MyClientId",
		"type":      "MyType",
	}
	expectedLabelNames = types.LabelNames{}
	expectedRecursive  = label.RecursiveReflector{
		Ln: testLabelNames("client_id", "name", "type"),
		T:  reflect.TypeOf(typed.Stats{}),
		Lr: &LabelReflector{
			T: reflect.TypeOf(typed.Stats{}),
			Generators: []label.KeyValueGenerator{
				{FieldIndex: 0, FieldType: reflect.TypeOf(""), LabelName: "name", T: reflect.TypeOf(typed.Stats{})},
				{FieldIndex: 1, FieldType: reflect.TypeOf(""), LabelName: "client_id", T: reflect.TypeOf(typed.Stats{})},
				{FieldIndex: 2, FieldType: reflect.TypeOf(""), LabelName: "type", T: reflect.TypeOf(typed.Stats{})},
			},
		},
		Fields: map[int]*label.RecursiveReflector{
			21: {
				Ln: testLabelNames("brokers_name", "brokers_nodeid", "brokers_nodename", "brokers_source", "brokers_state", "client_id", "name", "type"),
				T:  reflect.TypeOf(typed.BrokerStats{}),
				Fields: map[int]*label.RecursiveReflector{
					27: {
						Fields: map[int]*label.RecursiveReflector{},
						Ln:     testLabelNames("brokers_name", "brokers_nodeid", "brokers_nodename", "brokers_source", "brokers_state", "client_id", "name", "type"),
						Lr:     &LabelReflector{T: reflect.TypeOf(typed.WindowStats{})},
						T:      reflect.TypeOf(typed.WindowStats{}),
					},
					28: {
						Fields: map[int]*label.RecursiveReflector{},
						Ln:     testLabelNames("brokers_name", "brokers_nodeid", "brokers_nodename", "brokers_source", "brokers_state", "client_id", "name", "type"),
						Lr:     &LabelReflector{T: reflect.TypeOf(typed.WindowStats{})},
						T:      reflect.TypeOf(typed.WindowStats{}),
					},
					29: {
						Fields: map[int]*label.RecursiveReflector{},
						Ln:     testLabelNames("brokers_name", "brokers_nodeid", "brokers_nodename", "brokers_source", "brokers_state", "client_id", "name", "type"),
						Lr:     &LabelReflector{T: reflect.TypeOf(typed.WindowStats{})},
						T:      reflect.TypeOf(typed.WindowStats{}),
					},
					30: {
						Fields: map[int]*label.RecursiveReflector{},
						Ln:     testLabelNames("brokers_name", "brokers_nodeid", "brokers_nodename", "brokers_source", "brokers_state", "client_id", "name", "type"),
						Lr:     &LabelReflector{T: reflect.TypeOf(typed.WindowStats{})},
						T:      reflect.TypeOf(typed.WindowStats{}),
					},
				},
				Lr: &LabelReflector{
					T: reflect.TypeOf(typed.BrokerStats{}),
					Generators: []label.KeyValueGenerator{
						{FieldIndex: 0, FieldType: reflect.TypeOf(""), LabelName: "brokers_name", T: reflect.TypeOf(typed.BrokerStats{})},
						{FieldIndex: 1, FieldType: reflect.TypeOf(0), LabelName: "brokers_nodeid", T: reflect.TypeOf(typed.BrokerStats{})},
						{FieldIndex: 2, FieldType: reflect.TypeOf(""), LabelName: "brokers_nodename", T: reflect.TypeOf(typed.BrokerStats{})},
						{FieldIndex: 3, FieldType: reflect.TypeOf(""), LabelName: "brokers_source", T: reflect.TypeOf(typed.BrokerStats{})},
						{FieldIndex: 4, FieldType: reflect.TypeOf(""), LabelName: "brokers_state", T: reflect.TypeOf(typed.BrokerStats{})},
					},
				},
			},
			22: {
				Ln: testLabelNames("client_id", "name", "topics_topic", "type"),
				T:  reflect.TypeOf(typed.TopicStats{}),
				Lr: &LabelReflector{
					T: reflect.TypeOf(typed.TopicStats{}),
					Generators: []label.KeyValueGenerator{
						{FieldIndex: 0, FieldType: reflect.TypeOf(""), LabelName: "topics_topic", T: reflect.TypeOf(typed.TopicStats{})},
					},
				},
				Fields: map[int]*label.RecursiveReflector{
					3: {
						Fields: map[int]*label.RecursiveReflector{},
						Ln:     testLabelNames("client_id", "name", "topics_topic", "type"),
						Lr:     &LabelReflector{T: reflect.TypeOf(typed.WindowStats{})},
						T:      reflect.TypeOf(typed.WindowStats{}),
					},
					4: {
						Fields: map[int]*label.RecursiveReflector{},
						Ln:     testLabelNames("client_id", "name", "topics_topic", "type"),
						Lr:     &LabelReflector{T: reflect.TypeOf(typed.WindowStats{})},
						T:      reflect.TypeOf(typed.WindowStats{}),
					},
					5: {
						Fields: map[int]*label.RecursiveReflector{},
						Ln:     testLabelNames("client_id", "name", "topics_partitions_broker", "topics_partitions_fetch_state", "topics_partitions_leader", "topics_partitions_partition", "topics_topic", "type"),
						T:      reflect.TypeOf(typed.PartitionStats{}),
						Lr: &LabelReflector{
							T: reflect.TypeOf(typed.PartitionStats{}),
							Generators: []label.KeyValueGenerator{
								{FieldIndex: 0, FieldType: reflect.TypeOf(0), LabelName: "topics_partitions_partition", T: reflect.TypeOf(typed.PartitionStats{})},
								{FieldIndex: 1, FieldType: reflect.TypeOf(0), LabelName: "topics_partitions_broker", T: reflect.TypeOf(typed.PartitionStats{})},
								{FieldIndex: 2, FieldType: reflect.TypeOf(0), LabelName: "topics_partitions_leader", T: reflect.TypeOf(typed.PartitionStats{})},
								{FieldIndex: 11, FieldType: reflect.TypeOf(""), LabelName: "topics_partitions_fetch_state", T: reflect.TypeOf(typed.PartitionStats{})},
							},
						},
					},
				},
			},
			23: {
				Fields: map[int]*label.RecursiveReflector{},
				Ln:     testLabelNames("cgrp_join_state", "cgrp_rebalance_reason", "cgrp_state", "client_id", "name", "type"),
				T:      reflect.TypeOf(typed.CgrpStats{}),
				Lr: &LabelReflector{
					T: reflect.TypeOf(typed.CgrpStats{}),
					Generators: []label.KeyValueGenerator{
						{FieldIndex: 0, FieldType: reflect.TypeOf(""), LabelName: "cgrp_state", T: reflect.TypeOf(typed.CgrpStats{})},
						{FieldIndex: 2, FieldType: reflect.TypeOf(""), LabelName: "cgrp_join_state", T: reflect.TypeOf(typed.CgrpStats{})},
						{FieldIndex: 5, FieldType: reflect.TypeOf(""), LabelName: "cgrp_rebalance_reason", T: reflect.TypeOf(typed.CgrpStats{})},
					},
				},
			},
			24: {
				Fields: map[int]*label.RecursiveReflector{},
				Ln:     testLabelNames("client_id", "eos_idemp_state", "eos_producer_id", "eos_txn_state", "name", "type"),
				T:      reflect.TypeOf(typed.EosStats{}),
				Lr: &LabelReflector{
					T: reflect.TypeOf(typed.EosStats{}),
					Generators: []label.KeyValueGenerator{
						{FieldIndex: 0, FieldType: reflect.TypeOf(""), LabelName: "eos_idemp_state", T: reflect.TypeOf(typed.EosStats{})},
						{FieldIndex: 2, FieldType: reflect.TypeOf(""), LabelName: "eos_txn_state", T: reflect.TypeOf(typed.EosStats{})},
						{FieldIndex: 5, FieldType: reflect.TypeOf(0), LabelName: "eos_producer_id", T: reflect.TypeOf(typed.EosStats{})},
					},
				},
			},
		},
	}
	expectedSimpleLabelReflector = label.RecursiveReflector{
		T:  reflect.TypeOf(simpleStats{}),
		Ln: testLabelNames("name"),
		Lr: &LabelReflector{
			Generators: []label.KeyValueGenerator{
				{FieldIndex: 0, FieldType: reflect.TypeOf(""), LabelName: "name", T: reflect.TypeOf(simpleStats{})},
			},
			T: reflect.TypeOf(simpleStats{}),
		},
		Fields: map[int]*label.RecursiveReflector{
			2: {
				T: reflect.TypeOf(simpleBrokerStats{}),
				Lr: &LabelReflector{
					T: reflect.TypeOf(simpleBrokerStats{}),
					Generators: []label.KeyValueGenerator{
						{FieldIndex: 0, FieldType: reflect.TypeOf(""), LabelName: "brokers_name", T: reflect.TypeOf(simpleBrokerStats{})},
					},
				},
				Ln:     testLabelNames("brokers_name", "name"),
				Fields: map[int]*label.RecursiveReflector{},
			},
		},
	}
)

func init() {
	err := expectedLabelNames.AddStrings("client_id", "name", "type")
	if err != nil {
		panic(err)
	}
}

func TestMakeLabelReflector(t *testing.T) {
	tpe := reflect.TypeOf(typed.Stats{})
	lg, lns := MakeLabelReflector(tpe, "", types.LabelNames{})
	ls := lg.LabelsForValue(typed.Stats{
		Name:     "MyName",
		ClientId: "MyClientId",
		Type:     "MyType",
	})
	if !reflect.DeepEqual(ls, expectedLabels) {
		t.Fatal("Invalid labels generated:", ls)
	}
	lns.Sort()
	if !reflect.DeepEqual(lns, expectedLabelNames) {
		t.Fatal("Invalid label names generated:", lns)
	}
}

func TestLabelReflectorSimple(t *testing.T) {
	rlr := label.RecursiveReflector{}
	fillLabels(reflect.TypeOf(simpleStats{}), &rlr, "", types.LabelNames{}, func(value string) (labelName string) { return value })
	d := cmp.Diff(rlr, expectedSimpleLabelReflector, cmp.Comparer(internal.CompareType))
	if d != "" {
		t.Fatal("Diff", d)
	}
}

func TestLabelReflectorFull(t *testing.T) {
	rlr := label.RecursiveReflector{}
	fillLabels(reflect.TypeOf(typed.Stats{}), &rlr, "", types.LabelNames{}, func(value string) (labelName string) { return value })
	d := cmp.Diff(rlr, expectedRecursive, cmp.Comparer(internal.CompareType))
	if d != "" {
		t.Fatal("Diff", d)
	}
}

func testLabelNames(names ...string) types.LabelNames {
	lns := types.LabelNames{}
	err := lns.AddStrings(names...)
	if err != nil {
		panic(err)
	}
	return lns
}
