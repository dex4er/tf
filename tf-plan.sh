#!/bin/bash
grep="grep --line-buffered -v -P '^\s{4}(?!.*[~+/-]\e)|\(known after apply\)' | { IFS=''; while read line; do test \"\$line\" != \"\$prev\" || continue; case \"\$line\" in *': Refreshing state...'*) printf '.';; *) echo \"\$line\"; prev=\"\$line\"; esac; done; }"
args=()
for arg in "$@"; do
  case "$arg" in
    -compact) grep="grep --line-buffered -v -P '^\s\s[\s+~-]' | { IFS=''; while read line; do test \"\$line\" != \"\$prev\" || continue; case \"\$line\" in *': Refreshing state...'*) printf '.';; *) echo \"\$line\"; prev=\"\$line\"; esac; done; }";;
    -short) ;;
    -full) grep="cat";;
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
tf-init
terraform plan -detailed-exitcode ${args[*]} | eval "$grep"
exit ${PIPESTATUS[0]}
