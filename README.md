# SKGo voice recorder template

Reusable recorder sources and a development-time SKGo plugin, extracted from the
working Voice Notes application at `4bb9373840ca36a88c39b318da7dfb432c8196c2`.
The original Voice Notes checkout remains the working application baseline.

`recorder/` is the canonical recorder source. `template/` packages a reviewed
allowlist from that tree; generated payloads are build products, not a second
application maintained by hand. The plugin exports the `voice-notes` add-on and
the `voice-notes` recipe: SKGo's built-in `native-app`, then `voice-notes`.

Implementation uses released SKGo v0.27.0's actual `templateapi`, `nativeapp`
application slot and shared source-build helper. The recorder
exports its plugin factory and Apply implementation; see [source packaging](template/README.md)
for host-specific builds. A fresh FieldNotes consumer now builds and runs in
browser mode, on macOS and in the iPhone simulator with synthetic transcription.
Manual composition matches the template's sources and configuration. The released
foundation's signed iPhone app installs and opens. The user separately confirmed
real speech before backgrounding, on Home and while locked in the earlier
same-recorder-source development build. Generated web and macOS products rebuild
and run with this checkout and its plugins unavailable.
See [retained validation](ephemeral/validation/fresh-consumer/release-report.md)
for evidence and limits.

Released SKGo v0.28.1 also created the fresh **ReleaseFieldNotes** consumer through
ordinary `skgo new` with the released manual plugin helper and qualified VitePlus
1.0.0, without a private wrapper. It selected pnpm 12.9.1 despite ambient 12.10.1;
the production web/Go build, Go checks, Svelte diagnostics and all 11 frontend
tests passed, and regular browser computer use observed synthetic text.
This is separate from the v0.27.0 native/physical evidence;
no new native result is claimed. See the [v0.28.1 report](ephemeral/validation/released-v0.28.1/report.md).
To supply 1.0.0 alongside a newer global installation, follow the
[vendor's installation guide](https://viteplus.dev/guide/global-cli) with
`VP_VERSION=1.0.0` and an absolute `VP_HOME`, then select that home's `bin` on PATH
and keep the same `VP_HOME` for `skgo new`. This vendor-documented fresh-install
procedure was not exercised here; qualification reused an existing installation.

The recorder preserves microphone selection, native recording controls, timer,
five horizontal green bars stacked upward, live transcript delivery, synthetic
playback and active iPhone background recording. Its transcription package remains
independently testable. No transcript persistence or audio archive is introduced.

Read [VISION.md](VISION.md) and [AGENTS.md](AGENTS.md). Implementation decisions
and prerequisite gaps live under `ephemeral/`. Voice Notes Mission Control records
fresh-consumer evidence separately from evidence for the original application.

## Checks

`scripts/check-recorder.sh` exercises the existing Go transcript and synthetic
Swift package/controller checks. It does not record microphone audio or establish
physical-device behavior. Template packaging checks use `go test ./...` at this
repository's root.

## Create a consumer

### Source installation — pending SKGo release

SKGo's source installer is planned; released v0.27.0 and v0.28.0 do not provide
this command. After the installer ships, the intended workflow is:

```sh
skgo plugin install github.com/tylergannon/skgo-template-voice-recorder
skgo new --help
CI=1 skgo new --name FieldNotes --module example.com/field-notes \
  --template voice-notes --set bundle-id=com.example.fieldnotes field-notes
```

An optional `@version` selects a Go module query; omitting it means `latest`.
The root `skgo-plugin.json` names `./template`. SKGo will fetch the pinned full
repository and use its shared builder to run `go generate ./template` in a
writable stage with the actual host's shared dependency pins, compile the plugin,
and install it for that host with automatic discovery.
This runs trusted repository code. Git, the host's exact Go toolchain (at least
Go 1.27.1), and a working C toolchain are required.
Installing for a consumer's `go tool skgo` is a separate invocation against that
actual host; a changed host requires installation again.

Keep one discovered copy per host. Before using managed installation, remove the
manual recorder `.so` from the directories that host discovers. To keep a manual
copy instead, the planned `skgo plugin remove
github.com/tylergannon/skgo-template-voice-recorder` will remove only that host's
managed selection. Changing `SKGO_PLUGIN_DIRS` will not hide managed copies.

This packaging convention is prepared, but released source installation and
automatic discovery remain unqualified. The manual path below retains the historical
v0.27.0 qualification; v0.28.0 also predates the managed installer.

### Historical manual qualification — v0.27.0

Use Go 1.27.2, Node 24+, VitePlus 1.0.0 and pnpm 12.9.1 for the qualified
SKGo v0.27.0 creation path. Native builds need Xcode 27; macOS/simulator target
27.0, and the real iPhone preset supports iOS 26.0+. First install the standalone
host and build the development plugin for that exact host:

```sh
go install github.com/tylergannon/skgo/cmd/skgo@v0.27.0
scripts/build-plugin.sh --host /absolute/path/to/skgo --out /absolute/plugins/voice-notes.so
```

Create outside an ignored parent checkout. Point that same host at its plugin:

```sh
SKGO_PLUGIN_DIRS=/absolute/plugins CI=1 /absolute/path/to/skgo new \
  --name FieldNotes --module example.com/field-notes \
  --template voice-notes --set bundle-id=com.example.fieldnotes field-notes
```

The equivalent manual composition is ordinary `skgo new`, followed by
`skgo add --set bundle-id=com.example.fieldnotes native-app` and
`skgo add voice-notes`. A generated application owns its copied sources and does
not load the plugin at runtime. Signing inputs belong to the consumer developer.
This manual helper procedure is retained for reproducing that release's evidence;
it should not be combined with managed installation for the same host.
Rebuild the development plugin for each changed SKGo host or Go toolchain;
the shared helper verifies compatibility against the host that will load it.
See the [qualification limits](ephemeral/validation/fresh-consumer/release-report.md#limits)
for the v0.27.0 disposable SDK directory issue and separate update-release work.

From the created consumer, use `go tool skgo native build --platform macos
--preset synthetic` for microphone-free desktop playback, or change the platform
to `simulator`. Use `--platform iphone --preset real` for the real phone build.
Supply `SKGO_DEVELOPMENT_TEAM` and `SKGO_DEVICE_ID` from your own Xcode setup when
signing for a connected phone. Build products are under `native/build/`.
For browser playback, run the built Go server and submit synthetic text with
`scripts/feed-transcript.sh --origin http://127.0.0.1:8080`.
