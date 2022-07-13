package v0

import (
	"encoding/json"
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/gen"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter interface {
	UpdateWithStats(stats *kafka.Stats) error
	UpdateWithStatString(stats string) error
}

type exporter struct {
	registerer    prometheus.Registerer
	stats         typed.Stats
	cachedUpdater map[reflect.Type]struct {
		collector prometheus.Collector
		updater   gen.Updater
	}
}

func NewExporter(r prometheus.Registerer) *exporter {
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
