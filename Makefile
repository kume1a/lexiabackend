APP_NAME = lexia
MAIN = ./main.go
BUILD_DIR = bin

include .env.development

test:
	echo ${DB_CONNECTION_STRING}

run:
	air

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

clean:
	rm -rf $(TMP_DIR) $(BUILD_DIR)

docker-dev:
	docker-compose -f docker-compose.development.yaml up

schemagen:
	go generate ./ent

migration-status:
	atlas migrate status \
		--dir "file://ent/migrate/migrations" \
		--url ${DB_CONNECTION_URL}

migration-apply:
	atlas migrate apply \
		--dir "file://ent/migrate/migrations" \
		--url ${DB_CONNECTION_URL}

migration-create:
	atlas migrate diff init \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema" \
		--dev-url "docker://postgres/15/test?search_path=public"

migrate-lint:
	atlas migrate lint \
		--dev-url="docker://postgres/15/test?search_path=public" \
		--dir="file://ent/migrate/migrations" \
		--latest=1

help:
	@echo "Usage:"
	@echo "  make run              Run the backup service"
	@echo "  make build            Build the binary"
	@echo "  make clean            Clean up tmp files and binaries"
	@echo "  make docker-dev       Run the development Docker Compose file"
	@echo "  make test-e2e         Run all end-to-end tests"
	@echo "  make test-e2e-race    Run e2e tests with race detection"
	@echo "  make test-e2e-coverage Run e2e tests with coverage"

# E2E Testing targets
test-e2e:
	@echo "Running all E2E tests..."
	go test -v ./test/e2e/...

test-e2e-race:
	@echo "Running E2E tests with race detection..."
	go test -race -v ./test/e2e/...

test-e2e-coverage:
	@echo "Running E2E tests with coverage..."
	go test -cover -v ./test/e2e/...
