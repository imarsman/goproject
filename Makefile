# A makefile to set up base directory structure for a Go project

# http://stackoverflow.com/documentation/go/1020/cross-compilation#t=201703112045220092703

# Figure out dir for the makefile
MAKEFILE_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: help all
.DEFAULT_GOAL := help

# Format 2020:10:26T03:14:23Z
build_ts=`date -u '+%Y:%m:%dT%H:%M:%SZ'`
build_date=`date +%F`

# Get revision from git
version=`git rev-list --count HEAD`

# Various colours for help output
NO_C := \x1b[0m
GREEN := \x1b[32;01m
BLUE := \x1b[34;01m
TEAL := \033[36m
GODOC_PORT := 6060
# make foo=bar target

# Useful for making build information for use in version output by application
namespace := project.repository.domain/$(appname)/cmd/$(appname)

# http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
LDFLAGS=-ldflags="\
	-X $(namespace)/internal/build.Version=$(version) \
	-X $(namespace)/internal/build.Build=$(build_ts) \
	-X $(namespace)/internal/build.Platform=$(1) \
	-X $(namespace)/internal/build.Architecture=$(4)"


# From https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print this message
	@echo "version: $(version)"
	@echo "build ts: $(build_ts)"
	@echo "build date: $(build_date)"
	@echo "$(GREEN)Usage:$(NO_C)"
	@echo "$(TEAL)make appname=[appname] all$(NO_C)"
	@echo ""
	@echo "$(GREEN)List of targets:$(NO_C)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(TEAL)- %-15s$(NO_C) %s\n", $$1, $$2}'
	@echo
	@echo "$(GREEN)Examples$(NO_C)"
	@echo "$(TEAL)make all$(NO_C)"

# Check that given variables are set and all have non-empty values,
# die with an error otherwise.
#
# Params:
#   1. Variable name(s) to test.
#   2. (optional) Error message to print.
# https://stackoverflow.com/questions/10858261/how-to-abort-makefile-if-variable-not-set
check_defined = \
    $(strip $(foreach 1,$1, \
        $(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),, \
      $(error Undefined $1$(if $2, ($2))))

all: ## Make base directory structure and package in builds directory.
	$(call check_defined, appname)
	mkdir -p $(MAKEFILE_DIR)/bin/test/$(appname)
	mkdir -p $(MAKEFILE_DIR)/bin/test/$(appname)/config
	mkdir -p $(MAKEFILE_DIR)/bin/test/$(appname)/log
	mkdir -p $(MAKEFILE_DIR)/build
	mkdir -p $(MAKEFILE_DIR)/build/bin/basefiles/$(appname)
	mkdir -p $(MAKEFILE_DIR)/build/bin/basefiles/$(appname)/config
	mkdir -p $(MAKEFILE_DIR)/build/bin/basefiles/$(appname)/log
	mkdir -p $(MAKEFILE_DIR)/build/bin/build/$(appname)
	mkdir -p $(MAKEFILE_DIR)/build/bin/build/builds/$(appname)
	mkdir -p $(MAKEFILE_DIR)/build/cmd/$(appname)/internal/build
	mkdir -p $(MAKEFILE_DIR)/build/cmd/$(appname)/internal/common
	cp $(MAKEFILE_DIR)/basefiles/app/Makefile $(MAKEFILE_DIR)/build/cmd/$(appname)/
	cp $(MAKEFILE_DIR)/basefiles/app/internal/build/*go $(MAKEFILE_DIR)/build/cmd/$(appname)/internal/build
	cp $(MAKEFILE_DIR)/basefiles/app/internal/common/*go $(MAKEFILE_DIR)/build/cmd/$(appname)/internal/common
	find $(MAKEFILE_DIR)/build/ -type d -links 2 -exec touch {}/.keep \;
	mkdir -p builds
	cd $(MAKEFILE_DIR)/build && tar -zcvf project-$(appname).tar.gz .
	mv $(MAKEFILE_DIR)/build/project-$(appname).tar.gz $(MAKEFILE_DIR)/builds
	rm -rf $(MAKEFILE_DIR)/build/cmd/$(appname)
