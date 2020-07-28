#!/bin/bash
grep="grep --line-buffered -v -P '^\s{4}(?!.*[~+/-]\e)|\(known after apply\)'"
args=()
for arg in "$@"; do
  case "$arg" in
    -compact) grep="grep --line-buffered -v -P '^\s\s[\s+~-]' | uniq";;
    -short) ;;
    -full) grep="cat";;
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
terraform plan -detailed-exitcode ${args[*]} | eval $grep
exit ${PIPESTATUS[0]}
