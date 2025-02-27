#!/usr/bin/env bash

set -euo pipefail

if command -v gsed >/dev/null 2>&1; then
  sed="gsed"
else
  sed="sed"
fi

tr '\015' '\012' |
  $sed -e '/^ *$/d' \
    -e '/Your version of .* is out of date! The latest version/d' \
    -e '/You can update by downloading from/d' \
    -e '/: Creating\.\.\./d' \
    -e '/: Destroying\.\.\./d' \
    -e '/: Reading\.\.\./d' \
    -e '/: Still .*\.\.\./d' \
    -e 's/ \[id=[0-9a-zA-Z,:-]*\]\(\o033\)/\1/' \
    -e 's/ \[id=[0-9a-zA-Z,:-]*\]$//' \
    -e 's/ \[id=[0-9a-zA-Z,:-]*\(\o033\)/\1/' \
    -e 's/ \[id=[0-9a-zA-Z,:-]*$//' \
    -e 's/ \[id\(\o033\)/\1/' \
    -e 's/ \[id$//' \
    -e 's/ \[i\(\o033\)/\1/' \
    -e 's/ \[i$//' \
    -e 's/ \[\(\o033\)/\1/' \
    -e 's/ \[$//' \
    -e 's/ complete after [0-9][0-9]*s \(\o033\)/\1/' \
    -e 's/ complete after [0-9][0-9]*s $//' \
    -e 's/ complete after [0-9][0-9]*s\(\o033\)/\1/' \
    -e 's/ complete after [0-9][0-9]*s$//' \
    -e 's/ complete after [0-9]*\(\o033\)/\1/' \
    -e 's/ complete after [0-9]*s$//' \
    -e 's/ complete after\(\o033\)/\1/' \
    -e 's/ complete after$//' \
    -e 's/ complete afte\(\o033\)/\1/' \
    -e 's/ complete afte$//' \
    -e 's/ complete aft\(\o033\)/\1/' \
    -e 's/ complete aft$//' \
    -e 's/ complete af\(\o033\)/\1/' \
    -e 's/ complete af$//' \
    -e 's/ complete a\(\o033\)/\1/' \
    -e 's/ complete a$//' \
    -e 's/ complete \(\o033\)/\1/' \
    -e 's/ complete $//' \
    -e 's/ complete\(\o033\)/\1/' \
    -e 's/ complete$//' \
    -e 's/ complet\(\o033\)/\1/' \
    -e 's/ complet$//' \
    -e 's/ comple\(\o033\)/\1/' \
    -e 's/ comple$//' \
    -e 's/ compl\(\o033\)/\1/' \
    -e 's/ compl$//' \
    -e 's/ comp\(\o033\)/\1/' \
    -e 's/ comp$//' \
    -e 's/ com\(\o033\)/\1/' \
    -e 's/ com$//' \
    -e 's/ co\(\o033\)/\1/' \
    -e 's/ co$//' \
    -e 's/ c\(\o033\)/\1/' \
    -e 's/ c$//' \
    -e 's/\(\o033\)*\o033/\o033/g' \
    -e 's/\(\o033\[0m\)*\o033\[0m/\o033\[0m/g' \
    -e 's/\o033\[0\(\o033\)/\1/' \
    -e 's/\o033\[\(\o033\)/\1/' \
    -e 's/darwin_amd64/XXX/g' \
    -e 's/darwin_arm64/XXX/g' \
    -e 's/linux_amd64/XXX/g' \
    -e 's/linux_arm64/XXX/g' \
    -e 's/v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*[0-9a-zA-Z-]*/vX.X.X/g' \
    -e 's/key ID [0-9A-F][0-9A-F]*/key ID XXX/g' \
    -e 's/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z/XXXX-XX-XXTXX:XX:XXZ/g' \
    -e 's/\(hashicorp\/\)[a-z]* /\1XXX /g' \
    -e 's/\["[0-9][0-9]*s"\]:/["Xs"]:/g' \
    -e 's/\^[0-9][0-9]*/^X/g' \
    -e 's/\([=&+~-]\)[0-9][0-9]*\/[0-9][0-9]*/\1X\/X/g'
