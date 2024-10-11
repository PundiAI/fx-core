#!/usr/bin/make -f

VERSION := $(shell git describe --tags --always 2>/dev/null || echo 'unknown')
COMMIT := $(shell git log -1 --format='%H' 2>/dev/null || echo 'unknown')

LEDGER_ENABLED ?= true
BUILDDIR ?= $(CURDIR)/build
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)
GOPATH ?= '$(HOME)/go'
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

#ifeq (badgerdb,$(findstring badgerdb,$(FX_BUILD_OPTIONS)))
#  build_tags += badgerdb
#endif

#ifeq (cleveldb,$(findstring cleveldb,$(FX_BUILD_OPTIONS)))
#  build_tags += gcc cleveldb muslc
#endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

comma:= ,
empty:=
space:= $(empty) $(empty)
build_tags_comma_sep := $(subst $(space),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep) \
		  -X github.com/cosmos/cosmos-sdk/version.Name=fxcore \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=fxcored

ifeq (,$(findstring nostrip,$(FX_BUILD_OPTIONS)))
  ldflags += -w -s
endif

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags '$(build_tags)' -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(FX_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(FX_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

###############################################################################
###                                  Build                                  ###
###############################################################################

all: build lint test

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify
	go mod tidy
	@echo "--> Download go modules to local cache"
	go mod download

build: go.sum
	go build -mod=readonly -v $(BUILD_FLAGS) -o $(BUILDDIR)/bin/$(BINARYNAME) ./cmd/fxcored
	@echo "--> Done building."

build-win:
	@$(MAKE) build

build-linux:
	@GOOS=linux GOARCH=amd64 $(MAKE) build

INSTALL_DIR := $(shell go env GOPATH)/bin
install: build $(INSTALL_DIR)
	mv $(BUILDDIR)/bin/fxcored $(shell go env GOPATH)/bin/fxcored
	@echo "--> Run \"fxcored start\" or \"$(shell go env GOPATH)/bin/fxcored start\" to launch fxcored."

$(INSTALL_DIR):
	@echo "Folder $(INSTALL_DIR) does not exist"
	mkdir -p $@

docker:
	@echo "--> Building fxcore docker image"
	docker build --progress plain -t ghcr.io/functionx/fx-core:latest .

run-local: install
	@./local-node.sh init

.PHONY: build build-win install docker go.sum run-local

###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_version=v1.60.3

lint-install:
	@echo "--> Installing golangci-lint $(golangci_version)"
	@if golangci-lint version --format json | jq .version | grep -q $(golangci_version); then \
		echo "golangci-lint $(golangci_version) is already installed"; \
	else \
		echo "Installing golangci-lint $(golangci_version)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version); \
	fi

check-no-lint:
	@if [ $$(find . -name '*.go' -type f | xargs grep 'nolint\|#nosec' | wc -l) -ne 29 ]; then \
		echo "\033[91m--> increase or decrease nolint, please recheck them\033[0m"; \
		echo "\033[91m--> list nolint: \`find . -name '*.go' -type f | xargs grep 'nolint\|#nosec'\`\033[0m"; \
		exit 1;\
	fi

lint: check-no-lint lint-install
	@echo "--> Running linter"
	@golangci-lint run --build-tags=$(GO_BUILD) --out-format=tab

format: lint-install
	@golangci-lint run --build-tags=$(GO_BUILD) --out-format=tab --fix

shell-lint:
	# install shellcheck > https://github.com/koalaman/shellcheck
	grep -r '^#!/usr/bin/env bash' --exclude-dir={node_modules,build} . | cut -d: -f1 | xargs shellcheck

shell-format:
	# install shfmt > https://github.com/mvdan/sh
	#go install mvdan.cc/sh/v3/cmd/shfmt@v3.8.0
	grep -r '^#!/usr/bin/env bash' --exclude-dir={node_modules,build} . | cut -d: -f1 | xargs shfmt -l -w -i 2

.PHONY: format lint shell-lint shell-format

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@echo "--> Running tests"
	go test -mod=readonly ./...

test-count:
	go test -mod=readonly -cpu 1 -count 1 -cover ./... | grep -v 'types\|cli\|no test files'

test-nightly:
	@TEST_INTEGRATION=true go test -mod=readonly -timeout 20m -cpu 4 -v -run TestIntegrationTest ./tests
	@TEST_CROSSCHAIN=true go test -mod=readonly -cpu 4 -v -run TestCrosschainKeeperTestSuite ./x/crosschain/...

mocks:
	@go install go.uber.org/mock/mockgen@v0.4.0
	mockgen -source=x/crosschain/types/expected_keepers.go -package mock -destination x/crosschain/mock/expected_keepers_mocks.go

.PHONY: test test-count test-nightly mocks

###############################################################################
###                                Protobuf                                 ###
###############################################################################
protoVer=0.13.0
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=docker run --rm -v $(CURDIR):/workspace --user root --workdir /workspace $(protoImageName)

proto-all: proto-format proto-gen proto-swagger-gen

proto-format:
	@echo "Formatting Protobuf files"
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(protoImage) sh ./scripts/protoc-swagger-gen.sh
	$(MAKE) update-swagger-docs

proto-update-deps:
	@echo "Updating Protobuf dependencies"
	@docker run --rm -v $(CURDIR)/proto:/workspace --workdir /workspace $(protoImageName) buf mod update

.PHONY: proto-format proto-lint proto-gen proto-swagger-gen proto-update-deps

statik: $(STATIK)
$(STATIK):
	@echo "Installing statik..."
	@go install github.com/rakyll/statik@latest

update-swagger-docs: statik
	@if [ "$(shell sed -n '7p' docs/swagger-ui/swagger.yaml)" != "schemes:" ]; then \
		perl -pi -e "print \"host: fx-rest.functionx.io\nschemes:\n  - https\n\" if $$.==6 " ./docs/swagger-ui/swagger.yaml; \
	fi
	@$(GOPATH)/bin/statik -src=docs/swagger-ui -dest=docs -f -m

.PHONY: statik update-swagger-docs

###############################################################################
###                               Contracts                                 ###
###############################################################################

contract-abigen:
	@./contract/compile.sh

contract-publish:
	@./solidity/release.sh

.PHONY: contract-abigen contract-publish

###############################################################################
###                                Releasing                                ###
###############################################################################

PACKAGE_NAME := $(shell go list -m)
GOLANG_CROSS_VERSION := v1.23
release-dry-run:
	docker run --rm --privileged -e CGO_ENABLED=1 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip=validate --skip=publish --snapshot

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run --rm --privileged -e CGO_ENABLED=1 --env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip=validate --release-notes ./release-note.md

.PHONY: release-dry-run release
