#!/bin/bash
if [ $# -gt 0 ]; then
    terraform state show -no-color "$@" | sed 's/\x1b\[[01]m//g'
    exit ${PIPESTATUS[0]}
else
    terraform show -no-color
    exit $?
fi
