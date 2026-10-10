# Recorder source-install packaging review — round 01

Outcome: **no findings**.

## Target and requirements

Reviewed the recorder packaging worktree at
`/Users/tyler/src/.worktrees/recorder-plugin-install`, based on commit
`51553d793142f94d4a18b202aaa4baf6cb800749`, including the uncommitted
`README.md`, `template/README.md`, `skgo-plugin.json`, and packaging evidence.
The requested outcome is recorder-owned source-install metadata and documentation;
the shared installer is pending release.

Requirements were reconstructed from this repository's `AGENTS.md`, `VISION.md`,
`README.md`, and the proposed shared contract at
`/Users/tyler/.codex/worktrees/plugin-install-contract/skgo-project/ephemeral/plans/plugin-source-install.md`.
The caller's file ownership and read-only constraints were observed. No caller
instruction restricted the defects considered or predicted the outcome.

## Evidence inspected

- Full current documentation diff and the preceding packaging-constraint checkpoint.
- Root module, Git ignore rules, metadata, existing `template/main.go`,
  `template/apply.go`, `template/recorder.go`, staging implementation, allowlist,
  packaging tests, and source-build wrapper.
- `ephemeral/worklog/20261009-plugin-install-packaging.md` and
  `ephemeral/validation/plugin-install-packaging/{metadata.json,go-tests.log}`.
- Retained released fresh-consumer qualification report, including distinctions
  between synthetic runtime, signed physical launch, earlier real-speech
  observation, and checkout-independent rebuild evidence.
- Independent read-only assertions confirmed the exact one-field metadata,
  package-main directory, existing generation directive, version-injection
  variable, and all 47 allowlisted sources being Git-tracked regular files whose
  bytes match the staged payload.

## Assessment

No material defect or genuine nitpick was found. The root metadata is exactly
`{"package":"./template"}`, matching the proposed v1 schema. Its package already
exports the supported factory and `main.skgoVersion` injection variable.
`go generate ./template` delegates to the existing allowlisted stage; no alternate
SDK, preparation hook, payload tree, or recorder runtime change is introduced.

Both documentation additions correctly explain the full-repository requirement:
the nested `recorder/` module supplies canonical payload inputs that a root module
zip omits. They describe actual-host installation and distinguish planned automatic
discovery from the qualified released manual procedure. The proposed command is
explicitly marked pending, and its source-install/discovery behavior is explicitly
unproved. Existing native and fresh-consumer observations are not relabeled as
installation evidence.

The retained Go test output supports packaging regression only. This review did
not execute the pending installer, generate another consumer, or establish new
microphone/transcription behavior. Shared installation acceptance remains for the
shared SKGo implementation and subsequent fresh-caller evidence; that pending
work is not a recorder packaging defect.
