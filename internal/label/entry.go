package label

import (
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/internal/assert"
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
	return KeyValue{
		Key:   g.LabelName,
		Value: fv.String(),
	}
}
