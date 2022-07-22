#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  ifeq ($(OS),Windows_NT)
	VERSION := $(shell git describe --exact-match 2>$null)
  else
	VERSION := $(shell git describe --exact-match 2>/dev/null)
  endif
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

LEDGER_ENABLED ?= true
TM_VERSION := $(shell go list -m -f '{{ .Version }}' github.com/tendermint/tendermint)
BUILDDIR ?= $(CURDIR)/build

export GO111MODULE = on

ifeq ($(OS),Windows_NT)
  BINARYNAME := fxcored.exe
else
  BINARYNAME := fxcored
endif

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

ifneq (,$(FX_BUILD_OPTIONS))
  network=$(FX_BUILD_OPTIONS)
endif
ifeq (devnet,$(network))
  FX_BUILD_OPTIONS := devnet
endif
ifeq (testnet,$(network))
  FX_BUILD_OPTIONS := testnet
endif
ifeq (,$(network))
  FX_BUILD_OPTIONS := mainnet
endif

ifeq (cleveldb,$(findstring cleveldb,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif

ifeq (devnet,$(findstring devnet,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/functionx/fx-core/app.network=devnet
  ldflags += -X github.com/functionx/fx-core/app.ChainID=boonlay
endif

ifeq (testnet,$(findstring testnet,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/functionx/fx-core/app.network=testnet
  ldflags += -X github.com/functionx/fx-core/app.ChainID=dhobyghaut
endif

ifeq (mainnet,$(findstring mainnet,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/functionx/fx-core/app.network=mainnet
  ldflags += -X github.com/functionx/fx-core/app.ChainID=fxcore
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

all: install lint test

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy
	@echo "--> Download go modules to local cache"
	@go mod download

go-build: go.mod
	@go build -mod=readonly -v $(BUILD_FLAGS) -o $(BUILDDIR)/bin/$(BINARYNAME) ./cmd/fxcored

build: go.mod
	@echo "--> build mainnet <--"
	@FX_BUILD_OPTIONS=mainnet make go-build

build-devnet: go.mod
	@echo "--> build devnet <--"
	@FX_BUILD_OPTIONS=devnet make go-build

build-testnet: go.mod
	@echo "--> build testnet <--"
	@FX_BUILD_OPTIONS=testnet make go-build

build-linux:
	@CGO_ENABLED=0 TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 make build

build-linux-devnet:
	@CGO_ENABLED=0 TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 make build-devnet

build-linux-testnet:
	@CGO_ENABLED=0 TARGET_CC=clang LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 make build-testnet

build-win:
	@make go-build

install:
	@$(MAKE) build
	@mv $(BUILDDIR)/bin/fxcored $(GOPATH)/bin/fxcored

install-devnet:
	@$(MAKE) build-devnet
	@mv $(BUILDDIR)/bin/fxcored $(GOPATH)/bin/fxcored

install-testnet:
	@$(MAKE) build-testnet
	@mv $(BUILDDIR)/bin/fxcored $(GOPATH)/bin/fxcored

run-local: install
	@./develop/run_fxcore.sh init

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/functionx/fx-core/cmd/fxcored -d 2 | dot -Tpng -o dependency-graph.png

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@echo "--> Running linter"
	golangci-lint run -v --timeout 3m
	find . -name '*.go' -type f -not -path "./build*" -not -path "*.git*" -not -name '*.pb.*' | xargs gofmt -d -s

format:
	find . -name '*.go' -type f -not -path "./build*" -not -path "*.git*" -not -name '*.pb.*' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./build*" -not -path "*.git*" -not -name '*.pb.*' | xargs misspell -w
	find . -name '*.go' -type f -not -path "./build*" -not -path "*.git*" -not -name '*.pb.*' | xargs goimports -w -local github.com/functionx/fx-core

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

protoVer=v0.2
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(protoVer)
containerProtoGenSwagger=cosmos-sdk-proto-gen-swagger-$(protoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(protoVer)

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./develop/protocgen.sh; fi

###############################################################################
###                                 Other                                  ###
###############################################################################

.PHONY: build build-linux install go.sum format lint clean draw-deps protoc-gen \
	test test-cover test-unit test-race benchmark

