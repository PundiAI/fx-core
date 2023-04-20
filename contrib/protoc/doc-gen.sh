#!/usr/bin/env bash

set -eo pipefail

if [ ! -d build/docs/bin ]; then
  mkdir -p build/docs/bin
fi

if [ ! -f "./build/docs/bin/protoc-gen-doc" ]; then
  if [ ! -f "./build/docs/protoc-gen-doc.tar.gz" ]; then
    VERSION=1.5.1
    PLATFORM=amd64
    PL=$(uname -m)
    if [ "x86_64" != "$PL" ]; then
      PLATFORM=arm64
    fi
    echo "download protoc-gen-doc $VERSION $PLATFORM ..."
    wget -c https://github.com/pseudomuto/protoc-gen-doc/releases/download/v${VERSION}/protoc-gen-doc_${VERSION}_linux_${PLATFORM}.tar.gz -O ./build/docs/protoc-gen-doc.tar.gz
  fi
  (
    cd ./build/docs
    mkdir protoc-gen-doc-tmp
    tar -zxf protoc-gen-doc.tar.gz -C protoc-gen-doc-tmp
    mv protoc-gen-doc-tmp/protoc-gen-doc ./bin
    rm -rf protoc-gen-doc-tmp
  )
fi


export PATH=$PATH:$PWD/build/docs/bin

# generate fx proto doc
echo "generate fx proto doc"
buf generate --template proto/buf.gen.doc.yaml proto

if [ ! -f "./build/docs/cosmos-sdk/README.md" ]; then
  commit_hash=$(go list -m -f '{{.Version}}' github.com/cosmos/cosmos-sdk)
  if [ ! -f "./build/docs/cosmos-sdk-proto.zip" ]; then
    echo "download cosmos-sdk $commit_hash"
    wget -c "https://github.com/cosmos/cosmos-sdk/archive/$commit_hash.zip" -O "./build/docs/cosmos-sdk-proto.zip"
  fi
  (
    cd build/docs
    unzip -q -o "./cosmos-sdk-proto.zip"
    # shellcheck disable=SC2010
    mv "$(ls | grep cosmos-sdk | grep -v grep | grep -v zip)" cosmos-sdk
    rm -rf cosmos-sdk/.git
    # remove unused proto files
    rm -rf cosmos-sdk/proto/cosmos/group
    rm -rf cosmos-sdk/proto/cosmos/nft
  )
fi

if [ ! -f ./build/docs/ibc-go/README.md ]; then
  commit_hash=$(go list -m -f '{{.Version}}' github.com/cosmos/ibc-go/v6)
  if [ ! -f "./build/docs/ibc-go-proto.zip" ]; then
    echo "download ibc-go $commit_hash"
    wget -c "https://github.com/cosmos/ibc-go/archive/$commit_hash.zip" -O "./build/docs/ibc-go-proto.zip"
  fi
  (
    cd build/docs
    unzip -q -o "./ibc-go-proto.zip"
    # shellcheck disable=SC2010
    mv "$(ls | grep ibc-go | grep -v grep | grep -v zip)" ibc-go
    rm -rf ibc-go/.git
  )
fi

if [ ! -f ./build/docs/ethermint/README.md ]; then
  commit_hash=$(go list -m -f '{{.Version}}' github.com/evmos/ethermint | awk -F- '{print $1}')
  if [ ! -f "./build/docs/ethermint-proto.zip" ]; then
    echo "download ethermint $commit_hash"
    wget -c "https://github.com/evmos/ethermint/archive/$commit_hash.zip" -O "./build/docs/ethermint-proto.zip"
  fi
  (
    cd build/docs
    unzip -q -o "./ethermint-proto.zip"
    # shellcheck disable=SC2010
    mv "$(ls | grep ethermint | grep -v grep | grep -v zip)" ethermint
    rm -rf ethermint/.git
  )
fi


for item in "cosmos-sdk" "ibc-go" "ethermint"; do
  echo "generate $item proto doc"
  cp proto/buf.gen.doc.yaml build/docs/${item}/proto
  sed -i 's|fx|'"$item"'|g' "build/docs/${item}/proto/buf.gen.doc.yaml"
  sed -i 's|docs/proto|docs|g' "build/docs/${item}/proto/buf.gen.doc.yaml"

  (
    cd build/docs/${item}
    buf generate --template proto/buf.gen.doc.yaml proto
  )

  mv build/docs/${item}/docs/${item}-proto-docs.md docs/proto
  rm build/docs/${item}/proto/buf.gen.doc.yaml
done
