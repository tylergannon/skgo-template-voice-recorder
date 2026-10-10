# Adversarial review, round 02: released SKGo v0.28.1 ordinary recorder qualification

Date: 2026-10-10. Read-only review; no project file other than this artifact was
written. Round 01 is `ephemeral/reviews/20261010-released-v0281-qualification.md`.

**Target:** branch `codex/released-v0281-qualification` at commit `cd6161d`
("Complete tool provenance and qualification feedback"), reviewed as the whole
branch (`main...HEAD`, two commits, 42 files, additions only) against the same
sources as round 01. At review time the working tree also carried uncommitted
edits to `ephemeral/skgo-feedback.md` and
`ephemeral/worklog/released-v0281-qualification.md` (five lines, spacing
only). Both the committed and the working-tree text were inspected; findings
say which they refer to.

**Scope note:** scope was taken from AGENTS.md, VISION.md, README.md,
`template/README.md` and the artifacts. The caller asked only for a re-review
against the same sources and the artifact path; no narrowing was requested or
applied. The round 01 context (installer 299, updater 297 and native runtime are
separately owned) was again treated as context, and the report still claims
nothing in those areas.

## Evidence inspected

- Instructions: AGENTS.md, VISION.md, README.md (current and diff from
  `118f325`), `template/recorder.go` (template summary and option help).
- The `cd6161d` diff in full: README paragraph, `skgo-feedback.md`,
  `report.md`, worklog, `new-help.log`, and the six new receipts
  `vp-provenance.json`, `vp-original-install.log`, `svelte-sync.log`,
  `svelte-check.log`, `web-tests.log`, `web-check-results.json`.
- Carried-over receipts re-read where the new text depends on them:
  `creation-command.json`, `tool-versions.log`, `paths.json`, `new.log`
  lines 135-160, directory listing with timestamps.
- Source invariance: `git diff --quiet 118f325 HEAD -- recorder template
  internal scripts go.mod go.sum skgo-plugin.json` is clean. The only
  non-ephemeral change on the branch is the README paragraph (11 lines).
- Independent checks run this round:
  - GitHub release asset API `assets/594610033` returns
    `vp-aarch64-apple-darwin.tar.gz`, digest
    `sha256:309d4255…944b35`, size 4,399,239, updated 2026-09-28. This matches
    `publishedArchiveDigest` in `vp-provenance.json`.
  - Downloaded that archive into the scratchpad: SHA256 `309d4255…944b35`;
    its `vp` entry hashes to `d787c0a8…9e34`, equal to `vpSHA256` in
    `creation-command.json` and `selectedBinarySHA256` in
    `vp-provenance.json`. The round 01 provenance gap is closed by evidence I
    could reproduce, not by assertion.
  - Fetched `https://viteplus.dev/guide/global-cli`. It documents
    `VP_VERSION` ("Version to install", default latest, Unix example
    `curl -fsSL https://vite.plus | VP_VERSION=1.2.3 bash`) and `VP_HOME`
    ("Set it to an absolute path"), so the README's new instructions cite a
    real vendor procedure.
  - `template/recorder.go:28-29` carries the summary "A native voice recorder
    built on the standard SKGo project" and option help "Base application
    bundle identifier", Required true, matching the two lines appended to
    `new-help.log`.
  - `recorder/` ships two frontend spec files (`TranscriptScreen.svelte.spec.ts`,
    `native-status.spec.ts`), both allowlisted; `web-tests.log` reports four
    files and 11 tests, consistent with those two plus the scaffold's own
    examples. `web-check-results.json` records exit 0 for sync, check and tests.
  - `file` reports `web-synthetic.png` is JPEG data (1280x720), not PNG.
  - Grep of the committed and working-tree prose for digit-split identifiers
    (see finding 1).

## Status of round 01 findings

1. Provenance of the VitePlus 1.0.0 executable: **resolved.** Verified
   independently as above. `vp-original-install.log` additionally shows how
   the retained tree was made (vendor installer put 1.1.0 into the isolated
   `VP_HOME`, then an in-place update moved it to 1.0.0). The README now says
   how to supply 1.0.0 beside a newer global install.
2. Adapter `svelte-kit sync` ENOENT not recorded as feedback: **resolved.**
   New `skgo-feedback.md` entry names adapter 0.22.1, `prerender.js:330`,
   `new.log` 144-152, the reproduction command and the 24 deprecated
   subdependencies. `report.md` "Limits" now describes the same diagnostics.
3. `new-help.log` lacked template discovery: **resolved.** The capture now
   ends with the `voice-notes` template and its `bundle-id` option.
