package fuzzy

import (
	"fmt"
	"strings"
	"testing"
)

func TestNGramBasic(t *testing.T) {
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
	if len(results) < 2 {
		t.Errorf("Expected at least 2 results for 'hello', got %d", len(results))
	}
	
	// Check that hello-containing texts are found
	foundHello := false
	for _, r := range results {
		if strings.Contains(r.Text, "hello") {
			foundHello = true
			break
		}
	}
	if !foundHello {
		t.Error("Failed to find texts containing 'hello'")
	}
}

func TestNGramUnicode(t *testing.T) {
	ng := NewNGram(3)
	
	texts := []string{
		"こんにちは世界", // Japanese
		"你好世界",      // Chinese
		"مرحبا بالعالم", // Arabic
		"Привет мир",   // Russian
	}
	
	for i, text := range texts {
		ng.Add(text, i)
	}
	
	// Test Japanese search
	results := ng.Search("こんにちは", 0.5)
	if len(results) == 0 {
		t.Error("Failed to find Japanese text")
	}
}

func TestNGramNormalization(t *testing.T) {
	ng := NewNGram(3)
	
	// Add text with punctuation and mixed case
	ng.Add("Hello, World!", 0)
	ng.Add("HELLO WORLD", 1)
	ng.Add("hello world", 2)
	
	results := ng.Search("hello world", 0.8)
	if len(results) < 3 {
		t.Errorf("Expected all 3 variations to match, got %d", len(results))
	}
}

func TestNGramBatchAdd(t *testing.T) {
	ng := NewNGram(3)
	
	texts := []string{
		"apple pie",
		"apple juice",
		"orange juice",
		"banana split",
	}
	
	ng.BatchAdd(texts)
	
	if ng.Size() != 4 {
		t.Errorf("Expected 4 documents, got %d", ng.Size())
	}
	
	results := ng.Search("apple", 0.3)
	if len(results) < 2 {
		t.Error("Failed to find apple-related texts")
	}
}

func TestTrigramIndex(t *testing.T) {
	ti := NewTrigramIndex()
	
	words := []string{
		"algorithm",
		"logarithm",
		"rhythm",
		"arithmetic",
	}
	
	for i, word := range words {
		ti.Add(word, i)
	}
	
	results := ti.SearchWithDistance("algoritm", 2)
	
	found := false
	for _, r := range results {
		if r.Text == "algorithm" {
			found = true
			if r.Score <= 0 {
				t.Error("Expected positive similarity score")
			}
			break
		}
	}
	
	if !found {
		t.Error("Failed to find 'algorithm' when searching for 'algoritm'")
	}
}

func TestNGramConcurrency(t *testing.T) {
	ng := NewNGram(3)
	
	// Concurrent adds
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			text := fmt.Sprintf("document %d with some content", id)
			ng.Add(text, id)
			done <- true
		}(i)
	}
	
	for i := 0; i < 10; i++ {
		<-done
	}
	
	if ng.Size() != 10 {
		t.Errorf("Expected 10 documents, got %d", ng.Size())
	}
	
	// Concurrent searches
	for i := 0; i < 5; i++ {
		go func() {
			results := ng.Search("document", 0.3)
			if len(results) < 10 {
				t.Errorf("Expected 10 results, got %d", len(results))
			}
			done <- true
		}()
	}
	
	for i := 0; i < 5; i++ {
		<-done
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

func BenchmarkNGramBatchAdd(b *testing.B) {
	texts := make([]string, 100)
	for i := range texts {
		texts[i] = fmt.Sprintf("Document %d with some sample text content", i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ng := NewNGram(3)
		ng.BatchAdd(texts)
	}
}

func BenchmarkNGramUnicode(b *testing.B) {
	ng := NewNGram(3)
	texts := []string{
		"こんにちは世界、今日はいい天気ですね",
		"你好世界，今天天气很好",
		"مرحبا بالعالم، الطقس جميل اليوم",
		"Привет мир, сегодня хорошая погода",
	}
	
	for i, text := range texts {
		ng.Add(text, i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ng.Search("世界", 0.3)
	}
}

// Comparison benchmarks
func BenchmarkNGramOriginalAdd(b *testing.B) {
	ng := NewNGram(3)
	text := "The quick brown fox jumps over the lazy dog"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ng.Add(text, i)
	}
}

func BenchmarkNGramOriginalSearch(b *testing.B) {
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