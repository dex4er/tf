#!/bin/bash

trap 'rm -rf terraform.tfplan' EXIT

function filter_manifest_short() {
  grep --line-buffered -v -P '^\s{4}(?!.*[~+/-]\e)|\(known after apply\)'
}

function filter_manifest_compact() {
  grep --line-buffered -v -P '^\s\s[\s+~-]'
}

function filter_terraform_status() {
  declare -A statusline

  declare -A progress
  progress[.]="o"
  progress[o]="O"
  progress[O]="o"

  IFS=''
  while read line; do
    test "$line" != "$prev" || continue
    case "$line" in
      *': Refreshing state...'*)
        printf '.'
        ;;
      *': Modifying...'*)
        key="${line%: Modifying...*}"
        statusline[$key]="."
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Still modifying...'*)
        key="${line%: Still modifying...*}"
        statusline[$key]="${progress[${statusline[$key]}]}"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Modifications complete after '*)
        key="${line%: Modifications complete after *}"
        statusline[$key]="*"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      'Apply complete!'*)
        echo ""
        echo "$line"
        ;;
      *)
        echo "$line"
        prev="$line"
    esac
  done
}

filter="filter_manifest_short | filter_terraform_status"

args=()

for arg in "$@"; do
  case "$arg" in
    -compact) filter="filter_manifest_compact | filter_terraform_status";;
    -short) ;;
    -full) filter="cat";;
    -*) args+=("$arg");;
    *) args+=("-target=$arg")
  esac
done

terraform plan -destroy -detailed-exitcode ${args[*]} | eval $filter

test ${PIPESTATUS[0]} = 2 || exit $?

echo "[0m[1mDo you want to perform these actions?[0m"
echo "  Terraform will perform the actions described above."
echo "  Only 'yes' will be accepted to approve.:"
echo ""

read -p "  Enter a value: " VALUE

test "$VALUE" = "yes" || exit 0

echo ""

terraform destroy -auto-approve ${args[*]} | eval $filter
status=${PIPESTATUS[0]}

exit $status
