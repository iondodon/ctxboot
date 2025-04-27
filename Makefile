.PHONY: generate build

build:
	go build -o bin/ctxboot cmd/generate/main.go

generate: build
	go run cmd/generate/main.go examples/simple
	go run cmd/generate/main.go examples/di
	cd examples/simple && go run .
	cd examples/di && go run . 