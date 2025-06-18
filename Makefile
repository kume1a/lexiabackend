APP_NAME=lexia
MAIN=./main.go
BUILD_DIR=bin

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

help:
	@echo "Usage:"
	@echo "  make run        Run the backup service"
	@echo "  make build      Build the binary"
	@echo "  make clean      Clean up tmp files and binaries"
	@echo "  make docker-dev Run the development Docker Compose file"
