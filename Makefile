.DEFAULT_GOAL := check

.PHONY: build
build: generate
	go build -o your_name \
		cmd/main.go

.PHONY: check
check: clean generate lint test

FILES_TO_DELETE = 'mock_*.go' '*_string.go' '*.sql.go' '*.gen.go' 'copyfrom.go'
.PHONY: clean
clean:
	rm -f project_name coverage.out
	$(foreach file, $(FILES_TO_DELETE), find pkg -type f -name $(file) -delete;)

.PHONY: fmt
fmt:
	gofumpt -w cmd/ pkg/ scripts/ tests/ tools/
	goimports -w -local project_name/ cmd/ pkg/ scripts/ tests/ tools/
	goimports -w -local project_name/smt/project_name cmd/ pkg/ scripts/ tests/ tools/

.PHONY: generate
generate: tidy
	go generate ./tools/...
	go generate -run="sqlc|oapi-codegen" ./pkg/...
	go generate ./pkg/...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	@ ! gofumpt -l cmd/ pkg/ scripts/ tests/ tools/ | read || (echo "Bad format, run make fmt" && exit 1)
	go test -coverpkg=./pkg/... -coverprofile=coverage.out -race -timeout 10s ./...

.PHONY: tidy
tidy:
	go mod tidy

