#!/bin/bash
grep="grep --line-buffered -v -P '^\s\s[+-~\s]' | uniq"
args=()
for arg in "$@"; do
  case "$arg" in
    -compact);;
    -short) grep="grep --line-buffered -v -P '(known after apply)'";;
    -full) grep="cat";;
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
terraform plan -destroy -detailed-exitcode ${args[*]} | eval $grep
test ${PIPESTATUS[0]} = 2 || exit $?
echo "[0m[1mDo you want to perform these actions?[0m"
echo "  Terraform will perform the actions described above."
echo "  Only 'yes' will be accepted to approve.:"
echo ""
read -p "  Enter a value: " VALUE
test "$VALUE" = "yes" || exit 0
echo ""
terraform destroy -auto-approve ${args[*]}
exit $?
