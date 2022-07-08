package gen

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
)

type generated struct {
	collector prometheus.Collector
	index     int
	update    func(last, current int64, ls prometheus.Labels) int64
	last      int64
}

type UpdateCollectors struct {
	lg         LabelReflector
	collectors []generated
	children   []UpdatingCollector
}

func (u *UpdateCollectors) Update(v interface{}) {
	labels := u.lg(v)
	rv := reflect.ValueOf(v)
	for _, c := range u.collectors {
		fv := rv.FieldByIndex([]int{c.index})
		if !fv.CanInt() {
			fmt.Printf("Field %s %d\n", fv, c.index)
			panic("Only in update implemented yet")
		}
		current := fv.Int()
		c.update(c.last, current, labels)
		c.last = current
	}
}

func (u *UpdateCollectors) Describe(c chan<- *prometheus.Desc) {
	for _, g := range u.collectors {
		g.collector.Describe(c)
	}

	for _, child := range u.children {
		child.Describe(c)
	}
}

func (u *UpdateCollectors) Collect(c chan<- prometheus.Metric) {
	for _, g := range u.collectors {
		g.collector.Collect(c)
	}

	for _, child := range u.children {
		child.Collect(c)
	}
}

func NewRecursiveUpdaterFromTags(tagged interface{}) UpdatingCollector {
	t := reflect.TypeOf(tagged)
	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
	}

	u := &UpdateCollectors{}
	u.fill(t, "")
	return u
}

func (u *UpdateCollectors) fill(t reflect.Type, parent string) {
	var labelNames LabelNames
	u.lg, labelNames = MakeLabelReflector(t, parent)

	fields := reflect.VisibleFields(t)
	for i, f := range fields {
		tag := f.Tag.Get("kpromcol")
		if tag != "" {
			u.collectors = append(u.collectors, *makeGenerated(i, tag, f, parent, labelNames))
			continue
		}
		tag = f.Tag.Get("kprompnt")
		if tag != "" {
			cu := &UpdateCollectors{}
			if parent == "" {
				cu.fill(f.Type, tag+"_")
			} else {
				cu.fill(f.Type, parent+"_"+tag+"_")
			}
			u.children = append(u.children, cu)
			continue
		}
	}
}

func makeGenerated(i int, tag string, f reflect.StructField, parent string, labelNames LabelNames) *generated {
	prom := strings.SplitN(tag, ",", 2)
	help, err := url.QueryUnescape(prom[1])
	if err != nil {
		panic(err)
	}
	switch prom[0] {
	case "CounterVec":
		counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: parent + strcase.ToSnake(f.Name) + "_total",
			Help: help,
		}, labelNames)
		return &generated{
			collector: counterVec,
			index:     i,
			update: func(last, current int64, ls prometheus.Labels) int64 {
				return updateCounterVec(last, current, counterVec, ls)
			},
		}
	case "GaugeVec":
		gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: parent + strcase.ToSnake(f.Name),
			Help: help,
		}, labelNames)
		return &generated{
			collector: gaugeVec,
			index:     i,
			update: func(last, current int64, ls prometheus.Labels) int64 {
				return updateGaugeVec(current, gaugeVec, ls)
			},
		}
	case "":
		return nil
	default:
		panic(fmt.Sprintf("Unsupported prometheus Metric: %s", prom[0]))
	}
}

func updateCounterVec(last, current int64, counter *prometheus.CounterVec, ls prometheus.Labels) int64 {
	return updateCounter(last, current, counter.With(ls))
}

func updateCounter(last, current int64, counter prometheus.Counter) int64 {
	diff := current - last
	if diff < 0 {
		return last
	}
	counter.Add(float64(diff))
	return current
}

func updateGaugeVec(current int64, gauge *prometheus.GaugeVec, ls prometheus.Labels) int64 {
	return updateGauge(current, gauge.With(ls))
}

func updateGauge(current int64, gauge prometheus.Gauge) int64 {
	gauge.Set(float64(current))
	return current
}
