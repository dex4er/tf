#!/bin/bash
set -e
for r in "$@"; do
  terraform state rm "$r"
done
