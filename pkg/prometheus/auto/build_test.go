package auto

import (
	"testing"

	"github.com/abergmeier/kafka_stats_exporter/pkg/kafka/typed"
)

func TestBuildFromTags(t *testing.T) {
	s := typed.Stats{}
	BuildFromTags(&s)
}
