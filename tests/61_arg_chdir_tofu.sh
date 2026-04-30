#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"
test=$(basename "$0" .sh)

export TERRAFORM_PATH=${test##*_}

rm -rf tmp/$test
mkdir -p tmp/$test/subdir
cp ../demo.tf tmp/$test/subdir
pushd tmp/$test >/dev/null

{
  set -x
  ../../../tf -chdir=subdir init
  ../../../tf -chdir=subdir upgrade
  ../../../tf -chdir=subdir plan -parallelism=30
  ../../../tf -chdir=subdir apply -auto-approve -parallelism=30
  ../../../tf -chdir=subdir refresh -parallelism=30
  ../../../tf -chdir=subdir destroy -auto-approve -parallelism=30
} 2>&1 | ../../sanitize.sh >>tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
