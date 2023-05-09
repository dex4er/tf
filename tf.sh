#!/bin/bash

## https://github.com/dex4er/tf
##
## (c) 2020-2023 Piotr Roszatycki <piotr.roszatycki@gmail.com>
##
## MIT License

VERSION=1.5.0

if [[ -n ${BASH_VERSINFO[0]} ]] && [[ ${BASH_VERSINFO[0]} -le 3 ]] && command -v zsh >/dev/null; then
  exec zsh "$0" "$@"
fi

shopt -s inherit_errexit 2>/dev/null || true

if command -v ggrep >/dev/null; then
  alias grep=ggrep
fi

if command -v gsed >/dev/null; then
  alias sed=gsed
fi

function add_quotes() {
  echo "$1" | sed 's/\[\([a-z_][^]]*\)\]/["\1"]/gi'
}

function filter_manifest_short() {
  grep --line-buffered -v -P '= \(known after apply\)|\(\d+ unchanged \w+ hidden\)|\(config refers to values not yet known\)'
}

function filter_manifest_compact() {
  grep --line-buffered -v -P '^\s\s[\s+~-]|\(config refers to values not yet known\)'
}

function filter_outputs() {
  sed -u '/^Outputs:$/,$d'
}

function filter_progress() {
  local mode="$1"

  local fan="-"
  declare -A statusline

  declare -A progress
  progress=(["-"]='\\' ['\\']="|" ["|"]="/" ["/"]="-" [r]="R" [R]="r" [a]="A" [A]="a" [c]="C" [C]="c" [d]="D" [D]="d")

  local dot_ended=no

  local ignore=no

  local prev=''
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
      case "$mode" in
      fan)
        fan="${progress[$fan]}"
        echo "$fan ${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      dot)
        printf "."
        ;;
      esac
      ;;
    *': Reading...'*)
      case "$mode" in
      fan)
        fan="${progress[$fan]}"
        key="${line%: Reading...*}"
        statusline[$key]="r"
        echo "$fan ${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      dot)
        printf "r"
        ;;
      esac
      ;;
    *': Creating...'*)
      case "$mode" in
      fan)
        fan="${progress[$fan]}"
        key="${line%: Creating...*}"
        statusline[$key]="a"
        echo "$fan ${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      dot) printf "a" ;;
      esac
      ;;
    *': Modifying...'*)
      case "$mode" in
      fan)
        fan="${progress[$fan]}"
        key="${line%: Modifying...*}"
        statusline[$key]="c"
        echo "$fan ${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      dot) printf "c" ;;
      esac
      ;;
    *': Destroying...'*)
      case "$mode" in
      fan)
        fan="${progress[$fan]}"
        key="${line%: Destroying...*}"
        key="${key% (* ????????)}"
        statusline[$key]="d"
        echo "$fan ${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      dot) printf "d" ;;
      esac
      ;;
    *': Still '*'ing... '*)
      case "$mode" in
      fan)
        fan="${progress[$fan]}"
        key="${line%: Still *}"
        statusline[$key]="${progress[${statusline[$key]}]}"
        echo "$fan ${statusline[*]}" | xargs printf "%s"
        printf "\r"
        ;;
      dot) printf "." ;;
      esac
      ;;
    *': '*' complete after '*)
      case "$mode" in
      fan)
        fan="${progress[$fan]}"
        key="${line%: * complete after *}"
        statusline[$key]=" "
        echo "$fan ${statusline[*]}" | xargs printf "%s"
        printf " \r"
        ;;
      dot)
        case "$line" in
        *'Read complete after'*) printf "R" ;;
        *'Creation complete after'*) printf "A" ;;
        *'Destruction complete after'*) printf "D" ;;
        *'Modifications complete after'*) printf "C" ;;
        esac
        ;;
      esac
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
    *'Terraform used the selected providers to generate the following execution'*) echo ;;
    *'plan. Resource actions are indicated with the following symbols:') ;;
    *'Terraform will perform the following actions:'*) ;;
    *'Terraform has compared your real infrastructure against your configuration'*) ;;
    *'and found no differences, so no changes are needed.'*) ;;
    *'Unless you have made equivalent changes to your configuration, or ignored the'*) ;;
    *'relevant attributes using ignore_changes, the following plan may include'*) ;;
    *'actions to undo or respond to these changes.'*) ;;
    *'This is a refresh-only plan, so Terraform will not take any actions to undo'*) ;;
    *'these. If you were expecting these changes then you can apply this plan to'*) ;;
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
    *'Apply complete! Resources: 0 added, 0 changed, 0 destroyed.'*) ;;
    *'─────────────────────────────────────────────────────────────────────────────'*) ;;
    *)
      if [[ $mode == dot ]] && [[ $dot_ended == no ]]; then
        echo
        dot_ended=yes
      fi
      echo "$line"
      prev="$line"
      ;;
    esac
  done
}

