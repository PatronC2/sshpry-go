# Variables
BINARY_NAME=sshpry
OUTPUT_DIR=./bin
DOCKER_IMAGE=patron-sshpry

.PHONY: build extract clean

# Build the Go binary inside the container
build:
	docker build --target=build -t $(DOCKER_IMAGE) .

# Extract the built binary from the container to the host system
extract:
	mkdir -p $(OUTPUT_DIR)
	docker create --name $(TEST_CONTAINER) $(DOCKER_IMAGE)
	docker cp $(TEST_CONTAINER):/app/$(BINARY_NAME) $(OUTPUT_DIR)/
	docker rm $(TEST_CONTAINER)
	chmod +x $(OUTPUT_DIR)/$(BINARY_NAME)

# Run all steps: test-build, test, build, extract
all: build extract

# Clean up Docker images and built binary
clean:
	docker rmi -f $(DOCKER_IMAGE) || true
	rm -rf $(OUTPUT_DIR)