# Recorder source-install packaging review — round 03

Outcome: **only nitpicks remain**.

## Target and requirements

Reviewed `/Users/tyler/src/.worktrees/recorder-plugin-install` on branch
`codex/plugin-install-packaging` at `385ffbc` (clean tree), as the `main...HEAD`
diff of recorder PR #2 (draft, open). The diff touches `README.md`,
`template/README.md`, `skgo-plugin.json`, the root `go.mod` Go directive, the
packaging worklog, retained evidence under
`ephemeral/validation/plugin-install-packaging/`, and the two prior review rounds.
No Go source, allowlist, generator, canonical `recorder/` module or native code
changes.

Requirements came from `AGENTS.md`, `VISION.md`, `README.md` and the chosen
source-install contract. The caller named
`/Users/tyler/.codex/worktrees/plugin-install-contract/skgo-project/ephemeral/plans/plugin-source-install.md`;
that worktree directory now contains only a `.codex-worktree-name` marker and no
plan file. The same plan was read instead from the canonical checkout
`/Users/tyler/src/skgo-project/ephemeral/plans/plugin-source-install.md`, last
changed by commit `1d72fea` ("Define source plugin installation and host
discovery (#143)"), with no local modifications. Rounds 01 and 02 cite the now
absent path; future rounds should cite the merged skgo-project location.

The caller's statements that SKGo owns installer implementation and release
qualification, and that this PR prepares packaging and marks commands pending,
were treated as scope description, not as a limit on defects considered. The
read-only boundary and the artifact path were observed. No caller instruction
predicted a verdict.

## Evidence inspected

- Full `main...HEAD` diff and every branch commit message.
- `skgo-plugin.json`, root and `recorder/` `go.mod`, `.gitignore`,
  `template/recorder.go`, `template/main.go`, `template/stage/main.go`,
  `template/allowlist.json`, `scripts/build-plugin.sh`, both READMEs in full.
- Retained evidence: `metadata.json`, `go-tests.log`, `go1271-tests.log`,
  `tidy-before-generate.log`; packaging worklog; `ephemeral/skgo-feedback.md`;
  prior review rounds 01 and 02.
- Contract conformance checks run read-only against the tree:
  - Metadata is exactly `{"package":"./template"}`; `template` is `package main`
    with the `//go:generate go run ./stage` directive and the embed of
    `all:payload allowlist.json`.
  - All 47 allowlist sources are Git-tracked regular files (no mode 120000
    entries), none lives under `.git`, `node_modules`, `ephemeral`, `.local`,
    `.cache`, `build` or `.build`, and none is a `.so`. Executable bits are
    carried by the allowlist's explicit flag, not by filesystem mode, so the
    sanitized stage cannot lose them. The generator imports only the standard
    library and reads nothing from Git metadata.
  - `go.work` is absent for this tree; `GOWORK=off` has no hidden effect.
  - SKGo tags `v0.27.0` and `v0.28.0` exist; `cmd/skgo/main.go` at `v0.28.0`
    dispatches no `plugin` subcommand, so "pending release" is accurate for both.
  - Both SKGo tags declare `go 1.27.1`, matching the lowered plugin directive.
- Independent execution in a disposable copy of `git archive HEAD` with
  `ephemeral/` removed, mimicking a sanitized stage, under
  `GOTOOLCHAIN=go1.27.1 GOWORK=off GOFLAGS=-mod=mod` (Go reported
  `go1.27.1 darwin/arm64`):
  - `go mod tidy -e` with no payload: exit 0.
  - Pinned `go.mod` restored, `go generate ./template`: payload staged.
  - Strict `go mod tidy`: exit 0.
  - `go test ./...`: `internal/packagefiles` and `template` pass.
  - `go build -buildmode=plugin ./template`: compiles (7.2 MB). This shows the
    plugin package is buildable in plugin mode at Go 1.27.1 only; it is not host
    ABI qualification, load proof or installer proof.
  - `go run github.com/tylergannon/skgo/cmd/skgo-plugin-build --help` succeeds
    both with the committed `go.mod` and after strict tidy.

## Assessment

No material defect was found. The PR delivers exactly the recorder-owned part of
the contract: the one-field manifest, a package that already exports the factory
and host-version injection variable, author documentation reproducing the
contract's preparation order (host pins, preparatory `tidy -e`, complete pinned
`go.mod` restoration retaining `go.sum`, `go generate`, strict tidy, host graph
verification, `GOFLAGS=-mod=mod`, `GOWORK=off`, sanitized stage exclusions,
symlink rejection, no `git describe`), the nested `recorder/` module fetch
requirement, the Go 1.27.1 minimum required of the recorder owner with retained
exact-toolchain evidence, and the historical labelling of the manual helper with
the manual/managed same-host conflict warning and the planned
`skgo plugin remove MODULE` remedy. No installer behavior, automatic discovery,
new consumer or native runtime behavior is claimed. Retained v0.27.0 native and
physical-device evidence is not relabelled. The PR remains draft as the worklog
states.

## Findings

### Nitpick 1 — root `go.mod` is not tidy at Go 1.27.1

`go.mod:11-23` carries ten `// indirect` requirement lines (goja, goja_nodejs,
regexp2, sourcemap, pprof, devalue, polytype, x/net, x/sys, x/text). Strict
`go mod tidy` under Go 1.27.1 removes all of them, and `go mod tidy -diff` on the
committed file exits 1. The retained `tidy-before-generate.log` records the tidy
command but not that it rewrote `go.mod`. Impact: none for the planned installer,
whose contract explicitly tolerates redundant requirement lines and rewrites host
pins; the manual helper also still resolves after tidy. A tidy-check in CI would
fail, and a reader of the log may assume the committed file is already tidy.
Either tidy the file or note the rewrite in the log.

### Nitpick 2 — consumer-facing install section omits the plugin's Go minimum

`README.md:54-60` tells a consumer that "Git, Go and a working C toolchain are
required" for source installation. The contract has the installer reject a plugin
whose `go` directive exceeds the actual host toolchain, naming both versions, and
asks the recorder owner to qualify at or document the host Go minimum. The
qualified minimum (Go 1.27.1) appears only in `template/README.md:33-34`. Impact:
a consumer on an older host toolchain learns the minimum from an installer
diagnostic rather than from the consumer documentation. One clause in the root
README would close it.

### Nitpick 3 — prior review rounds cite a plan path that no longer exists

`ephemeral/reviews/20261009-plugin-packaging-round-01.md` and `round-02.md` cite
the `.codex/worktrees/plugin-install-contract/...` plan path, which is now empty.
The merged plan lives in skgo-project at `1d72fea`. Impact: a reader of the
retained reviews cannot reopen the contract they were judged against without
locating it elsewhere. Record the merged location in the worklog or PR body.

## Limits

This review did not run the unreleased installer, build a plugin against an
actual host, load a plugin, create a consumer, or exercise microphone or native
behavior. Those remain SKGo's installer qualification and later fresh-caller
evidence, and the documentation under review does not claim them.
