#!/bin/bash
terraform state list "$@" | sed 's/\x1b\[[01]m//g'
exit ${PIPESTATUS[0]}