declare command="$1"
test $# -gt 0 && shift

case "$command" in

apply | destroy | plan | refresh)
  declare auto_approve=no
  declare tf_log_file
  declare logging="cat"
  if [[ -n $TF_LOG_FILE ]]; then
    tf_log_file=$(LC_ALL=C date "+$TF_LOG_FILE")
    logging="cat 2> >(tee -a \"$tf_log_file\" >&2) | tee -a \"$tf_log_file\""
  fi

  declare filter_manifest="cat"
  declare filter_progress="cat"
  declare filter_outputs="cat"

  declare args=()

  declare default_args=(-short -fan)
  if [[ ${TF_IN_AUTOMATION:-0} == 1 ]]; then
    default_args=(-short -verbose)
  fi

  for arg in "${default_args[@]}" "$@"; do
    case "$arg" in
    -auto-approve) auto_approve=yes ;;
    -compact) filter_manifest="filter_manifest_compact" ;;
    -dot) filter_progress="filter_progress dot" ;;
    -fan) filter_progress="filter_progress fan" ;;
    -full) filter_manifest="cat" ;;
    -no-outputs) filter_outputs="filter_outputs" ;;
    -quiet) filter_progress="filter_progress quiet" ;;
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

  declare filter="$logging | $filter_manifest | $filter_outputs | $filter_progress"
  filter=${filter//| cat /}

  trap '' INT

  case "$command" in
  plan)
    terraform plan -compact-warnings "${args[@]}" | eval "$filter"
    ;;
  apply | destroy)
    trap 'rm -rf terraform.tfplan' EXIT

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

  exec terraform init | eval "$logging" |
    grep --line-buffered -v -P 'Finding .* versions matching|Initializing Terraform|Initializing (modules|the backend|provider plugins)...|Upgrading modules...|Using previously-installed|Reusing previous version of|from the shared cache directory|in modules/' |
    sed -u '/You may now begin working with Terraform/,$d' |
    uniq |
    sed -u '1{/^$/d}'
  ;;

list)
  # trunk-ignore(shellcheck/SC2046)
  terraform state list $(
    for r in "$@"; do
      add_quotes "$r"
    done
  ) | sed -u 's/\x1b\[[01]m//g'
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
      terraform state show -no-color "$r" | sed -u 's/\x1b\[[01]m//g'
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

upgrade)
  declare logging="cat"
  if [[ -n $TF_LOG_FILE ]]; then
    logging="tee -a \"$TF_LOG_FILE\""
  fi

  set -e

  exec terraform init -upgrade | eval "$logging" |
    grep --line-buffered -v -P 'Finding .* versions matching|Initializing Terraform|Initializing (modules|the backend|provider plugins)...|Upgrading modules...|Using previously-installed|Reusing previous version of|from the shared cache directory|in modules/' |
    sed -u '/You may now begin working with Terraform/,$d' |
    uniq |
    sed -u '1{/^$/d}'
  ;;

version)
  echo "tf ${VERSION}"
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
