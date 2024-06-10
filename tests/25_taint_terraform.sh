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
  ../../../tf import 'time_sleep.this["2s"]' 2s,2s
  ../../../tf import time_sleep.this["3s"] 3s,3s
  ../../../tf import time_sleep.this["4s"] 4s,4s
  ../../../tf taint time_sleep.this["1s"]
  ../../../tf taint 'time_sleep.this["2s"]' time_sleep.this["3s"]
  ../../../tf taint time_sleep.this["4s"] time_sleep.this["foo"] || true
  ../../../tf untaint time_sleep.this["1s"]
  ../../../tf untaint 'time_sleep.this["2s"]' time_sleep.this["3s"]
  ../../../tf untaint time_sleep.this["4s"] time_sleep.this["foo"] || true
} 2>&1 | ../../sanitize.sh >>tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
