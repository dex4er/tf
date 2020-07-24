#!/bin/bash
grep="grep --line-buffered -v -P '^\s\s[+-~\s]' | uniq"
args=()
for arg in "$@"; do
  case "$arg" in
    -compact);;
    -short) grep="grep --line-buffered -v -P '(known after apply)'";;
    -full) grep="cat";;
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
terraform plan -detailed-exitcode ${args[*]} | eval $grep
exit ${PIPESTATUS[0]}
