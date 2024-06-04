#!/usr/bin/env bash

set -euo pipefail

cd $(dirname $0)
test=$(basename $0 .sh)

export TERRAFORM_PATH=${test##*_}

rm -rf tmp/$test
mkdir -p tmp/$test
cp ../demo.tf tmp/$test
pushd tmp/$test >/dev/null

{
  ../../../tf version
} | ../../sanitize.sh 2>&1 >> tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
