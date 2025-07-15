generate-mocks:
	@echo "Generating mocks..."
	@go install github.com/golang/mock/mockgen@v1.6.0
	@go generate ./...


test:
	@echo "Running tests..."
	@go test -race -failfast ./...
