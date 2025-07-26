package fuzzy

import (
	"testing"
)

func TestBKTree(t *testing.T) {
	tree := NewBKTree()
	words := []string{"book", "books", "cake", "boo", "boon", "cook", "cape", "cart"}
	
	for _, word := range words {
		tree.Add(word)
	}
	
	results := tree.Search("book", 2)
	expected := map[string]bool{"book": true, "books": true, "boo": true, "boon": true, "cook": true}
	
	for _, result := range results {
		if !expected[result] {
			t.Errorf("Unexpected result: %s", result)
		}
		delete(expected, result)
	}
	
	if len(expected) > 0 {
		t.Errorf("Missing results: %v", expected)
	}
}

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "def", 3},
		{"kitten", "sitting", 3},
		{"saturday", "sunday", 3},
	}
	
	for _, tt := range tests {
		got := LevenshteinDistance(tt.s1, tt.s2)
		if got != tt.want {
			t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestDamerauLevenshteinDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "acb", 1}, // transposition
		{"ca", "abc", 3},
	}
	
	for _, tt := range tests {
		got := DamerauLevenshteinDistance(tt.s1, tt.s2)
		if got != tt.want {
			t.Errorf("DamerauLevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func TestMyersDistance(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "a", 1},
		{"abc", "abc", 0},
		{"abc", "def", 6},
		{"ABCABBA", "CBABAC", 5},
	}
	
	for _, tt := range tests {
		got := MyersDistance(tt.s1, tt.s2)
		if got != tt.want {
			t.Errorf("MyersDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, got, tt.want)
		}
	}
}

func BenchmarkLevenshtein(b *testing.B) {
	s1 := "The quick brown fox jumps over the lazy dog"
	s2 := "The quick brown fox jumped over the lazy dogs"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}

func BenchmarkDamerauLevenshtein(b *testing.B) {
	s1 := "The quick brown fox jumps over the lazy dog"
	s2 := "The quick brown fox jumped over the lazy dogs"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DamerauLevenshteinDistance(s1, s2)
	}
}

func BenchmarkMyers(b *testing.B) {
	s1 := "The quick brown fox jumps over the lazy dog"
	s2 := "The quick brown fox jumped over the lazy dogs"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MyersDistance(s1, s2)
	}
}