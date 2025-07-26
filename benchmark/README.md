# Fuzzy Search Benchmark Suite

Comprehensive benchmarks comparing three Go fuzzy search libraries:
- github.com/sahilm/fuzzy
- github.com/lithammer/fuzzysearch  
- github.com/raphamorim/fuzzy

## Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark patterns
go test -bench=BenchmarkExactMatch -benchmem
go test -bench=BenchmarkTypoMatch -benchmem
go test -bench=BenchmarkMemory -benchmem
go test -bench=BenchmarkConcurrent -benchmem
go test -bench=BenchmarkRealWorld -benchmem

# Run with more iterations for stable results
go test -bench=. -benchtime=10s -benchmem

# Run accuracy tests
go test -v -run TestAccuracy

# Run latency percentile tests
go test -v -run TestLatencyPercentiles

# Generate CPU profile
go test -bench=BenchmarkTypoMatch_Sahilm_Medium -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Generate memory profile
go test -bench=BenchmarkMemory -memprofile=mem.prof
go tool pprof mem.prof
```

## Benchmark Categories

### 1. Exact Matching
Tests performance when searching for exact matches in datasets of varying sizes.

### 2. Typo Matching
Tests fuzzy matching performance with queries containing typos.

### 3. Memory Usage
Measures memory allocation and usage patterns.

### 4. Concurrent Access
Tests thread-safety and performance under concurrent load.

### 5. Real-World Scenarios
- File path fuzzy search
- Autocomplete functionality
- Large-scale text search

### 6. Index Building
Measures the time and memory required to build search indices.

### 7. Scaling Tests
Shows how performance scales with dataset size.

## Dataset Sizes

- Small: 100 items
- Medium: 10,000 items
- Large: 100,000 items
- Huge: 1,000,000 items (used selectively)

## Key Metrics

- **ns/op**: Nanoseconds per operation
- **B/op**: Bytes allocated per operation
- **allocs/op**: Number of allocations per operation
- **P50/P90/P99**: Latency percentiles

## Interpreting Results

### sahilm/fuzzy
- Optimized for editor-style filename matching
- Simple character-by-character algorithm
- Good performance for small to medium datasets
- No index building required

### lithammer/fuzzysearch
- Uses Levenshtein distance for ranking
- Good balance of accuracy and performance
- Suitable for general-purpose fuzzy matching
- No index building required

### raphamorim/fuzzy
- Multiple advanced algorithms for different use cases:
  - **BK-Tree**: Best for spell checking with bounded edit distance
  - **N-gram/Trigram**: Fast partial matching and similarity search
  - **LSH**: Excellent for finding similar documents
  - **Suffix Array**: Powerful substring search
  - **Wu-Manber**: Bit-parallel approximate matching
- Requires index building but provides better query performance
- More memory usage but better scalability

## Recommendations

1. **For simple fuzzy matching without indices**: Use lithammer/fuzzysearch
2. **For filename/path matching in editors**: Use sahilm/fuzzy
3. **For spell checking**: Use raphamorim/fuzzy with BK-Tree
4. **For large-scale similarity search**: Use raphamorim/fuzzy with LSH
5. **For substring search**: Use raphamorim/fuzzy with Suffix Array
6. **For autocomplete**: Use raphamorim/fuzzy with N-gram or Trigram index