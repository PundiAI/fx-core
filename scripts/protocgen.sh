#!/usr/bin/env bash

set -eo pipefail

# generate proto files
echo "Generating gogo proto code"
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $proto_files; do
    # Check if the go_package in the file is pointing to functionx
    if grep -q "option go_package.*functionx" "$file"; then
      buf generate --template proto/buf.gen.gogo.yaml "$file"
    fi
  done
done

# move proto files to the right places
cp -r github.com/functionx/fx-core/* ./
rm -rf github.com

sh ./scripts/protocgen-pulsar.sh
