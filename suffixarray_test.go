package fuzzy

import (
	"testing"
)

func TestSuffixArray(t *testing.T) {
	text := "banana"
	sa := NewSuffixArray(text)
	
	tests := []struct {
		pattern string
		want    []int
	}{
		{"ana", []int{1, 3}},
		{"ban", []int{0}},
		{"nan", []int{2}},
		{"xyz", nil},
	}
	
	for _, tt := range tests {
		got := sa.Search(tt.pattern)
		if len(got) != len(tt.want) {
			t.Errorf("Search(%q) returned %d results, want %d", tt.pattern, len(got), len(tt.want))
			continue
		}
		
		gotMap := make(map[int]bool)
		for _, pos := range got {
			gotMap[pos] = true
		}
		
		for _, pos := range tt.want {
			if !gotMap[pos] {
				t.Errorf("Search(%q) missing position %d", tt.pattern, pos)
			}
		}
	}
}

func TestSuffixArrayFuzzySearch(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	sa := NewSuffixArray(text)
	
	results := sa.FuzzySearch("quik", 1)
	found := false
	for _, pos := range results {
		if pos == 4 { // "quick" starts at position 4
			found = true
			break
		}
	}
	
	if !found {
		t.Error("FuzzySearch failed to find 'quick' when searching for 'quik' with 1 error")
	}
}

func TestFMIndex(t *testing.T) {
	text := "mississippi"
	fm := NewFMIndex(text, 2)
	
	tests := []struct {
		pattern string
		want    int
	}{
		{"si", 2},
		{"ssi", 2},
		{"iss", 2},
		{"i", 4},
		{"xyz", 0},
	}
	
	for _, tt := range tests {
		got := fm.Count(tt.pattern)
		if got != tt.want {
			t.Errorf("Count(%q) = %d, want %d", tt.pattern, got, tt.want)
		}
	}
}

func BenchmarkSuffixArrayBuild(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog. " +
		"Pack my box with five dozen liquor jugs. " +
		"How vexingly quick daft zebras jump!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewSuffixArray(text)
	}
}

func BenchmarkSuffixArraySearch(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog. " +
		"Pack my box with five dozen liquor jugs. " +
		"How vexingly quick daft zebras jump!"
	sa := NewSuffixArray(text)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sa.Search("quick")
	}
}

func BenchmarkFMIndexBuild(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog. " +
		"Pack my box with five dozen liquor jugs. " +
		"How vexingly quick daft zebras jump!"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFMIndex(text, 4)
	}
}

func BenchmarkFMIndexCount(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog. " +
		"Pack my box with five dozen liquor jugs. " +
		"How vexingly quick daft zebras jump!"
	fm := NewFMIndex(text, 4)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fm.Count("quick")
	}
}