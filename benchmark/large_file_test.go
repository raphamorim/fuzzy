package benchmark

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	raphamorim "github.com/raphamorim/fuzzy"
)

func loadLargeFile(filename string, limit int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	count := 0
	
	for scanner.Scan() && (limit == 0 || count < limit) {
		lines = append(lines, scanner.Text())
		count++
	}
	
	return lines, scanner.Err()
}

func BenchmarkBKTree10GB(b *testing.B) {
	// Check if 10GB file exists
	if _, err := os.Stat("testdata/10gb_words.txt"); os.IsNotExist(err) {
		b.Skip("10GB test file not found. Run: go run generate_10gb.go")
	}

	// Load a subset for indexing (e.g., first 1M lines)
	fmt.Println("Loading data for indexing...")
	data, err := loadLargeFile("testdata/10gb_words.txt", 1000000)
	if err != nil {
		b.Fatal(err)
	}

	// Build BK-tree
	fmt.Printf("Building BK-tree with %d items...\n", len(data))
	startTime := time.Now()
	tree := raphamorim.NewBKTree()
	for _, word := range data {
		tree.Add(word)
	}
	fmt.Printf("BK-tree built in %v\n", time.Since(startTime))

	// Memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory usage: %d MB\n", m.Alloc/1024/1024)

	queries := []string{"algorithm", "databse", "netwrk", "securty", "performnce"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		_ = tree.Search(query, 2)
	}
}

func BenchmarkNGram10GB(b *testing.B) {
	if _, err := os.Stat("testdata/10gb_words.txt"); os.IsNotExist(err) {
		b.Skip("10GB test file not found. Run: go run generate_10gb.go")
	}

	fmt.Println("Loading data for indexing...")
	data, err := loadLargeFile("testdata/10gb_words.txt", 1000000)
	if err != nil {
		b.Fatal(err)
	}

	fmt.Printf("Building N-gram index with %d items...\n", len(data))
	startTime := time.Now()
	ngram := raphamorim.NewNGram(3)
	for i, text := range data {
		ngram.Add(text, i)
	}
	fmt.Printf("N-gram index built in %v\n", time.Since(startTime))

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory usage: %d MB\n", m.Alloc/1024/1024)

	queries := []string{"algorithm", "database", "network", "security", "performance"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		_ = ngram.Search(query, 0.7)
	}
}

func BenchmarkStreamSearch10GB(b *testing.B) {
	if _, err := os.Stat("testdata/10gb_words.txt"); os.IsNotExist(err) {
		b.Skip("10GB test file not found. Run: go run generate_10gb.go")
	}

	queries := []string{"algorithm", "database", "network", "security", "performance"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		
		// Stream search through file
		file, err := os.Open("testdata/10gb_words.txt")
		if err != nil {
			b.Fatal(err)
		}
		
		scanner := bufio.NewScanner(file)
		matches := 0
		for scanner.Scan() {
			line := scanner.Text()
			if raphamorim.LevenshteinDistance(query, line) <= 2 {
				matches++
			}
		}
		
		file.Close()
	}
}