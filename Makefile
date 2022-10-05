DOCKER:=docker

CONTOSO_DIR?=./contoso
AVAILABILITY_DIR?=./wherebuy/availability
PRODUCT_DIR?=./wherebuy/product

DOCKER_REGISTRY?=ghcr.io/shubham1172/daprcon
DOCKER_TAG?=latest

docker-build:
	$(info Building Docker images...)
	$(DOCKER) build -t $(DOCKER_REGISTRY)/contoso:$(DOCKER_TAG) -f $(CONTOSO_DIR)/Dockerfile $(CONTOSO_DIR)
	$(DOCKER) build -t $(DOCKER_REGISTRY)/availability:$(DOCKER_TAG) -f $(AVAILABILITY_DIR)/Dockerfile $(AVAILABILITY_DIR)

docker-push: docker-build
	$(info Pushing Docker images...)
	$(DOCKER) push $(DOCKER_REGISTRY)/contoso:$(DOCKER_TAG)
	$(DOCKER) push $(DOCKER_REGISTRY)/availability:$(DOCKER_TAG)