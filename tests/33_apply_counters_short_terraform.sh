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
  ../../../tf upgrade
  ../../../tf plan -parallelism=30 -counters -short
  ../../../tf apply -auto-approve -parallelism=30 -counters -short
  ../../../tf refresh -parallelism=30 -counters -short
  ../../../tf destroy -auto-approve -parallelism=30 -counters -short
} 2>&1 | ../../sanitize.sh >>tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
