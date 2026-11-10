# Agent Guide: tf

## Overview

`tf` is a wrapper of `terraform` and `tofu` in Go. It runs one of them,
filters what it prints, and takes resource addresses as plain arguments.
`README.md` is for its users; this file describes how it is done and why.

- `main.go` dispatches the subcommand. Only `-chdir=` is accepted before it;
  anything else there becomes the "command" and goes to terraform or tofu,
  which reject it. Unknown subcommands go through unchanged.
- `run/` has a file per subcommand, and a helper for each kind of them:
  `terraform_with_progress.go` for the ones with a plan and a progress
  indicator (`plan`, `apply`, `destroy`, `refresh`, `import`), with the regular
  expressions of what is dropped; `terraform_with_filter.go` for `init` and
  `upgrade`; `terraform_without_colors.go` for `list` and `show`;
  `terraform_with_args_quotes.go` for `mv` and `rm`; and
  `terraform_with_args_iteration.go` for `taint` and `untaint`, which run once
  per address.
  `exec_terraform_command.go` chooses `TERRAFORM_PATH`, or `tofu` or
  `terraform` from the nearest `.opentofu-version` or `.terraform-version`
  (`util/find.go`).
- `util/strings.go` `AddQuotes` puts back the quotes of an index that holds a
  letter or `_`, so `foo.bar[baz]` becomes `foo.bar["baz"]` and `foo.bar[0]`
  stays.
- `util/output.go` opens the log of `TF_OUTPUT_PATH` for appending, its name
  through `strftime`.
- `progress/` holds the progress formats, `console/` the printing that a
  progress line written with `\r` can overwrite.

What `tf` prints, its options and its variables are documented in
`README.md` and in the skill `plugins/tf/skills/tf/SKILL.md`; a change of them
changes both.

## Build and Test

The tools come from `mise.toml`: GoReleaser and yq pinned there, and Go,
Terraform and OpenTofu from `.go-version`, `.terraform-version` and
`.opentofu-version` (`idiomatic_version_file_enable_tools`), so that Renovate
and the workflows read the same versions.

```sh
make build   # goreleaser build --snapshot --single-target, the binary ./tf
make test    # tests/test.sh against ./tf
```

Each `tests/NN_<name>_<terraform|tofu>.sh` copies `demo.tf` into
`tests/tmp`, runs `tf` there with `TERRAFORM_PATH` from the suffix of its
name, and compares the output, after `tests/sanitize.sh`, with its
`tests/NN_<name>_<tool>.out`. A test of a tool that is not installed is
skipped. The tests download the providers of `demo.tf` and sleep in it, so
they need the network and take a while.

The `.out` files keep the ANSI codes that terraform and tofu print, so a
change of the filtering or of the colors changes them. `TEST_UPDATE=1 make
test` rewrites the `.out` of each failing test with its current output; read
the diff of those files before committing it, as it is the specification of
what `tf` prints. A new release of Terraform or OpenTofu that prints
otherwise fails the tests the same way.

## Linting and Formatting

