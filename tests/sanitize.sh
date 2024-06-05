#!/usr/bin/env bash

tr '\015' '\012' | \
sed -e '/^ *$/d' \
    -e '/Your version of .* is out of date! The latest version/d' \
    -e '/You can update by downloading from/d' \
    -e '/: Creating\.\.\./d' \
    -e '/: Destroying\.\.\./d' \
    -e '/: Reading\.\.\./d' \
    -e '/: Still .*\.\.\./d' \
    -e 's/: Creation complete a/: Creation complete/' \
    -e 's/: Creation complete /: Creation complete/' \
    -e 's/: Destruction complete/: Destruction co/' \
    -e 's/: Destruction complet/: Destruction co/' \
    -e 's/: Destruction comple/: Destruction co/' \
    -e 's/: Destruction compl/: Destruction co/' \
    -e 's/: Destruction comp/: Destruction co/' \
    -e 's/: Destruction com/: Destruction co/' \
    -e 's/: Refreshing state\.\.\./: Refreshing state../' \
    -e 's/: Read complete/: Read compl/' \
    -e 's/: Read complet/: Read compl/' \
    -e 's/: Read comple/: Read compl/' \
    -e 's/: Read compl /: Read compl/' \
    -e 's/: Refreshing state\.\. /: Refreshing state../' \
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
    -e 's/\^[0-9][0-9]*/^X/g' \
    -e 's/\([=&+~-]\)[0-9][0-9]*\/[0-9][0-9]*/\1X\/X/g' \
    -e 's/ complete after [0-9][0-9]*s/ complete after Xs/g'
