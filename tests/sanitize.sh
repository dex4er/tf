#!/usr/bin/env bash

tr '\015' '\012' | \
sed -e '/^ *$/d' \
    -e '/: Creating\.\.\./d' \
    -e '/: Destroying\.\.\./d' \
    -e '/: Reading\.\.\./d' \
    -e '/: Still .*\.\.\./d' \
    -e 's/darwin_amd64/XXX/g' \
    -e 's/darwin_arm64/XXX/g' \
    -e 's/linux_amd64/XXX/g' \
    -e 's/linux_arm64/XXX/g' \
    -e 's/v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*[0-9a-zA-Z-]*/vX.X.X/g' \
    -e 's/key ID [0-9A-F][0-9A-F]*/key ID XXX/g' \
    -e 's/[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z/XXXX-XX-XXTXX:XX:XXZ/g' \
    -e 's/\[id=.*/[id=X]/g' \
    -e 's/\(hashicorp\/\)[a-z]* /\1XXX /g' \
    -e 's/\["[0-9][0-9]*s"\]:/["Xs"]:/g' \
    -e 's/\([=&+~-]\)[0-9][0-9]*\/[0-9][0-9]*/\1X\/X/g' \
    -e 's/ complete after [0-9][0-9]*s/ complete after Xs/g'
