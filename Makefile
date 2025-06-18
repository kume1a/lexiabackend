APP_NAME = lexia
MAIN = ./main.go
BUILD_DIR = bin

include .env.development

test:
	echo ${DB_CONNECTION_STRING}

run:
	go run $(MAIN)

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
	atlas migrate diff migration_name \
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
	@echo "  make run        Run the backup service"
	@echo "  make build      Build the binary"
	@echo "  make clean      Clean up tmp files and binaries"
	@echo "  make docker-dev Run the development Docker Compose file"
