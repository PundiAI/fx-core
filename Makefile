#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

LEDGER_ENABLED ?= true
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::')
BUILDDIR ?= $(CURDIR)/build

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(FX_BUILD_OPTIONS)))
  build_tags += gcc cleveldb muslc

endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
BUILD_TAGS_COMMA_SEP := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(BUILD_TAGS_COMMA_SEP)" \
		  -X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Name=fxcore \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=fxcored \

ifeq (cleveldb,$(findstring cleveldb,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif

ifeq (devnet,$(findstring devnet,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/functionx/fx-core/app.network=devnet
endif

ifeq (testnet,$(findstring testnet,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/functionx/fx-core/app.network=testnet
endif

ifeq (mainnet,$(findstring mainnet,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/functionx/fx-core/app.network=mainnet
endif

ifeq (,$(findstring nostrip,$(FX_BUILD_OPTIONS)))
  ldflags += -w -s
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(FX_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

###############################################################################
###                              Documentation                              ###
###############################################################################

build:
	@go build -mod=readonly -v $(BUILD_FLAGS) -o $(BUILDDIR)/bin/fxcored ./cmd/fxcored

build-devnet:
	@echo "--> build devnet <--"
	@FX_BUILD_OPTIONS=devnet make build

build-testnet:
	@echo "--> build testnet <--"
	@echo "replace cosmos-sdk to github.com/cosmos/cosmos-sdk=github.com/functionx/cosmos-sdk@v0.42.5-0.20211015120647-6c0e91f2e952"
	@go mod edit --replace=github.com/cosmos/cosmos-sdk=github.com/functionx/cosmos-sdk@v0.42.5-0.20211015120647-6c0e91f2e952
	@go mod tidy -v
	@FX_BUILD_OPTIONS=testnet make build
	@echo "recover cosmos-sdk to github.com/cosmos/cosmos-sdk=github.com/functionx/cosmos-sdk@v0.42.5-0.20210927070625-89306d0caf62"
	@go mod edit --replace=github.com/cosmos/cosmos-sdk=github.com/functionx/cosmos-sdk@v0.42.5-0.20210927070625-89306d0caf62
	@go mod tidy -v

build-mainnet:
	@echo "--> build mainnet <--"
	@FX_BUILD_OPTIONS=mainnet make build

build-linux-devnet:
	@TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 make build-devnet

build-linux-testnet:
	@TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 make build-testnet

build-linux:
	@TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 make build-mainnet

install-devnet:
	@$(MAKE) build-devnet
	@mv $(BUILDDIR)/bin/fxcored $(GOPATH)/bin/fxcored

install-testnet:
	@$(MAKE) build-testnet
	@mv $(BUILDDIR)/bin/fxcored $(GOPATH)/bin/fxcored

install:
	@$(MAKE) build-mainnet
	@mv $(BUILDDIR)/bin/fxcored $(GOPATH)/bin/fxcored

docker-devnet: build-linux-devnet
	@docker build --no-cache -f ./cmd/fxcored/Dockerfile -t functionx/fx-core:latest .
	@docker tag functionx/fx-core:latest functionx/fx-core:dev

docker-testnet: build-linux-testnet
	@docker build --no-cache -f ./cmd/fxcored/Dockerfile -t functionx/fx-core:testnet-1.0 .

docker: build-linux
	@docker build --no-cache -f ./cmd/fxcored/Dockerfile -t functionx/fx-core:mainnet-1.0 .

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy
	@echo "--> Download go modules to local cache"
	@go mod download

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/functionx/fx-core/cmd/fxcored -d 2 | dot -Tpng -o dependency-graph.png

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@echo "--> Running linter"
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.*' | xargs gofmt -d -s

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.*' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.*' | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.*' | xargs goimports -w -local github.com/functionx/fx-core

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	go test -mod=readonly -short $(shell go list ./...)

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -short -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -short -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -short -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

benchmark:
	@go test -mod=readonly -bench=. ./...

###############################################################################
###                                Protobuf                                 ###
###############################################################################

# The below include contains the tools target.
include develop/devtools.mk

proto-gen:
	@echo "Generating Protobuf files"
	@./develop/protocgen.sh

###############################################################################
###                                 Other                                  ###
###############################################################################

.PHONY: build build-linux install go.sum format lint clean draw-deps protoc-gen \
	test test-cover test-unit test-race benchmark

