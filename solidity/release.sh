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

package_name=$(jq -r '.name' "$solidity_dir/package.json")
last_version=$(npm info "$package_name" version)
latest_version=$(jq -r '.version' "$solidity_dir/package.json")
if [ "$last_version" == "$latest_version" ]; then
  echo "The latest version of $package_name is already published." && exit 0
fi

if [ ! -f "$solidity_dir/hardhat.config.ts" ]; then
  echo "This script must be run from the root of the repository." && exit 1
fi

cp "$solidity_dir/package.json" "$solidity_dir/contracts/"
cp "$solidity_dir/README.md" "$solidity_dir/contracts/"
cp "$project_dir/LICENSE" "$solidity_dir/contracts/"

(
  cd "$solidity_dir" || exit 1
  yarn install
  yarn clean
  yarn compile
)

mkdir -p "$solidity_dir"/contracts/build/contracts
find "$solidity_dir"/artifacts/contracts -name '*.json' -exec cp {} "$solidity_dir"/contracts/build/contracts \;
rm "$solidity_dir"/contracts/build/contracts/*.dbg.json

for key in "scripts" "dependencies" "devDependencies"; do
  jq -r "del(.$key)" "$solidity_dir/contracts/package.json" >"tmp.json"
  mv "tmp.json" "$solidity_dir/contracts/package.json"
done

# publish contracts package
(
  cd "$solidity_dir/contracts"
  if [ -z "$NODE_AUTH_TOKEN" ]; then
    echo "WARN: NODE_AUTH_TOKEN is required to publish the package."
    npm publish --dry-run
  else
    npm publish --access public
  fi
)
