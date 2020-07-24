#!/bin/bash
set -e
terraform state rm "$@"
