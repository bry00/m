MODULE_ROOT ?= $(shell git rev-parse --show-toplevel)
MODULE = $(shell basename $(MODULE_ROOT))
MODULES = buffers config controller utl view view/tv
SOURCE_DIR = src
GOLANG_MODULES_SOURCES=$(foreach dir,$(addprefix $(SOURCE_DIR)/,$(MODULES)),$(wildcard $(dir)/*.go))
GOLANG_SOURCES=$(wildcard $(SOURCE_DIR)/*.go)

GO111MODULE = on

BUILD_DIR = _target
TARGET_DIRS = $(BUILD_DIR)
INSTALL_DIR = /usr/local/bin

.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
	sed 's/{version}/$(VERSION)/'

.PHONY: build
build: $(BUILD_DIR)/$(MODULE) ## Compile and build the Golang module

.PHONY: compile
compile: $(BUILD_DIR)/$(MODULE) ## An alias to the build

.PHONY: clean
clean: ## Clean build
	@echo "Cleaning $(CURDIR)..." &&\
	rm -rf $(BUILD_DIR) &&\
	echo "Cleaned."

.PHONY: rebuild
rebuild: clean build ## Rebuild the module

.PHONY: install
install: $(INSTALL_DIR)/$(MODULE) ## Install the module

.PHONY: run
run: build ## Test run
	@$(BUILD_DIR)/$(MODULE) Makefile


$(TARGET_DIRS):
	@$(info Creating target directory $@...)
	@mkdir -p $@

$(BUILD_DIR)/$(MODULE): $(GOLANG_SOURCES) $(GOLANG_MODULES_SOURCES) $(TARGET_DIRS)
	@echo "Building module $(MODULE)" &&\
	cd $(SOURCE_DIR) &&\
	go build -o $(abspath $@) $(notdir $(GOLANG_SOURCES)) &&\
	echo "Success: the module build: $@"

$(INSTALL_DIR)/$(MODULE): $(BUILD_DIR)/$(MODULE)
	@sudo cp -v $< $@ &&\
	sudo chmod 755 $@ &&\
	echo "Installed as $@"
