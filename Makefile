.PHONY: examples test pkg internal cmd

examples:
	go run ./examples/main.go

test:
	go test ./...