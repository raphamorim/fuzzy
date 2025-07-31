package benchmark

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"

	raphamorim "github.com/raphamorim/fuzzy"
)

var (
	// Test data
	smallDataset   []string
	mediumDataset  []string
	largeDataset   []string
	hugeDataset    []string
	
	// Search queries
	exactQueries      []string
	typoQueries       []string
	partialQueries    []string
	
	// Pre-built indices
	raphamorimBKTree     *raphamorim.BKTree
	raphamorimNGram      *raphamorim.NGram
	raphamorimLSH        *raphamorim.LSH
	raphamorimTrigram    *raphamorim.TrigramIndex
	raphamorimSuffixArr  *raphamorim.SuffixArray
)

func init() {
	rand.Seed(time.Now().UnixNano())
	
	// Generate datasets
	smallDataset = generateDataset(100)
	mediumDataset = generateDataset(10000)
	largeDataset = generateDataset(100000)
	hugeDataset = generateDataset(1000000)
	
	// Generate queries
	exactQueries = selectRandomItems(mediumDataset, 100)
	typoQueries = generateTypos(exactQueries)
	partialQueries = generatePartials(exactQueries)
	
	// Build indices for raphamorim/fuzzy
	buildRaphamorimIndices()
}

func generateDataset(size int) []string {
	// Try to load dictionary file first
	words := loadDictionary()
	if len(words) > 0 {
		return expandDataset(words, size)
	}
	
	// Fallback to generated data
	dataset := make([]string, size)
	for i := 0; i < size; i++ {
		dataset[i] = generateRandomString()
	}
	return dataset
}

func loadDictionary() []string {
	// Try common dictionary locations
	paths := []string{
		"/usr/share/dict/words",
		"/usr/dict/words",
		"./testdata/words.txt",
	}
	
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()
		
		var words []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			if len(word) > 0 {
				words = append(words, word)
			}
		}
		
		if len(words) > 0 {
			return words
		}
	}
	
	return nil
}

func expandDataset(baseWords []string, targetSize int) []string {
	dataset := make([]string, 0, targetSize)
	
	// Add original words
	for _, word := range baseWords {
		if len(dataset) >= targetSize {
			break
		}
		dataset = append(dataset, word)
	}
	
	// Generate variations
	for len(dataset) < targetSize {
		base := baseWords[rand.Intn(len(baseWords))]
		variation := generateVariation(base)
		dataset = append(dataset, variation)
	}
	
	return dataset
}

func generateVariation(base string) string {
	variations := []func(string) string{
		func(s string) string { return s + generateSuffix() },
		func(s string) string { return generatePrefix() + s },
		func(s string) string { return s + "_" + generateRandomString() },
		func(s string) string { return strings.Title(s) },
		func(s string) string { return strings.ToUpper(s) },
		func(s string) string { return strings.Replace(s, "e", "3", -1) },
		func(s string) string { return strings.Replace(s, "a", "@", -1) },
	}
	
	return variations[rand.Intn(len(variations))](base)
}

func generatePrefix() string {
	prefixes := []string{"pre", "post", "anti", "super", "ultra", "mega", "micro", "nano"}
	return prefixes[rand.Intn(len(prefixes))]
}

func generateSuffix() string {
	suffixes := []string{"ing", "ed", "er", "est", "ly", "ness", "ment", "ful", "less"}
	return suffixes[rand.Intn(len(suffixes))]
}

func generateRandomString() string {
	length := rand.Intn(15) + 5
	chars := "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func selectRandomItems(dataset []string, count int) []string {
	if count > len(dataset) {
		count = len(dataset)
	}
	
	selected := make([]string, count)
	indices := rand.Perm(len(dataset))
	for i := 0; i < count; i++ {
		selected[i] = dataset[indices[i]]
	}
	return selected
}

func generateTypos(queries []string) []string {
	typos := make([]string, len(queries))
	for i, query := range queries {
		typos[i] = introduceTypo(query)
	}
	return typos
}

func introduceTypo(s string) string {
	if len(s) < 2 {
		return s
	}
	
	runes := []rune(s)
	typoType := rand.Intn(4)
	
	switch typoType {
	case 0: // Deletion
		pos := rand.Intn(len(runes))
		return string(append(runes[:pos], runes[pos+1:]...))
	case 1: // Insertion
		pos := rand.Intn(len(runes) + 1)
		char := rune('a' + rand.Intn(26))
		return string(append(runes[:pos], append([]rune{char}, runes[pos:]...)...))
	case 2: // Substitution
		pos := rand.Intn(len(runes))
		runes[pos] = rune('a' + rand.Intn(26))
		return string(runes)
	case 3: // Transposition
		if len(runes) > 1 {
			pos := rand.Intn(len(runes) - 1)
			runes[pos], runes[pos+1] = runes[pos+1], runes[pos]
		}
		return string(runes)
	}
	
	return s
}

func generatePartials(queries []string) []string {
	partials := make([]string, len(queries))
	for i, query := range queries {
		if len(query) > 4 {
			start := rand.Intn(len(query) / 2)
			end := start + len(query)/2 + rand.Intn(len(query)/2)
			if end > len(query) {
				end = len(query)
			}
			partials[i] = query[start:end]
		} else {
			partials[i] = query
		}
	}
	return partials
}

func buildRaphamorimIndices() {
	// Build BK-Tree
	raphamorimBKTree = raphamorim.NewBKTree()
	for _, word := range mediumDataset {
		raphamorimBKTree.Add(word)
	}
	
	// Build N-gram index
	raphamorimNGram = raphamorim.NewNGram(3)
	for i, text := range mediumDataset {
		raphamorimNGram.Add(text, i)
	}
	
	// Build LSH
	raphamorimLSH = raphamorim.NewLSH(10, 5, 3)
	for _, text := range mediumDataset {
		raphamorimLSH.Add(text)
	}
	
	// Build Trigram index
	raphamorimTrigram = raphamorim.NewTrigramIndex()
	for i, text := range mediumDataset {
		raphamorimTrigram.Add(text, i)
	}
	
	// Build Suffix Array (for smaller dataset due to memory)
	if len(mediumDataset) > 0 {
		combined := strings.Join(mediumDataset[:min(1000, len(mediumDataset))], " ")
		raphamorimSuffixArr = raphamorim.NewSuffixArray(combined)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}