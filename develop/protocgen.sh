#!/usr/bin/env bash

set -e -x

set -eo pipefail

protoc_gen_gocosmos() {
  if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null; then
    echo -e "\tPlease run this command from somewhere inside the cosmos-sdk folder."
    return 1
  fi

  go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
}

protoc_gen_gocosmos

if [ ! -d build ]; then
  mkdir -p build
fi

if [ ! -f "./build/cosmos-sdk/README.md" ]; then
  commit_hash=$(go list -m -f '{{.Version}}' github.com/cosmos/cosmos-sdk)
  if [ ! -f "./build/cosmos-sdk-42-proto.zip" ]; then
    wget -c "https://github.com/cosmos/cosmos-sdk/archive/$commit_hash.zip" -O "./build/cosmos-sdk-42-proto.zip"
  fi
  (
    cd build
    unzip -q -o "./cosmos-sdk-42-proto.zip"
    mv $(ls | grep cosmos-sdk | grep -v grep | grep -v zip) cosmos-sdk
    rm -rf cosmos-sdk/.git
  )
fi

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  buf protoc \
    -I "proto" \
    -I "build/cosmos-sdk/proto" \
    -I "build/cosmos-sdk/third_party/proto" \
    --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
    --grpc-gateway_out=logtostderr=true:. \
    $(find "${dir}" -maxdepth 1 -name '*.proto')
done

# command to generate docs using protoc-gen-doc
buf protoc \
  -I "proto" \
  -I "build/cosmos-sdk/proto" \
  -I "build/cosmos-sdk/third_party/proto" \
  --doc_out=./docs/proto \
  --doc_opt=./docs/proto/proto-doc-markdown.tmpl,fx-proto-docs.md \
  $(find "$(pwd)/proto" -maxdepth 5 -name '*.proto')

buf protoc \
    -I "build/cosmos-sdk/proto" \
    -I "build/cosmos-sdk/third_party/proto" \
    --doc_out=./docs/proto \
    --doc_opt=./docs/proto/proto-doc-markdown.tmpl,cosmos-sdk-proto-docs.md \
    $(find "$(pwd)/build/cosmos-sdk/proto" -maxdepth 5 -name '*.proto')

cp -r github.com/functionx/fx-core/* ./
rm -rf github.com