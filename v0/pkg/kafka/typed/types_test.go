package typed

import (
	"encoding/json"
	"os"
	"testing"
)

func TestStats(t *testing.T) {
	f, err := os.Open("testdata/example.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	stats := Stats{}
	err = dec.Decode(&stats)
	if err != nil {
		t.Fatal("Decode failed:", err)
	}
	if stats.Type != "producer" {
		t.Fatalf("Type `%s` is wrong. Expected: producer", stats.Type)
	}
}
