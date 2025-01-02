BIN_DIR=$(PWD)/bin

bin-deps:
	GOBIN=$(BIN_DIR) go install github.com/golang/mock/mockgen@latest
test:
	go test ./...
mocks: bin-deps
	$(BIN_DIR)/mockgen -package=mock_calendar -destination ./internal/mocks/mock_calendar/generator.go --source=internal/clients/calendar/client.go Client
	$(BIN_DIR)/mockgen -package=mock_drive -destination ./internal/mocks/mock_drive/generator.go --source=internal/clients/drive/client.go Client
	$(BIN_DIR)/mockgen -package=mock_mm -destination ./internal/mocks/mock_mm/generator.go --source=internal/clients/mm/client.go Client


