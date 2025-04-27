.PHONY: generate

generate:
	go run cmd/generate/main.go examples/simple
	go run cmd/generate/main.go examples/di
	cd examples/simple && go run .
	cd examples/di && go run . 