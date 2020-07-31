#!/bin/bash
trap 'rm -rf terraform.tfplan' EXIT
grep="grep --line-buffered -v -P '^\s{4}(?!.*[~+/-]\e)|\(known after apply\)' | uniq"
args=()
for arg in "$@"; do
  case "$arg" in
    -compact) grep="grep --line-buffered -v -P '^\s\s[\s+~-]' | uniq";;
    -short) ;;
    -full) grep="cat";;
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done
terraform plan -detailed-exitcode ${args[*]} -out=terraform.tfplan | eval $grep
test ${PIPESTATUS[0]} = 2 || exit $?
echo "[0m[1mDo you want to perform these actions?[0m"
echo "  Terraform will perform the actions described above."
echo "  Only 'yes' will be accepted to approve.:"
echo ""
read -p "  Enter a value: " VALUE
test "$VALUE" = "yes" || exit 0
echo ""
terraform apply -auto-approve -refresh=false ${args[*]} terraform.tfplan
status=$?
exit $status
