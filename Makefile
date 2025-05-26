APP_NAME := middleware
FILES_TO_DELETE := '*.sql.go' '*.pb.go'
ENV_FILE := .env

.PHONY: all
all: build

.PHONY: generate
generate: tidy
	protoc \
		--go_out=pkg/messages \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/messages \
		--go-grpc_opt=paths=source_relative \
		proto/messages.proto
	sqlc generate --file internal/repositories/sqlc.yaml

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run: generate
	go run ./cmd/app/main.go

.PHONY: build
build: generate
	go build -o bin/$(APP_NAME) ./cmd/app

.PHONY: build-docker
build-docker: generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/middleware ./cmd/app/main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: fmt vet

.PHONY: clean
clean:
	@echo "Cleaning up generated files..."
	@for file in $(FILES_TO_DELETE); do \
		find . -name "$$file" -delete -print; \
	rm -f bin/$(APP_NAME); \
	done
	@echo "Cleanup complete."

.PHONY: env
env:
	@echo "Loading env vars from $(ENV_FILE)..."
	set -a; source $(ENV_FILE); set +a

.PHONY: docker-build
docker-build:
	docker build -t middleware-app .

.PHONY: docker-run
docker-run:
	docker run --rm -it \
		--env-file .env \
		middleware-app