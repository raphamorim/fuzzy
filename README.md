# Advanced Fuzzy Search Library for Go

A comprehensive fuzzy search library implementing state-of-the-art algorithms and data structures for approximate string matching, text search, and similarity detection.

## Features

### Distance Metrics
- **Levenshtein Distance**: Classic edit distance algorithm
- **Damerau-Levenshtein Distance**: Supports transpositions
- **Myers' Algorithm**: Efficient diff algorithm for edit distance

### Data Structures
- **BK-Tree**: Metric tree for efficient similarity search
- **Suffix Array**: For substring search and pattern matching
- **FM-Index**: Compressed full-text index based on Burrows-Wheeler Transform

### Indexing Methods
- **N-gram Indexing**: Character n-gram based search with Jaccard similarity
- **Trigram Index**: Optimized 3-gram indexing for approximate matching
- **Q-gram Distance**: Distance metric based on q-gram profiles

### Locality-Sensitive Hashing
- **MinHash LSH**: For finding similar documents
- **SimHash**: Near-duplicate detection using hamming distance

### Advanced Algorithms
- **Wu-Manber**: Bit-parallel approximate string matching
- **Ukkonen A***: A* search for approximate pattern matching
- **Bitonic Sort**: Parallel-friendly sorting for fuzzy matches

## Installation

```bash
go get github.com/raphamorim/fuzzy
```

## Usage Examples

### BK-Tree for Spell Checking

```go
tree := fuzzy.NewBKTree()
words := []string{"book", "books", "cake", "boo", "boon", "cook"}
for _, word := range words {
    tree.Add(word)
}

// Find words within edit distance 2 of "bok"
suggestions := tree.Search("bok", 2)
// Returns: ["book", "boo", "cook"]
```

### N-gram Search

```go
ng := fuzzy.NewNGram(3)
ng.Add("hello world", 0)
ng.Add("hello there", 1)
ng.Add("goodbye world", 2)

// Search with similarity threshold
results := ng.Search("helo wrld", 0.3)
// Returns indices of matching documents
```

### LSH for Document Similarity

```go
lsh := fuzzy.NewLSH(10, 5, 3) // 10 hash tables, 5 hash functions, 3-shingles
lsh.Add("The quick brown fox jumps over the lazy dog")
lsh.Add("The fast brown fox jumps over the lazy cat")

similar := lsh.Query("The quick brown fox leaps over the lazy dog", 0.7)
// Returns indices of similar documents
```

### Wu-Manber Approximate Search

```go
wm := fuzzy.NewWuManber("pattern")
text := "This text contains a patern with an error"
matches := wm.Search(text, 1) // Allow 1 error

for _, match := range matches {
    fmt.Printf("Found at position %d-%d with %d errors\n", 
        match.Start, match.End, match.Distance)
}
```

## Benchmarks

Run benchmarks with:
```bash
go test -bench=. ./...
```

## Algorithm Complexity

| Algorithm | Build Time | Search Time | Space |
|-----------|------------|-------------|-------|
| BK-Tree | O(n log n) | O(log n) | O(n) |
| Suffix Array | O(n log n) | O(log n + m) | O(n) |
| FM-Index | O(n log n) | O(m) | O(n) |
| N-gram Index | O(n) | O(m + k) | O(n) |
| LSH | O(n) | O(1) | O(n) |
| Wu-Manber | O(m) | O(n) | O(m) |

Where:
- n = text/corpus size
- m = pattern size
- k = number of results

## Comparison with Other Libraries

- **sahilm/fuzzy**: Uses a simple character-by-character matching algorithm optimized for editor-style filename matching
- **lithammer/fuzzysearch**: Implements basic fuzzy matching with Levenshtein distance ranking

This library implements advanced algorithms not found in those libraries:
- BK-trees and other metric tree structures
- Suffix arrays and FM-indexes for large text search
- N-gram based indexing
- Locality-sensitive hashing approaches
- Advanced edit distance algorithms like Myers' algorithm and Wu-Manber

## License

MIT