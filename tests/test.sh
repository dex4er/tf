#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

status=0

for t in [0-9]*.sh; do
  echo -n "Testing $t... "
  if bash $t </dev/null 1>&2; then
    echo "OK"
  else
    echo "FAILED"
    status=$((status + 1))
  fi
done

exit $status
