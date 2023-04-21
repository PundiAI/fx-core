#!/usr/bin/env bash
SWAGGER_DIR=./swagger-proto

set -eo pipefail
#set -x

# prepare swagger generation
mkdir -p "$SWAGGER_DIR/proto"
mkdir -p "$SWAGGER_DIR/third_party"

printf "version: v1\ndirectories:\n  - proto\n  - third_party" > "$SWAGGER_DIR/buf.work.yaml"
printf "version: v1\nname: buf.build/functionx/fx-core\n" > "$SWAGGER_DIR/proto/buf.yaml"
cp ./proto/buf.gen.swagger.yaml "$SWAGGER_DIR/proto/buf.gen.swagger.yaml"

# copy existing proto files
cp -r ./proto/fx "$SWAGGER_DIR/proto"

# download fx proto deps
fx_deps=$(awk '/^[^ ]/{ f=/^deps:/; next } f{ if (sub(/:$/,"")) deps=$2; else print $2 }' proto/buf.yaml)
# add fork proto deps
fx_deps="${fx_deps} buf.build/functionx/cosmos-sdk:$(go list -m -f '{{.Version}}' github.com/cosmos/cosmos-sdk)"
fx_deps="${fx_deps} buf.build/functionx/ethermint:$(go list -m -f '{{.Version}}' github.com/evmos/ethermint)"
fx_deps="${fx_deps} buf.build/functionx/ibc:$(go list -m -f '{{.Version}}' github.com/cosmos/ibc-go/v6)"
echo "download fx-core proto deps ..."
for dep in $fx_deps ; do
  echo "$dep downloading..."
  buf export "$dep" --output "$SWAGGER_DIR/third_party" --exclude-imports
done


# create temporary folder to store intermediate results from `buf generate`
mkdir -p ./tmp-swagger-gen

# step into swagger folder
cd "$SWAGGER_DIR"
echo ""
# create swagger files on an individual basis  w/ `buf build` and `buf generate` (needed for `swagger-combine`)
proto_dirs=$(find ./proto ./third_party -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ -n "$query_file" ]]; then
    buf generate --template proto/buf.gen.swagger.yaml "$query_file"
  fi
done

cd ..

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./docs/config.json -o ./docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen
rm -rf "$SWAGGER_DIR"
