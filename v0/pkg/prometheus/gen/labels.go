package gen

import (
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/types"
)

// Kept for compat only.
type LabelNames = types.LabelNames

type LabelReflector = label.Reflector

// MakeLabelReflector uses reflection Information to build LabelReflector.
// for structs. LabelReflector can then later on be used to performantly
// extract Labels from the same struct type.
// Internally LabelReflector specifically looks at field tag `kpromlbl` for
// Label names.
func MakeLabelReflector(t reflect.Type, parent string, parentLabelNames types.LabelNames) (*LabelReflector, types.LabelNames) {

	switch t.Kind() {
	case reflect.Struct:
	default:
		panic("LabelReflector can only be constructed from Struct")
	}

	var generators []label.KeyValueGenerator
	labelNames := types.LabelNames{}
	err := labelNames.AddLabelNames(&parentLabelNames)
	if err != nil {
		panic(err)
	}

	fields := reflect.VisibleFields(t)
	for i, f := range fields {
		tag := f.Tag.Get("kpromlbl")
		if tag == "" {
			continue
		}
		var key string
		if parent == "" {
			key = tag
		} else {
			key = parent + "_" + tag
		}

		err := labelNames.AddStrings(key)
		if err != nil {
			panic(err)
		}
		generators = append(generators, label.KeyValueGenerator{
			FieldIndex: i,
			FieldType:  f.Type,
			LabelName:  key,
			T:          t,
		})
	}

	labelNames.Sort()
	return &LabelReflector{
		Generators: generators,
		T:          t,
	}, labelNames
}

func fillLabels(t reflect.Type, rlr *label.RecursiveReflector, parent string, parentLabelNames types.LabelNames) {

	rlr.T = t

	rlr.Lr, rlr.Ln = MakeLabelReflector(t, parent, parentLabelNames)
	rlr.Fields = map[int]*label.RecursiveReflector{}

	for _, ft := range reflect.VisibleFields(t) {
		tag := ft.Tag.Get("kprommap")
		if tag != "" {
			switch ft.Type.Kind() {
			case reflect.Map:
			default:
				panic("Invalid code path")
			}
			frlr := label.RecursiveReflector{}
			rlr.Fields[ft.Index[0]] = &frlr
			if parent == "" {
				fillLabels(ft.Type.Elem(), &frlr, tag, rlr.Ln)
			} else {
				fillLabels(ft.Type.Elem(), &frlr, parent+"_"+tag, rlr.Ln)
			}
		}
		tag = ft.Tag.Get("kprompnt")
		if tag != "" {
			switch ft.Type.Kind() {
			case reflect.Struct:
			default:
				panic("Invalid code path")
			}
			frlr := label.RecursiveReflector{}
			rlr.Fields[ft.Index[0]] = &frlr
			if parent == "" {
				fillLabels(ft.Type, &frlr, tag, rlr.Ln)
			} else {
				fillLabels(ft.Type, &frlr, parent+"_"+tag, rlr.Ln)
			}
		}
	}

}
