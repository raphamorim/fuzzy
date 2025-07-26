package benchmark

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/lithammer/fuzzysearch/fuzzy"
	raphamorim "github.com/raphamorim/fuzzy"
	sahilm "github.com/sahilm/fuzzy"
)

func TestMain(m *testing.M) {
	fmt.Println("Fuzzy Search Library Benchmark Suite")
	fmt.Println("====================================")
	fmt.Printf("Dataset sizes: Small=%d, Medium=%d, Large=%d\n", 
		len(smallDataset), len(mediumDataset), len(largeDataset))
	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Println()
	
	code := m.Run()
	
	// Print summary
	printSummary()
	
	os.Exit(code)
}

func printSummary() {
	fmt.Println("\nBenchmark Summary")
	fmt.Println("=================")
	fmt.Println("\nKey Findings:")
	fmt.Println("1. sahilm/fuzzy: Optimized for editor-style filename matching")
	fmt.Println("2. lithammer/fuzzysearch: Good balance with Levenshtein distance")
	fmt.Println("3. raphamorim/fuzzy: Advanced algorithms for different use cases")
	fmt.Println("   - BK-Tree: Efficient for spell checking with edit distance")
	fmt.Println("   - N-gram/Trigram: Fast for partial matches")
	fmt.Println("   - LSH: Excellent for document similarity")
	fmt.Println("   - Suffix Array: Powerful for substring search")
	fmt.Println("\nRecommendations:")
	fmt.Println("- For simple fuzzy matching: lithammer/fuzzysearch")
	fmt.Println("- For filename/path matching: sahilm/fuzzy")
	fmt.Println("- For advanced use cases: raphamorim/fuzzy")
}

// Benchmark different edit distances
func BenchmarkEditDistance_Comparison(b *testing.B) {
	s1 := "The quick brown fox jumps over the lazy dog"
	s2 := "The quick brown fox jumped over the lazy dogs"
	
	b.Run("raphamorim_levenshtein", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			raphamorim.LevenshteinDistance(s1, s2)
		}
	})
	
	b.Run("raphamorim_damerau", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			raphamorim.DamerauLevenshteinDistance(s1, s2)
		}
	})
	
	b.Run("raphamorim_myers", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			raphamorim.MyersDistance(s1, s2)
		}
	})
}

// Benchmark advanced algorithms
func BenchmarkAdvanced_WuManber(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog. The fox is quick and brown."
	pattern := "quick"
	
	wm := raphamorim.NewWuManber(pattern)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wm.Search(text, 1)
	}
}

func BenchmarkAdvanced_SimHash(b *testing.B) {
	sh := raphamorim.NewSimHash(64)
	
	// Add documents
	for i := 0; i < 1000; i++ {
		sh.Add(mediumDataset[i%len(mediumDataset)])
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sh.Query("The quick brown fox", 5)
	}
}

// Scaling benchmark
func BenchmarkScaling(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}
	
	for _, size := range sizes {
		var dataset []string
		if size <= 10000 {
			dataset = mediumDataset[:min(size, len(mediumDataset))]
		} else {
			dataset = largeDataset[:min(size, len(largeDataset))]
		}
		
		b.Run(fmt.Sprintf("sahilm_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sahilm.Find("test", dataset)
			}
		})
		
		b.Run(fmt.Sprintf("lithammer_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fuzzy.Find("test", dataset)
			}
		})
		
		b.Run(fmt.Sprintf("raphamorim_bktree_%d", size), func(b *testing.B) {
			tree := raphamorim.NewBKTree()
			for _, word := range dataset {
				tree.Add(word)
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tree.Search("test", 2)
			}
		})
		
		b.Run(fmt.Sprintf("raphamorim_ngram_%d", size), func(b *testing.B) {
			ngram := raphamorim.NewNGram(3)
			for i, word := range dataset {
				ngram.Add(word, i)
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ngram.Search("test", 0.5)
			}
		})
	}
}

// Benchmark with different query lengths
func BenchmarkQueryLength(b *testing.B) {
	queries := map[string]string{
		"short":  "cat",
		"medium": "elephant",
		"long":   "supercalifragilisticexpialidocious",
	}
	
	for name, query := range queries {
		b.Run(fmt.Sprintf("sahilm_%s", name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sahilm.Find(query, mediumDataset)
			}
		})
		
		b.Run(fmt.Sprintf("lithammer_%s", name), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fuzzy.Find(query, mediumDataset)
			}
		})
	}
}