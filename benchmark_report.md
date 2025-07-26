# Fuzzy Search Benchmark Report
Generated on: Sat Jul 26 14:44:32 CEST 2025

## Quick Benchmark Results (5s)
```
goos: darwin
goarch: arm64
pkg: github.com/raphamorim/fuzzy
cpu: Apple M4 Pro
BenchmarkLevenshtein-14            	 2980075	      2034 ns/op	     704 B/op	       2 allocs/op
BenchmarkDamerauLevenshtein-14     	  910986	      6525 ns/op	   18048 B/op	      45 allocs/op
BenchmarkMyers-14                  	35907814	       174.2 ns/op	    1536 B/op	       1 allocs/op
BenchmarkLSHAdd-14                 	 2261422	      2644 ns/op	    1557 B/op	      12 allocs/op
BenchmarkLSHQuery-14               	 1000000	      6681 ns/op	    6352 B/op	      33 allocs/op
BenchmarkSimHashCompute-14         	 7129040	       821.5 ns/op	    1136 B/op	      18 allocs/op
BenchmarkSimHashQuery-14           	 6793328	       872.8 ns/op	    1136 B/op	      18 allocs/op
BenchmarkNGramAdd-14               	 4904584	      1336 ns/op	    3019 B/op	      44 allocs/op
BenchmarkNGramSearch-14            	20176755	       288.9 ns/op	     296 B/op	      11 allocs/op
BenchmarkTrigramSearch-14          	 5352698	      1126 ns/op	    1056 B/op	      39 allocs/op
BenchmarkSuffixArrayBuild-14       	 1956535	      3062 ns/op	    1080 B/op	       3 allocs/op
BenchmarkSuffixArraySearch-14      	143961518	        41.56 ns/op	      16 B/op	       1 allocs/op
BenchmarkFMIndexBuild-14           	  139498	     42450 ns/op	   44616 B/op	      68 allocs/op
BenchmarkFMIndexCount-14           	146757506	        41.01 ns/op	       0 B/op	       0 allocs/op
BenchmarkWuManberExact-14          	33255105	       172.4 ns/op	     424 B/op	       3 allocs/op
BenchmarkWuManberApproximate-14    	 1738185	      3467 ns/op	    3186 B/op	      37 allocs/op
BenchmarkUkkonenAStar-14           	   40221	    149591 ns/op	  163344 B/op	    3627 allocs/op
BenchmarkBitonicSort-14            	 2483631	      2407 ns/op	    6752 B/op	       8 allocs/op
PASS
ok  	github.com/raphamorim/fuzzy	137.895s
```

## Memory Profile Summary
```
File: fuzzy.test
Type: alloc_space
Time: 2025-07-26 14:44:29 CEST
Showing nodes accounting for 1025.05kB, 100% of 1025.05kB total
      flat  flat%   sum%        cum   cum%
     513kB 50.05% 50.05%      513kB 50.05%  runtime.allocm
  512.05kB 49.95%   100%   512.05kB 49.95%  runtime.(*scavengerState).init
         0     0%   100%   512.05kB 49.95%  runtime.bgscavenge
         0     0%   100%      513kB 50.05%  runtime.mstart
         0     0%   100%      513kB 50.05%  runtime.mstart0
         0     0%   100%      513kB 50.05%  runtime.mstart1
         0     0%   100%      513kB 50.05%  runtime.newm
         0     0%   100%      513kB 50.05%  runtime.resetspinning
         0     0%   100%      513kB 50.05%  runtime.schedule
         0     0%   100%      513kB 50.05%  runtime.startm
         0     0%   100%      513kB 50.05%  runtime.wakep
```

## Accuracy Test Results
```
testing: warning: no tests to run
PASS
ok  	github.com/raphamorim/fuzzy	0.275s
```

## Latency Percentiles
```
testing: warning: no tests to run
PASS
ok  	github.com/raphamorim/fuzzy	0.141s
```
