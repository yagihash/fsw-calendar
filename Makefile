.PHONY: setup
setup:
	go mod download

.PHONY: test
test:
	go test -v -cover -race ./...


