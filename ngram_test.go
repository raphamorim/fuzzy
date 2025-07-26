package fuzzy

import (
	"testing"
)

func TestNGram(t *testing.T) {
	ng := NewNGram(3)
	
	texts := []string{
		"hello world",
		"hello there",
		"world peace",
		"goodbye world",
	}
	
	for i, text := range texts {
		ng.Add(text, i)
	}
	
	results := ng.Search("hello", 0.3)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'hello', got %d", len(results))
	}
	
	if results[0] != 0 && results[0] != 1 {
		t.Errorf("Expected results to contain indices 0 or 1, got %v", results)
	}
}

func TestJaccardSimilarity(t *testing.T) {
	ng := NewNGram(2)
	
	tests := []struct {
		text1, text2 string
		minSim       float64
	}{
		{"hello", "hello", 1.0},
		{"hello", "hallo", 0.3},
		{"abc", "def", 0.0},
		{"hello world", "hello there", 0.2},
	}
	
	for _, tt := range tests {
		sim := ng.JaccardSimilarity(tt.text1, tt.text2)
		if sim < tt.minSim {
			t.Errorf("JaccardSimilarity(%q, %q) = %f, want >= %f", tt.text1, tt.text2, sim, tt.minSim)
		}
	}
}

func TestTrigramIndex(t *testing.T) {
	ti := NewTrigramIndex()
	
	words := []string{
		"hello",
		"hallo",
		"hullo",
		"world",
		"word",
		"work",
	}
	
	for _, word := range words {
		ti.Add(word)
	}
	
	results := ti.Search("helo", 1)
	found := false
	for _, result := range results {
		if result == "hello" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("TrigramIndex failed to find 'hello' when searching for 'helo' with distance 1")
	}
}

func TestQGram(t *testing.T) {
	qg1 := NewQGram("hello world", 2)
	qg2 := NewQGram("hello there", 2)
	qg3 := NewQGram("goodbye world", 2)
	
	dist12 := qg1.Distance(qg2)
	dist13 := qg1.Distance(qg3)
	
	if dist12 >= dist13 {
		t.Error("Expected 'hello world' to be closer to 'hello there' than to 'goodbye world'")
	}
	
	sim12 := qg1.CosineSimilarity(qg2)
	sim13 := qg1.CosineSimilarity(qg3)
	
	if sim12 <= sim13 {
		t.Error("Expected higher cosine similarity between 'hello world' and 'hello there'")
	}
}

func BenchmarkNGramAdd(b *testing.B) {
	ng := NewNGram(3)
	text := "The quick brown fox jumps over the lazy dog"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ng.Add(text, i)
	}
}

func BenchmarkNGramSearch(b *testing.B) {
	ng := NewNGram(3)
	texts := []string{
		"The quick brown fox jumps over the lazy dog",
		"Pack my box with five dozen liquor jugs",
		"How vexingly quick daft zebras jump",
		"The five boxing wizards jump quickly",
	}
	
	for i, text := range texts {
		ng.Add(text, i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ng.Search("quick", 0.3)
	}
}

func BenchmarkTrigramSearch(b *testing.B) {
	ti := NewTrigramIndex()
	words := []string{
		"algorithm", "logarithm", "rhythm", "arithmetic",
		"programming", "program", "diagram", "grammar",
		"computer", "compiler", "interpreter", "terminal",
	}
	
	for _, word := range words {
		ti.Add(word)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ti.Search("algoritm", 2)
	}
}