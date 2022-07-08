package gen

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var (
	stats             = typed.Stats{}
	expectedCollected io.Reader
)

func init() {
	f, err := os.Open("testdata/example.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&stats)
	if err != nil {
		panic(err)
	}

	expectedCollected, err = os.Open("testdata/updated.txt")
	if err != nil {
		panic(err)
	}
}

func TestNewRecursiveUpdaterFromTags(t *testing.T) {
	_ = NewRecursiveUpdaterFromTags(typed.Stats{})
	_ = NewRecursiveUpdaterFromTags(&typed.Stats{})
}

func TestUpdate(t *testing.T) {
	upd := NewRecursiveUpdaterFromTags(&stats)
	upd.Update(stats)
	err := testutil.CollectAndCompare(upd, expectedCollected)
	if err != nil {
		t.Fatal("CollectAndCompare failed:", err)
	}
}
