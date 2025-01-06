#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

pwd=$(pwd)
export TF_PLUGIN_CACHE_DIR=$pwd/tmp/cache
mkdir -p $TF_PLUGIN_CACHE_DIR

status=0

for t in [0-9]*.sh; do
  echo -n "Testing $t... "
  if bash $t </dev/null 1>&2; then
    echo "OK"
  else
    echo "FAILED"
    status=$((status + 1))
    if [[ -n ${TEST_UPDATE-} ]]; then
      n=$(basename "$t" .sh)
      cat tmp/$n/tf.out >$n.out
    fi
  fi
done

rm -rf $TF_PLUGIN_CACHE_DIR

exit $status
