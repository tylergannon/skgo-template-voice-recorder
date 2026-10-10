# Adversarial review: released SKGo v0.28.2 updater safe refusal

Date: 2026-10-10. Read-only review; no project file other than this artifact was
written.

**Target:** branch `codex/released-v0282-updater` at commit `1b118fe` ("Record
released v0.28.2 updater safe refusal"), reviewed as `main...HEAD` (one commit,
24 files, all under `ephemeral/`) against AGENTS.md, VISION.md, README.md,
`ephemeral/validation/released-update-v0.28.2/report.md` and its receipts, the
new `ephemeral/skgo-feedback.md` entry and the new worklog.

**Scope note:** the caller named the target, the coordinator's authorization
(one ordinary released updater attempt on a fresh unpatched copy of the held
edited v0.27.0 baseline) and the artifact path. No narrowing was requested or
applied; scope was taken from the repository instructions and artifacts.

## Evidence inspected

- The full branch diff: report, 22 receipt files, feedback entry, worklog.
- AGENTS.md, VISION.md, README.md; the prior round-03 review for format.
- The original v0.28.0 refusal report and `after.json` at `c1cc76d`
  (`ephemeral/validation/update-v0.28.0/`), which is on the remote branch
  `codex/released-update-qualification` and **not** on `main`.
- Upstream: issue tylergannon/skgo#297 (still OPEN), the v0.28.0...v0.28.2
  compare (three commits: #301, #303 "preserve authored apps when aligning
  VitePlus", #304), PR #303's description, and `internal/update/frontend.go`
  at tag v0.28.2 (the guard compares `package.json` and `pnpm-workspace.yaml`
  and names the changed key).
- Independent checks run this round:
  - Tag `v0.28.2` is an annotated tag whose object is
    `458aa5d204af8217afbed126981bf4154d52b238`; release-gate run 38035420960
    concluded `success` on that SHA. Both match the report.
  - `host-buildinfo.json` equals `host-buildinfo-after.json`; the live host at
    `/private/tmp/recorder-update-v0282-iup2x7f8/host/skgo` hashes to
    `300a6ee3…48c3`, matching `invocation.json` and `preservation.json`.
  - `baseline-before.json` == `consumer-before.json` == `baseline-after.json`
    == `consumer-after.json` (120 entries each). Every entry of the live copy
    and the live baseline rehashes to the inventory; no file is missing or
    extra; file modes are identical across both trees; neither tree contains a
    symlink.
  - The four `prior-authored-receipts.json` hashes equal the `c1cc76d`
    `after.json` hashes and the current inventory entries.
  - All six retained `input-*` files rehash to their inventory entries.
  - `update.log`, `update-result.json`, `update-help.log` and
    `tool-versions-before.log` are consistent with the report's narrative; the
    caller's argv contains no `--complete-version` or override flags.
  - Leak scan of the diff (excluding the lockfile): no home paths, team IDs or
    the user's email; only documented `/private/tmp/...` qualification paths.
  - Timing: the host directory was created 2026-10-10T02:07:39-0600 (08:07Z)
    and the commit landed 02:10:36-0600. npm reports
    `@skgo/sveltekit-adapter@0.28.2` published 2026-10-10T07:53:32Z, roughly
    14 minutes before the run. `pnpm config get minimumReleaseAge` is unset in
    the user's environment.

The preservation claim, the identity of the host, the exact command, and the
"no patch, no retry" constraint from the coordinator are all substantiated by
the receipts. The report correctly refrains from any updated native-app or
runtime claim and leaves #297 open.

**Working-tree note:** during this review, uncommitted edits (mtime
2026-10-10 02:11:39, after the `1b118fe` commit) appeared in `report.md:72-75`
and the worklog, adding a "limitation" that the updater's temporary stage was
removed after refusal, so no staged after-bytes of `pnpm-workspace.yaml` are
retained and the pnpm log line is the only record of the added entry. Those
edits were not made by this review. They are accurate against `update.log`
and do not change the findings below; finding 1 still applies to that text.

## Findings

### 1. Issue — the refusal is time-dependent and neither the report nor the feedback entry says so

`update.log:9-10` shows pnpm adding `@skgo/sveltekit-adapter@0.28.2` to
`minimumReleaseAgeExclude` because the package was newer than pnpm 12.9.1's
release-age policy at install time; the adapter was published at 07:53Z and the
run began at about 08:07Z. The guard then refused that key
(`update.log:40`). PR #303 reports that disposable copies of this same edited
v0.27 input completed under development qualification, so the difference between
that result and this refusal is plausibly the age of the adapter at run time,
not the input. AGENTS.md requires feedback entries to carry a reproduction;
`skgo-feedback.md:106-120` and `report.md:52-60` give the command and the
consumer but omit the condition that makes it reproduce (a fresh adapter
release inside pnpm's release-age window) and the observation that a rerun
after the window would likely take a different path. A framework owner
reproducing this tomorrow may not see it and may close it as unreproducible,
and the repository would then record a "normal input" refusal that was really a
publish-timing interaction. Impact: the recorded observation is correct but its
reproduction is incomplete and its characterization as a refusal of "untouched
recorder input" overstates what was shown.

### 2. Issue — `host-install.log` is empty but is cited as host-installation evidence

`ephemeral/validation/released-update-v0.28.2/host-install.log` is 0 bytes.
`report.md:14` links it as "[Host installation](host-install.log)" and states the
updater was "independently installed". The buildinfo (`Main.Sum
h1:4Vro+…`, `Version v0.28.2`) does establish the executable's module identity,
so provenance is not lost, but the report presents an empty receipt as if it
recorded the install. Either the capture failed or the install was done
elsewhere; the report should say which, and the link should not imply content
that is not there. README.md's "Distinguish compile/tests from running
consumer evidence" discipline (AGENTS.md:27-28) applies equally to an evidence
link that points at nothing.

### 3. Nitpick — report describes inventory coverage the inventories do not record

`report.md:21-22` says all entries were copied "with modes and symlink targets
retained" and `report.md:63-64` says the comparison includes "regular-file
SHA256 and symlink targets". The inventory files record only `sha256` per path;
there is no mode or symlink field, and the trees contain no symlinks. The claims
happen to be true (verified live above), but they are not verifiable from the
committed receipts alone, and the live `/private/tmp` trees will not survive.

### 4. Nitpick — the feedback entry does not name the upstream repair or the shift in the refusal point

`skgo-feedback.md:110-111` says "Expected the repaired ordinary released
updater to preserve authored inputs and complete". The repair is PR #303 (in
v0.28.1), and the observable change versus the v0.28.0 run is that the
`package.json` refusal no longer occurs and the guard now names the changed key,
which is exactly what #297 asked for. Stating this would tell the framework
owner that #303 worked as far as it went and that the remaining gap is the
workspace policy file, not a regression.

### 5. Nitpick — the "original refusal report" link depends on an unmerged branch commit

`report.md:17-19` and the feedback entry cite the v0.28.0 report at GitHub
commit `c1cc76d`. That commit exists on the remote but is not on `main`;
`ephemeral/validation/update-v0.28.0/` is absent from this checkout. The branch
is self-contained for hashes via `prior-authored-receipts.json`, but the
narrative context for the "original edited v0.27 baseline" lives only on
`codex/released-update-qualification`. If that branch is deleted without a
merge, the link becomes a dangling commit reference.

## Outcome

material findings remain
