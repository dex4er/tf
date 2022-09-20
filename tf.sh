#!/bin/bash

## https://github.com/dex4er/tf
##
## (c) 2020-2022 Piotr Roszatycki <piotr.roszatycki@gmail.com>
##
## MIT License

shopt -s inherit_errexit

function add_quotes() {
  echo "$1" | sed 's/\[\([a-z_][^]]*\)\]/["\1"]/g'
}

function filter_manifest_short() {
  grep --line-buffered -v -P '= \(known after apply\)|\(\d+ unchanged \w+ hidden\)|\(config refers to values not yet known\)'
}

function filter_manifest_compact() {
  grep --line-buffered -v -P '^\s\s[\s+~-]|\(config refers to values not yet known\)'
}

function filter_progress_fan() {
  declare fan="-"
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
      fan="${progress[$fan]}"
      echo "$fan ${statusline[*]}" | xargs printf "%s"
      printf "\r"
      ;;
    *': Reading...'*)
      fan="${progress[$fan]}"
      key="${line%: Reading...*}"
      statusline[$key]="r"
      echo "$fan ${statusline[*]}" | xargs printf "%s"
      printf "\r"
      ;;
    *': Creating...'*)
      fan="${progress[$fan]}"
      key="${line%: Creating...*}"
      statusline[$key]="a"
      echo "$fan ${statusline[*]}" | xargs printf "%s"
      printf "\r"
      ;;
    *': Modifying...'*)
      fan="${progress[$fan]}"
      key="${line%: Modifying...*}"
      statusline[$key]="c"
      echo "$fan ${statusline[*]}" | xargs printf "%s"
      printf "\r"
      ;;
    *': Destroying...'*)
      fan="${progress[$fan]}"
      key="${line%: Destroying...*}"
      key="${key% (* ????????)}"
      statusline[$key]="d"
      echo "$fan ${statusline[*]}" | xargs printf "%s"
      printf "\r"
      ;;
    *': Still '*'ing... '*)
      fan="${progress[$fan]}"
      key="${line%: Still *}"
      statusline[$key]="${progress[${statusline[$key]}]}"
      echo "$fan ${statusline[*]}" | xargs printf "%s"
      printf "\r"
      ;;
    *': '*' complete after '*)
      fan="${progress[$fan]}"
      key="${line%: * complete after *}"
      statusline[$key]=" "
      echo "$fan ${statusline[*]}" | xargs printf "%s"
      printf " \r"
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
    *'Unless you have made equivalent changes to your configuration, or ignored the'*) ;;
    *'relevant attributes using ignore_changes, the following plan may include'*) ;;
    *'actions to undo or respond to these changes.'*) ;;
    *'This is a refresh-only plan, so Terraform will not take any actions to undo'*) ;;
    *'these. If you were expecting these changes then you can apply this plan to
