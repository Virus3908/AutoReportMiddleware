
FILES_TO_DELETE = '*.sql.go' '*.pb.go' 
.PHONY: clean
clean:
	@echo "Cleaning up generated files..."
	@for file in $(FILES_TO_DELETE); do \
		if [ -f "$$file" ]; then \
			echo "Deleting $$file..."; \
			rm -f "$$file"; \
		else \
			echo "$$file not found."; \
		fi; \
	done
	@echo "Cleanup complete."


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
