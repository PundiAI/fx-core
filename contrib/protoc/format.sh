#!/usr/bin/env bash

set -eo pipefail

# export third_party proto files
mkdir -p third_party/proto
echo "export fork third_party proto files..."
fx_deps=""
fx_deps="${fx_deps} buf.build/functionx/cosmos-sdk:$(go list -m -f '{{.Version}}' github.com/cosmos/cosmos-sdk)"
fx_deps="${fx_deps} buf.build/functionx/ibc:$(go list -m -f '{{.Version}}' github.com/cosmos/ibc-go/v7)"
for dep in $fx_deps; do
  echo "$dep downloading..."
  buf export "$dep" --output third_party/proto --exclude-imports
done

echo "buf format proto files..."
buf format -w

rm -rf third_party
