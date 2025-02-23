BINARY_NAME=sshpry
OUTPUT_DIR=./bin
DOCKER_IMAGE=patron-sshpry
CONTAINER_NAME=sshpry_container

.PHONY: build extract clean all

# Build the Go binary inside the container
build:
	docker build -t $(DOCKER_IMAGE) .

# Extract the built binary from the container to the host system
extract:
	mkdir -p $(OUTPUT_DIR)
	# Create a container from the image without starting it
	docker create --name $(CONTAINER_NAME) $(DOCKER_IMAGE)
	# Copy the binary from the container's /app directory to the host's bin folder
	docker cp $(CONTAINER_NAME):/app/$(BINARY_NAME) $(OUTPUT_DIR)/
	# Remove the temporary container
	docker rm $(CONTAINER_NAME)
	# Make sure the binary is executable
	chmod +x $(OUTPUT_DIR)/$(BINARY_NAME)

# Run all steps
all: build extract

# Clean up Docker images and built binary
clean:
	docker rmi -f $(DOCKER_IMAGE) || true
	rm -rf $(OUTPUT_DIR)
