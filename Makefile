##############################
#### BEGIN GENERAL PART ######


BASEPATH = $(shell pwd)
LOCALBIN = $(BASEPATH)/bin
PATH := $(LOCALBIN):$(PATH)
SHELL := env PATH=$(PATH) /bin/bash

# GIT_TAG do not set it manually
# GIT_TAG    ?= $$(git describe --tags)
GIT_BRANCH ?= $$(git rev-parse --abbrev-ref HEAD)
GIT_COMMIT ?= $$(git rev-parse --short HEAD)

BUILD_DATE = $$(date -u +'%FT%T.%NZ')

# GO
GOCMD       = go
GOTEST      = gotestsum


# all src packages without generated code
PKGS = $(shell go list ./...)

# Colors
BLUE_COLOR    = "\033[0;34m"
DEFAULT_COLOR = "\033[m"




install-test:
	@echo -e $(BLUE_COLOR)[install-test]$(DEFAULT_COLOR)
	@$(ISTALLTEST)

test: install-test
	@echo -e $(BLUE_COLOR)[test]$(DEFAULT_COLOR)
	@go clean -testcache
	@mkdir -p report
	@$(GOTEST) --junitfile report/test-junit.xml --format testname --jsonfile report/test.json -- $(GOBUILDFLAG) -race -coverprofile=report/coverage.out ./... 
	@gocover-cobertura < report/coverage.out > report/coverage.xml


#### END GENERAL PART ########
############################## 


BINNAME     = "./homeTask"
GOOSARCH    = GOOS=linux GOARCH=amd64


# GO
GOBUILD     = GOPRIVATE=$(GOPRIVATE) GOINSECURE=$(GOINSECURE) $(GOOSARCH) $(GOCMD) build $(GOBUILDFLAG)
GOINSTALL   = GOPATH=$(GOPATH) GOBIN=$(LOCALBIN) $(GOCMD) install
GOGENERATE  = $(GOCMD) generate
GORUN       = $(GOCMD) run --race

# INSTALL
ISTALLTEST      	=   $(GOINSTALL) github.com/boumenot/gocover-cobertura@$(VCOUBERTURA) && \
						$(GOINSTALL) gotest.tools/gotestsum@$(VGOTESTSUM)
INSTALLALL  		= 	$(ISTALLTEST) && \
						$(GOINSTALL) github.com/golang/mock/mockgen@$(VMOCKGEN) && \
						$(GOINSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@$(VPROTOCGENGO) && \
						$(GOINSTALL) google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(VPROTOCGENGOGRPC)
					


# Versions

VCOUBERTURA        = "v1.2.0"
VMOCKGEN           = "v1.6.0"
VPROTOCGENGO       = "v1.27.1"
VPROTOCGENGOGRPC   = "v1.1"
VGOTESTSUM         = "v1.8.2"
# Colors
PURPLE_COLOR   = "\033[0;35m"

all: generate ci

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Help screen.'
	@echo '    test               Run unit tests.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    generate           Perform go generate'
	@echo ''

install: 
	@echo -e $(BLUE_COLOR)[install]$(DEFAULT_COLOR)
	@$(INSTALLALL)


build: 
	@echo -e $(BLUE_COLOR)[build]$(DEFAULT_COLOR)
	@CGO_ENABLED=0 $(GOBUILD) -ldflags=$(LDFLAGS) -o . ./...

generate: install
	@echo -e $(PURPLE_COLOR)[generate]$(DEFAULT_COLOR)
	@$(GOGENERATE) $(PKGS)

run: generate build
	@echo -e $(PURPLE_COLOR)[run]$(DEFAULT_COLOR)
	@LOGGER_LEVEL=INFO $(BINNAME)

