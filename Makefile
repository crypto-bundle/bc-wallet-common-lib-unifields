default: lint

lint:
	golangci-lint run --config .golangci.yml -v ./...

test:
	go test -v -race -cover ./...

.PHONY: lint test