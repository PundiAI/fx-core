#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
TM_VERSION := $(shell go list -m -f '{{ .Version }}' github.com/tendermint/tendermint)

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
BUILDDIR ?= $(CURDIR)/build
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
STATIK = $(GOPATH)/bin/statik

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

#ifeq (cleveldb,$(findstring cleveldb,$(FX_BUILD_OPTIONS)))
#  build_tags += gcc cleveldb muslc
#endif
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

ifeq (cleveldb,$(findstring cleveldb,$(FX_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
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
###                                  Build                                  ###
###############################################################################

all: build lint test

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy
	@echo "--> Download go modules to local cache"
	@go mod download

build: go.mod
	@go build -mod=readonly -v $(BUILD_FLAGS) -o $(BUILDDIR)/bin/$(BINARYNAME) ./cmd

build-win:
	@$(MAKE) build

build-linux:
	@GOOS=linux GOARCH=amd64 $(MAKE) build

install:
	$(MAKE) build
	@mv $(BUILDDIR)/bin/fxcored $(GOPATH)/bin/fxcored

docker:
	@docker build --no-cache -f Dockerfile -t functionx/fx-core:latest .

run-local: install
	@./develop/run_fxcore.sh init

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/functionx/fx-core/app -d 2 | dot -Tpng -o dependency-graph.png

.PHONY: build build-win install docker go.sum run-local draw-deps

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@echo "--> Running linter"
	golangci-lint run -v --timeout 5m
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.*' -not -name "statik.go" | xargs gofmt -d -s

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./build*" -not -path "*.git*" -not -name '*.pb.*' -not -name "statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./build*" -not -path "*.git*" -not -name '*.pb.*' -not -name "statik.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "./build*" -not -path "*.git*" -not -name '*.pb.*' -not -name "statik.go" | xargs goimports -w -local github.com/functionx/fx-core

.PHONY: format lint

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@go test -mod=readonly -cover -short $(shell go list ./...)

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -short -tags='ledger test_ledger_mock' ./...

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race -short -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -race -short -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

benchmark:
	@go test -mod=readonly -short -bench=. ./...

.PHONY: test test-cover test-unit test-race benchmark

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=v0.2
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(protoVer)
containerProtoGenSwagger=cosmos-sdk-proto-gen-swagger-$(protoVer)
containerProtoFmt=cosmos-sdk-proto-fmt-$(protoVer)

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./ -name "*.proto" -exec clang-format -i {} \; ; fi

proto-gen: proto-format
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./develop/protocgen.sh; fi

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./develop/protoc-swagger-gen.sh; fi

.PHONY: proto-format proto-gen proto-swagger-gen

# Install the runsim binary with a temporary workaround of entering an outside
# directory as the "go get" command ignores the -mod option and will polute the
# go.{mod, sum} files.
#
# ref: https://github.com/golang/go/issues/30515
statik: $(STATIK)
$(STATIK):
	@echo "Installing statik..."
	@(cd /tmp && go install github.com/rakyll/statik@latest)

update-swagger-docs: statik
	$(GOPATH)/bin/statik -src=docs/swagger-ui -dest=docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "\033[91mSwagger docs are out of sync!!!\033[0m";\
        exit 1;\
    else \
        echo "\033[92mSwagger docs are in sync\033[0m";\
    fi

.PHONY: statik update-swagger-docs

###############################################################################
###                                Releasing                                ###
###############################################################################

PACKAGE_NAME:=github.com/functionx/fx-core
GOLANG_CROSS_VERSION  = v1.18
GOPATH ?= '$(HOME)/go'
release-dry-run:
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist --skip-validate

.PHONY: release-dry-run release