'*) ;;
    *'record the updated values in the Terraform state without changing any remote'*) ignore=next ;;
    *'Terraform has checked that the real remote objects still match the result of'*) ;;
    *'your most recent changes, and found no differences.'*) ;;
    *'To perform exactly these actions, run the following command to apply:'*) ;;
    *"Saved the plan to: terraform.tfplan"*) ignore=next ;;
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
  declare auto_approve=no
  declare tf_log_file
  declare logging="cat"
  if [[ -n $TF_LOG_FILE ]]; then
    tf_log_file=$(LC_ALL=C date "+$TF_LOG_FILE")
    logging="tee -a \"$tf_log_file\""
  fi

  declare filter_manifest="cat"
  declare filter_progress="cat"

  declare args=()

  declare default_args=(-short -fan)
  if [[ ${TF_IN_AUTOMATION:-0} == 1 ]]; then
    default_args=(-short -verbose)
  fi

  for arg in "${default_args[@]}" "$@"; do
    case "$arg" in
    -auto-approve) auto_approve=yes ;;
    -compact) filter_manifest="filter_manifest_compact" ;;
    -fan) filter_progress="filter_progress_fan" ;;
    -full) filter_manifest="cat" ;;
    -short) filter_manifest="filter_manifest_short" ;;
    -verbose) filter_progress="cat" ;;
    -*) args+=("$arg") ;;
    *)
      declare r
      r=$(add_quotes "$arg")
      args+=("-target=$r")
      ;;
    esac
  done

  declare filter="$logging | $filter_manifest | $filter_progress"
  filter=${filter//| cat /}

  case "$command" in
  plan)
    terraform plan -compact-warnings "${args[@]}" | eval "$filter"
    ;;
  apply | destroy)
    trap 'rm -rf terraform.tfplan' EXIT
    trap '' INT

    declare workspace
    workspace=$(terraform workspace show 2>/dev/null || true)

    case "$command" in
    apply)
      terraform plan -compact-warnings -detailed-exitcode "${args[@]}" -out=terraform.tfplan | eval "$filter"
      ;;
    destroy)
      terraform plan -destroy -compact-warnings -detailed-exitcode "${args[@]}" -out=terraform.tfplan | eval "$filter"
      ;;
    esac

    test "${PIPESTATUS[0]}" = 2 || exit 0

    if [[ $auto_approve == no ]]; then
      if [[ -n $workspace ]]; then
        echo -n "[0m[1mDo you want to perform these actions in workspace \"$workspace\"?[0m "
      else
        echo -n "[0m[1mDo you want to perform these actions?[0m "
      fi

      trap - INT
      read -r VALUE
      trap '' INT

      test "$VALUE" = "yes" || exit 0

      echo ""
    fi

    terraform apply -compact-warnings -auto-approve -refresh=false "${args[@]}" terraform.tfplan | eval "$filter"
    ;;
  refresh)
    terraform apply -compact-warnings -auto-approve "${args[@]}" -refresh-only | eval "$filter"
    ;;
  esac

  exit "${PIPESTATUS[0]}"
  ;;

import)
  declare args=()
  declare show=no
  declare src

  while [[ ${1#-} != "$1" ]]; do
    if [[ $1 == "-show" ]]; then
      show=yes
    else
      args+=("$1")
    fi
    shift
  done

  if [[ $# -lt 2 ]]; then
    terraform import "${args[@]}"
    exit $?
  fi

  src=$(add_quotes "$1")
  shift

  terraform import "${args[@]}" "$src" "$*" | grep --line-buffered -v -P '(The resources that were imported are shown above. These resources are now in|your Terraform state and will henceforth be managed by Terraform.)'

  if [[ $show == "yes" ]]; then
    terraform state show -no-color "$src"
  fi
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
  # trunk-ignore(shellcheck/SC2046)
  terraform state list $(
    for r in "$@"; do
      add_quotes "$r"
    done
  ) | sed 's/\x1b\[[01]m//g'
  exit "${PIPESTATUS[0]}"
  ;;

mv)
  set -e
  # trunk-ignore(shellcheck/SC2046)
  terraform state mv $(
    for r in "$@"; do
      add_quotes "$r"
    done
  )
  ;;

rm)
  set -e
  # trunk-ignore(shellcheck/SC2046)
  terraform state rm $(
    for r in "$@"; do
      add_quotes "$r"
    done
  )
  ;;

show)
  if [[ $# -gt 0 ]]; then
    for r in "$@"; do
      declare r
      r="$(add_quotes "$r")"
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
    declare r
    r=$(add_quotes "$r")
    terraform taint "$r"
  done
  ;;

untaint)
  set -e
  for r in "$@"; do
    declare r
    r=$(add_quotes "$r")
    terraform untaint "$r"
  done
  ;;

version)
  echo tf 1.0.0
  terraform version
  ;;

*)
  if command -v "tf-$command" >/dev/null; then
    exec "tf-$command" "$@"
  else
    exec terraform "$command" "$@"
  fi
  ;;
esac
