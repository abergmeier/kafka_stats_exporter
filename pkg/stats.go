package pkg

import (
	"encoding/json"
	"reflect"

	"github.com/abergmeier/kafka_stats_exporter/pkg/kafka/typed"
	"github.com/abergmeier/kafka_stats_exporter/pkg/prometheus/auto"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Exporter struct {
	stats         typed.Stats
	cachedUpdater map[reflect.Type]interface{}
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
		upd = auto.BuildFromTags(e.stats)
		e.cachedUpdater[t] = upd
	}

	// TODO: Call update

	return nil
}
