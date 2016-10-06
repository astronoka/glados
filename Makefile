PROJECT_PATH:=github.com/astronoka/glados
CUR_DIR:=$(shell pwd)
BUILD_DIR:=$(CUR_DIR)/_build
BIN_DIR:=$(BUILD_DIR)/bin
BIN_DIR_ON_CONTAINER:=/glados
CONTAINER_FILE_DIR:=$(BUILD_DIR)/container
VENDORING_FILE:=$(BUILD_DIR)/vendoring

BUILD_CONTAINER_FILE:=$(CONTAINER_FILE_DIR)/go-build
BUILD_CONTAINER_IMAGE:=astronoka/go:glados

DOCKER_NETWORK?=bridge
DOCKER_PUBLISH_HOST_PORT?=7000
WITH_BUILD_CONTAINER:=docker run --rm \
	--volume $(CUR_DIR):/go/src/$(PROJECT_PATH) \
	--volume $(BIN_DIR):$(BIN_DIR_ON_CONTAINER) \
	--workdir /go/src/$(PROJECT_PATH) \
	--network $(DOCKER_NETWORK) \
	--publish $(DOCKER_PUBLISH_HOST_PORT):7000 \
	$(BUILD_CONTAINER_IMAGE)

GLADOS_SERVER_CONTAINER:=$(CONTAINER_FILE_DIR)/glados-server
GLADOS_SERVER_BIN:=$(BIN_DIR)/glados-server
GLADOS_SRCS:=$(shell find $(CUR_DIR) -type f -name '*.go' -not -path "$(CUR_DIR)/vendor/*" -not -path "$(CUR_DIR)/.glide/*")

#
# glados-server container
#
$(GLADOS_SERVER_CONTAINER): $(GLADOS_SERVER_BIN)
	docker build -t astronoka/glados -f Dockerfile.glados-server . \
		&& mkdir -p $(@D) \
		&& touch $(GLADOS_SERVER_CONTAINER)

$(GLADOS_SERVER_BIN): $(BUILD_CONTAINER_FILE) $(VENDORING_FILE) $(GLADOS_SRCS)
	mkdir -p $(BIN_DIR)
	$(WITH_BUILD_CONTAINER) go build -v -o $(BIN_DIR_ON_CONTAINER)/glados-server $(PROJECT_PATH)/example/cmd/glados-server

run-glados-server:
	$(WITH_BUILD_CONTAINER) reflex -r '\.go$$' -s -- sh -c 'go run -v -race example/cmd/glados-server/main.go'

#
# vendoring
#
$(VENDORING_FILE): glide.yaml
	$(WITH_BUILD_CONTAINER) glide install -v \
		&& mkdir -p $(@D) \
		&& touch $(VENDORING_FILE)

$(BUILD_CONTAINER_FILE): Dockerfile.build
	docker build -t $(BUILD_CONTAINER_IMAGE) -f Dockerfile.build . \
		&& mkdir -p $(@D) \
		&& touch $(BUILD_CONTAINER_FILE)

clean:
	rm -rf $(BUILD_DIR)

.PHONY: clean run-glados-server
