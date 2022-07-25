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

for item in "cosmos-sdk" "ibc-go" "ethermint"; do
  proto_dirs=$(find "./build/${item}/proto" -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
  for dir in $proto_dirs; do
    # generate swagger files (filter query files)
    query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
    if [[ ! -z "$query_file" ]]; then
      buf protoc \
        -I "build/${item}/proto" \
        -I "build/${item}/third_party/proto" \
        "$query_file" \
        --swagger_out=./tmp-swagger-gen \
        --swagger_opt=logtostderr=true \
        --swagger_opt=fqn_for_swagger_name=true \
        --swagger_opt=simple_operation_ids=true
    fi
  done
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
