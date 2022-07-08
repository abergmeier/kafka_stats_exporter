package auto

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type generated struct {
	collector prometheus.Collector
	update    func(last int, current int) int
	last      int
}

func BuildFromTags(tagged interface{}) interface{} {
	t := reflect.TypeOf(tagged)
	switch t.Kind() {
	case reflect.Pointer:
		t = t.Elem()
	}

	var collectors []generated

	fields := reflect.VisibleFields(t)
	for _, f := range fields {
		tag := f.Tag.Get("kpromauto")
		if tag == "" {
			continue
		}

		prom := strings.SplitN(tag, ",", 2)
		switch prom[0] {
		case "CounterVec":
			counterVec := promauto.NewCounterVec(prometheus.CounterOpts{
				Name: strcase.ToSnake(f.Name),
			}, nil)
			collectors = append(collectors, generated{
				collector: counterVec,
				update: func(last int, current int) int {
					return updateCounterVec(last, current, counterVec)
				},
			})
		case "GaugeVec":
			gaugeVec := promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name: strcase.ToSnake(f.Name),
			}, nil)
			collectors = append(collectors, generated{
				collector: gaugeVec,
				update: func(last int, current int) int {
					updateGaugeVec(current, gaugeVec)
					return current
				},
			})
		case "":
		default:
			panic(fmt.Sprintf("Unsupported prometheus Metric: %s", prom[0]))
		}
	}

	return collectors
}

func calculateLabels() prometheus.Labels {
	return prometheus.Labels{}
}

func updateCounterVec(last, current int, counter *prometheus.CounterVec) int {
	return updateCounter(last, current, counter.With(calculateLabels()))
}

func updateCounter(last, current int, counter prometheus.Counter) int {
	diff := current - last
	if diff < 0 {
		return last
	}
	counter.Add(float64(diff))
	return current
}

func updateGaugeVec(current int, gauge *prometheus.GaugeVec) {
	updateGauge(current, gauge.With(calculateLabels()))
}

func updateGauge(current int, gauge prometheus.Gauge) {
	gauge.Set(float64(current))
}
