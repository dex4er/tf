#!/bin/bash
set -e
for r in "$@"; do
  terraform taint "$r"
done
