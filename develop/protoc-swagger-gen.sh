#!/usr/bin/env bash

set -eo pipefail

if [ -d "build/cosmos-sdk/proto/ibc/applications" ]; then
  rm -rf "build/cosmos-sdk/proto/ibc/applications"
fi

if [ -d ./tmp-swagger-gen ]; then
  rm -rf ./tmp-swagger-gen
fi
mkdir -p ./tmp-swagger-gen

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf protoc \
      -I "proto" \
      -I "build/cosmos-sdk/proto" \
      -I "build/cosmos-sdk/third_party/proto" \
      "$query_file" \
      --swagger_out=./tmp-swagger-gen \
      --swagger_opt=logtostderr=true \
      --swagger_opt=fqn_for_swagger_name=true \
      --swagger_opt=simple_operation_ids=true
  fi
done

proto_dirs=$(find ./build/cosmos-sdk/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf protoc \
      -I "build/cosmos-sdk/proto" \
      -I "build/cosmos-sdk/third_party/proto" \
      "$query_file" \
      --swagger_out=./tmp-swagger-gen \
      --swagger_opt=logtostderr=true \
      --swagger_opt=fqn_for_swagger_name=true \
      --swagger_opt=simple_operation_ids=true
  fi
done

proto_dirs=$(find ./build/ethermint/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf protoc \
      -I "build/ethermint/proto" \
      -I "build/ethermint/third_party/proto" \
      "$query_file" \
      --swagger_out=./tmp-swagger-gen \
      --swagger_opt=logtostderr=true \
      --swagger_opt=fqn_for_swagger_name=true \
      --swagger_opt=simple_operation_ids=true
  fi
done

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./docs/config.json \
  -o ./docs/swagger-ui/swagger.yaml \
  -f yaml \
  --continueOnConflictingPaths true \
  --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen
