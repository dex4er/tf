---
name: tf
description: Run Terraform and OpenTofu through tf, the wrapper of dex4er/tf, whose output takes a fraction of the tokens of theirs and whose resource addresses are plain arguments that need no shell quoting. Use when tf is installed or a project asks for it, to read a plan or the result of an apply in fewer tokens, for init, plan, apply, destroy and refresh, the state commands list, show, mv, rm, taint, untaint and import, addresses with indexes such as foo.bar["baz"], composing them with xargs, its plan and progress formats, the full log of TF_OUTPUT_PATH, and which of terraform and tofu it runs.
---

# tf

`tf` runs `terraform` or `tofu` with input and output made for a reader that
pays per token, as an agent does:

- **Output in fewer tokens.** It drops what carries no decision: the refresh
  and progress lines, `(known after apply)`, the plan of data sources, the
  advice to save the plan, the outputs of an apply, and with `-compact` the
  attributes of the resources too. On the demo of tf, 20 resources of
  OpenTofu, `tofu plan`
  prints 20.9 kB, `tf plan -quiet -no-color` 5.2 kB and `tf plan -quiet
  -compact -no-color` 1.1 kB, a twentieth; `tofu apply` prints 25.8 kB and
  `tf apply -quiet -no-color` 5.3 kB.
- **Plain text where it is read by a program.** `tf list` and `tf show` have
  no ANSI codes, so their lines are addresses and attributes, ready for `grep`
  and `xargs`. With `-quiet` there is no spinner, whose redraws would be noise
  in a captured output.
- **Input with no quoting to get wrong.** An address is a positional argument
  instead of `-target=...`, and `tf` puts back the quotes of an index that the
  shell removed, so a command written by an agent needs no nested quoting of
  `foo.bar["baz"]`, where a targeted command is easiest to get wrong.
- **The full output stays available.** `TF_OUTPUT_PATH` writes the verbatim
  output to a file, so the filtered output costs nothing when the details turn
  out to be needed, as after a failed apply.

Prefer `tf` over `terraform` and `tofu` wherever it is installed, and call them
directly only for what `tf` cannot do.

## For an agent

```sh
tf plan -quiet -compact -no-color                  # what changes, one line per resource
tf plan -quiet -no-color foo.bar                   # the attributes of the ones that matter
TF_OUTPUT_PATH=tf.log tf apply -quiet -no-color    # filtered output, verbatim log
```

`tf` keeps the colors of `plan` and `apply`, which are escape codes in a
captured output, so pass `-no-color`; it goes through to terraform or tofu.
Read a big plan with `-compact` first, which lists each resource and its action
without its attributes, and then the resources that matter in full.

## Commands

```sh
tf init                      # init, quiet
tf upgrade                   # init -upgrade
tf validate                  # any other command goes to terraform or tofu as it is
tf plan -quiet               # plan, without the progress indicator
tf plan -quiet foo.bar       # plan -target=foo.bar
tf plan -quiet -out=tfplan   # options other than its own go through unchanged
tf apply -plan=tfplan        # apply the saved plan
tf apply foo.bar foo.baz     # apply -target=foo.bar -target=foo.baz
tf destroy foo.bar           # destroy -target=foo.bar
tf refresh foo.bar           # refresh -target=foo.bar
tf list                      # state list, without colors
tf show                      # show of the whole state
tf show foo.bar foo.baz      # state show of each address
tf mv foo.bar foo.baz        # state mv
tf rm foo.bar foo.baz        # state rm of each address
tf taint foo.bar foo.baz     # taint of each address
tf untaint foo.bar           # untaint of each address
tf import foo.bar some id    # import "foo.bar" "some id", then show of foo.bar
tf version                   # version of terraform or tofu, and of tf
```

## Addresses

