# Released v0.28.1 ordinary recorder creation and web runtime

Outcome: released SKGo v0.28.1 created a fresh independently named recorder
consumer through ordinary `skgo new`, using the released manual plugin helper
and an unmodified qualified VitePlus executable. Its production web/Go build
and Go tests passed. Regular browser computer use observed the synthetic feeder's
final text: **“This is synthetic text.”** No private VitePlus wrapper was used.

Recorder source is canonical `118f325610f924881abcaf15c3ca2ae6597da040`.
SKGo v0.28.1 resolves to `bc267b284c19bab3ae47aaf78bc550c1cbbfcc50`;
the [release gate](https://github.com/tylergannon/skgo/actions/runs/38029289700)
succeeded. The independently installed standalone host and generated application's
Go-tool host both identify v0.28.1 without a local replacement.
[Host buildinfo](host-buildinfo.json), [project host](project-tool-buildinfo.json),
[identities and hashes](resolved-identities.json).

## Prerequisites and normal creation

The ambient global VitePlus was 1.1.0 and pnpm was 12.10.1. The first ordinary
invocation stopped before scaffolding with an explicit requirement for qualified
VitePlus 1.0.0. [Prerequisite refusal](ambient-vp-prerequisite.log).
One corrected attempt selected the already retained native VitePlus 1.0.0
installation by its ordinary PATH and VP_HOME. This is the actual vendor binary,
not the development managed-pnpm wrapper. Its resolved executable and hash,
complete invocation and environment overrides are in [creation-command.json](creation-command.json).
The user's ordinary global installation was not changed.

Ambient `pnpm --version` still reported 12.10.1. Released creation selected its
qualified 12.9.1 companion; `new.log` records installation of pnpm@12.9.1, and
both generated `devEngines.packageManager` and installed `.modules.yaml` identify
12.9.1. [Tool versions](tool-versions.log), [creation](new.log),
[exit status](new-result.json), [resolved identities](resolved-identities.json).
The direct shell Node reported v26.11.0, while VitePlus reported its managed
Node v24.21.0; those observations are preserved separately.

The existing recorder `scripts/build-plugin.sh` built and loaded the plugin for
the actual released standalone host (Go1.27.2); this was manual discovery through
`SKGO_PLUGIN_DIRS`, not the unshipped managed installer.
[Helper output](plugin-build.log), [host installation](host-install.log).
The consumer is `/private/tmp/recorder-v0281-7iiyiinf/ReleaseFieldNotes`, module
`example.com/release-field-notes`, bundle `com.example.releasefieldnotes`.
Creation applied builtin native-app and then voice-notes without editing or
patching the generated application. The resolved identities retain hashes for
all 47 recorder destinations and the module/toolchain/native configuration.

## Build and browser observation

The generated `just build` recipe completed frontend generation/build and
production Go compilation. `go test ./...` passed, including transcript and
feeder tests. [Build](app-build.log), [Go checks](consumer-go-tests.log),
[exit statuses](build-results.json), [app buildinfo](app-buildinfo.txt).

The identified production binary served `http://127.0.0.1:8197` using its normal
`--listen` and `--origin` flags. The initial browser view identified ReleaseFieldNotes
and explained native microphone capture. Its documented external feeder submitted
six provisional/final updates and closed the synthetic session. Regular browser
computer use then observed Stopped, Synthetic transcript and the final sentence.
[Feeder](feeder.log), [browser accessibility](web-synthetic-ax.txt),
[screenshot](web-synthetic.png), [HTTP headers](web-headers.txt),
[initial response](web-response.html), [server](server.log),
[owned process](server-process.txt). This establishes final synthetic browser
delivery; it does not assert that every intermediate revision was visually observed.
The owned server was stopped after the screenshot.

![Released consumer synthetic text](web-synthetic.png)

## Limits and Mission Control attribution

This qualifies released ordinary creation under its documented VitePlus 1.0.0
prerequisite, not automatic adoption of an incompatible ambient VitePlus version
or a from-scratch VitePlus installation. No private wrapper, source replacement,
managed plugin install, updater retry, native build, native runtime, physical
device or microphone check ran. The desktop remains unavailable for native UI
qualification; no locked-screen retry or screenshot workaround was attempted.

Creation emitted two nonfatal upstream scaffold lint warnings (empty lib index
and an unawaited Welcome test render), retained in `new.log`; no consumer patch
was made. Build/test success does not erase them.

Voice Notes Mission Control's existing `native-starter` assessment remains bound
to recorder a481333 and released v0.27.0 web/Mac/simulator evidence. This bounded
v0.28.1 web result does not replace that native assessment or alter original
Voice Notes promise statuses. Previous physical-launch and development real-speech
observations stay attached to their original products. No website status change
is needed to preserve that attribution.