[Trunk](https://docs.trunk.io) runs the linters of `.trunk/trunk.yaml`:
golangci-lint and gofmt for Go, Biome for JSON, Markdown and YAML, yamllint
for what Biome does not check (duplicated keys, values that YAML 1.1 reads
otherwise), shellcheck and shfmt for the scripts, actionlint for the
workflows, taplo for TOML, and the Terraform linters for `demo.tf`.

```sh
trunk fmt
trunk check --no-fix   # without --fix or --no-fix it asks about each autofix and hangs a non-interactive shell
```

- Biome reads `biome.jsonc`. The definition of `biome` in `.trunk/trunk.yaml`
  runs it in the sandbox of Trunk with a copy of `biome.jsonc`, as Biome takes
  the options of its parser, as `markdown.parser.frontmatter` for the front
  matter of `SKILL.md`, only for the files under the directory of its
  configuration.
- Its version is pinned with `!` in `.trunk/trunk.yaml`, which keeps it from
  `trunk upgrade`, and Renovate bumps it. Renovate cannot run `biome migrate`
  or `biome format`, so the workflow Biome Migrate does on its branch and
  commits the result; elsewhere it fails when a file is not migrated or
  formatted. It cannot change a workflow file, which is then formatted by
  hand.
- shfmt runs without `-s`, as any formatting flag makes it ignore
  `.editorconfig` and indent with tabs; `simplify` and `switch_case_indent`
  are in `.editorconfig`.
- actionlint comes from the maintained fork `kjanat/actionlint`, through
  `downloads` in `.trunk/trunk.yaml`, as `rhysd/actionlint` stopped at 1.7.12.
- The workflow Automated trunk upgrade runs `trunk upgrade` and `trunk fmt
  --all` daily in the pull request of `automated/trunk-upgrade`.

## Dependencies

Renovate (`.renovaterc.json`) updates Go, the Go modules, the tools and the
actions.

- A version in a workflow or in `.trunk/trunk.yaml` has a hint in a comment,
  `# datasource=... depName=...`, which a regex manager reads. That manager
  drops the `v` of a tag, so it writes the version without it, whatever the
  file had: `GRYPE_VERSION` is kept without the `v`, which the step that
  installs Grype adds. Keep any version read this way in that form.
- The versions of Go, Terraform and OpenTofu live in their version files, and
  their updates do not touch `mise.lock`; lock file maintenance, at any time
  and automerged, runs `mise lock` after them. The workflow Mise Lock Check
  fails a lock whose checksum changed while its URL did not.
- Terraform is one dependency read from two places, the hints and
  `terraform_version` of `setup-terraform`, so Renovate groups it as a
  monorepo; a package rule titles its pull request "Update dependency
  terraform".
- `gitIgnoredAuthors` lets Renovate keep updating its branches after Biome
  Migrate commits to them as github-actions.
- Running a mise shim in this directory can let mise install the pinned tool
  and rewrite `mise.lock` on its own. Check `git status` after it, and do not
  commit such a change by accident.

## Releases

A tag `vX.Y.Z` on `main` starts the workflow Release; Tag semver creates the
next one from the input `major`, `minor` or `patch`. GoReleaser builds the
binaries of six platforms, as bare binaries, with an SPDX SBOM each, publishes
the release, writes the cask of `dex4er/homebrew-tap`, and the workflow then
attests the build provenance of `checksums.txt`.

- The cask is written with `HOMEBREW_GITHUB_TOKEN` of the environment
  `publishing`. When the token is expired, GoReleaser fails after the release
  is published: the cask stays at the previous version and the attestation
  step does not run. Releasing the same tag again brought back the generated
  notes (v2.14.2), so edit the notes after a run that succeeded.
- The cask clears the quarantine attribute with `install_steps`, which
  generates `postflight_steps`; the raw Ruby of `install` generates
  `postflight`, which Homebrew deprecated and warns about on every command.
  `.goreleaser.yaml` explains why the whole staged directory is cleared.
- The changelog of GoReleaser drops by their titles the updates of CI, of the
  development tools and of Renovate, and keeps those of Go and the Go modules,
  which change the binaries. A commit that changes only CI under another title
  stays in the generated notes and is removed by hand; `changelog.paths`,
  which could tell it, is a feature of GoReleaser Pro. The notes are for the
  users of `tf`: what changes in the binaries, and the vulnerabilities that an
  update of Go fixes, checked with Grype on the SBOMs and with `govulncheck`
  on the code for whether `tf` reaches them.
- Main Scan and Release Scan scan the SBOMs of a snapshot and of a release
  with Grype and report to their discussions.

## Skill for Coding Agents

The skill that teaches an agent to use `tf` is the plugin `tf` of Claude Code
in `plugins/tf`, a skill in each `skills/<name>/SKILL.md`, which the
marketplace of `.claude-plugin/marketplace.json` lists, so that
`claude plugin install tf --marketplace dex4er/tf` installs it. The plugin is
not in the root of the repository, as Claude Code warns of the `CLAUDE.md` of
the root of a plugin. It has no version, so that Claude Code updates it with
each commit, instead of a version that a change of the skill would have to
bump; `claude plugin validate` warns about it, which is expected.

A skill holds the idioms, the commands and the mistakes of its subject, not a
copy of `README.md`, and its description says when it applies. Its front
matter has only the fields of the format of Agent Skills, `name` and
`description`, which other agents read too. The sizes of the output that it
quotes were measured on `demo.tf`; measure them again when the filtering
changes.

```sh
claude plugin validate .
claude plugin validate ./plugins/tf
```
