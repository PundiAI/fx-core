#!/usr/bin/env bash

set -eo pipefail

VERSION=$(buf --version)
echo "buf version: $VERSION"

if [ -z "$BUF_NAME" ]; then
  echo "buf name not found, please set BUF_NAME"
  exit 1
fi

if [ -z "$BUF_TOKEN" ]; then
  echo "buf token not found, please set BUF_TOKEN"
  exit 1
fi

if [ -z "$BUF_ORG" ]; then
  echo "buf org not found, please set BUF_ORG"
  exit 1
fi

echo "buf registry login $BUF_NAME with ******"
echo "$BUF_TOKEN" | buf registry login --username "$BUF_NAME" --token-stdin

echo "USER $BUF_NAME push proto to $BUF_ORG ..."

read -rp "Input want to push proto to $BUF_ORG: " input
if [ "$input" != "cosmos-sdk" ] && [ "$input" != "ethermint" ] && [ "$input" != "ibc" ]; then
  echo "input '$input' error, please input 'cosmos-sdk' or 'ethermint' or 'ibc'" && exit 1
fi

if [ ! -d build/fork ]; then
  mkdir -p build/fork
fi

if [ "$input" == "cosmos-sdk" ]; then
  # download cosmos-sdk proto
  commit_hash=$(go list -m -f '{{.Version}}' github.com/cosmos/cosmos-sdk)
  if [ ! -f "./build/fork/cosmos-sdk-proto.zip" ]; then
    echo "download cosmos-sdk $commit_hash"
    wget -c "https://github.com/cosmos/cosmos-sdk/archive/$commit_hash.zip" -O "./build/fork/cosmos-sdk-proto.zip"
  fi

  (
    cd build/fork
    rm -rf cosmos-sdk
    unzip -q -o "./cosmos-sdk-proto.zip"
    # shellcheck disable=SC2010
    mv "$(ls | grep cosmos-sdk | grep -v grep | grep -v zip)" cosmos-sdk
    rm -rf cosmos-sdk/.git

    # buf push
    cd cosmos-sdk/proto
    # replace buf.yaml buf.build/cosmos/cosmos-sdk => buf.build/functionx/cosmos-sdk
    perl -pi -e 's|buf.build/cosmos/cosmos-sdk|buf.build/'"$BUF_ORG"'/cosmos-sdk|g' buf.yaml

    echo "buf push cosmos-sdk proto with tag $commit_hash ..."
    buf push --tag "$commit_hash"
  )
fi

if [ "$input" == "ethermint" ]; then
  # download ethermint proto
  commit_hash=$(go list -m -f '{{.Replace.Version}}' github.com/evmos/ethermint | awk -F '-' '{print $NF}')
  if [ ! -f "./build/fork/ethermint-proto.zip" ]; then
    echo "download ethermint $commit_hash"
    wget -c "https://github.com/functionx/ethermint/archive/$commit_hash.zip" -O "./build/fork/ethermint-proto.zip"
  fi

  (
    cd build/fork
    rm -rf ethermint
    unzip -q -o "./ethermint-proto.zip"
    # shellcheck disable=SC2010
    mv "$(ls | grep ethermint | grep -v grep | grep -v zip)" ethermint
    rm -rf ethermint/.git

    # buf push
    cd ethermint/proto
    # replace buf.yaml buf.build/evmos/ethermint => buf.build/functionx/ethermint
    perl -pi -e 's|buf.build/evmos/ethermint|buf.build/'"$BUF_ORG"'/ethermint|g' buf.yaml

    echo "buf push ethermint proto with tag $commit_hash ..."
    buf push --tag "$commit_hash"
  )
fi

if [ "$input" == "ibc" ]; then
  # download ibc-go proto
  commit_hash=$(go list -m -f '{{.Version}}' github.com/cosmos/ibc-go/v7)
  if [ ! -f "./build/fork/ibc-go-proto.zip" ]; then
    echo "download ibc-go $commit_hash"
    wget -c "https://github.com/cosmos/ibc-go/archive/$commit_hash.zip" -O "./build/fork/ibc-go-proto.zip"
  fi
  (
    cd build/fork
    rm -rf ibc-go
    unzip -q -o "./ibc-go-proto.zip"
    # shellcheck disable=SC2010
    mv "$(ls | grep ibc-go | grep -v grep | grep -v zip)" ibc-go
    rm -rf ibc-go/.git

    cd ibc-go/proto
    perl -pi -e 's|buf.build/cosmos/ibc|buf.build/'"$BUF_ORG"'/ibc|g' buf.yaml

    buf mod update
    echo "buf push ibc-go proto with tag $commit_hash ..."
    buf push --tag "$commit_hash"
  )
fi
