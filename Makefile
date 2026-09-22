.PHONY: fmt tidy vet lint test test-race cover clean run build migrate-status migrate-up migrate-down migrate-reset migrate-create

fmt:
	go fmt ./...
	
tidy:
	go mod tidy
	
vet:
	go vet ./...
	
lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint tidak terinstall"
	
test:
	go test ./...
	
test-race:
	go test -race ./...
	
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	
clean:
	rm -f coverage.out coverage.html
	go clean -testcache
	
build:
	go build ./...
