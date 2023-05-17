#!/usr/bin/env bash

set -euo pipefail

if command -v tput &>/dev/null && tty -s; then
  RED=$(tput setaf 1)
  GREEN=$(tput setaf 2)
  MAGENTA=$(tput setaf 5)
  NORMAL=$(tput sgr0)
  BOLD=$(tput bold)
else
  RED=$(echo -en "\e[31m")
  GREEN=$(echo -en "\e[32m")
  MAGENTA=$(echo -en "\e[35m")
  NORMAL=$(echo -en "\e[00m")
  BOLD=$(echo -en "\e[01m")
fi

function log_header() {
  printf "\n${BOLD}${MAGENTA}==========  %s  ==========${NORMAL}\n" "$@" >&2
}

function log_success() {
  printf "${GREEN}✔ %s${NORMAL}\n" "$@" >&2
}

function log_failure() {
  printf "${RED}✖ %s${NORMAL}\n" "$@" >&2
}

function assert_eq() {
  local expected="$1"
  local actual="$2"
  local msg="${3-}"

  if [ "$expected" == "$actual" ]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "FAIL: Expected '$expected' but got '$actual'. $msg"
    return 1
  fi
}

function assert_not_eq() {
  local expected="$1"
  local actual="$2"
  local msg="${3-}"

  if [ ! "$expected" == "$actual" ]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "FAIL: Expected not '$expected' but got '$actual'. $msg"
    return 1
  fi
}

function assert_true() {
  local actual="$1"
  local msg="${2-}"

  assert_eq true "$actual" "$msg"
  return "$?"
}

function assert_false() {
  local actual="$1"
  local msg="${2-}"

  assert_eq false "$actual" "$msg"
  return "$?"
}

function assert_empty() {
  local actual=$1
  local msg="${2-}"

  assert_eq "" "$actual" "$msg"
  return "$?"
}

function assert_not_empty() {
  local actual=$1
  local msg="${2-}"

  assert_not_eq "" "$actual" "$msg"
  return "$?"
}

function assert_contain() {
  local haystack="$1"
  local needle="${2-}"
  local msg="${3-}"

  if [ -z "${needle:+x}" ]; then
    return 0
  fi

  if [ -z "${haystack##*"${needle}"*}" ]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "$haystack doesn't contain $needle. $msg"
    return 1
  fi
}

function assert_not_contain() {
  local haystack="$1"
  local needle="${2-}"
  local msg="${3-}"

  if [ -z "${needle:+x}" ]; then
    return 0
  fi

  if [ "${haystack##*"${needle}"*}" ]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "$haystack contains $needle. $msg"
    return 1
  fi
}

function assert_gt() {
  local first="$1"
  local second="$2"
  local msg="${3-}"

  if [[ $(echo "$first > $second" | bc) -eq 1 ]]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "$first > $second. $msg"
    return 1
  fi
}

function assert_ge() {
  local first="$1"
  local second="$2"
  local msg="${3-}"

  if [[ $(echo "$first >= $second" | bc) -eq 1 ]]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "$first >= $second. $msg"
    return 1
  fi
}

function assert_lt() {
  local first="$1"
  local second="$2"
  local msg="${3-}"

  if [[ $(echo "$first < $second" | bc) -eq 1 ]]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "$first < $second. $msg"
    return 1
  fi
}

function assert_le() {
  local first="$1"
  local second="$2"
  local msg="${3-}"

  if [[ $(echo "$first <= $second" | bc) -eq 1 ]]; then
    log_success "PASS: $msg"
    return 0
  else
    log_failure "$first <= $second. $msg"
    return 1
  fi
}

# shellcheck source=/dev/null
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/footer.sh"
