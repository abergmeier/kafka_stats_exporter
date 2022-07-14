package v0

import (
	"encoding/json"
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/gen"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter interface {
	UpdateWithStatString(stats string) error
	// Returns the last updated stats
	Stats() *typed.Stats
}

type exporter struct {
	registerer    prometheus.Registerer
	stats         typed.Stats
	cachedUpdater map[reflect.Type]struct {
		collector prometheus.Collector
		updater   gen.Updater
	}
}

func NewExporter(r prometheus.Registerer) Exporter {
	return &exporter{
		registerer: r,
		cachedUpdater: map[reflect.Type]struct {
			collector prometheus.Collector
			updater   gen.Updater
		}{},
	}
}

func (e *exporter) UpdateWithStatString(stats string) error {
	err := json.Unmarshal([]byte(stats), &e.stats)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(e.stats)
	ce, ok := e.cachedUpdater[t]
	if !ok {
		ce.collector, ce.updater = gen.NewRecursiveMetricsFromTags(e.stats)
		err := e.registerer.Register(ce.collector)
		if err != nil {
			return err
		}
		e.cachedUpdater[t] = ce
	}

	ce.updater.Update(e.stats, prometheus.Labels{})
	return nil
}

func (e *exporter) Stats() *typed.Stats {
	return &e.stats
}
