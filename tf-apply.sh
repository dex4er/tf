#!/bin/bash

trap 'rm -rf terraform.tfplan' EXIT
trap '' INT

function filter_manifest_short() {
  grep --line-buffered -v -P '^\s{4}(?!.*[~+/-]\e)|\(known after apply\)'
}

function filter_manifest_compact() {
  grep --line-buffered -v -P '^\s\s[\s+~-]'
}

function filter_terraform_status() {
  declare -A statusline

  declare -A progress
  progress=(["-"]="\\\\" ["\\\\"]="|" ["|"]="/" ["/"]="-" [r]="R" [R]="r" [a]="A" [A]="a" [c]="C" [C]="c" [d]="D" [D]="d")

  IFS=''
  while read line; do
    test "$line" != "$prev" || continue
    case "$line" in
      *': Refreshing state...'*)
        key="-"
        currentstate="${statusline[$key]:-/}"
        statusline[$key]="${progress[$currentstate]}"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Reading...'*)
        statusline["-"]="*"
        key="${line%: Reading...*}"
        statusline[$key]="r"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Creating...'*)
        statusline["-"]="*"
        key="${line%: Creating...*}"
        statusline[$key]="a"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Modifying...'*)
        statusline["-"]="*"
        key="${line%: Modifying...*}"
        statusline[$key]="c"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Destroying...'*)
        statusline["-"]="*"
        key="${line%: Destroying...*}"
        statusline[$key]="d"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': Still '*'ing...'*)
        statusline["-"]="*"
        key="${line%: Still *ing...*}"
        statusline[$key]="${progress[${statusline[$key]}]}"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      *': '*' complete after '*)
        statusline["-"]="*"
        key="${line%: * complete after *}"
        statusline[$key]="*"
        echo "${statusline[*]}" | xargs printf "%s"
        printf "\r"
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

tf-init

terraform plan -detailed-exitcode ${args[*]} -out=terraform.tfplan | eval $filter

test ${PIPESTATUS[0]} = 2 || exit $?

echo "[0m[1mDo you want to perform these actions?[0m"
echo "  Terraform will perform the actions described above."
echo "  Only 'yes' will be accepted to approve.:"
echo ""

read -p "  Enter a value: " VALUE

test "$VALUE" = "yes" || exit 0

echo ""

terraform apply -auto-approve -refresh=false ${args[*]} terraform.tfplan | eval $filter
status=${PIPESTATUS[0]}

exit $status
