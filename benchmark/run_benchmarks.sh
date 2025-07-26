#!/bin/bash

# Generate test data if needed
if [ ! -d "testdata" ]; then
    mkdir -p testdata
    
    # Try to copy system dictionary
    if [ -f "/usr/share/dict/words" ]; then
        cp /usr/share/dict/words testdata/words.txt
    elif [ -f "/usr/dict/words" ]; then
        cp /usr/dict/words testdata/words.txt
    else
        echo "No system dictionary found. Using generated data."
    fi
fi

echo "Running Fuzzy Search Benchmark Suite"
echo "===================================="
echo

# Quick benchmarks
echo "1. Running quick benchmarks (5s each)..."
go test -bench=. -benchtime=5s -benchmem -run=^$ | tee results_quick.txt

# Detailed benchmarks
echo
echo "2. Running detailed benchmarks (10s each)..."
go test -bench=. -benchtime=10s -benchmem -run=^$ | tee results_detailed.txt

# Memory profiling
echo
echo "3. Running memory profiling..."
go test -bench=BenchmarkMemory -benchmem -memprofile=mem.prof -run=^$
go tool pprof -text mem.prof > memory_profile.txt

# CPU profiling for hot paths
echo
echo "4. Running CPU profiling..."
go test -bench=BenchmarkTypoMatch_Sahilm_Medium -benchtime=10s -cpuprofile=cpu_sahilm.prof -run=^$
go test -bench=BenchmarkTypoMatch_Lithammer_Medium -benchtime=10s -cpuprofile=cpu_lithammer.prof -run=^$
go test -bench=BenchmarkTypoMatch_Raphamorim_BKTree -benchtime=10s -cpuprofile=cpu_raphamorim.prof -run=^$

# Accuracy tests
echo
echo "5. Running accuracy tests..."
go test -v -run TestAccuracy | tee accuracy_results.txt

# Latency percentiles
echo
echo "6. Running latency percentile tests..."
go test -v -run TestLatencyPercentiles | tee latency_results.txt

# Generate summary report
echo
echo "7. Generating summary report..."
cat > benchmark_report.md << EOF
# Fuzzy Search Benchmark Report
Generated on: $(date)

## Quick Benchmark Results (5s)
\`\`\`
$(tail -n 50 results_quick.txt)
\`\`\`

## Memory Profile Summary
\`\`\`
$(head -n 20 memory_profile.txt)
\`\`\`

## Accuracy Test Results
\`\`\`
$(cat accuracy_results.txt)
\`\`\`

## Latency Percentiles
\`\`\`
$(cat latency_results.txt)
\`\`\`
EOF

echo
echo "Benchmark suite completed!"
echo "Results saved to:"
echo "  - results_quick.txt"
echo "  - results_detailed.txt"
echo "  - memory_profile.txt"
echo "  - accuracy_results.txt"
echo "  - latency_results.txt"
echo "  - benchmark_report.md"