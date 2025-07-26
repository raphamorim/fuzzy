package fuzzy

import (
	"testing"
)

func TestLSH(t *testing.T) {
	lsh := NewLSH(5, 3, 3)
	
	texts := []string{
		"The quick brown fox",
		"The quick brown dog",
		"A slow green turtle",
		"The fast brown fox",
	}
	
	for _, text := range texts {
		lsh.Add(text)
	}
	
	results := lsh.Query("The quick brown fox", 0.3)
	if len(results) < 1 {
		t.Errorf("Expected at least 1 similar result, got %d", len(results))
	}
}

func TestSimHash(t *testing.T) {
	sh := NewSimHash(64)
	
	texts := []string{
		"The quick brown fox jumps over the lazy dog",
		"The quick brown fox jumped over the lazy dog",
		"A completely different sentence with no similarity",
		"The fast brown fox jumps over the lazy cat",
	}
	
	for _, text := range texts {
		sh.Add(text)
	}
	
	results := sh.Query("The quick brown fox jumps over the lazy dog", 10)
	if len(results) < 2 {
		t.Errorf("Expected at least 2 similar results, got %d", len(results))
	}
	
	found := false
	for _, idx := range results {
		if idx == 0 {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("SimHash failed to find exact match")
	}
}

func TestHammingDistance(t *testing.T) {
	tests := []struct {
		a, b uint64
		want int
	}{
		{0, 0, 0},
		{0xFF, 0x00, 8},
		{0b1010, 0b0101, 4},
		{0b1111, 0b1110, 1},
	}
	
	for _, tt := range tests {
		got := hammingDistance(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("hammingDistance(%b, %b) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func BenchmarkLSHAdd(b *testing.B) {
	lsh := NewLSH(10, 5, 3)
	text := "The quick brown fox jumps over the lazy dog"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lsh.Add(text)
	}
}

func BenchmarkLSHQuery(b *testing.B) {
	lsh := NewLSH(10, 5, 3)
	
	texts := []string{
		"The quick brown fox jumps over the lazy dog",
		"Pack my box with five dozen liquor jugs",
		"How vexingly quick daft zebras jump",
		"The five boxing wizards jump quickly",
		"Sphinx of black quartz, judge my vow",
		"Waltz, bad nymph, for quick jigs vex",
	}
	
	for _, text := range texts {
		lsh.Add(text)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lsh.Query("The quick brown cat jumps over the lazy dog", 0.5)
	}
}

func BenchmarkSimHashCompute(b *testing.B) {
	sh := NewSimHash(64)
	text := "The quick brown fox jumps over the lazy dog"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sh.computeHash(text)
	}
}

func BenchmarkSimHashQuery(b *testing.B) {
	sh := NewSimHash(64)
	
	texts := []string{
		"The quick brown fox jumps over the lazy dog",
		"Pack my box with five dozen liquor jugs",
		"How vexingly quick daft zebras jump",
		"The five boxing wizards jump quickly",
		"Sphinx of black quartz, judge my vow",
		"Waltz, bad nymph, for quick jigs vex",
	}
	
	for _, text := range texts {
		sh.Add(text)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sh.Query("The quick brown cat jumps over the lazy dog", 5)
	}
}