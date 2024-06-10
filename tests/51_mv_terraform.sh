#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"
test=$(basename "$0" .sh)

export TERRAFORM_PATH=${test##*_}

rm -rf tmp/$test
mkdir -p tmp/$test
cp ../demo.tf tmp/$test
pushd tmp/$test >/dev/null

{
  set -x
  ../../../tf init
  ../../../tf destroy -auto-approve
  ../../../tf import time_sleep.this["1s"] 1s,1s
  ../../../tf list
  ../../../tf mv time_sleep.this["1s"] time_sleep.this["2s"]
  ../../../tf list
} 2>&1 | ../../sanitize.sh >>tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
