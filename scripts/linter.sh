#!/usr/bin/env bash

set -eo pipefail

patternLimits=(
  "nolint:21"
  "#nosec:5"
  "CrossChain:4"
  "cross chain:0"
)

if ! command -v rg &>/dev/null; then
  echo "rg command not found, please install rg first: https://github.com/BurntSushi/ripgrep?tab=readme-ov-file#installation" && exit 1
fi

check_pattern_count() {
  local pattern=$1
  local allowed_count=$2
  local file_type=$3

  # Default values if not provided
  file_type=${file_type:-go}

  rg_args=(--type go --glob '!*.pb.go' --glob '!*.pulsar.go' --glob '!*.sol.go' --glob '!legacy.go')

  if [[ "$allowed_count" -eq 0 ]]; then
    if rg "${rg_args[@]}" "$pattern" ./ >/dev/null; then
      echo "Warning: Matches found for '$pattern'."
      exit 1
    fi
    return
  fi

  # Count the actual number of 'pattern' in specified file types
  actual_count=$(rg "${rg_args[@]}" "$pattern" ./ | wc -l)
  echo "Actual count of '$pattern': $actual_count"

  # Compare with the allowed count
  if [ "$actual_count" -eq "$allowed_count" ]; then
    echo "The count matches the allowed number of suppressions."
  elif [ "$actual_count" -lt "$allowed_count" ]; then
    echo "The actual count is less than $allowed_count. Consider further reducing the allowed count."
    exit 1
  else
    echo "Warning: The actual count is higher than $allowed_count. Please review and update suppressions or adjust the allowed count."
    rg "${rg_args[@]}" "$pattern" ./
    exit 1
  fi
}

for pattern_limit in "${patternLimits[@]}"; do
  pattern=$(echo "$pattern_limit" | cut -d: -f1)
  allowed_count=$(echo "$pattern_limit" | cut -d: -f2)
  check_pattern_count "$pattern" "$allowed_count"
done
