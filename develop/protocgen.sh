#!/usr/bin/env bash

set -eo pipefail

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
  cosmos_sdk_commit_hash=$(go list -m -f '{{.Replace.Version}}' github.com/cosmos/cosmos-sdk | awk -F- '{print $3}')
  if [ ! -f "./build/cosmos-proto.zip" ]; then
    wget -c "https://github.com/functionx/cosmos-sdk/archive/$cosmos_sdk_commit_hash.zip" -O "./build/cosmos-proto.zip"
  fi
  (
    cd build
    unzip -q -o "./cosmos-proto.zip"
    for dir in *; do
      if [[ -d $dir && "$dir" == "cosmos-sdk-$cosmos_sdk_commit_hash"* ]]; then
        mv "./$dir" cosmos-sdk
      fi
    done
    rm -rf ./cosmos-sdk/proto/ibc/applications
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
  --doc_out=./docs \
  --doc_opt=./docs/protodoc-markdown.tmpl,proto-docs.md \
  $(find "$(pwd)/proto" -maxdepth 5 -name '*.proto')

#go mod tidy

cp -r github.com/functionx/fx-core/* ./
rm -rf github.com
