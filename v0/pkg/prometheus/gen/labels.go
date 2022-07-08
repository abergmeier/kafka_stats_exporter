package gen

import (
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
)

type labelKeyValue struct {
	key   string
	value string
}

type LabelNames []string
type LabelReflector func(interface{}) prometheus.Labels
type oneLabelGenerator []func(interface{}) labelKeyValue

// MakeLabelReflector uses reflection Information to build LabelReflector.
// for structs.
// It specifically looks at field tag `kpromlbl` for Keys.
func MakeLabelReflector(t reflect.Type, parent string) (LabelReflector, LabelNames) {

	switch t.Kind() {
	case reflect.Struct:
	default:
		panic("LabelReflector can only be constructed from Struct")
	}

	var generators oneLabelGenerator
	var labelNames LabelNames

	fields := reflect.VisibleFields(t)
	for i, f := range fields {
		tag := f.Tag.Get("kpromlbl")
		if tag == "" {
			continue
		}
		key := parent + tag
		labelNames = append(labelNames, key)
		// This is the ONE thing that is absolutely insane about Go
		fi := i
		generators = append(generators, func(v interface{}) labelKeyValue {
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Struct:
				fv := rv.FieldByIndex([]int{fi})
				if !fv.CanInterface() {
					panic("Can not interface")
				}
				return labelKeyValue{
					key:   key,
					value: fv.String(),
				}
			default:
				panic("Not implemented yet")
			}
		})
	}

	return func(v interface{}) prometheus.Labels {
		ls := make(prometheus.Labels, len(generators))
		for _, g := range generators {
			kv := g(v)
			ls[kv.key] = kv.value
		}
		return ls
	}, labelNames
}
