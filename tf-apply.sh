#!/bin/bash
terraform plan -detailed-exitcode "$@" | grep --line-buffered -v -P '^\s\s[+-~\s]' | uniq
test ${PIPESTATUS[0]} = 2 || exit $?
echo "Do you want to perform these actions?"
echo "  Terraform will perform the actions described above."
echo "  Only 'yes' will be accepted to approve.:"
echo ""
read -p "  Enter a value: " VALUE
test "$VALUE" = "yes" || exit 0
echo ""
terraform apply -auto-approve "$@"
exit $?
