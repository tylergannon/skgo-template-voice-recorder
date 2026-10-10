# Adversarial review, round 03: released SKGo v0.28.1 ordinary recorder qualification

Date: 2026-10-10. Read-only review; no project file other than this artifact was
written. Rounds 01 and 02 are
`ephemeral/reviews/20261010-released-v0281-qualification.md` and
`ephemeral/reviews/20261010-released-v0281-qualification-round-02.md`.

**Target:** branch `codex/released-v0281-qualification` at commit `ca2be13`
("Correct qualification identifiers and screenshot format"), reviewed as the
whole branch (`main...HEAD`, three commits) against the same sources. The
working tree was clean at review time; the uncommitted edits noted in round 02
were superseded by this commit.

**Scope note:** scope was taken from AGENTS.md, VISION.md, README.md,
`template/README.md` and the artifacts. The caller asked only for a re-review
and the artifact path; no narrowing was requested or applied.

## Evidence inspected

- The `ca2be13` diff in full: README paragraph, the rewritten v0.28.1
  `skgo-feedback.md` entry, `report.md` screenshot links, the worklog, and the
  `web-synthetic.png` to `web-synthetic.jpg` rename.
- Current `README.md:25-37`, `report.md` head, `skgo-feedback.md:73-104`,
  worklog in full.
- Independent checks run this round:
  - `git diff --quiet 118f325 HEAD -- recorder template internal scripts
    go.mod go.sum skgo-plugin.json` is clean; the branch still changes only the
    README paragraph outside `ephemeral/`.
  - Renamed screenshot: SHA256 of `cd6161d:web-synthetic.png` equals the
    current `web-synthetic.jpg` (`b1f3569f…3db6`), and `file` confirms JPEG
    1280x720. No remaining `web-synthetic.png` reference in any `.md` or
    `.json` outside `ephemeral/reviews/`.
  - Identifier grep: `bc267b284c19bab3ae47aaf78bc550c1cbbfcc50` appears intact
    at `skgo-feedback.md:76`; `118f325` and `a481333` are intact in the worklog;
    no `v 0.`, `118f 325`, `a 481333` or `released-v 0` tokens remain in the
    branch's prose.
  - Leak scan of non-ephemeral tracked files for home paths, team IDs and the
    user's email addresses: only the documented `SKGO_DEVELOPMENT_TEAM` variable
    name in README.
  - The round 02 provenance verification (GitHub asset digest and downloaded
    binary hash) was not repeated; the receipts it verified are unchanged in
    `ca2be13`.

## Status of round 02 findings

1. Corrupted identifiers in the feedback file and worklog: **resolved.** The
   entry was rewritten explicitly rather than re-regexed, the SKGo commit and
   receipt path are correct, and the worklog records the cause as friction with
   an instruction not to format identifiers with letter-to-digit regexes.
2. README 1.0.0 procedure unexercised: **resolved.** `README.md:35-37` now says
   the vendor-documented fresh-install procedure was not exercised and that
   qualification reused an existing installation.
3. Screenshot format: **resolved.** Byte-identical rename to `.jpg`; both report
   links updated.
4. README understating frontend coverage: **resolved.** `README.md:28-29` now
   lists Go checks, Svelte diagnostics and the 11 frontend tests.

The worklog additionally corrects its earlier claim that re-review ran in the
same reviewer session, and records the tooling syntax gap that caused a fresh
session. This matches what happened (round 02 and this round ran in session
`3e82c7c7`), and the worklog explicitly says not to claim session continuation.

## Findings

No material findings remain. Two genuine nitpicks:

### 1. Nitpick: the pre-existing feedback entry still has digit-joined versions

`ephemeral/skgo-feedback.md:28-37` (the "Confirmed setup friction" entry) reads
`vp1.1.0/pnpm12.10.1`, `qualifies1.0.0`, `vp1.0.0/pnpm12.9.1`, `pnpm12.10.1`
and `pnpm12.9.1`. These lines exist unchanged on `main`, so they are outside
this branch's diff, but they sit in the same file the framework peer reads and
are the pattern round 01 first flagged. Fix by hand, as `ca2be13` did for the
v0.28.1 entry; do not apply a regex.

### 2. Nitpick: README paragraph has an orphaned short line

`README.md:29-30` wraps as "tests passed, and regular browser computer use" /
"observed synthetic text." after the edit, leaving a short line mid-paragraph.
Cosmetic only.

## What held up

- Sources, template, scripts and module files are byte-identical to `118f325`.
- The v0.28.1 feedback entry carries version, expectation, observation and
  reproduction with correct identifiers, as AGENTS.md requires.
- The README states what was qualified, what was reused, and what was not
  exercised, and keeps v0.28.1 web evidence separate from v0.27.0 native and
  physical evidence.
- Evidence receipts were not altered beyond the byte-identical rename; the
  report's links resolve.
- No credentials, team IDs or home-directory paths entered non-ephemeral files.

## Outcome

only nitpicks remain
