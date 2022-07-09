package gen

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/abergmeier/kafka_stats_exporter/internal/assert"
	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
)

type dynamicMap struct {
	indexInStruct int
	structParent  string
	fieldName     string // Field names get saved in snake cased prometheus format

	mapped map[reflect.Value]*UpdateCollectors
}

type generatedUpdator struct {
	collector prometheus.Collector
	index     int
	update    func(last, current int64, ls prometheus.Labels) int64
	last      int64
}

type UpdateCollectors struct {
	rlr              *recursiveLabelReflector
	staticCollectors []generatedUpdator
	children         []UpdatingCollector
	maps             []dynamicMap
	t                reflect.Type
}

func (u *UpdateCollectors) Update(v interface{}, labels prometheus.Labels) {
	for k, v := range u.rlr.Lr.LabelsForValue(v) {
		labels[k] = v
	}

	u.update(v, labels)
}

func (u *UpdateCollectors) update(v interface{}, labels prometheus.Labels) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		rv = rv.Elem()
	}
	assert.AssertType(rv, u.t)
	for _, c := range u.staticCollectors {
		fv := rv.FieldByIndex([]int{c.index})
		if !fv.CanInt() {
			panic("Only in update implemented yet")
		}
		current := fv.Int()
		c.update(c.last, current, labels)
		c.last = current
	}
	// Up until here we could do statically initialize
	// all data. Here map keys can change while runtime
	// thus we need to handle
	for _, m := range u.maps {
		fv := rv.FieldByIndex([]int{m.indexInStruct})
		assertMap(fv)
		ln := LabelNames{}
		for k := range labels {
			ln = append(ln, k)
		}
		updateMapped(&m, u.rlr.Fields[m.indexInStruct], fv, ln)
		for mk, mc := range m.mapped {
			// While we only support strings and ints we usually use an alias to
			// make semantics clear. Thus we need to convert from string or int
			// to alias here.
			aliased := mk.Convert(fv.Type().Key())
			mv := fv.MapIndex(aliased)
			mc.rlr = u.rlr.Fields[m.indexInStruct]
			mc.Update(mv.Interface(), labels)
		}
	}
}

func updateMapped(d *dynamicMap, rlr *recursiveLabelReflector, fv reflect.Value, labelNames LabelNames) {

	var keysToDelete map[reflect.Value]struct{}
	for k := range d.mapped {
		keysToDelete[k] = struct{}{}
	}

	iter := fv.MapRange()
	for iter.Next() {
		key := iter.Key()
		_, ok := d.mapped[key]
		if ok {
			// We can reuse map entry
			delete(keysToDelete, key)
			continue
		}
		vt := iter.Value().Type()
		cu := &UpdateCollectors{}
		if d.structParent == "" {
			cu.fill(vt, rlr, d.fieldName+"_"+key.String()+"_", labelNames)
		} else {
			cu.fill(vt, rlr, d.structParent+"_"+d.fieldName+"_"+key.String()+"_", labelNames)
		}
		// We need new map entry
		d.mapped[key] = cu
	}
	// Delete superfluous keys
	for k := range keysToDelete {
		delete(d.mapped, k)
	}
}

func (u *UpdateCollectors) Describe(c chan<- *prometheus.Desc) {
	for _, g := range u.staticCollectors {
		g.collector.Describe(c)
	}

	for _, child := range u.children {
		child.Describe(c)
	}

	for _, m := range u.maps {
		for _, collectors := range m.mapped {
			collectors.Describe(c)
		}
	}
}

func (u *UpdateCollectors) Collect(c chan<- prometheus.Metric) {
	for _, g := range u.staticCollectors {
		g.collector.Collect(c)
	}

	for _, child := range u.children {
		child.Collect(c)
	}

	for _, m := range u.maps {
		for _, collectors := range m.mapped {
			collectors.Collect(c)
		}
	}
}

func NewRecursiveUpdaterFromTags(tagged interface{}) UpdatingCollector {
	t := reflect.TypeOf(tagged)
	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
	}

	rlr := recursiveLabelReflector{}
	fillLabels(t, &rlr, "", nil)

	u := &UpdateCollectors{}
	u.fill(t, &rlr, "", LabelNames{})
	return u
}

func (u *UpdateCollectors) fill(t reflect.Type, rlr *recursiveLabelReflector, parent string, labelNames LabelNames) {
	u.t = t
	u.rlr = rlr
	if u.t != u.rlr.T {
		panic(fmt.Sprintf("LabelReflector type `%s` does not match collected type `%s`", u.rlr.T, u.t))
	}

	fields := reflect.VisibleFields(t)
	for i, f := range fields {
		tag := f.Tag.Get("kpromcol")
		if tag != "" {
			ln := labelNames
			ln = append(ln, rlr.Ln...)
			u.staticCollectors = append(u.staticCollectors, *makeGenerated(i, tag, f, parent, ln))
			continue
		}
		tag = f.Tag.Get("kprommap")
		if tag != "" {
			switch f.Type.Kind() {
			case reflect.Map:
			default:
				panic("Only supported on maps")
			}
			u.maps = append(u.maps, dynamicMap{
				indexInStruct: i,
				structParent:  parent,
				mapped:        map[reflect.Value]*UpdateCollectors{},
				fieldName:     tag,
			})
			continue
		}
		tag = f.Tag.Get("kprompnt")
		if tag != "" {
			cu := &UpdateCollectors{}
			if parent == "" {
				cu.fill(f.Type, rlr.Fields[i], tag+"_", labelNames)
			} else {
				cu.fill(f.Type, rlr.Fields[i], parent+"_"+tag+"_", labelNames)
			}
			u.children = append(u.children, cu)
			continue
		}
	}
}

func makeGenerated(i int, tag string, f reflect.StructField, parent string, labelNames LabelNames) *generatedUpdator {
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
		return &generatedUpdator{
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
		return &generatedUpdator{
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
