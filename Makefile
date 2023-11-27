.PHONY: setup
setup:
	go mod download

.PHONY: test
test:
	go test -v -cover -race ./...

.PHONY: coverage
coverage:
	go test -v -cover -race ./... -coverprofile=/tmp/cover.out
	go tool cover -html=/tmp/cover.out -o /tmp/cover.html
	open /tmp/cover.html

