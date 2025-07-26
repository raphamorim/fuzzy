package fuzzy

import (
	"testing"
)

func TestWuManber(t *testing.T) {
	wm := NewWuManber("hello")
	
	tests := []struct {
		text      string
		maxErrors int
		wantMin   int
		wantMax   int
	}{
		{"hello world", 0, 1, 1},
		{"helo world", 1, 1, 10},
		{"hllo world", 1, 1, 10},
		{"goodbye", 2, 0, 10},
		{"hello hello", 0, 2, 2},
	}
	
	for _, tt := range tests {
		matches := wm.Search(tt.text, tt.maxErrors)
		if len(matches) < tt.wantMin || len(matches) > tt.wantMax {
			t.Errorf("Search(%q, %d) returned %d matches, want %d-%d", tt.text, tt.maxErrors, len(matches), tt.wantMin, tt.wantMax)
		}
	}
}

func TestUkkonenAStar(t *testing.T) {
	ua := NewUkkonenAStar("pattern", 2)
	
	text := "This is a patern in the text with pattern too"
	matches := ua.Search(text)
	
	if len(matches) < 2 {
		t.Errorf("Expected at least 2 matches, got %d", len(matches))
	}
	
	foundExact := false
	foundApprox := false
	
	for _, match := range matches {
		if match.Distance == 0 {
			foundExact = true
		} else if match.Distance == 1 {
			foundApprox = true
		}
	}
	
	if !foundExact {
		t.Error("Failed to find exact match for 'pattern'")
	}
	
	if !foundApprox {
		t.Error("Failed to find approximate match for 'patern'")
	}
}

func TestBitonicSort(t *testing.T) {
	matches := []FuzzyMatch{
		{Text: "match1", Score: 0.5, Distance: 2},
		{Text: "match2", Score: 0.8, Distance: 1},
		{Text: "match3", Score: 0.3, Distance: 3},
		{Text: "match4", Score: 0.9, Distance: 0},
		{Text: "match5", Score: 0.6, Distance: 2},
	}
	
	bs := NewBitonicSort(matches)
	sorted := bs.Sort()
	
	for i := 1; i < len(sorted); i++ {
		if sorted[i-1].Score < sorted[i].Score {
			t.Errorf("BitonicSort failed: %f < %f at positions %d and %d", 
				sorted[i-1].Score, sorted[i].Score, i-1, i)
		}
	}
}

func BenchmarkWuManberExact(b *testing.B) {
	wm := NewWuManber("algorithm")
	text := "This text contains the word algorithm multiple times. The algorithm is efficient."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wm.Search(text, 0)
	}
}

func BenchmarkWuManberApproximate(b *testing.B) {
	wm := NewWuManber("algorithm")
	text := "This text contains the word algoritm and algorythm with errors. The algorthm is efficient."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wm.Search(text, 2)
	}
}

func BenchmarkUkkonenAStar(b *testing.B) {
	ua := NewUkkonenAStar("pattern", 2)
	text := "This is a long text with multiple occurrences of patern, pattern, and pattrn scattered throughout."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ua.Search(text)
	}
}

func BenchmarkBitonicSort(b *testing.B) {
	matches := make([]FuzzyMatch, 64)
	for i := range matches {
		matches[i] = FuzzyMatch{
			Text:     "match",
			Score:    float64(i) / 64.0,
			Distance: i % 5,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bs := NewBitonicSort(matches)
		bs.Sort()
	}
}