#!/usr/bin/env bash

set -eo pipefail

# export third_party proto files
mkdir -p third_party/proto
echo "export fork third_party proto files..."
fx_deps=""
fx_deps="${fx_deps} buf.build/functionx/cosmos-sdk:$(go list -m -f '{{.Version}}' github.com/cosmos/cosmos-sdk)"
fx_deps="${fx_deps} buf.build/functionx/ethermint:$(go list -m -f '{{.Version}}' github.com/evmos/ethermint)"
fx_deps="${fx_deps} buf.build/functionx/ibc:$(go list -m -f '{{.Version}}' github.com/cosmos/ibc-go/v6)"
for dep in $fx_deps; do
  echo "$dep downloading..."
  buf export "$dep" --output third_party/proto --exclude-imports
done

# generate proto files
echo "generate proto files..."
proto_dirs=$(find ./proto ./third_party/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $proto_files; do
    # Check if the go_package in the file is pointing to evmos
    if grep -q "option go_package.*fx" "$file"; then
      buf generate --template proto/buf.gen.gogo.yaml "$file"
    fi
  done
done

# move proto files to the right places
cp -r github.com/functionx/fx-core/* ./
rm -rf github.com
rm -rf third_party