4. Whitespace stripped before numbers: **regressed.** See finding 1.
5. Frontend checks not run on the fresh consumer: **resolved.** Svelte sync,
   svelte-check (0 errors, 0 warnings) and `just test` (11 passed) were run on
   the unchanged consumer at 00:32 and are linked from the report.

## Findings

### 1. Issue: the spacing repair corrupted the identifiers the feedback file exists to carry

Evidence (committed `cd6161d`):

- `ephemeral/skgo-feedback.md:76`:
  `Version: v 0.28.1/bc 267b 284c 19bab 3ae 47aaf 78bc 550c 1cbbfcc 50;`
  The SKGo commit `bc267b284c19bab3ae47aaf78bc550c1cbbfcc50` is split into
  ten tokens and the version reads `v 0.28.1`.
- `ephemeral/skgo-feedback.md:84`: the receipt path reads
  `ephemeral/validation/released-v 0.28.1/`, which does not exist.
- `ephemeral/worklog/released-v0281-qualification.md:1,3,7`: `v 0.28.1`,
  `recorder 118f 325`, `v 0.27/a 481333`.

The uncommitted working tree repairs `v 0.28.1` and the path but still has
`bc 267b 284c 19bab 3ae 47aaf 78bc 550c 1cbbfcc 50` on line 76 and
`118f 325` and `a 481333` in the worklog.

Requirement: AGENTS.md, "Record shared-SKGo gaps in ephemeral/skgo-feedback.md
with version, expectation, observation and reproduction." The version line of
the v0.28.1 entry is now wrong in the committed tree, and the recorder commit in
the worklog is wrong in both trees. Round 01 flagged dropped spaces as a
nitpick because the digits were still contiguous and greppable; a blanket
"insert a space before digits" pass has turned a cosmetic problem into
incorrect identifiers in the handoff document, and the worklog's "review:"
paragraph reports the spacing as "corrected". This is the same authoring
pipeline defect seen in round 01, now in the opposite direction. The fix is to
restore the exact strings `v0.28.1/bc267b284c19bab3ae47aaf78bc550c1cbbfcc50`,
`118f325` and `a481333`, and to stop applying mechanical whitespace rewrites
to hashes, versions and paths.

### 2. Nitpick: the README's 1.0.0 installation procedure is vendor-documented but not the one exercised

`README.md:31-34` instructs readers to install with `VP_VERSION=1.0.0` and an
absolute `VP_HOME`. The vendor guide does document both variables, so the
instruction is sound. But `vp-original-install.log` shows the retained
installation was made differently: the installer placed 1.1.0, then an
in-place update moved it to 1.0.0. The report's "Limits" says a from-scratch
1.0.0 installation is not qualified; the README paragraph does not carry that
caveat, so a reader may take the `VP_VERSION` route as qualified when it was
not run. One clause in the README ("not exercised here") would close this.

### 3. Nitpick: `web-synthetic.png` is a JPEG

`file` identifies `ephemeral/validation/released-v0.28.1/web-synthetic.png` as
JFIF JPEG data, 1280x720. Browsers and the Markdown image link in `report.md`
render it anyway, so the evidence is intact, but the extension misstates the
format and a strict tool or checksum-by-type comparison would reject it. Rename
to `.jpg` and update the two report links, or re-export as PNG.

### 4. Nitpick: README still understates the frontend coverage it now has

`README.md:28` says "the production web/Go build and Go checks passed". After
`cd6161d` the fresh consumer also passed Svelte diagnostics and the 11
frontend tests, and `report.md` says so. The README is accurate but stale
relative to the evidence it links; one phrase ("Go, Svelte and frontend
checks") would align them.

## What held up

- Recorder, template, internal, scripts and module files are byte-identical to
  `118f325`; the branch changes only prose and ephemeral evidence.
- The VitePlus 1.0.0 executable used for creation is byte-identical to the
  official v1.0.0 darwin-arm64 release binary, verified here from GitHub's
  published digest and a fresh download.
- The adapter ordering defect and dependency deprecations are now recorded with
  version, expectation, observation and reproduction.
- `new-help.log` content matches the template definition in `template/recorder.go`.
- The frontend check receipts are consistent with the recorder's allowlisted
  spec files and record exit status 0 for each step.
- The report and README still keep v0.28.1 web evidence separate from v0.27.0
  native and physical evidence and claim no native, installer or updater result.
- No credentials, team IDs or home-directory paths entered non-ephemeral files.

## Outcome

material findings remain
