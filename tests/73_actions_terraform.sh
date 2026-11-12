#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"
test=$(basename "$0" .sh)

export TERRAFORM_PATH=${test##*_}

rm -rf tmp/$test
mkdir -p tmp/$test
cp actions.tf tmp/$test
pushd tmp/$test >/dev/null

{
  set -x
  ../../../tf init
  ../../../tf plan
  ../../../tf apply -auto-approve -counters
  ../../../tf apply -auto-approve -verbose -invoke=action.local_command.hello
  ../../../tf apply -auto-approve -quiet -invoke action.local_command.hello
} 2>&1 | ../../sanitize.sh >>tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
