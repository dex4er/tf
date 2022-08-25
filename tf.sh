#!/bin/bash

## https://github.com/dex4er/tf
##
## (c) 2020-2022 Piotr Roszatycki <piotr.roszatycki@gmail.com>
##
## MIT License

function filter_manifest_short() {
  grep --line-buffered -v -P '\(known after apply\)|\(\d+ unchanged \w+ hidden\)'
}

function filter_manifest_compact() {
  grep --line-buffered -v -P '^\s\s[\s+~-]'
}

function filter_terraform_status() {
  declare -A statusline

  declare -A progress
  # trunk-ignore(shellcheck/SC1003)
  progress=(["-"]='\\' ['\\']="|" ["|"]="/" ["/"]="-" [r]="R" [R]="r" [a]="A" [A]="a" [c]="C" [C]="c" [d]="D" [D]="d")

  ignore=no

  IFS=''
  while read -r line; do
    if [[ $ignore == "yes" ]]; then
      continue
    fi
    if [[ $ignore == "next" ]]; then
      ignore=no
      continue
    fi
    if [[ $line == "$prev" ]]; then
      continue
    fi
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
      key="${key% (* ????????)}"
      statusline[$key]="d"
      echo "${statusline[*]}" | xargs printf "%s"
      printf "\r"
      ;;
    *': Still '*'ing... '*)
      statusline["-"]="*"
      key="${line%: Still *}"
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
    *'Warning:'*'Applied changes may be incomplete'*) ignore=yes ;;
    *'Warning:'*'Resource targeting is in effect'*) ignore=yes ;;
    'This plan was saved to: terraform.tfplan') ignore=yes ;;
    *'suggests to use it as part of an error message'*)
      ignore=no
      continue
      ;;
    *'exceptional situations such as recovering from errors or mistakes'*)
      ignore=no
      continue
      ;;
    *'terraform apply "terraform.tfplan"')
      ignore=no
      continue
      ;;
    *'Terraform used the selected providers to generate the following execution'*)
      echo
      ;;
    *'plan. Resource actions are indicated with the following symbols:') ;;
    *'Terraform will perform the following actions:'*) ;;
    *'Terraform has compared your real infrastructure against your configuration'*) ;;
    *'and found no differences, so no changes are needed.'*) ;;
    *'To see the full warning notes, run Terraform without -compact-warnings.'*) ;;
    *'Acquiring state lock. This may take a few moments...'*) ;;
    *'Releasing state lock. This may take a few moments...'*) ;;
    *'Warnings:'*) ;;
    *'- Experimental feature '*' is active'*) ignore=next ;;
    *'Note: You '*' use the -out option to save this plan, so Terraform'*) ;;
    *'guarantee to take exactly these actions if you run "terraform apply" now.'*) ;;
    *'─────────────────────────────────────────────────────────────────────────────'*) ;;
    *)
      echo "$line"
      prev="$line"
      ;;
    esac
  done
}

declare command="$1"
shift

case "$command" in

apply | destroy | plan | refresh)
  declare tf_log_file
  declare logging=""
  if [[ -n $TF_LOG_FILE ]]; then
    tf_log_file=$(LC_ALL=C date "+$TF_LOG_FILE")
    logging="tee -a \"$tf_log_file\" | "
  fi

  declare filter="${logging}filter_manifest_short | filter_terraform_status"

  declare args=()

  for arg in "$@"; do
    case "$arg" in
    -compact) filter="${logging}filter_manifest_compact | filter_terraform_status" ;;
    -short) ;;
    -full) filter="${logging}" ;;
    -*) args+=("$arg") ;;
    *) args+=("-target=$arg") ;;
    esac
  done

  case "$command" in
  plan)
    terraform plan -compact-warnings "${args[@]}" | eval "$filter"
    ;;
  apply | destroy)
    trap 'rm -rf terraform.tfplan' EXIT
    trap '' INT

    case "$command" in
    apply)
      terraform plan -compact-warnings -detailed-exitcode "${args[@]}" -out=terraform.tfplan | eval "$filter"
      ;;
    destroy)
      terraform plan -destroy -compact-warnings -detailed-exitcode "${args[@]}" -out=terraform.tfplan | eval "$filter"
      ;;
    esac

    test "${PIPESTATUS[0]}" = 2 || exit 0

    echo "[0m[1mDo you want to perform these actions?[0m"
    echo "  Terraform will perform the actions described above."
    echo "  Only 'yes' will be accepted to approve.:"
    echo ""

    read -r -p "  Enter a value: " VALUE

    test "$VALUE" = "yes" || exit 0

    echo ""

    terraform apply -compact-warnings -auto-approve -refresh=false "${args[@]}" terraform.tfplan | eval "$filter"
    ;;
  refresh)
    terraform apply -compact-warnings "${args[@]}" -refresh-only | eval "$filter"
    ;;
  esac

  exit "${PIPESTATUS[0]}"
  ;;

init)
  declare logging="cat"
  if [[ -n $TF_LOG_FILE ]]; then
    logging="tee -a \"$TF_LOG_FILE\""
  fi

  set -e

  exec terraform init -upgrade | eval "$logging" |
    grep -v -P 'Finding .* versions matching|Initializing (modules|the backend|provider plugins)...|Upgrading modules...|Using previously-installed|Reusing previous version of|from the shared cache directory|in modules/' |
    sed '/Terraform has been successfully initialized/,$d' |
    uniq |
    sed '1{/^$/d}'
  ;;

list)
  terraform state list "$@" | sed 's/\x1b\[[01]m//g'
  exit "${PIPESTATUS[0]}"
  ;;

mv)
  set -e
  terraform state mv "$@"
  ;;

rm)
  set -e
  terraform state rm "$@"
  ;;

show)
  if [[ $# -gt 0 ]]; then
    for r in "$@"; do
      terraform state show -no-color "$r" | sed 's/\x1b\[[01]m//g'
    done
  else
    terraform show -no-color
    exit $?
  fi
  ;;

taint)
  set -e
  for r in "$@"; do
    terraform taint "$r"
  done
  ;;

untaint)
  set -e
  for r in "$@"; do
    terraform untaint "$r"
  done
  ;;

version)
  echo tf 1.0.0
  terraform version
  ;;

*)
  exec "tf-$command" "$@"
  ;;
esac
