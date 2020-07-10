#!/bin/bash
args=()
for arg in "$@"; do
  case "$arg" in
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
terraform plan -detailed-exitcode ${args[*]} | grep --line-buffered -v -P '^\s\s[+-~\s]' | uniq
test ${PIPESTATUS[0]} = 2 || exit $?
echo "Do you want to perform these actions?"
echo "  Terraform will perform the actions described above."
echo "  Only 'yes' will be accepted to approve.:"
echo ""
read -p "  Enter a value: " VALUE
test "$VALUE" = "yes" || exit 0
echo ""
terraform apply -auto-approve ${args[*]}
exit $?
