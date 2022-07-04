.PHONY: build clean tool lint help
GO_CMD=go
GO_BUILD=$(GO_CMD) build -ldflags "-s -w"
GO_BUILD_PLUGIN=$(GO_CMD) build --buildmode=plugin
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_TEST_C=$(GO_CMD) test -c
GO_GET=$(GO_CMD) get
GO_FMT=gofmt
GO_LINT=golint
GO_VET=$(GO_CMD) vet
REMOVE=rm -rf
COPY = cp -rf
TARGET= build
TARGET_BIN= build/faas
TARGET_CONFIG= build/config
MAIN_FILE= manager.go
PLUGIN_PATH=./plugin
PLUGIN_MODULE=$(shell find ./plugin -type f -name "*.go" | xargs grep "package main" | awk -F. '{print $$2}' | sed "s/^\///g")


all: clean tidy build-plugin build

tidy:
	$(GO_CMD) mod tidy
build:
	$(GO_BUILD) -o ${TARGET_BIN} -v ./${MAIN_FILE}

release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o ${TARGET_BIN} -v ./${MAIN_FILE}

build-plugin:
	$(foreach plugin, $(PLUGIN_MODULE), $(GO_BUILD_PLUGIN) -o $(TARGET)/$(plugin).so  $(plugin).go;)
	$(foreach plugin, $(PLUGIN_MODULE), $(COPY) $(TARGET)/$(plugin).so $(plugin).so;)

tool:
	$(GO_VET) ./...
	$(GO_FMT) -w .

lint:
	$(GO_LINT) .

clean:
	@echo ++++ clean start ++++
	$(foreach plugin, $(PLUGIN_MODULE), $(REMOVE) ${plugin}.so;)
	$(REMOVE) ${TARGET}
	$(GO_CLEAN) -i .
	@echo ---- clean end ----

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint ./..."
	@echo "make clean: remove object files and cached files"