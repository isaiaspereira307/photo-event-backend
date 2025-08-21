.PHONY: default run build clean docs
APP_NAME = "photo-backend"
BIN_DIR = "bin"

default: run

run:
	@swag init --parseDependency --parseInternal
	@go run main.go
build:
	@mkdir -p $(BIN_DIR)
	@swag init --parseDependency --parseInternal
	@go build -o $(BIN_DIR)/$(APP_NAME) main.go
test:
	@go test -v ./...
docs:
	@swag init --parseDependency --parseInternal
clean:
	@rm -f $(APP_NAME)
	@rm -f ./docs/swagger.json