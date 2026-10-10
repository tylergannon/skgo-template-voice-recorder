# Source packaging

`allowlist.json` names each canonical source and its generated consumer path.
Prepare the embedded payload before building or checking the plugin:

```sh
go generate ./template
go test ./...
```

The stage copies only named sources, keeps a `.txt` suffix so Go does not compile
application payload files as plugin packages, and is ignored by Git. The renderer
removes that suffix and derives imports, display name and probe identity from the
destination project. No application source is maintained under the payload tree.

The module-root `skgo-plugin.json` contains only `{"package":"./template"}`.
The planned SKGo source installer uses this package path and ordinary Go
generation before invoking its existing host-matched builder. No custom
preparation command is needed. Its repository fetch must preserve the nested
`recorder/` module: a root Go module zip omits those canonical payload inputs.
The installer command remains pending release; the following explicit helper
procedure is the qualified v0.27.0 path.

The source-build wrapper stages those assets and invokes SKGo's shared
`skgo-plugin-build`. Build separately for the actual standalone or project-tool
host; matching version text alone does not establish Go plugin compatibility:

```sh
scripts/build-plugin.sh --host /absolute/path/to/skgo --out /absolute/plugins/voice-notes.so
# Or, for a consumer's Go tool:
scripts/build-plugin.sh --project /absolute/consumer --out /absolute/tool-plugins/voice-notes.so
```

Set `SKGO_PLUGIN_DIRS` to the selected output directory when invoking that host.
The helper reproduces host toolchain/dependencies/settings, verifies actual plugin
loading and installs atomically. For an unreleased checkout host, its explicit
`--skgo-source` option identifies the matching framework sources; this is a
development override and does not become a consumer dependency.

The plugin exports the real `SKGoPluginV1` factory. Apply requires an existing
native-app foundation, replaces only its exact default ApplicationView and the
ordinary minimal SvelteKit page, and merges recorder contributions into
`native/skgo-native.json`. It rejects edited recorder destinations or conflicting
configuration before writing any files; identical reapplication changes nothing.
Unknown configuration keys are retained, though the current core native builder
rejects unsupported keys. No commands run during Describe or Apply.
