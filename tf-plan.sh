#!/bin/bash

function filter_manifest_short() {
  grep --line-buffered -v -P '\(known after apply\)|\(\d+ unchanged \w+ hidden\)'
}

function filter_manifest_compact() {
  grep --line-buffered -v -P '^\s\s[\s+~-]'
}

function filter_terraform_status() {
  declare -A statusline

  declare -A progress
  progress=(["-"]="\\\\" ["\\\\"]="|" ["|"]="/" ["/"]="-" [r]="R" [R]="r" [a]="A" [A]="a" [c]="C" [C]="c" [d]="D" [D]="d")

  IFS=''
  while read line; do
    test "$line" != "$prev" || continue
    case "$line" in
      *': Refreshing state...'*)
        key="-"
        currentstate="${statusline[$key]:-/}"
        statusline[$key]="${progress[$currentstate]}"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Reading...'*)
        statusline["-"]="*"
        key="${line%: Reading...*}"
        statusline[$key]="r"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *)
        echo "$line"
        prev="$line"
    esac
  done
}

filter="filter_manifest_short | filter_terraform_status"

args=()

for arg in "$@"; do
  case "$arg" in
    -compact) filter="filter_manifest_compact | filter_terraform_status";;
    -short) ;;
    -full) filter="cat";;
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done

terraform plan -detailed-exitcode ${args[*]} | eval $filter

exit ${PIPESTATUS[0]}
