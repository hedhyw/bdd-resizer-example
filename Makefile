GOLANG_CI_LINT_VER:=v1.45.2

all: lint test run
.PHONY: all

run:
	go run cmd/bdd-resizer-example/main.go
.PHONY: run

test:
	go test \
		-coverpkg ./... \
		-covermode=count \
		-coverprofile=coverage.out \
		./...
	go tool cover -func=coverage.out
.PHONY: test

lint: bin/golangci-lint
	./bin/golangci-lint run --timeout=10m ./...
.PHONY: lint

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor

bin/golangci-lint:
	curl \
		-sSfL \
		https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s $(GOLANG_CI_LINT_VER)
