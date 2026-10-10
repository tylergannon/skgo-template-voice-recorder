# Adversarial review: released SKGo v0.28.1 ordinary recorder qualification

Date: 2026-10-10. Read-only review; no project file other than this artifact was
written.

**Target:** branch `codex/released-v0281-qualification`, commit `628ee76`
(`main...HEAD`, one commit, 35 files, additions only): the README paragraph on
ReleaseFieldNotes, the new `ephemeral/skgo-feedback.md` entry, the worklog, and
`ephemeral/validation/released-v0.28.1/` (report plus 33 receipts). No source,
template, script or configuration file changed; `recorder/` is byte-identical to
`118f325` (`git diff --quiet 118f325 HEAD -- recorder`).

**Scope note:** scope was taken from AGENTS.md, VISION.md, README.md,
`template/README.md` and the artifacts themselves. The caller's summary of the
coordinator-authorized outcome (ordinary creation without the private vp wrapper,
independent name, production build and web runtime, exact identities) and its
statement that installer 299, updater 297 and native runtime are separately owned
were treated as context. They did not limit what was inspected; the review
confirms that the report claims nothing in those areas (all four
`nativeBuildRun`/`nativeRuntimeObserved`/`managedInstallRun`/`updaterRun` flags
in `resolved-identities.json` are `false`).

## Evidence inspected

- Instructions: AGENTS.md, VISION.md, README.md (before/after), `template/README.md`,
  `template/allowlist.json`, `template/recorder.go` (rewrite rules),
  `scripts/build-plugin.sh`, `scripts/check-recorder.sh`, `.gitignore`.
- Report and receipts: `report.md`, `paths.json`, `initial-identities.json`,
  `resolved-identities.json`, `creation-command.json`,
  `ambient-vp-prerequisite-{command.json,result.json,log}`, `new.log` (all 313
  lines), `new-result.json`, `new-help.log`, `tool-versions.log`,
  `plugin-build.log`, `host-install.log`, `host-buildinfo.json`,
  `project-tool-buildinfo.json`, `consumer-go.mod`, `consumer-web-package.json`,
  `consumer-web-pnpm-workspace.yaml`, `consumer-native-skgo-native.json`,
  `app-build.log`, `consumer-go-tests.log`, `build-results.json`,
  `app-buildinfo.txt`, `server.log`, `server-process.txt`, `feeder.log`,
  `feeder-result.json`, `web-headers.txt`, `web-response.html`,
  `web-synthetic-ax.txt`, `web-synthetic.png` (viewed).
- Prior context: earlier `skgo-feedback.md` entries, `released-new.log`,
  `field-notes-create.log`, the `fresh-consumer` evidence directory listing and
  `ephemeral/reviews/release-qualification-review.md`.
- Independent checks run:
  - `recorderDestinations` has 47 entries, equal to the 47 allowlist entries;
    `new.log` lists 48 `voice-notes: changed` lines (the 47 plus
    `native/skgo-native.json`, which native-app owns and voice-notes extends).
  - SHA256 of eight canonical `recorder/` files that the renderer does not
    rewrite (both licenses, `synthetic-transcript.json`,
    `burst-transcript.json`, `feed-transcript.sh`, `Package.swift`,
    `TranscriptScreen.svelte`) equals the recorded destination hash.
    `check-recorder.sh` differs, as expected from the `cd` rewrite at
    `template/recorder.go:90-91`.
  - `host-buildinfo.json` and `project-tool-buildinfo.json` differ only in the
    `executable` line; both report `github.com/tylergannon/skgo v0.28.1`,
    sum `h1:faO0BK…`, Go 1.27.2, and contain no `Replace` entries.
    `app-buildinfo.txt` pins the same skgo sum; `consumer-go.mod` requires
    v0.28.1 and declares `tool github.com/tylergannon/skgo/cmd/skgo`.
  - `devEngines.packageManager` 12.9.1 in the saved package.json and
    `installedPackageManager: pnpm@12.9.1`; ambient `pnpm --version` 12.10.1.
  - `server-process.txt` shows the identified `bin/ReleaseFieldNotes` listening
    on 127.0.0.1:8197; the feeder accepted six updates and closed; the AX dump
    and screenshot agree. The screenshot lacks the "Stopped" label only because
    `.status` is visually hidden by design (`TranscriptScreen.svelte:81`).
  - Leak scan of the non-ephemeral diff for `/Users/`, `DEVELOPMENT_TEAM`, the
    team ID and the user's email: nothing found.

## Findings

### 1. Issue: the qualified VitePlus 1.0.0 executable has no provenance beyond a local hash

Evidence: `creation-command.json` sets `VP_HOME` and prepends PATH to
`/private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp`, a
disposable temp directory left by earlier qualification work
(`skgo-feedback.md:35` records that tree as a throwaway consumer parent). The
report calls this "the actual vendor binary, not the development managed-pnpm
wrapper" and "already retained native VitePlus 1.0.0 installation", but the only
identity recorded is a self-reported `vp v1.0.0` and a SHA256
(`d787c0a8…`) with no download source, installer record or upstream checksum.
`tool-versions.log` also shows that executable reporting "Local vite-plus: Not
found" and managed Node v24.21.0, i.e. an isolated home, not the user's ordinary
installation.

