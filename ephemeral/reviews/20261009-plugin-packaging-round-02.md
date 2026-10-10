# Recorder source-install packaging review — round 02

Outcome: **no findings**.

## Target and requirements

Reviewed `/Users/tyler/src/.worktrees/recorder-plugin-install` at
`e69d4e7`, including the current uncommitted Go directive, documentation,
metadata and Go 1.27.1 test evidence. Requirements came from repository
`AGENTS.md`, `VISION.md`, the updated source-install contract at
`/Users/tyler/.codex/worktrees/plugin-install-contract/skgo-project/ephemeral/plans/plugin-source-install.md`,
and the authoritative request to support actual Go 1.27.1 hosts while preserving
native evidence, explaining manual/managed conflict remedies, and keeping the PR
draft pending the second shared contract review.

The caller's read-only boundary and report ownership were observed. No caller
instruction restricted the defects considered or predicted a verdict.

## Evidence inspected

- Full current diff, cumulative change statistics and canonical-source diff
  against the released baseline `178135e`.
- Root and nested recorder module declarations, metadata, current consumer and
  source-packaging documentation, existing generation implementation and the
  prior review's payload/factory inspection.
- Packaging metadata, `go1271-tests.log`, and packaging worklog.
- Read-only GitHub PR inspection: recorder PR #2 is open and remains a draft.
- Independently executed `GOTOOLCHAIN=go1.27.1 GOWORK=off go version` and
  `GOTOOLCHAIN=go1.27.1 GOWORK=off go test ./...`: exact Go 1.27.1 was reported;
  root tests passed (cached).

## Assessment

No material defect or genuine nitpick was found. The sole module change lowers
the plugin root directive from Go 1.27.2 to Go 1.27.1. Dependencies, canonical
application source, staging code, payload allowlist and retained native
qualification artifacts are unchanged. The nested recorder module already uses
Go 1.27.1. Recorded generation and tests use the explicitly selected Go 1.27.1
toolchain with workspace resolution disabled; the independent test invocation
confirms the root packages work under that toolchain.

The documentation distinguishes the plugin's minimum Go version from the Go
1.27.2 used for historical native qualification. It marks source installation and
removal as pending release, labels manual qualification historical, warns against
duplicate discovered recorder copies, and correctly explains that changing
`SKGO_PLUGIN_DIRS` does not hide managed installations. Removing the manual file
or invoking the planned current-host managed removal are consistent with the
updated shared contract. The PR remains draft as requested.

The review does not establish actual-host Go-plugin ABI qualification, automatic
source installation or new consumer/native runtime behavior. Those observations
remain dependent on the shared installer and its acceptance work. The current
documentation and metadata do not claim them. Retained native and physical-device
evidence remains attached to its original identified products.
