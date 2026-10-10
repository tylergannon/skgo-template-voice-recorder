# SKGo extraction feedback

## Open: native application extension seam

Version: templateapi SDK commit `4badfe9` (prerelease).
Reference: merged application-template contract and framework peer handoff.
Expected: recorder replaces the application view, contributes its Settings and
lifecycle callbacks, and adds package/resource/plist inputs without copying a
generic host or build pipeline.
Resolved: nativeapp config/default slot arrived in 322fca3; actual builtin native
shell/build arrived in 2bde811. Recorder Apply uses those actual interfaces.

## Confirmed: fixture codec tied to generated root ordering

Version: Voice Notes source `4bb9373`, SKGo `v0.26.1`.
Reference: original `app/cmd/transcript-fixture/main.go`.
Observed: feeder calls generated `EncodeRoot0`, `EncodeRoot2` and similar numbered
functions. A fresh consumer's different type graph may reorder these identifiers.
Expected: the reusable feeder encodes named transcript contracts without assuming
the original graph order. Resolving this coupling belongs to extraction; no
fresh-consumer failure has been observed yet. Resolved in extraction: named
generated encoders/DecodeAck and Remotes registration lookup replace numbered
functions and the fixed ID. Fresh-consumer feeder Go race tests and browser
delivery now pass.

## Confirmed setup friction: VitePlus qualification and ignored parent

Version: actual installed SKGo 2bde811; global vp1.1.0/pnpm12.10.1.
Expected: ordinary skgo new prerequisites are explicit and creation does not
inherit irrelevant parent checkout ignore rules.
Observed: vp1.1.0 rejected before scaffold because core qualifies1.0.0. After
installing separate vp1.0.0/pnpm12.9.1, creating under the plugin checkout's ignored
build/ failed at vp check --fix: every matched target was excluded. This prevented
the recipe from running. Moved fresh consumers outside any checkout to
/private/tmp/skgo-recorder-consumer.GthvK4; normal creation/recipe then succeeded.
The generated devEngines still selected pnpm12.10.1 for VitePlus-managed installs
despite the explicit shell pnpm12.9.1. Shared peer confirmed VitePlus does not pin
pnpm and is correcting managed invocation/metadata. These are setup findings,
not app behavior failures.

## Confirmed: disposable native SDK enters Go package discovery

Version: nativebuild at2bde811. Repro: build macos/simulator, then go test -race ./...
in the fresh consumer. Observed tests for native/build/macos-synthetic/sdk and
simulator-synthetic/sdk, ~15 seconds each. Expected normal consumer regression
to exclude copied disposable SDK code. Shared peer supplied ebe668a, staging
under .sdk instead; not yet adopted/retested in this evidence snapshot.

## Physical install prerequisite: account/provisioning discrepancy

Version: Xcode27, nativebuild2bde811, fresh com.example.fieldnotes.real bundle.
Signed build failed once with No Accounts and no provisioning profile. Computer
use showed the current Apple account and Personal Team in Xcode Settings, so
account presence alone does not demonstrate usable command-line provisioning.
Unsigned iphone real build passed. The account-refresh request was premature. Selecting the existing Personal Team
in the generated project's Signing & Capabilities created the required profile
without a password, and signed retry/installation passed. The user subsequently
confirmed the fresh FieldNotes speech/background session. Diagnose team/profile
selection before asking for sign-in again; the CLI error was not a visible popup.


## Released adoption and subsequent framework fixes

Version: released v0.27.0/e64e350. Standalone and project-tool helpers, fresh new,
manual add-on composition, native Mac/simulator/phone builds and detached
source/plugin-unavailable builds/runs passed. Native discovery still scans sdk
staging as noted above; PR296's .sdk/update changes are separately merged and
await next-release qualification. No branch replacement entered these consumers.
The framework peer also removed VitePlus installer shell-init entries pointing
into its archived qualification worktree; those entries emitted login-shell
missing-file errors after cleanup. Preserve normal VitePlus configuration and
avoid installer side effects referencing disposable worktrees.

## v0.28.0 full explicit vp baseline is still unstable under updater

Confirmed result: NEW disposable27 consumer adopts all four files of exact reviewed full vp1.0 migration, including widened Vitest peer rules; every supplied before/post SHA matches. One released28 update then refuses another Welcome example source change. Stage/dependency/Kit patch pass, app remains27 and native build not reached. Full source inventory/authored/native config preserved; no guard relaxed or retry. Precise next rewrite unknown because stage removed; do not guess codemod/formatter cause. Reproduction/diff/log/hash evidence: validation/full-migration-update-v0.28.0/report.md. Both earlier guarded refusals retained; normal-original-input and conditional successful update remain gaps.
