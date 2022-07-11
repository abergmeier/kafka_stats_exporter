package collector

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/abergmeier/kafka_stats_exporter/internal/label"
	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
)

type DynamicMap struct {
	IndexInStruct int
	StructParent  string
	FieldName     string // Field names get saved in snake cased prometheus format

	Mapped map[reflect.Value]*Collectors
}

type GeneratedUpdator struct {
	Collector prometheus.Collector
	Index     int
	Update    func(last, current int64, ls prometheus.Labels) int64
	Last      int64
}

func makeGenerated(i int, tag string, f reflect.StructField, parent string, labelNames label.Names) *GeneratedUpdator {
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
		return &GeneratedUpdator{
			Collector: counterVec,
			Index:     i,
			Update: func(last, current int64, ls prometheus.Labels) int64 {
				return updateCounterVec(last, current, counterVec, ls)
			},
		}
	case "GaugeVec":
		gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: parent + strcase.ToSnake(f.Name),
			Help: help,
		}, labelNames)
		return &GeneratedUpdator{
			Collector: gaugeVec,
			Index:     i,
			Update: func(last, current int64, ls prometheus.Labels) int64 {
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
