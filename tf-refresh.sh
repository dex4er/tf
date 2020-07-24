#!/bin/bash
args=()
for arg in "$@"; do
  case "$arg" in
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
terraform refresh ${args[*]}
exit ${PIPESTATUS[0]}
