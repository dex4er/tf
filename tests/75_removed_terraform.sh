#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"
test=$(basename "$0" .sh)

export TERRAFORM_PATH=${test##*_}

rm -rf tmp/$test
mkdir -p tmp/$test
cp removed.tf tmp/$test
pushd tmp/$test >/dev/null

{
  set -x
  ../../../tf init
  ../../../tf apply -auto-approve -quiet
  perl -0 -i -pe 's/resource "terraform_data" "this" \{\n  input = "this"\n\}/removed {\n  from = terraform_data.this\n  lifecycle {\n    destroy = false\n  }\n}/' removed.tf
  ../../../tf plan -quiet
  ../../../tf apply -auto-approve -quiet
} 2>&1 | ../../sanitize.sh >>tf.out

diff -u ../../$test.out tf.out

popd >/dev/null

rm -rf tmp/$test
