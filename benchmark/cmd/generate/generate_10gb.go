package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	targetSize = 10 * 1024 * 1024 * 1024 // 10GB
	bufferSize = 1024 * 1024             // 1MB buffer
)

var words = []string{
	"algorithm", "database", "network", "security", "performance",
	"optimization", "distributed", "scalability", "reliability", "efficiency",
	"architecture", "microservice", "kubernetes", "container", "deployment",
	"monitoring", "logging", "debugging", "testing", "automation",
	"integration", "continuous", "delivery", "pipeline", "infrastructure",
	"cloud", "computing", "storage", "processing", "analytics",
	"machine", "learning", "artificial", "intelligence", "neural",
	"blockchain", "cryptocurrency", "decentralized", "consensus", "protocol",
	"encryption", "authentication", "authorization", "certificate", "firewall",
	"vulnerability", "penetration", "compliance", "governance", "audit",
}

func generateSentence(r *rand.Rand) string {
	length := r.Intn(10) + 5
	sentence := make([]string, length)
	for i := 0; i < length; i++ {
		sentence[i] = words[r.Intn(len(words))]
	}
	return strings.Join(sentence, " ")
}

func main() {
	fmt.Println("Generating 10GB test file...")
	startTime := time.Now()

	file, err := os.Create("testdata/10gb_words.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, bufferSize)
	defer writer.Flush()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var written int64
	var lineCount int64

	for written < targetSize {
		sentence := generateSentence(r) + "\n"
		n, err := writer.WriteString(sentence)
		if err != nil {
			panic(err)
		}
		written += int64(n)
		lineCount++

		if lineCount%1000000 == 0 {
			progress := float64(written) / float64(targetSize) * 100
			fmt.Printf("\rProgress: %.2f%% (%d MB written, %d lines)", 
				progress, written/(1024*1024), lineCount)
		}
	}

	fmt.Printf("\nCompleted! Generated %d MB in %d lines\n", 
		written/(1024*1024), lineCount)
	fmt.Printf("Time taken: %v\n", time.Since(startTime))
}