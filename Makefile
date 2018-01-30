# some metadata
BUILD_DIR=.build
PWD=$(shell pwd)
BUILD_TIME=$(shell date)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

# information about the go package to build
GO_PACKAGE=github.com/faryon93/hlswatch
GO_ARTIFACT=rootfs/usr/sbin/hlswatch
GO_LD_FLAGS=-X "main.BUILD_TIME=$(BUILD_TIME)" -X "main.GIT_COMMIT=$(GIT_COMMIT)"
GO_SRC=src

# information about the docker image
DOCKER_IMAGE=faryon93/nginx-hls:latest

all: hlswatch
	docker build -t $(DOCKER_IMAGE) .

hlswatch:
	# setup go build environment
	if [ ! -d $(PWD)/$(BUILD_DIR) ]; then \
		mkdir -p $(BUILD_DIR)/src/$(GO_PACKAGE); \
		rm -r $(BUILD_DIR)/src/$(GO_PACKAGE); \
		mkdir -p $(BUILD_DIR)/pkg; \
		ln -s $(PWD)/$(GO_SRC) $(BUILD_DIR)/src/$(GO_PACKAGE); \
    fi

	# get the dependencies
	GOPATH=$(PWD)/$(BUILD_DIR) go get $(GO_PACKAGE)

	# build a statically linked golang binary
	GOPATH=$(PWD)/$(BUILD_DIR) go build \
									-tags netgo -a \
									-ldflags '$(GO_LD_FLAGS)' \
									-o $(GO_ARTIFACT) \
									$(GO_PACKAGE)

clean:
	rm -rf $(BUILD_DIR)/