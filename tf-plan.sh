#!/bin/bash
terraform plan -detailed-exitcode "$@" | grep --line-buffered -v -P '^\s\s[+-~\s]' | uniq
exit ${PIPESTATUS[0]}
