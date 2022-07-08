package pkg

import (
	"encoding/json"
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/prometheus/gen"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Exporter struct {
	stats         typed.Stats
	cachedUpdater map[reflect.Type]gen.UpdatingCollector
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
	upd, ok := e.cachedUpdater[t]
	if !ok {
		upd = gen.NewRecursiveUpdaterFromTags(e.stats)
		e.cachedUpdater[t] = upd
	}

	upd.Update(e.stats)
	return nil
}