Impact: the stated outcome "exact source/host/toolchain identities" is met for
recorder, SKGo host, Go and pnpm, but the toolchain piece most at issue (the
component the private wrapper used to replace) rests on an assertion. Nothing in
the repository lets a reader confirm the binary is unmodified vendor 1.0.0, and
nothing lets them reproduce the run: the README's new paragraph says creation
used "qualified VitePlus 1.0.0" without saying how to supply it when the
ambient install is 1.1.0, and the report explicitly declines to qualify a
from-scratch 1.0.0 installation. The prerequisite refusal
(`ambient-vp-prerequisite.log`) tells the user to "Install that version"; the
retained path instead re-used a leftover from `/private/tmp`, which the previous
feedback entry warned against relying on ("avoid installer side effects
referencing disposable worktrees").

Suggested resolution: record where that `vp/1.0.0` tree came from (installer
command or npm package and version, and the upstream checksum for the
darwin-arm64 binary), or re-run creation once against a 1.0.0 installed by the
documented vendor procedure into a non-disposable `VP_HOME`, and state in the
README how a reader supplies 1.0.0 alongside an ambient 1.1.0.

### 2. Issue: a recurring shared-SKGo scaffold error is not recorded as feedback

Evidence: `new.log:144-152`. The adapter's `prepare` script
(`svelte-kit sync || echo ''`) throws
`ENOENT: no such file or directory, open '…/web/skgo.remotes.json'` from
`@skgo/sveltekit-adapter@0.22.1/skgo-adapter/prerender.js:330` during the
managed install, prints a full stack trace, and is only hidden by the `|| echo ''`
guard; SKGo writes `skgo.remotes.json` several lines later. The same trace is
present in `released-new.log` and `field-notes-create.log`, so it reproduces on
every fresh creation. `new.log:139` also reports 24 deprecated subdependencies.
The report's "Limits" section says creation emitted "two nonfatal upstream
scaffold lint warnings" and nothing else; the new `skgo-feedback.md` entry
records only the `vp check --fix` warnings.

Requirement: AGENTS.md, "Record shared-SKGo gaps in ephemeral/skgo-feedback.md
with version, expectation, observation and reproduction."

Impact: the report under-describes creation output, and a shared-framework
ordering defect (adapter prepare runs before the remotes manifest exists) stays
invisible to the framework peer. Creation still exits 0, so this is not a
qualification failure.

### 3. Nitpick: `new-help.log` does not evidence template discovery

`new-help.log` ends with `Available templates:` followed by nothing, so it was
captured without `SKGO_PLUGIN_DIRS` (or before `plugin-build.log`). The report
does not link it, but it sits beside receipts that do show `applied voice-notes`
and reads as if the host found no template. Re-capture with the plugin directory
set, or remove it.

### 4. Nitpick: whitespace stripped before numbers in committed prose

`skgo-feedback.md` (new section) and `worklog/released-v0281-qualification.md`
contain "VitePlus1.0.0", "vp1.1.0", "install1.0.0", "vendor1.0.0",
"pnpm12.10.1", "identify12.9.1", "lines7-8", "observed0errors/2warnings",
"Node26.11.0", "Node24.21.0", "installer299", "updater297", "recorder118f325".
`report.md` has "Go1.27.2". Earlier feedback entries show the same pattern, so
the authoring pipeline is dropping spaces systematically. These are the
versions a framework peer will grep for.

### 5. Nitpick: the generated consumer's web checks were not exercised

The v0.27.0 qualification kept `field-notes-web-tests.log` and
`field-notes-svelte-check.log`. For v0.28.1 only `go test ./...` ran; the
delivered `TranscriptScreen.svelte.spec.ts` and `native-status.spec.ts` and
`svelte-check` were not run on the fresh consumer, although `just test` was
available. The report and README word their claims accurately ("Go checks"), so
this is a coverage gap, not a misstatement.

## What held up

- Independent name, module and bundle: ReleaseFieldNotes,
  `example.com/release-field-notes`, `com.example.releasefieldnotes`.
- Ordinary `skgo new` with `--template voice-notes`, no wrapper
  (`privateVPWrapper: false`), no `-skgo-replace`/`-adapter` overrides, exit 0.
- Recorder source, SKGo host and project tool identities are consistent across
  all receipts; destination hashes match canonical sources where applicable.
- Production build (`just build`, `-ldflags -X main.origin=…`) and Go tests
  exit 0; the production binary served the page and displayed the feeder's
  final text.
- The README and report keep v0.28.1 web evidence separate from v0.27.0 native
  and physical evidence and claim no native, installer or updater result.
- No credentials, team IDs or home paths entered non-ephemeral files.

## Outcome

material findings remain
