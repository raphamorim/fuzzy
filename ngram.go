package fuzzy

import (
	"math"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// NGram is an optimized n-gram index with better performance and features
type NGram struct {
	n          int
	grams      map[string][]int
	corpus     []string
	corpusSize []int // Store corpus text sizes for better scoring
	mu         sync.RWMutex
	
	// Options
	normalize bool // Enable/disable text normalization
}

// NewNGram creates an optimized n-gram index
func NewNGram(n int) *NGram {
	return &NGram{
		n:          n,
		grams:      make(map[string][]int),
		corpus:     make([]string, 0, 1024),
		corpusSize: make([]int, 0, 1024),
		normalize:  true, // Default to true for better matching
	}
}

// SetNormalization enables or disables text normalization
func (ng *NGram) SetNormalization(enabled bool) {
	ng.normalize = enabled
}

// Add adds a text to the index with optimizations
func (ng *NGram) Add(text string, id int) {
	ng.mu.Lock()
	defer ng.mu.Unlock()
	
	ng.corpus = append(ng.corpus, text)
	ng.corpusSize = append(ng.corpusSize, utf8.RuneCountInString(text))
	
	// Process text based on normalization setting
	processedText := text
	if ng.normalize {
		processedText = ng.normalizeTextFast(text)
	} else {
		processedText = strings.ToLower(text)
	}
	
	// Generate n-grams
	grams := ng.generateNGramsFast(processedText)
	
	// Add to index without deduplication for speed
	for _, gram := range grams {
		ng.grams[gram] = append(ng.grams[gram], id)
	}
}

// normalizeTextFast performs faster text normalization
func (ng *NGram) normalizeTextFast(text string) string {
	// Fast path for ASCII-only text
	if isASCII(text) {
		return normalizeASCII(text)
	}
	
	// Unicode path
	var result []rune
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result = append(result, unicode.ToLower(r))
		}
	}
	return string(result)
}

// normalizeASCII fast normalization for ASCII text
func normalizeASCII(text string) string {
	result := make([]byte, 0, len(text))
	for i := 0; i < len(text); i++ {
		c := text[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == ' ' {
			result = append(result, c)
		} else if c >= 'A' && c <= 'Z' {
			result = append(result, c+32) // Convert to lowercase
		}
	}
	return string(result)
}

// generateNGramsFast generates n-grams with minimal overhead
func (ng *NGram) generateNGramsFast(text string) []string {
	n := ng.n
	
	// Fast path for ASCII
	if isASCII(text) {
		if len(text) < n {
			return []string{text}
		}
		
		count := len(text) - n + 1
		grams := make([]string, count)
		for i := 0; i < count; i++ {
			grams[i] = text[i:i+n]
		}
		return grams
	}
	
	// Unicode path
	runes := []rune(text)
	if len(runes) < n {
		return []string{text}
	}
	
	count := len(runes) - n + 1
	grams := make([]string, count)
	for i := 0; i < count; i++ {
		grams[i] = string(runes[i:i+n])
	}
	return grams
}

// NGramResult contains search result with additional metadata
type NGramResult struct {
	ID    int
	Score float64
	Text  string
}

// Search performs optimized search with better scoring
func (ng *NGram) Search(query string, threshold float64) []NGramResult {
	ng.mu.RLock()
	defer ng.mu.RUnlock()
	
	// Process query
	processedQuery := query
	if ng.normalize {
		processedQuery = ng.normalizeTextFast(query)
	} else {
		processedQuery = strings.ToLower(query)
	}
	
	queryGrams := ng.generateNGramsFast(processedQuery)
	if len(queryGrams) == 0 {
		return nil
	}
	
	// Count matching grams per document
	candidates := make(map[int]int, 32)
	
	for _, gram := range queryGrams {
		if ids, exists := ng.grams[gram]; exists {
			for _, id := range ids {
				candidates[id]++
			}
		}
	}
	
	// Calculate scores
	results := make([]NGramResult, 0, len(candidates))
	queryLen := float64(len(queryGrams))
	
	for id, matchCount := range candidates {
		score := float64(matchCount) / queryLen
		
		if score >= threshold {
			results = append(results, NGramResult{
				ID:    id,
				Score: score,
				Text:  ng.corpus[id],
			})
		}
	}
	
	// Sort by score descending
	if len(results) > 1 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Score > results[j].Score
		})
	}
	
	return results
}

// SearchIDs performs search returning only IDs for backward compatibility
func (ng *NGram) SearchIDs(query string, threshold float64) []int {
	results := ng.Search(query, threshold)
	ids := make([]int, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}
	return ids
}

