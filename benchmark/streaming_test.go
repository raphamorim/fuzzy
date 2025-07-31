package benchmark

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	raphamorim "github.com/raphamorim/fuzzy"
)

// StreamingBKTree builds index in chunks to handle large files
type StreamingBKTree struct {
	trees []*raphamorim.BKTree
	mu    sync.RWMutex
}

func NewStreamingBKTree() *StreamingBKTree {
	return &StreamingBKTree{
		trees: make([]*raphamorim.BKTree, 0),
	}
}

func (s *StreamingBKTree) AddChunk(words []string) {
	tree := raphamorim.NewBKTree()
	for _, word := range words {
		tree.Add(word)
	}
	
	s.mu.Lock()
	s.trees = append(s.trees, tree)
	s.mu.Unlock()
}

func (s *StreamingBKTree) Search(word string, maxDist int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	seen := make(map[string]bool)
	var results []string
	
	for _, tree := range s.trees {
		matches := tree.Search(word, maxDist)
		for _, match := range matches {
			if !seen[match] {
				seen[match] = true
				results = append(results, match)
			}
		}
	}
	
	return results
}

func BenchmarkStreaming10GB(b *testing.B) {
	if _, err := os.Stat("testdata/10gb_words.txt"); os.IsNotExist(err) {
		b.Skip("10GB test file not found. Run: go run generate_10gb.go")
	}

	fmt.Println("Building streaming index for 10GB file...")
	startTime := time.Now()
	
	streamTree := NewStreamingBKTree()
	chunkSize := 100000 // Process 100k lines at a time
	
	file, err := os.Open("testdata/10gb_words.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	chunk := make([]string, 0, chunkSize)
	totalLines := 0
	
	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		
		if len(chunk) >= chunkSize {
			streamTree.AddChunk(chunk)
			totalLines += len(chunk)
			chunk = make([]string, 0, chunkSize)
			
			if totalLines%1000000 == 0 {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				fmt.Printf("Processed %d lines, Memory: %d MB\n", totalLines, m.Alloc/1024/1024)
			}
		}
	}
	
	// Add remaining chunk
	if len(chunk) > 0 {
		streamTree.AddChunk(chunk)
		totalLines += len(chunk)
	}
	
	fmt.Printf("Index built in %v for %d lines\n", time.Since(startTime), totalLines)
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Final memory usage: %d MB\n", m.Alloc/1024/1024)
	
	queries := []string{"algorithm", "databse", "netwrk", "securty", "performnce"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		_ = streamTree.Search(query, 2)
	}
}

func BenchmarkParallelSearch10GB(b *testing.B) {
	if _, err := os.Stat("testdata/10gb_words.txt"); os.IsNotExist(err) {
		b.Skip("10GB test file not found. Run: go run generate_10gb.go")
	}

	numWorkers := runtime.NumCPU()
	queries := []string{"algorithm", "database", "network", "security", "performance"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		
		file, err := os.Open("testdata/10gb_words.txt")
		if err != nil {
			b.Fatal(err)
		}
		
		// Parallel search
		resultChan := make(chan string, 100)
		var wg sync.WaitGroup
		
		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if raphamorim.LevenshteinDistance(query, line) <= 2 {
						resultChan <- line
					}
				}
			}()
		}
		
		// Collect results
		go func() {
			wg.Wait()
			close(resultChan)
		}()
		
		matches := 0
		for range resultChan {
			matches++
		}
		
		file.Close()
	}
}