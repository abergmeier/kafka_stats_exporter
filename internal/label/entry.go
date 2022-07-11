package label

import (
	"reflect"
	"strconv"

	"github.com/abergmeier/kafka_stats_exporter/internal/assert"
	"github.com/prometheus/client_golang/prometheus"
)

type KeyValue struct {
	Key   string
	Value string
}

type KeyValueGenerator struct {
	FieldIndex int
	FieldType  reflect.Type
	LabelName  string
	T          reflect.Type
}

//Names are all possible names for Labels
type Names = []string

//LabelReflector uses a struct values to extract Labels
type Reflector struct {
	Generators []KeyValueGenerator
	T          reflect.Type
}

type RecursiveReflector struct {
	Lr     *Reflector
	Ln     Names
	Fields map[int]*RecursiveReflector
	T      reflect.Type
}

func (lr *Reflector) LabelsForValue(v interface{}) prometheus.Labels {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		rv = rv.Elem()
	}
	assert.AssertType(rv, lr.T)
	ls := make(prometheus.Labels, len(lr.Generators))
	for _, g := range lr.Generators {
		kv := g.GenerateFromValue(v)
		ls[kv.Key] = kv.Value
	}
	return ls
}

func (g *KeyValueGenerator) GenerateFromValue(v interface{}) KeyValue {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		rv = rv.Elem()
	}
	assert.AssertType(rv, g.T)
	rvk := rv.Kind()
	switch rvk {
	case reflect.Struct:
	case reflect.Pointer:
		rv = rv.Elem()
	default:
		panic("Unhandled code path")
	}
	fv := rv.FieldByIndex([]int{g.FieldIndex})
	assert.AssertType(fv, g.FieldType)
	if fv.CanInt() {
		return KeyValue{
			Key:   g.LabelName,
			Value: strconv.FormatInt(fv.Int(), 10),
		}
	}
	return KeyValue{
		Key:   g.LabelName,
		Value: fv.String(),
	}
}
