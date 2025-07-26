package fuzzy

import (
	"math"
	"sort"
	"strings"
)

type NGram struct {
	n      int
	grams  map[string][]int
	corpus []string
}

func NewNGram(n int) *NGram {
	return &NGram{
		n:      n,
		grams:  make(map[string][]int),
		corpus: make([]string, 0),
	}
}

func (ng *NGram) Add(text string, id int) {
	ng.corpus = append(ng.corpus, text)
	grams := ng.generateNGrams(text)
	
	for _, gram := range grams {
		ng.grams[gram] = append(ng.grams[gram], id)
	}
}

func (ng *NGram) generateNGrams(text string) []string {
	text = strings.ToLower(text)
	runes := []rune(text)
	
	if len(runes) < ng.n {
		return []string{text}
	}
	
	grams := make([]string, 0, len(runes)-ng.n+1)
	for i := 0; i <= len(runes)-ng.n; i++ {
		grams = append(grams, string(runes[i:i+ng.n]))
	}
	
	return grams
}

func (ng *NGram) Search(query string, threshold float64) []int {
	queryGrams := ng.generateNGrams(query)
	candidates := make(map[int]int)
	
	for _, gram := range queryGrams {
		if ids, exists := ng.grams[gram]; exists {
			for _, id := range ids {
				candidates[id]++
			}
		}
	}
	
	type result struct {
		id    int
		score float64
	}
	
	results := make([]result, 0)
	for id, count := range candidates {
		score := float64(count) / float64(len(queryGrams))
		if score >= threshold {
			results = append(results, result{id: id, score: score})
		}
	}
	
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})
	
	ids := make([]int, len(results))
	for i, r := range results {
		ids[i] = r.id
	}
	
	return ids
}

func (ng *NGram) JaccardSimilarity(text1, text2 string) float64 {
	grams1 := ng.generateNGrams(text1)
	grams2 := ng.generateNGrams(text2)
	
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)
	
	for _, g := range grams1 {
		set1[g] = true
	}
	for _, g := range grams2 {
		set2[g] = true
	}
	
	intersection := 0
	for g := range set1 {
		if set2[g] {
			intersection++
		}
	}
	
	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

type TrigramIndex struct {
	trigrams map[string]map[int]bool
	texts    []string
}

func NewTrigramIndex() *TrigramIndex {
	return &TrigramIndex{
		trigrams: make(map[string]map[int]bool),
		texts:    make([]string, 0),
	}
}

func (ti *TrigramIndex) Add(text string) int {
	id := len(ti.texts)
	ti.texts = append(ti.texts, text)
	
	trigrams := extractTrigrams(text)
	for _, trigram := range trigrams {
		if ti.trigrams[trigram] == nil {
			ti.trigrams[trigram] = make(map[int]bool)
		}
		ti.trigrams[trigram][id] = true
	}
	
	return id
}

func (ti *TrigramIndex) Search(query string, maxDistance int) []string {
	queryTrigrams := extractTrigrams(query)
	candidates := make(map[int]int)
	
	for _, trigram := range queryTrigrams {
		if ids, exists := ti.trigrams[trigram]; exists {
			for id := range ids {
				candidates[id]++
			}
		}
	}
	
	type match struct {
		text     string
		distance int
	}
	
	matches := make([]match, 0)
	
	for id, trigramCount := range candidates {
		text := ti.texts[id]
		
		minLen := len(query) - maxDistance
		maxLen := len(query) + maxDistance
		if len(text) < minLen || len(text) > maxLen {
			continue
		}
		
		textTrigrams := len(extractTrigrams(text))
		queryTrigramCount := len(queryTrigrams)
		
		maxPossibleTrigrams := max(textTrigrams, queryTrigramCount)
		minRequiredTrigrams := maxPossibleTrigrams - maxDistance*3
		
		if trigramCount >= minRequiredTrigrams {
			dist := LevenshteinDistance(query, text)
			if dist <= maxDistance {
				matches = append(matches, match{text: text, distance: dist})
			}
		}
	}
	
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].distance < matches[j].distance
	})
	
	results := make([]string, len(matches))
	for i, m := range matches {
		results[i] = m.text
	}
	
	return results
}

func extractTrigrams(text string) []string {
	text = strings.ToLower(text)
	runes := []rune(text)
	
	if len(runes) < 3 {
		return []string{text}
	}
	
	trigrams := make([]string, 0, len(runes)-2)
	for i := 0; i <= len(runes)-3; i++ {
		trigrams = append(trigrams, string(runes[i:i+3]))
	}
	
	return trigrams
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type QGram struct {
	q     int
	grams map[string]float64
}

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

func (qg *QGram) Distance(other *QGram) float64 {
	distance := 0.0
	
	allGrams := make(map[string]bool)
	for g := range qg.grams {
		allGrams[g] = true
	}
	for g := range other.grams {
		allGrams[g] = true
	}
	
	for g := range allGrams {
		diff := qg.grams[g] - other.grams[g]
		distance += math.Abs(diff)
	}
	
	return distance
}

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