package benchmark

import (
	"testing"

	"github.com/lithammer/fuzzysearch/fuzzy"
	raphamorim "github.com/raphamorim/fuzzy"
	sahilm "github.com/sahilm/fuzzy"
)

// Benchmark exact matching
func BenchmarkExactMatch_Sahilm_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			sahilm.Find(query, smallDataset)
		}
	}
}

func BenchmarkExactMatch_Sahilm_Medium(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			sahilm.Find(query, mediumDataset)
		}
	}
}

func BenchmarkExactMatch_Sahilm_Large(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			sahilm.Find(query, largeDataset)
		}
	}
}

func BenchmarkExactMatch_Lithammer_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			fuzzy.Find(query, smallDataset)
		}
	}
}

func BenchmarkExactMatch_Lithammer_Medium(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			fuzzy.Find(query, mediumDataset)
		}
	}
}

func BenchmarkExactMatch_Lithammer_Large(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			fuzzy.Find(query, largeDataset)
		}
	}
}

func BenchmarkExactMatch_Raphamorim_BKTree(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			raphamorimBKTree.Search(query, 0)
		}
	}
}

func BenchmarkExactMatch_Raphamorim_NGram(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range exactQueries[:10] {
			raphamorimNGram.Search(query, 0.9)
		}
	}
}

// Benchmark fuzzy matching with typos
func BenchmarkTypoMatch_Sahilm_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range typoQueries[:10] {
			sahilm.Find(query, smallDataset)
		}
	}
}

func BenchmarkTypoMatch_Sahilm_Medium(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range typoQueries[:10] {
			sahilm.Find(query, mediumDataset)
		}
	}
}

func BenchmarkTypoMatch_Lithammer_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range typoQueries[:10] {
			fuzzy.Find(query, smallDataset)
		}
	}
}

func BenchmarkTypoMatch_Lithammer_Medium(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range typoQueries[:10] {
			fuzzy.Find(query, mediumDataset)
		}
	}
}

func BenchmarkTypoMatch_Raphamorim_BKTree(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range typoQueries[:10] {
			raphamorimBKTree.Search(query, 2)
		}
	}
}

func BenchmarkTypoMatch_Raphamorim_Trigram(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range typoQueries[:10] {
			raphamorimTrigram.Search(query, 2)
		}
	}
}

// Benchmark ranking/scoring
func BenchmarkRanking_Lithammer_Medium(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range partialQueries[:10] {
			fuzzy.RankFind(query, mediumDataset)
		}
	}
}

func BenchmarkRanking_Lithammer_Large(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, query := range partialQueries[:5] {
			fuzzy.RankFind(query, largeDataset)
		}
	}
}

// Benchmark index building
func BenchmarkIndexBuild_Raphamorim_BKTree_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := raphamorim.NewBKTree()
		for _, word := range smallDataset {
			tree.Add(word)
		}
	}
}

func BenchmarkIndexBuild_Raphamorim_BKTree_Medium(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree := raphamorim.NewBKTree()
		for _, word := range mediumDataset[:1000] {
			tree.Add(word)
		}
	}
}

func BenchmarkIndexBuild_Raphamorim_NGram_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ng := raphamorim.NewNGram(3)
		for j, text := range smallDataset {
			ng.Add(text, j)
		}
	}
}

func BenchmarkIndexBuild_Raphamorim_LSH_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lsh := raphamorim.NewLSH(10, 5, 3)
		for _, text := range smallDataset {
			lsh.Add(text)
		}
	}
}