#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"
test=$(basename "$0" .sh)

export TERRAFORM_PATH=${test##*_}

rm -rf tmp/$test
mkdir -p tmp/$test
cp ephemeral.tf tmp/$test
./lock.sh ephemeral.tf tmp/$test
pushd tmp/$test >/dev/null

{
  set -x
  ../../../tf init
  ../../../tf plan -verbose
  ../../../tf plan -counters
  ../../../tf plan -dots
  ../../../tf plan -quiet
  ../../../tf apply -auto-approve -verbose
  ../../../tf apply -auto-approve -counters
} 2>&1 | ../../sanitize.sh >>tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
