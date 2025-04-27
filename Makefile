.PHONY: generate build

generate: build
	go run cmd/ctxboot/main.go examples/simple
	go run cmd/ctxboot/main.go examples/di
	cd examples/simple && go run .
	cd examples/di && go run . 