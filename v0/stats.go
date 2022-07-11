package pkg

import (
	"encoding/json"
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/gen"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	stats         typed.Stats
	cachedUpdater map[reflect.Type]struct {
		collector prometheus.Collector
		updater   gen.Updater
	}
}

func (e *Exporter) UpdateWithStats(stats *kafka.Stats) error {
	s := stats.String()
	return e.UpdateWithStatString(s)
}

func (e *Exporter) UpdateWithStatString(stats string) error {
	err := json.Unmarshal([]byte(stats), &e.stats)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(e.stats)
	ce, ok := e.cachedUpdater[t]
	if !ok {
		ce.collector, ce.updater = gen.NewRecursiveMetricsFromTags(e.stats)
		e.cachedUpdater[t] = ce
	}

	ce.updater.Update(e.stats, prometheus.Labels{})
	return nil
}
