.PHONY: all build clean test test-unit test-integration run

# Default Go build flags
GO_BUILD_FLAGS=-v

# Default Go test flags
GO_TEST_FLAGS=-v

# Binary name
BINARY_NAME=github-api-service

all: test build

build:
	go build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) .

clean:
	go clean
	rm -f $(BINARY_NAME)

test: test-unit test-integration

test-unit:
	go test $(GO_TEST_FLAGS) ./internal/...

test-integration:
	go test $(GO_TEST_FLAGS) -tags=integration .

test-coverage:
	go test $(GO_TEST_FLAGS) -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

run:
	go run .

# Run with race detector
run-race:
	go run -race .

# Tidy dependencies
tidy:
	go mod tidy

# Lint code
lint:
	go vet ./...
	test -z $(gofmt -l .)

# Check for vulnerabilities in dependencies
vuln:
	go list -json -m all | go run golang.org/x/vuln/cmd/govulncheck@latest

# Docker targets
docker-build:
	docker build -t $(BINARY_NAME):latest .

docker-run:
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest