# Adversarial review, round 02: released SKGo v0.28.2 updater safe refusal

Date: 2026-10-10. Read-only review; no project file other than this artifact was
written. Round 01 is `ephemeral/reviews/20261010-released-v0282-updater.md`.

**Target:** branch `codex/released-v0282-updater` at commit `f933fe5`
("Clarify publication timing and silent installation evidence"), reviewed as
the whole branch (`main...HEAD`, two commits) against AGENTS.md, VISION.md,
README.md, `ephemeral/validation/released-update-v0.28.2/report.md` and its
receipts, the `skgo-feedback.md` entry and the worklog. The working tree was
clean at review time; the uncommitted edits noted in round 01 are now part of
`f933fe5`.

**Scope note:** the caller asked for a re-review of the entire current work
against the same sources and gave the artifact path. No narrowing was requested
or applied.

## Evidence inspected

- The `f933fe5` diff in full: report rewording, feedback entry additions,
  worklog "limitation" and "review" lines, and six new receipt files
  (`install-receipt.json`, `publication-timing.json`,
  `file-modes-after.json`, and the `prior-v0.28.0/` copies).
- The round 01 verification results, which remain valid: tag `v0.28.2` ->
  `458aa5d2`, release gate success, host SHA matches the live binary, the four
  120-entry inventories are identical and match a live rehash of both trees,
  authored receipts match `c1cc76d`, and all `input-*` files match the
  inventory.
- Independent checks run this round:
  - All seven `prior-v0.28.0/` files are byte-identical to
    `c1cc76d:ephemeral/validation/update-v0.28.0/`.
  - `file-modes-after.json` holds two 120-entry mode maps (baseline, copy)
    that are identical to each other and agree with the live `stat` comparison
    from round 01.
  - `install-receipt.json` records the `go install ...@v0.28.2` argv, GOBIN,
    exit 0 and zero captured bytes; the report now says the empty
    `host-install.log` is intentional and no longer links it as evidence.
  - `publication-timing.json`: npm reports `@skgo/sveltekit-adapter@0.28.2`
    published 2026-10-10T07:53:32Z; the invocation receipt time
    08:07:39.395Z matches the live host directory mtime (02:07:39 -0600).
    Costa Rica is UTC-6, so the report's 01:53:32 local figure is correct.
  - Upstream attribution: the v0.28.0...v0.28.1 compare contains only #301;
    v0.28.1...v0.28.2 contains #303 and #304. The feedback entry's statement
    that repair PR #303 is "included in v0.28.2" is correct, and the worklog's
    note that round 01 misattributed #303 to v0.28.1 is also correct.
  - `git diff --stat main...HEAD -- . ':!ephemeral'` is empty: the branch
    changes nothing outside `ephemeral/`.
  - Leak scan of the new diff: the copied `prior-v0.28.0/` receipts and
    `update.log` contain the user's home path (`/Users/tyler/go/bin/skgo`,
    shell startup file paths). These were already present at `c1cc76d`, live
    under `ephemeral/`, and contain no credentials, team IDs or email.

## Disposition of round 01 findings

1. Timing dependence (issue): addressed. `report.md:70-79` and
   `skgo-feedback.md:119-125` now give the publication and receipt times, the
   likely release-age dependence, a reproduction condition, and an explicit
   refusal to assert a later outcome. The outcome line is qualified as
   "incomplete during fresh-package publication".
2. Empty `host-install.log` (issue): addressed by `install-receipt.json` and
   the report's explicit statement that the log is intentionally empty.
3. Inventory coverage wording (nitpick): addressed; the report now says
   regular-file SHA256 only, no symlinks exist, and a separate after-mode
   comparison was added with its "no before-mode" limit stated.
4. Upstream repair not named (nitpick): addressed; #303 and the moved refusal
   point are now in the feedback entry.
5. Unmerged-branch link (nitpick): addressed by copying the `c1cc76d`
   receipts verbatim into `prior-v0.28.0/`.

## Findings

No material findings remain. Genuine nitpicks:

### 1. Nitpick — report line 3 breaks the file's wrap width

`report.md:3` is now about 110 characters while the rest of the document wraps
near 80. Cosmetic only.

### 2. Nitpick — home-directory paths carried into the repository by the copied receipts

`prior-v0.28.0/normal-startup-before.json:2-4` and
`prior-v0.28.0/update.log:3` (and the current `update.log:3`) embed
`/Users/tyler/...` paths. AGENTS.md's exclusion of absolute checkout paths is
aimed at packaged template content, and these files are validation artifacts
under `ephemeral/`, so this is not a violation; it is noted because the branch
newly duplicates them from an unmerged commit.

### 3. Nitpick — the Costa Rica framing is unexplained

`report.md:70-71` states the publication time "in Costa Rica" without saying why
that zone is the reference. The machine zone is -0600 and the figures are
right; a reader outside the project may not know the zone is the operator's.

## Outcome

only nitpicks remain
