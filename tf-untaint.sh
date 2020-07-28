#!/bin/bash
set -e
for r in "$@"; do
  terraform untaint "$r"
done
