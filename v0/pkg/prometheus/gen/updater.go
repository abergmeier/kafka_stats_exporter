package gen

import (
	"reflect"
	"strconv"

	"github.com/abergmeier/kafka_stats_exporter/internal/assert"
	"github.com/abergmeier/kafka_stats_exporter/internal/collector"
	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/prometheus/client_golang/prometheus"
)

//Updater allows for updating prometheus Metrics from Value
type Updater interface {
	Update(v interface{}, labels prometheus.Labels)
}

type updater struct {
	c *collector.Collectors
}

func (u *updater) Update(v interface{}, labels prometheus.Labels) {
	newLabels := make(prometheus.Labels, len(labels))
	for k, v := range labels {
		newLabels[k] = v
	}
	for k, v := range u.c.Rlr.Lr.LabelsForValue(v) {
		newLabels[k] = v
	}

	u.update(v, newLabels)
}

func (u *updater) update(v interface{}, labels prometheus.Labels) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		rv = rv.Elem()
	}
	assert.AssertType(rv, u.c.T)
	for _, c := range u.c.StaticCollectors {
		fv := rv.FieldByIndex([]int{c.Index})
		if !fv.CanInt() {
			panic("Only in update implemented yet")
		}
		current := fv.Int()
		c.Update(c.Last, current, labels)
		c.Last = current
	}
	// Up until here we could do statically initialize
	// all data. Here map keys can change while runtime
	// thus we need to handle
	for _, m := range u.c.Maps {
		fv := rv.FieldByIndex([]int{m.IndexInStruct})
		assert.AssertMap(fv)
		updateMapped(&m, u.c.Rlr.Fields[m.IndexInStruct], fv)
		for mk, mc := range m.Mapped {
			// While we only support strings and ints we usually use an alias to
			// make semantics clear. Thus we need to convert from string or int
			// to alias here.
			aliased := mk.Convert(fv.Type().Key())
			mv := fv.MapIndex(aliased)
			(&updater{
				c: mc,
			}).Update(mv.Interface(), labels)
		}
	}
}

func updateMapped(d *collector.DynamicMap, rlr *label.RecursiveReflector, fv reflect.Value) {

	keysToDelete := make(map[reflect.Value]struct{}, len(d.Mapped))
	for k := range d.Mapped {
		keysToDelete[k] = struct{}{}
	}

	iter := fv.MapRange()
	for iter.Next() {
		rk := iter.Key()
		_, ok := d.Mapped[rk]
		if ok {
			// We can reuse map entry
			delete(keysToDelete, rk)
			continue
		}
		var keyString string
		if rk.CanInt() {
			keyString = strconv.FormatInt(rk.Int(), 10)
		} else {
			keyString = rk.String()
		}
		vt := iter.Value().Type()
		cu := &collector.Collectors{}
		if d.StructParent == "" {
			cu.Fill(vt, rlr, d.FieldName+"_"+keyString)
		} else {
			cu.Fill(vt, rlr, d.StructParent+"_"+d.FieldName+"_"+keyString)
		}
		// We need new map entry
		d.Mapped[rk] = cu
	}
	// Delete superfluous keys
	for k := range keysToDelete {
		delete(d.Mapped, k)
	}
}
