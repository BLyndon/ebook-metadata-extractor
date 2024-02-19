BINARY_NAME=./bin/extractor
DOCKER_IMAGE_NAME=extractor-image
DOCKER_CONTAINER_NAME=extractor-container
CURRENT_DIR=$(shell pwd)
DOCKER_REGISTRY=blyndon
IMAGE_NAME=extractor
IMAGE_TAG=latest
IMAGE=$(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
K8S_DEPLOYMENT_FILE=k8s/deployment.yml
K8S_SERVICE_FILE=k8s/service.yml
K8S_INGRESS_FILE=k8s/ingress.yml
NAMESPACE=default

build:
	@go build -o $(BINARY_NAME) ./cmd/extractor

run: build
	@$(BINARY_NAME)

stop:
	@pkill -f $(BINARY_NAME) || echo "Process not found."

reload: stop run

docker-build:
	@docker build -t $(DOCKER_IMAGE_NAME) .

docker-login:
	@echo "Logging in to Docker registry $(DOCKER_REGISTRY)"
	@docker login $(DOCKER_REGISTRY)

docker-push: docker-build
	@echo "Pushing Docker image to registry"
	@docker tag $(DOCKER_IMAGE_NAME) $(IMAGE)
	@docker push $(IMAGE)

docker-run: docker-build
	@docker run -p 8080:8080 --name $(DOCKER_CONTAINER_NAME) -v $(CURRENT_DIR)/data:/data -e OPENAI_API_KEY=$(OPENAI_API_KEY) -d $(DOCKER_IMAGE_NAME)

docker-clean:
	@docker stop $(DOCKER_CONTAINER_NAME)
	@docker rm $(DOCKER_CONTAINER_NAME)

docker-rmi:
	@docker rmi $(DOCKER_IMAGE_NAME)

clean:
	@rm $(BINARY_NAME)

create-secret:
	@kubectl get secret openai-api-secret -n $(NAMESPACE) >/dev/null 2>&1 || \
    kubectl create secret generic openai-api-secret --namespace $(NAMESPACE) --from-literal=api_key=$(OPENAI_API_KEY)
	@kubectl patch secret openai-api-secret -n $(NAMESPACE) -p "{\"data\":{\"api_key\":\"$(shell echo -n $(OPENAI_API_KEY) | base64)\"}}"


k8s-deploy: create-secret docker-push
	@echo "Deploying to Kubernetes"
	@kubectl apply -f $(K8S_DEPLOYMENT_FILE) -n $(NAMESPACE)
	@kubectl apply -f $(K8S_SERVICE_FILE) -n $(NAMESPACE)
	@kubectl apply -f $(K8S_INGRESS_FILE) -n $(NAMESPACE)

k8s-delete:
	@echo "Deleting deployment from Kubernetes"
	@kubectl delete -f $(K8S_DEPLOYMENT_FILE) -n $(NAMESPACE)
	@kubectl delete -f $(K8S_SERVICE_FILE) -n $(NAMESPACE)
	@kubectl delete -f $(K8S_INGRESS_FILE) -n $(NAMESPACE)
