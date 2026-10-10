# Released v0.28.2 updater: safe refusal on untouched recorder input

Outcome: **update incomplete**. One ordinary released updater attempt refused a
staged change to authored `pnpm-workspace.yaml` at `minimumReleaseAgeExclude`.
The copied application and held original baseline remained byte-identical across
all 120 input entries, including the four previously recorded authored receipts.
The application still requires SKGo v0.27.0. No application/web/native completion
or browser runtime was reached; no consumer patch or retry followed.

## Released command and exact input

The independently installed updater is SKGo v0.28.2, released source
`458aa5d204af8217afbed126981bf4154d52b238`. The annotated tag resolved to that
commit, and [release gate 38035420960](https://github.com/tylergannon/skgo/actions/runs/38035420960)
completed successfully. [Host installation](host-install.log),
[initial host buildinfo](host-buildinfo.json), [final host buildinfo](host-buildinfo-after.json).
Go is 1.27.2/darwin-arm64; the executable SHA256 before and after is
`300a6ee364df20734104424727dfbae8494926e79dd28b19db0e041e6148a0c3`.

The original normal edited v0.27.0 input is
`/private/tmp/skgo-recorder-consumer.GthvK4/edited-update-consumer`, documented in
the [original refusal report](https://github.com/tylergannon/skgo-template-voice-recorder/blob/c1cc76d/ephemeral/validation/update-v0.28.0/report.md).
The fresh copy is `/private/tmp/recorder-update-v0282-iup2x7f8/UpdatedFieldNotes`.
All 120 files/symlinks were copied with modes and symlink targets retained; no
source, manifest or prerequisite patch was applied. Its product identity remains
FieldNotes, `example.com/field-notes`, `com.example.fieldnotes`; the directory
name does not rename the application. [Paths](paths.json),
[held baseline before](baseline-before.json), [copy before](consumer-before.json).

```sh
env VP_HOME=/private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp \
  GOWORK=off CI=1 GOBIN=/private/tmp/recorder-update-v0282-iup2x7f8/host \
  /private/tmp/recorder-update-v0282-iup2x7f8/host/skgo update \
  --root /private/tmp/recorder-update-v0282-iup2x7f8/UpdatedFieldNotes \
  --vp /private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp/bin/vp \
  --native-platform macos --native-preset synthetic
```

[Exact argv/environment](invocation.json), [public help](update-help.log),
[complete output](update.log), [exit status 1](update-result.json).
The caller used the normal top-level command. The internal `--complete-version`
appears only in SKGo's own subprocess diagnostic; the caller did not invoke that
phase or supply a buildinfo/source replacement override. The command installed
the same released v0.28.2 host, then reported project completion as incomplete.

The explicit `--vp` selects the retained native vendor VitePlus 1.0.0 executable,
whose SHA256 matches the [previous independently verified official binary](https://github.com/tylergannon/skgo-template-voice-recorder/blob/73010529030d23ae1cc3347d0f04b05fd4a55c3a/ephemeral/validation/released-v0.28.1/vp-provenance.json).
No private managed-pnpm wrapper was used. Direct shell pnpm is 12.10.1 and Node
is 26.11.0; VitePlus reports managed Node 24.21.0. [Tool observations](tool-versions-before.log).

## What the command reached

`update.log` identifies adapter 0.28.2 and sv add-on 0.17.1. It reports an existing
VitePlus project and aligns dependencies without source migration. Its disposable
stage installed VitePlus 1.0.0 and dependencies using pnpm 12.9.1, and verified
Kit 3.0.0 with the installed prerender queue patch. These are staged dependency
operations, not updated-application completion.

The same log reports pnpm adding one `minimumReleaseAgeExclude` entry:
`@skgo/sveltekit-adapter@0.28.2`. The authored-configuration guard then stopped:

```text
VitePlus update would change authored configuration in pnpm-workspace.yaml at ["minimumReleaseAgeExclude"]; application configuration is untouched
```

This identifies the observed staged edit and refusal; it does not establish a
reviewed core fix or authorize changing the consumer's release-age policy.
The original [workspace policy](input-web-pnpm-workspace.yaml),
[package manifest](input-web-package.json), [lockfile](input-web-pnpm-lock.yaml),
[Go module](input-go.mod), [Go sums](input-go.sum) and
[native configuration](input-native-skgo-native.json) are retained verbatim.

## Preservation and limits

[Preservation result](preservation.json) compares the entire 120-entry input
inventories, including regular-file SHA256 and symlink targets. Both
[copy after](consumer-after.json) and [baseline after](baseline-after.json) match
their respective before inventories exactly. README.md and
native/Sources/ApplicationView.swift still contain `recorder-update-20261009`;
their hashes, native/skgo-native.json and web/src/lib/TranscriptScreen.svelte
also match the immutable [original authored receipts](prior-authored-receipts.json).
No template was reapplied and the held source was not mutated.

The forwarded macOS/synthetic options are invocation evidence only. Native build,
native UI, browser behavior, frontend tests and real microphone capture were not
reached or attempted. The desktop's pending unlock did not cause this failure;
the authored-policy guard stopped earlier. No locked UI retry was attempted.

Issue #297 remains without successful released updater qualification for this
normal input. Historical v0.27.0 native/physical evidence and released v0.28.1
fresh-creation/web evidence remain attached to their identified products. Managed
installer #299 is still a separate unshipped feature. This failure changes no
Voice Notes promise verdict and does not justify an updated native-app claim.
