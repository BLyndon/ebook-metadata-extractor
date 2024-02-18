BINARY_NAME=./bin/extractor

DOCKER_IMAGE_NAME=extractor-image

DOCKER_CONTAINER_NAME=extractor-container

CURRENT_DIR=$(shell pwd)

build:
	@go build -o $(BINARY_NAME) ./cmd/extractor

run: build
	@$(BINARY_NAME)

stop:
	@pkill -f $(BINARY_NAME) || echo "Process not found."

reload: stop run

docker-build:
	@docker build -t $(DOCKER_IMAGE_NAME) .

docker-run: docker-build
	@docker run -p 8080:8080 --name $(DOCKER_CONTAINER_NAME) -v ${CURRENT_DIR}/data/:/data/ -e OPENAI_API_KEY=$(OPENAI_API_KEY) -d $(DOCKER_IMAGE_NAME)

docker-clean:
	@docker stop $(DOCKER_CONTAINER_NAME)
	@docker rm $(DOCKER_CONTAINER_NAME)

docker-rmi:
	@docker rmi $(DOCKER_IMAGE_NAME)

clean:
	@rm $(BINARY_NAME)
