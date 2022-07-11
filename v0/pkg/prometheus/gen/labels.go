package gen

import (
	"reflect"
	"sort"

	"github.com/abergmeier/kafka_stats_exporter/internal/label"
)

//LabelNames are all possible names for Labels
type LabelNames = label.Names

type labelNameCache map[reflect.Type]LabelNames

type LabelReflector = label.Reflector

// MakeLabelReflector uses reflection Information to build LabelReflector.
// for structs. LabelReflector can then later on be used to performantly
// extract Labels from the same struct type.
// Internally LabelReflector specifically looks at field tag `kpromlbl` for
// Label names.
func MakeLabelReflector(t reflect.Type, parent string, parentLabelNames LabelNames) (*LabelReflector, LabelNames) {

	switch t.Kind() {
	case reflect.Struct:
	default:
		panic("LabelReflector can only be constructed from Struct")
	}

	var generators []label.KeyValueGenerator
	labelNames := LabelNames{}
	for _, ln := range parentLabelNames {
		labelNames = append(labelNames, ln)
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

		labelNames = append(labelNames, key)
		generators = append(generators, label.KeyValueGenerator{
			FieldIndex: i,
			FieldType:  f.Type,
			LabelName:  key,
			T:          t,
		})
	}

	sort.Strings(labelNames)
	return &LabelReflector{
		Generators: generators,
		T:          t,
	}, labelNames
}

func fillLabels(t reflect.Type, rlr *label.RecursiveLabelReflector, parent string, parentLabelNames LabelNames) {

	rlr.T = t

	rlr.Lr, rlr.Ln = MakeLabelReflector(t, parent, parentLabelNames)
	rlr.Fields = map[int]*label.RecursiveLabelReflector{}

	for _, ft := range reflect.VisibleFields(t) {
		tag := ft.Tag.Get("kprommap")
		if tag != "" {
			switch ft.Type.Kind() {
			case reflect.Map:
			default:
				panic("Invalid code path")
			}
			frlr := label.RecursiveLabelReflector{}
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
			frlr := label.RecursiveLabelReflector{}
			rlr.Fields[ft.Index[0]] = &frlr
			if parent == "" {
				fillLabels(ft.Type, &frlr, tag, rlr.Ln)
			} else {
				fillLabels(ft.Type, &frlr, parent+"_"+tag, rlr.Ln)
			}
		}
	}

}