An address is a positional argument, never `-target=`, and `tf` adds the
quotes of a string index that the shell has removed: `tf apply foo.bar[baz]`
and `tf apply foo.bar["baz"]` both reach `-target='foo.bar["baz"]'`. A numeric
index, as `foo.bar[0]`, stays as it is. So the addresses need no shell quoting,
and the output of `tf list` feeds the other commands:

```sh
tf list | grep module.old | xargs tf show
tf list | grep random_password | xargs tf rm
tf apply $(tf list | grep aws_vpc)
tf list | grep module.one | while read -r r; do tf mv "$r" "${r/module.one./module.two.}"; done
```

`tf apply` and `tf destroy` ask for a confirmation, which `xargs` does not
pass on, so their addresses go in `$(...)`.

In zsh, brackets are globs, and an unquoted `foo.bar[baz]` fails with `no
matches found` before `tf` runs: use `alias tf='noglob tf'` there, or quote
the address.

## Options

The options of `tf` and of `terraform` go **after** the subcommand: `tf plan
-quiet`, never `tf -quiet plan`, which goes to terraform or tofu as it is, and
which they reject as an unknown option. The only option before the subcommand
is `-chdir=DIR`.

- Plan format: `-short` (the default) drops the noise, `-compact` drops the
  attributes of the resources too, `-full` keeps the original plan. A
  `destroy` with `-short` still lists every attribute that goes to `null`, so
  it saves little there; `-compact` is the one that pays off for it.
- Progress format: `-fan` (the default) is a spinner, `-dots` a character per
  line, `-counters` the counts of the resources, `-verbose` the original lines
  with counters (the default when `TF_IN_AUTOMATION` is set), `-verbatim` the
  original lines, `-quiet` no progress at all. Use `-quiet` when the output is
  read by a program or an agent rather than watched in a terminal.
- `-no-outputs=false` shows the outputs, which `apply`, `refresh` and `show`
  hide by default.
- `TF_PLAN_FORMAT` and `TF_PROGRESS_FORMAT` set the defaults, as
  `TF_PROGRESS_FORMAT=quiet`.
- `-codesign` of `init` and `upgrade` re-signs the providers and drops their
  checksums from the lock file, for a macOS that refuses to run them.

## The full log

What `tf` prints is filtered, so it does not show how far a failed `apply`
went. `TF_OUTPUT_PATH` names a file that gets the verbatim output of `init`,
`upgrade`, `plan`, `apply`, `destroy` and `refresh`:

```sh
: > tf.log
TF_OUTPUT_PATH=tf.log tf apply -quiet foo.bar
```

The file is opened for appending, so truncate it first, or a later read mixes
in an earlier run. Its name goes through `strftime`, so `tf-%Y%m%d.log` is a
file per day and a literal `%` is written `%%`.

## terraform or tofu

`TERRAFORM_PATH` names the command to run, as `TERRAFORM_PATH=tofu`. Without
it, `tf` looks for `.opentofu-version` and `.terraform-version` in the current
directory and its parents: the nearest directory with either decides, and
`.opentofu-version` wins when it has both. The search starts in the directory
that `tf` runs in, not in the one of `-chdir`. Without either file, `tf` runs
`terraform`.

## Exit status

`tf` exits with the status of terraform or tofu, so `tf plan
-detailed-exitcode` still returns 2 for a plan with changes, and a script can
test it as it would test terraform.

## Mistakes

- `apply`, `destroy`, `mv`, `rm`, `taint`, `untaint` and `import` change the
  infrastructure or its state. Run them only when they were asked for, and
  plan before an apply.
- `tf apply foo.bar` is a targeted apply, not a whole one, as an address is
  `-target=`; leave the addresses out to apply everything.
- A summary line of `tf` is not the result of a failed `apply`: read the log of
  `TF_OUTPUT_PATH` or run `tf plan` again before deciding what is left.
- An option before the subcommand fails, as above, and so does a forgotten
  `noglob` in zsh.

The documentation of every command is in the
[README](https://github.com/dex4er/tf#usage) of tf.
