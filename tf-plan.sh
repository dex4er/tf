#!/bin/bash
args=()
for arg in "$@"; do
  case "$arg" in
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
terraform plan -detailed-exitcode ${args[*]} | grep --line-buffered -v -P '^\s\s[+-~\s]' | uniq
exit ${PIPESTATUS[0]}
