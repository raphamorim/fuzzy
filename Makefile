.PHONY: all build test bench clean fmt vet lint

all: build

build:
	go build -v ./...

test:
	go test -v ./...

test-race:
	go test -race -v ./...

bench:
	go test -bench=. -benchmem ./...

bench-all:
	cd benchmark && go test -bench=. -benchmem ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

clean:
	go clean
	rm -f coverage.out coverage.html

deps:
	go mod download
	go mod tidy

update-deps:
	go get -u ./...
	go mod tidy

run-benchmarks:
	cd benchmark && ./run_benchmarks.sh

plot:
	cd benchmark/plot && go run main.go