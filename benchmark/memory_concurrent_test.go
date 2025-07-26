package benchmark

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
	raphamorim "github.com/raphamorim/fuzzy"
	sahilm "github.com/sahilm/fuzzy"
)

// Memory benchmarks
func BenchmarkMemory_Sahilm_Large(b *testing.B) {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	before := m.Alloc
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sahilm.Find("test", largeDataset)
	}
	
	runtime.GC()
	runtime.ReadMemStats(&m)
	after := m.Alloc
	
	b.ReportMetric(float64(after-before)/float64(b.N), "bytes/op")
}

func BenchmarkMemory_Lithammer_Large(b *testing.B) {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	before := m.Alloc
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fuzzy.Find("test", largeDataset)
	}
	
	runtime.GC()
	runtime.ReadMemStats(&m)
	after := m.Alloc
	
	b.ReportMetric(float64(after-before)/float64(b.N), "bytes/op")
}

func BenchmarkMemory_Raphamorim_BKTree_Build(b *testing.B) {
	var m runtime.MemStats
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runtime.GC()
		runtime.ReadMemStats(&m)
		before := m.Alloc
		
		tree := raphamorim.NewBKTree()
		for _, word := range mediumDataset[:1000] {
			tree.Add(word)
		}
		
		runtime.GC()
		runtime.ReadMemStats(&m)
		after := m.Alloc
		
		b.ReportMetric(float64(after-before), "bytes/tree")
	}
}

// Concurrent benchmarks
func BenchmarkConcurrent_Sahilm_Medium(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			query := exactQueries[i%len(exactQueries)]
			sahilm.Find(query, mediumDataset)
			i++
		}
	})
}

func BenchmarkConcurrent_Lithammer_Medium(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			query := exactQueries[i%len(exactQueries)]
			fuzzy.Find(query, mediumDataset)
			i++
		}
	})
}

func BenchmarkConcurrent_Raphamorim_BKTree(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			query := exactQueries[i%len(exactQueries)]
			raphamorimBKTree.Search(query, 2)
			i++
		}
	})
}

// Latency percentile benchmarks
func measureLatencies(name string, searchFunc func(string), queries []string) {
	latencies := make([]time.Duration, len(queries))
	
	for i, query := range queries {
		start := time.Now()
		searchFunc(query)
		latencies[i] = time.Since(start)
	}
	
	// Sort latencies
	for i := 0; i < len(latencies); i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}
	
	p50 := latencies[len(latencies)*50/100]
	p90 := latencies[len(latencies)*90/100]
	p99 := latencies[len(latencies)*99/100]
	
	fmt.Printf("%s - P50: %v, P90: %v, P99: %v\n", name, p50, p90, p99)
}

func TestLatencyPercentiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping latency test in short mode")
	}
	
	fmt.Println("\nLatency Percentiles (Medium Dataset):")
	
	measureLatencies("sahilm/fuzzy", func(q string) {
		sahilm.Find(q, mediumDataset)
	}, exactQueries)
	
	measureLatencies("lithammer/fuzzysearch", func(q string) {
		fuzzy.Find(q, mediumDataset)
	}, exactQueries)
	
	measureLatencies("raphamorim/fuzzy (BK-Tree)", func(q string) {
		raphamorimBKTree.Search(q, 2)
	}, exactQueries)
	
	measureLatencies("raphamorim/fuzzy (NGram)", func(q string) {
		raphamorimNGram.Search(q, 0.5)
	}, exactQueries)
}