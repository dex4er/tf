#!/usr/bin/env bash

# Usage: ./lock.sh FIXTURE DIR, from tests/, with TERRAFORM_PATH set

set -euo pipefail

fixture=$1
dir=$2
tool=${TERRAFORM_PATH:?TERRAFORM_PATH is not set}

name=$(basename "$fixture" .tf)
lock=tmp/locks/${name}_${tool}.lock.hcl

if [[ ! -f ${lock} || ${fixture} -nt ${lock} ]]; then
  mkdir -p tmp/locks
  workdir=$(mktemp -d tmp/locks/init.XXXXXX)
  cp "${fixture}" "${workdir}"
  (cd "${workdir}" && "${tool}" init -backend=false -input=false >/dev/null)
  # a fixture with no providers but the built-in one gets no lock file
  if [[ -f ${workdir}/.terraform.lock.hcl ]]; then
    cp "${workdir}/.terraform.lock.hcl" "${lock}"
  else
    : >"${lock}"
  fi
  rm -rf "${workdir}"
fi

if [[ -s ${lock} ]]; then
  cp "${lock}" "${dir}/.terraform.lock.hcl"
fi
