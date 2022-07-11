package gen

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/abergmeier/kafka_stats_exporter/v0/pkg/kafka/typed"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var (
	simple                       = simpleStats{}
	full                         = typed.Stats{}
	expectedFull, expectedSimple io.Reader
)

type simpleStats struct {
	Name    string                                 `json:"name"     kpromlbl:"name"` //Handle instance name
	RxBytes int                                    `json:"rx_bytes" kpromcol:"CounterVec,Total number of bytes received from Kafka brokers"`
	Brokers map[typed.BrokerName]simpleBrokerStats `json:"brokers"  kprommap:"brokers"`
}

type simpleBrokerStats struct {
	Name    string `json:"name"    kpromlbl:"name"` //Broker hostname, port and broker id
	Rxbytes int    `json:"rxbytes" kpromcol:"CounterVec,Total number of bytes received"`
}

func init_full() {
	f, err := os.Open("testdata/full.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&full)
	if err != nil {
		panic(err)
	}

	expectedFull, err = os.Open("testdata/full_expected.txt")
	if err != nil {
		panic(err)
	}
}

func init_simple() {
	f, err := os.Open("testdata/simple.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&simple)
	if err != nil {
		panic(err)
	}

	expectedSimple, err = os.Open("testdata/simple_expected.txt")
	if err != nil {
		panic(err)
	}
}

func init() {
	init_simple()
	init_full()
}

func TestNewRecursiveUpdaterFromTags(t *testing.T) {
	_, _ = NewRecursiveUpdaterFromTags(typed.Stats{})
	_, _ = NewRecursiveUpdaterFromTags(&typed.Stats{})
}

func TestSimple(t *testing.T) {
	col, upd := NewRecursiveUpdaterFromTags(simpleStats{})
	upd.Update(&simple, prometheus.Labels{})
	err := testutil.CollectAndCompare(col, expectedSimple)
	if err != nil {
		t.Fatal("CollectAndCompare failed:", err)
	}
}

func TestUpdate(t *testing.T) {
	col, upd := NewRecursiveUpdaterFromTags(&full)
	upd.Update(full, prometheus.Labels{})
	err := testutil.CollectAndCompare(col, expectedFull)
	if err != nil {
		t.Fatal("CollectAndCompare failed:", err)
	}
}
