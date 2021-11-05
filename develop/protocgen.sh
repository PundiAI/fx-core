#!/usr/bin/env bash

set -e -x

set -eo pipefail

if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null ; then
  echo -e "\tPlease run this command from somewhere inside the fx-core folder."
  exit 1
fi

if [ ! -d "${GOPATH}/src/github.com/cosmos/cosmos-sdk" ]; then
  echo -e "\tPlease clone cosmos-sdk to ${GOAPTH}/src/."
  exit 1
fi

(cd ${GOPATH}/src/github.com/cosmos/cosmos-sdk && git stash && git checkout v0.42.1)

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  protoc \
  -I "proto" \
  -I ${GOPATH}/src/github.com/cosmos/cosmos-sdk/proto \
  -I ${GOPATH}/src/github.com/cosmos/cosmos-sdk/third_party/proto \
  --gocosmos_out=plugins=interfacetype+grpc,Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done

cp -r github.com/functionx/fx-core/* ./
rm -rf github.com
