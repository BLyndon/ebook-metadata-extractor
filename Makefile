# Constants
BINARY_NAME=./bin/ebook-metadata-extractor-be
IMAGE_NAME=blyndon/ebook-metadata-extractor-be
IMAGE_TAG=latest
IMAGE=$(IMAGE_NAME):$(IMAGE_TAG)
CONTAINER_NAME=ebook-metadata-extractor-be-container
CURRENT_DIR=$(shell pwd)

# Build and Run Commands
build:
	@go build -o $(BINARY_NAME) ./cmd/ebook-metadata-extractor-be

run: build
	@$(BINARY_NAME)

stop:
	@pkill -f $(BINARY_NAME) || echo "Process not found."

reload: stop run

# Docker Commands
docker-build:
	@docker build -t $(IMAGE) .

docker-login:
	@echo "Logging in to Docker registry"
	@docker login $(shell echo $(IMAGE_NAME) | cut -d'/' -f1)

docker-push: docker-build
	@echo "Pushing Docker image to registry"
	@docker push $(IMAGE)

docker-run: docker-build
	@docker run -p 8080:8080 --name $(CONTAINER_NAME) -v $(CURRENT_DIR)/data:/data -e OPENAI_API_KEY=$(OPENAI_API_KEY) -d $(IMAGE)

docker-clean:
	@docker stop $(CONTAINER_NAME) || echo "Container not running"
	@docker rm $(CONTAINER_NAME) || echo "Container not found"

docker-rmi:
	@docker rmi $(IMAGE)

clean:
	@rm $(BINARY_NAME)
