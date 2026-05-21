#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

pwd=$(pwd)
export TF_PLUGIN_CACHE_DIR=$pwd/tmp/cache
mkdir -p $TF_PLUGIN_CACHE_DIR

status=0

has_terraform=$(command -v terraform >/dev/null && echo yes || echo no)
has_tofu=$(command -v terraform >/dev/null && echo yes || echo no)

for t in [0-9]*.sh; do
  echo -n "Testing $t... "
  if [[ ${t} == *terraform.out && $has_terraform == no ]]; then
    echo -e "\033[0mSKIPPED (terraform not found)"
    continue
  fi
  if [[ ${t} == *tofu.out && $has_tofu == no ]]; then
    echo -e "\033[0mSKIPPED (tofu not found)"
    continue
  fi
  if bash $t </dev/null 1>&2; then
    echo -e "\033[0mOK"
  else
    echo -e "\033[0mFAILED"
    status=$((status + 1))
    if [[ -n ${TEST_UPDATE-} ]]; then
      n=$(basename "$t" .sh)
      cat tmp/$n/tf.out >$n.out
    fi
  fi
done

rm -rf $TF_PLUGIN_CACHE_DIR

exit $status
