package benchmark

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lithammer/fuzzysearch/fuzzy"
	raphamorim "github.com/raphamorim/fuzzy"
	sahilm "github.com/sahilm/fuzzy"
)

// Accuracy tests
func TestAccuracy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping accuracy test in short mode")
	}
	
	testCases := []struct {
		name     string
		query    string
		dataset  []string
		expected []string
	}{
		{
			name:     "Exact match",
			query:    "hello",
			dataset:  []string{"hello", "hallo", "hullo", "help", "hell"},
			expected: []string{"hello"},
		},
		{
			name:     "One character difference",
			query:    "helo",
			dataset:  []string{"hello", "halo", "help", "hero", "helm"},
			expected: []string{"hello", "halo", "help", "hero"},
		},
		{
			name:     "Transposition",
			query:    "hlelo",
			dataset:  []string{"hello", "hallo", "hullo", "help", "hell"},
			expected: []string{"hello"},
		},
		{
			name:     "Prefix match",
			query:    "hel",
			dataset:  []string{"hello", "help", "helmet", "world", "held"},
			expected: []string{"hello", "help", "helmet", "held"},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test sahilm/fuzzy
			sahilmResults := sahilm.Find(tc.query, tc.dataset)
			sahilmMatches := make([]string, len(sahilmResults))
			for i, r := range sahilmResults {
				sahilmMatches[i] = tc.dataset[r.Index]
			}
			
			// Test lithammer/fuzzysearch
			lithammerMatches := fuzzy.Find(tc.query, tc.dataset)
			
			// Test raphamorim/fuzzy BK-Tree
			tree := raphamorim.NewBKTree()
			for _, word := range tc.dataset {
				tree.Add(word)
			}
			raphamorimMatches := tree.Search(tc.query, 2)
			
			fmt.Printf("\n%s:\n", tc.name)
			fmt.Printf("  Query: %q\n", tc.query)
			fmt.Printf("  Expected: %v\n", tc.expected)
			fmt.Printf("  sahilm: %v\n", sahilmMatches)
			fmt.Printf("  lithammer: %v\n", lithammerMatches)
			fmt.Printf("  raphamorim: %v\n", raphamorimMatches)
		})
	}
}

// Real-world use case benchmarks
func BenchmarkRealWorld_FilenameFuzzySearch(b *testing.B) {
	// Simulate searching through file paths
	filePaths := make([]string, 10000)
	for i := range filePaths {
		depth := 3 + i%5
		parts := make([]string, depth)
		for j := 0; j < depth; j++ {
			parts[j] = generateRandomString()
		}
		filePaths[i] = strings.Join(parts, "/") + ".txt"
	}
	
	queries := []string{
		"config",
		"test",
		"main",
		"index",
		"util",
	}
	
	b.Run("sahilm", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, q := range queries {
				sahilm.Find(q, filePaths)
			}
		}
	})
	
	b.Run("lithammer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, q := range queries {
				fuzzy.Find(q, filePaths)
			}
		}
	})
	
	b.Run("raphamorim_ngram", func(b *testing.B) {
		ng := raphamorim.NewNGram(3)
		for i, path := range filePaths {
			ng.Add(path, i)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, q := range queries {
				ng.Search(q, 0.3)
			}
		}
	})
}

func BenchmarkRealWorld_AutoComplete(b *testing.B) {
	// Simulate autocomplete scenarios
	words := []string{
		"javascript", "java", "python", "typescript", "golang",
		"react", "angular", "vue", "svelte", "ember",
		"docker", "kubernetes", "terraform", "ansible", "jenkins",
		"postgresql", "mysql", "mongodb", "redis", "elasticsearch",
	}
	
	// Expand with variations
	expandedWords := make([]string, 0, len(words)*100)
	for _, word := range words {
		for i := 0; i < 100; i++ {
			expandedWords = append(expandedWords, generateVariation(word))
		}
	}
	
	prefixes := []string{"jav", "pyth", "dock", "kub", "post"}
	
	b.Run("sahilm", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, prefix := range prefixes {
				sahilm.Find(prefix, expandedWords)
			}
		}
	})
	
	b.Run("lithammer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, prefix := range prefixes {
				fuzzy.Find(prefix, expandedWords)
			}
		}
	})
	
	b.Run("raphamorim_trigram", func(b *testing.B) {
		ti := raphamorim.NewTrigramIndex()
		for i, word := range expandedWords {
			ti.Add(word, i)
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, prefix := range prefixes {
				ti.Search(prefix, 2)
			}
		}
	})
}