// BatchAdd adds multiple texts efficiently
func (ng *NGram) BatchAdd(texts []string) {
	ng.mu.Lock()
	defer ng.mu.Unlock()
	
	startID := len(ng.corpus)
	
	// Pre-allocate space
	ng.corpus = append(ng.corpus, texts...)
	
	for i, text := range texts {
		ng.corpusSize = append(ng.corpusSize, utf8.RuneCountInString(text))
		
		processedText := text
		if ng.normalize {
			processedText = ng.normalizeTextFast(text)
		} else {
			processedText = strings.ToLower(text)
		}
		
		grams := ng.generateNGramsFast(processedText)
		
		for _, gram := range grams {
			ng.grams[gram] = append(ng.grams[gram], startID+i)
		}
	}
}

// Size returns the number of documents in the index
func (ng *NGram) Size() int {
	ng.mu.RLock()
	defer ng.mu.RUnlock()
	return len(ng.corpus)
}

// Clear removes all data from the index
func (ng *NGram) Clear() {
	ng.mu.Lock()
	defer ng.mu.Unlock()
	
	ng.grams = make(map[string][]int)
	ng.corpus = ng.corpus[:0]
	ng.corpusSize = ng.corpusSize[:0]
}

// Helper functions

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// TrigramIndex is an optimized trigram index
type TrigramIndex struct {
	*NGram
}

// NewTrigramIndex creates an optimized trigram index
func NewTrigramIndex() *TrigramIndex {
	return &TrigramIndex{
		NGram: NewNGram(3),
	}
}

// SearchWithDistance searches using edit distance filtering
func (ti *TrigramIndex) SearchWithDistance(query string, maxDistance int) []NGramResult {
	// Use trigram similarity to filter candidates
	minSimilarity := 1.0 - float64(maxDistance)*0.3
	candidates := ti.Search(query, minSimilarity)
	
	// Filter by actual edit distance
	results := make([]NGramResult, 0, len(candidates))
	for _, candidate := range candidates {
		dist := LevenshteinDistance(query, candidate.Text)
		if dist <= maxDistance {
			// Update score based on edit distance
			candidate.Score = 1.0 - float64(dist)/float64(maxInt(len(query), len(candidate.Text)))
			results = append(results, candidate)
		}
	}
	
	// Re-sort by score
	if len(results) > 1 {
		sort.Slice(results, func(i, j int) bool {
			return results[i].Score > results[j].Score
		})
	}
	
	return results
}

// JaccardSimilarity calculates Jaccard similarity between two texts
func (ng *NGram) JaccardSimilarity(text1, text2 string) float64 {
	// Process texts
	if ng.normalize {
		text1 = ng.normalizeTextFast(text1)
		text2 = ng.normalizeTextFast(text2)
	} else {
		text1 = strings.ToLower(text1)
		text2 = strings.ToLower(text2)
	}
	
	grams1 := ng.generateNGramsFast(text1)
	grams2 := ng.generateNGramsFast(text2)
	
	// Use map for set operations
	set1 := make(map[string]struct{}, len(grams1))
	for _, g := range grams1 {
		set1[g] = struct{}{}
	}
	
	intersection := 0
	set2Size := 0
	seen := make(map[string]struct{}, len(grams2))
	
	for _, g := range grams2 {
		if _, exists := seen[g]; !exists {
			seen[g] = struct{}{}
			set2Size++
			if _, inSet1 := set1[g]; inSet1 {
				intersection++
			}
		}
	}
	
	union := len(set1) + set2Size - intersection
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

// QGram represents a q-gram profile for a text
type QGram struct {
	q     int
	grams map[string]float64
}

// NewQGram creates a new q-gram profile
func NewQGram(text string, q int) *QGram {
	qg := &QGram{
		q:     q,
		grams: make(map[string]float64),
	}
	
	text = strings.ToLower(text)
	runes := []rune(text)
	
	for i := 0; i <= len(runes)-q; i++ {
		gram := string(runes[i : i+q])
		qg.grams[gram]++
	}
	
	return qg
}

// Distance calculates the q-gram distance
func (qg *QGram) Distance(other *QGram) float64 {
	distance := 0.0
	
	// Process all grams from first profile
	for g, count1 := range qg.grams {
		count2 := other.grams[g]
		distance += math.Abs(count1 - count2)
	}
	
	// Add grams only in second profile
	for g, count2 := range other.grams {
		if _, exists := qg.grams[g]; !exists {
			distance += count2
		}
	}
	
	return distance
}

// CosineSimilarity calculates cosine similarity between q-gram profiles
func (qg *QGram) CosineSimilarity(other *QGram) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0
	
	for g, count1 := range qg.grams {
		magnitude1 += count1 * count1
		if count2, exists := other.grams[g]; exists {
			dotProduct += count1 * count2
		}
	}
	
	for _, count2 := range other.grams {
		magnitude2 += count2 * count2
	}
	
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}
	
	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}