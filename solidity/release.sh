#!/usr/bin/env bash

set -eo pipefail

commands=(git yarn jq)
for cmd in "${commands[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "$cmd command not found, please install $cmd first" && exit 1
  fi
done

project_dir="$(git rev-parse --show-toplevel)"
solidity_dir="$project_dir"
if [ -d "$project_dir/solidity" ]; then
  solidity_dir="$project_dir/solidity"
fi

if [ ! -f "$solidity_dir/hardhat.config.ts" ]; then
  echo "This script must be run from the root of the repository." && exit 1
fi

cp "$solidity_dir/README.md" "$solidity_dir/contracts/"
cp "$project_dir/LICENSE" "$solidity_dir/contracts/"

(
  cd "$solidity_dir" || exit 1
  yarn clean
  yarn compile
)

mkdir -p "$solidity_dir"/contracts/build/contracts
find "$solidity_dir"/artifacts/contracts -name '*.json' -exec cp {} "$solidity_dir"/contracts/build/contracts \;
rm "$solidity_dir"/contracts/build/contracts/*.dbg.json

# read current version from solidity/package.json
cur_version=$(jq -r .version "$solidity_dir/package.json")
# increment version
next_version=$(echo "$cur_version" | awk -F. -v OFS=. '{$NF++;print}')
next_version=${SOLIdITY_VERSION:-$next_version}
# update version
for file in "$solidity_dir/package.json" "$solidity_dir/contracts/package.json"; do
  jq ".version = \"$next_version\"" "$file" >"$file.tmp"
  mv "$file.tmp" "$file"
done

# publish contracts package
(
  cd "$solidity_dir/contracts"
  npm publish --access public
)
