module github.com/raphamorim/fuzzy/benchmark

go 1.21

require (
	github.com/lithammer/fuzzysearch v1.1.8
	github.com/raphamorim/fuzzy v0.0.0
	github.com/sahilm/fuzzy v0.1.1
)

replace github.com/raphamorim/fuzzy => ../

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
