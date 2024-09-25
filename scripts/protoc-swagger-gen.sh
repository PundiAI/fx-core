#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./tmp-swagger-gen ./third_party ./build
trap 'rm -rf ./tmp-swagger-gen ./third_party' EXIT

commit_hash=$(go list -m -f '{{.Replace.Version}}' github.com/evmos/ethermint | awk -F '-' '{print $NF}')
if [[ ! -f "./build/$commit_hash.zip" ]]; then
  wget -c "https://github.com/functionx/ethermint/archive/$commit_hash.zip" -O "./build/$commit_hash.zip"
fi
unzip -q -o "./build/$commit_hash.zip" -d "./build"
# shellcheck disable=SC2010
cp -r "./build/$(ls ./build | grep ethermint | grep -v grep | grep -v zip)/proto" ./third_party/
rm -rf ./build/ethermint-*

mkdir -p "./third_party/cosmos/ics23/v1"
wget https://raw.githubusercontent.com/cosmos/ics23/master/proto/cosmos/ics23/v1/proofs.proto -O "./third_party/cosmos/ics23/v1/proofs.proto"

buf generate --template ./proto/buf.gen.swagger.yaml "$(grep cosmos/cosmos-sdk proto/buf.yaml | awk '{print $2}')"
buf generate --template ./proto/buf.gen.swagger.yaml "$(grep cosmos/ibc proto/buf.yaml | awk '{print $2}')"

# create swagger files on an individual basis  w/ `buf build` and `buf generate` (needed for `swagger-combine`)
proto_dirs=$(find ./proto ./third_party -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ -n "$query_file" && -f "$query_file" ]]; then
    buf generate --template ./proto/buf.gen.swagger.yaml "$query_file"
  fi
done

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./docs/config.json -o ./docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true
