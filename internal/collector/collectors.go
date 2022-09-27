package collector

import (
	"fmt"
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/types"
	"github.com/prometheus/client_golang/prometheus"
)

type Collectors struct {
	Rlr              *label.RecursiveReflector
	StaticCollectors []GeneratedUpdator
	children         []prometheus.Collector
	Maps             []DynamicMap
	T                reflect.Type
}

func (u *Collectors) Describe(c chan<- *prometheus.Desc) {
	for _, g := range u.StaticCollectors {
		g.Collector.Describe(c)
	}

	for _, child := range u.children {
		child.Describe(c)
	}

	for _, m := range u.Maps {
		for _, collectors := range m.Mapped {
			collectors.Describe(c)
		}
	}
}

func (u *Collectors) Collect(c chan<- prometheus.Metric) {
	for _, g := range u.StaticCollectors {
		g.Collector.Collect(c)
	}

	for _, child := range u.children {
		child.Collect(c)
	}

	for _, m := range u.Maps {
		for _, collectors := range m.Mapped {
			collectors.Collect(c)
		}
	}
}

func (u *Collectors) Fill(t reflect.Type, rlr *label.RecursiveReflector, parent string, metricNameTransform types.MetricNameTransformer) {
	u.T = t
	u.Rlr = rlr
	if u.T != u.Rlr.T {
		panic(fmt.Sprintf("LabelReflector type `%s` does not match collected type `%s`", u.Rlr.T, u.T))
	}

	fields := reflect.VisibleFields(t)
	for i, f := range fields {
		tag := f.Tag.Get("kpromcol")
		if tag != "" {
			u.StaticCollectors = append(u.StaticCollectors, *makeGenerated(i, tag, f, parent, u.Rlr.Ln, metricNameTransform))
			continue
		}
		tag = f.Tag.Get("kprommap")
		if tag != "" {
			switch f.Type.Kind() {
			case reflect.Map:
			default:
				panic("Only supported on maps")
			}
			u.Maps = append(u.Maps, DynamicMap{
				IndexInStruct: i,
				StructParent:  parent,
				Mapped:        map[reflect.Value]*Collectors{},
				FieldName:     tag,
			})
			continue
		}
		tag = f.Tag.Get("kprompnt")
		if tag != "" {
			cu := &Collectors{}
			if parent == "" {
				cu.Fill(f.Type, rlr.Fields[i], tag, metricNameTransform)
			} else {
				cu.Fill(f.Type, rlr.Fields[i], parent+"_"+tag, metricNameTransform)
			}
			u.children = append(u.children, cu)
			continue
		}
	}
}
