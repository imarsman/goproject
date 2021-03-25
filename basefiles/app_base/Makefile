# A makefile to build FeedFetcher for Linux, Windows, and MacOs

# http://stackoverflow.com/documentation/go/1020/cross-compilation#t=201703112045220092703

.PHONY: help all
.DEFAULT_GOAL := help

# Various colours for help output
NO_C := \x1b[0m
GREEN := \x1b[32;01m
BLUE := \x1b[34;01m
TEAL := \033[36m
GODOC_PORT := 6060

# Figure out dir for the makefile
MAKEFILE_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Figure out the rood dir for the project
PROJECT_ROOT_DIR:=$(shell cd $(MAKEFILE_DIR); cd ./../..; pwd)

# Format 2020:10:26T03:14:23Z
build_ts=`date -u '+%Y:%m:%dT%H:%M:%SZ'`
build_date=`date +%F`

# Get revision from git
version=`git rev-list --count HEAD`

# The base for building the app
build_base := $(PROJECT_ROOT_DIR)/build/$(appname)

# The dir where builds go
app_builds_base := $(PROJECT_ROOT_DIR)/build/builds/$(appname)

check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

basefiles_dir := $(PROJECT_ROOT_DIR)/bin/basefiles/$(appname)
namespace := test/cmd/$(appname)

# http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
LDFLAGS=-ldflags="\
	-X $(namespace)/internal/build.Version=$(version) \
	-X $(namespace)/internal/build.Build=$(build_ts) \
	-X $(namespace)/internal/build.Platform=$(1) \
	-X $(namespace)/internal/build.Architecture=$(4)"

# Build architecture specific to specific build file name, make it executable,
# then copy base files such as config to the build directory.
build = GOOS=$(1) GOARCH=$(2) go build \
	$(LDFLAGS) \
	-o $(build_base)/bin/$(appname)_$(3)_$(4)$(5) &&\
	chmod 755 $(build_base)/bin/$(appname)_$(3)_$(4)$(5) &&\
	find $(basefiles_dir)/config/ -type f -exec cp {} $(build_base)/config/ \; && \
	find $(basefiles_dir)/log/ -type f -exec cp {} $(build_base)/log/ \; && \
	find $(build_base)/ -type f -name ".gitignore" -delete

godoc: ## Run the godoc server
	@echo "$(GREEN)Running godoc server locally on port $(GODOC_PORT)$(NO_C)"
	@echo "$(TEAL)http://localhost:6060$(NO_C)"
	@echo "$(TEAL)Press CONTROL-C to exit.$(NO_C)"
	godoc -http=localhost:$(GODOC_PORT) -notes 'BUG|TODO' -goroot=. -play

# From https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this message
	@echo "version: $(version)"
	@echo "build ts: $(build_ts)"
	@echo "build date: $(build_date)"
	@echo "$(GREEN)Usage:$(NO_C)"
	@echo "$(TEAL)make [target]$(NO_C)"
	@echo ""
	@echo "$(GREEN)List of targets:$(NO_C)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(TEAL)- %-15s$(NO_C) %s\n", $$1, $$2}'
	@echo
	@echo "$(GREEN)Examples$(NO_C)"
	@echo "$(TEAL)make appname=[appname] all$(NO_C)"
	@echo "$(BLUE)Make darwin (MacOS) and Linux builds:$(NO_C)"
	@echo "$(TEAL)- Builds placed in$(NO_C) ./builds/bin/[architecture]/$(appname)"

make_base_dirs: # Create directories needed later
	@echo "$(GREEN)Making base directories$(NO_C)"
	mkdir -p $(build_base)/ && \
	mkdir -p $(build_base)/bin && \
	mkdir -p $(build_base)/config && \
	mkdir -p $(build_base)/log && \
	mkdir -p $(app_builds_base)

dependencies: ## Update go.mod dependencies in vendor dir
	go mod vendor

package: ## Package into tar.gz zip build, both dated then move
	@echo "$(GREEN)Packaging files and moving them to builds dir for app$(NO_C)"
	cd $(build_base)/.. && \
	tar --exclude '*gz' --exclude '*zip' -zcv -f $(appname)-build-$(build_date).tar.gz $(appname) && \
	zip -x *gz *zip -r $(appname)-build-$(build_date).zip $(appname) && \
	mv $(appname)-build-$(build_date).tar.gz $(app_builds_base)/ && \
	mv $(appname)-build-$(build_date).zip $(app_builds_base)/

clean: ## Clean things up
	@echo "$(GREEN)Cleaning out files in preparation for build$(NO_C)"
	find $(build_base)/config -type f -delete && \
	find $(build_base)/log -type f -delete

check:
	$(call check_defined, appname)

all: check make_base_dirs clean buildall package  ## Build for all supported architectures

buildall: linuxamd64 linuxarm64 darwinamd64 darwinarm64 windows

##### LINUX BUILDS #####
linuxamd64: ## Build linux
	$(call check_defined, appname)
	@echo "$(GREEN)Building Linux$(NO_C)"
	$(call build,linux,amd64,lin,x64,)

##### LINUX BUILDS #####
linuxarm64: ## Build linux
	$(call check_defined, appname)
	@echo "$(GREEN)Building Linux$(NO_C)"
	$(call build,linux,arm64,lin,arm64,)

##### WINDOWS BUILDS #####
windows: ## Build Windows
	$(call check_defined, appname)
	@echo "$(GREEN)Building Windows$(NO_C)"
	$(call build,windows,amd64,win,x64,.exe)

##### DARWIN (MAC) BUILDS #####
darwinamd64: ## Build MacOs (Darwin)
	$(call check_defined, appname)
	@echo "$(GREEN)Building MacOS (Darwin)$(NO_C)"
	$(call build,darwin,amd64,mac,x64,)

##### DARWIN (MAC) BUILDS #####
darwinarm64: ## Build MacOs (Darwin)
	$(call check_defined, appname)
	@echo "$(GREEN)Building MacOS (Darwin)$(NO_C)"
	$(call build,darwin,arm64,mac,arm64,)
