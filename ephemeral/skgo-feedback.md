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

## Released v0.28.1 bootstrap selection and scaffold warnings

Version: v 0.28.1/bc 267b 284c 19bab 3ae 47aaf 78bc 550c 1cbbfcc 50;
reference: released `skgo new --help` and internal/toolchain's VitePlus 1.0.0 pin.
The ambient global vp 1.1.0 was correctly rejected before scaffolding with an
explicit install 1.0.0 diagnostic. Selecting the existing unmodified vendor 1.0.0
binary through PATH/VP_HOME, with ambient pnpm 12.10.1 unchanged, let ordinary new
choose pnpm 12.9.1 and create the recorder without the development wrapper.
Generated devEngines and installed .modules.yaml both identify 12.9.1. This closes
the previously observed managed-companion selection issue for this released
noninteractive path; it does not qualify fresh VitePlus installation or interactive
creation. Exact commands and output: ephemeral/validation/released-v 0.28.1/.

Confirmed nonblocking scaffold issue: `skgo new` emitted `vp check --fix` warnings
for an empty src/lib/index.ts and an unawaited render promise in generated
src/lib/vitest-examples/Welcome.svelte.spec.ts (lines 7-8). Expected clean generated
scaffold checks; observed 0 errors / 2 warnings with creation exit 0, followed by passing
production build and Go tests. No patch was applied. The upstream scaffold/test
warnings remain separate from recorder behavior and updater 297. Both direct
shell Node 26.11.0 and VitePlus-reported managed Node 24.21.0 are retained rather
than presenting one inferred toolchain identity.


Confirmed nonfatal adapter setup ordering issue: during ordinary v0.28.1 new,
managed install ran package.json's `svelte-kit sync || echo ''` prepare script
before SKGo generated web/skgo.remotes.json. Adapter 0.22.1 createPrerenderOwner
(prerender.js:330) threw ENOENT for that path; the existing shell guard tolerated
it, and generation subsequently wrote the manifest. Expected a fresh scaffold to
complete setup without this stack trace. Reproduce through the ordinary command
in released-v0.28.1/creation-command.json; new.log lines 144-152 retain the trace.
The same diagnostic appears in prior released-new.log. Creation exits 0 and the
later build, Svelte diagnostics and all 11 frontend tests pass; this is an ordering
gap, not an observed recorder failure. new.log also reports 24 deprecated
subdependencies; dependency hygiene remains a separate upstream concern.
