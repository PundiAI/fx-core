#!/usr/bin/env bash

set -eo pipefail

export GOPROXY=goproxy.cn

protoc_gen_gocosmos() {
  if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null; then
    echo -e "\tPlease run this command from somewhere inside the cosmos-sdk folder."
    return 1
  fi

  go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
}

protoc_gen_gocosmos

if [ ! -f ./build/cosmos-sdk/README.md ]; then
  mkdir -p build
  commit_hash=$(go list -m -f '{{.Version}}' github.com/cosmos/cosmos-sdk)
    if [ ! -f "./build/cosmos-proto.zip" ]; then
      wget -c "https://github.com/cosmos/cosmos-sdk/archive/$commit_hash.zip" -O "./build/cosmos-proto.zip"
    fi
  (
    cd build
    unzip -q -o "./cosmos-proto.zip"
    for dir in *; do
      if [[ -d $dir && "$dir" == "cosmos-sdk-"* ]]; then
        mv "./$dir" cosmos-sdk
      fi
    done
  )
fi

if [ ! -f ./build/ibc-go/README.md ]; then
  mkdir -p build
  commit_hash=$(go list -m -f '{{.Version}}' github.com/cosmos/ibc-go/v3)
    if [ ! -f "./build/ibc-proto.zip" ]; then
      wget -c "https://github.com/cosmos/ibc-go/archive/$commit_hash.zip" -O "./build/ibc-proto.zip"
    fi
  (
    cd build
    unzip -q -o "./ibc-proto.zip"
    for dir in *; do
      if [[ -d $dir && "$dir" == "ibc-go-"* ]]; then
        mv "./$dir" ibc-go
      fi
    done
  )
    rm -rf ./build/ibc-go/proto/ibc/applications/transfer
fi

if [ ! -f ./build/ethermint/README.md ]; then
  mkdir -p build
  commit_hash=$(go list -m -f '{{.Replace.Version}}' github.com/evmos/ethermint | awk -F- '{print $3}')
    if [ ! -f "./build/ethermint-proto.zip" ]; then
      wget -c "https://github.com/evmos/ethermint/archive/$commit_hash.zip" -O "./build/ethermint-proto.zip"
    fi
  (
    cd build
    unzip -q -o "./ethermint-proto.zip"
    for dir in *; do
      if [[ -d $dir && "$dir" == "ethermint-"* ]]; then
        mv "./$dir" ethermint
      fi
    done
  )
fi

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  buf protoc \
    -I "proto" \
    -I "build/ibc-go/proto" \
    -I "build/cosmos-sdk/proto" \
    -I "build/cosmos-sdk/third_party/proto" \
    --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
    --grpc-gateway_out=logtostderr=true:. \
    $(find "${dir}" -maxdepth 1 -name '*.proto')
done

# command to generate docs using protoc-gen-doc
buf protoc \
  -I "proto" \
  -I "build/ibc-go/proto" \
  -I "build/cosmos-sdk/proto" \
  -I "build/cosmos-sdk/third_party/proto" \
  --doc_out=./docs/proto \
  --doc_opt=./docs/proto/proto-doc-markdown.tmpl,fx-proto-docs.md \
  $(find "$(pwd)/proto" -maxdepth 5 -name '*.proto')

buf protoc \
  -I "build/cosmos-sdk/proto" \
  -I "build/cosmos-sdk/third_party/proto" \
  --doc_out=./docs/proto \
  --doc_opt=./docs/proto/proto-doc-markdown.tmpl,cosmos-proto-docs.md \
  $(find "$(pwd)/build/cosmos-sdk/proto" -maxdepth 5 -name '*.proto')

buf protoc \
  -I "build/ethermint/proto" \
  -I "build/ethermint/third_party/proto" \
  --doc_out=./docs/proto \
  --doc_opt=./docs/proto/proto-doc-markdown.tmpl,ethermint-proto-docs.md \
  $(find "$(pwd)/build/ethermint/proto" -maxdepth 5 -name '*.proto')

cp -r github.com/functionx/fx-core/* ./
rm -rf github.com
