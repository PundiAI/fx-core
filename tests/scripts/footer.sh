#!/usr/bin/env bash

set -eo pipefail

if [[ "$1" == "help" || "$#" -eq 0 ]]; then
  help && exit 0
fi

if [[ "$#" -gt 0 && "$(type -t "$1")" != "function" ]]; then
  echo "invalid command: $1" && help && exit 1
fi

if ! "$@"; then
  echo "failed: $0" "$@" && exit 1
fi
