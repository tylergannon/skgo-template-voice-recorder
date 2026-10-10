# Preparation review

Model: Claude Opus 5.5. Session: ce1aead3-357d-4d82-8cd9-481f6bb8091c

No blocking defects in the extraction itself, but I found two problems that need fixing before this can count as a reproducible package, plus one native requirement that isn't tracked anywhere. Nothing here shows that a template was delivered or that a fresh consumer works. I didn't mark the missing SKGoPluginV1 plugin as a defect, since you've stated it as pending.

## What I ran (read-only, except the gitignored `template/payload/`)
- `go run ./template/stage`, then at the root: `go vet ./...` and `go test ./...` pass.
- `go build ./...` **fails**: `template: function main is undeclared in the main package`.
- In `recorder/`: `go test ./internal/transcript`, with and without `-tags synthetic`, both pass.
- I compared each allowlisted file against `git show 4bb9373:` in the source app, using its `app/` prefix for the Go and web files.

## Material defects
1. **The package can't be built from a clean clone.** `.gitignore:1` ignores `template/payload/`, but `template/recorder.go:17-18` embeds it with `//go:embed all:payload`. A fresh clone or module-proxy fetch fails to compile until someone runs `go generate`, so a source distribution isn't reproducible from a commit. The plugin build recipe would need a generate step, or the staged payload must be committed with the existing staleness test as a guard.
2. **Required iOS/macOS settings were dropped and nothing tracks them.** The original sets these in `native/Info.plist`, plus `UIBackgroundModes=audio` in `scripts/build-native.sh:55-56`:
   - `NSMicrophoneUsageDescription`
   - `NSSpeechRecognitionUsageDescription`
   - `UIBackgroundModes=audio`

   None are in the allowlist, the worklog or `skgo-feedback.md`. Yet README says it preserves "active iPhone background recording". Without them, iOS refuses microphone access and background recording can't work. They need to be an explicit item in the native config handoff.
3. **The root `go build ./...` fails.** `template/` is `package main` with no `main`. The tests pass anyway, so nothing in the checks catches it. Please confirm this is part of the pending plugin entry point and record it.

## Smaller defects
- **The allowlist's `executable` flag is ignored.** File modes are actually chosen by the `.sh` suffix (`template/recorder.go:63`), and embedding discards modes, so the 0755 written by `template/stage/main.go` has no effect. The two sources of truth can drift apart.
- **Reapplying doesn't fix a wrong file mode.** `internal/packagefiles/files.go:59` skips a file whose content matches, even if its mode differs, so a script installed as non-executable stays that way.
- **`Replace` is never filled in**, so the "replace one known scaffold file" path has no test with real data.
- **A project name containing `{` or `}` would break the page.** It's inserted into Svelte markup (`template/recorder.go:51`) with only HTML escaping. The bundle ID isn't validated either.
- **`nameLiteral`** (`template/recorder.go:76`) is unused and discards its error.
- **Go and SKGo versions differ.** `recorder/go.mod` uses Go 1.27.1 and skgo v0.26.1; the root uses Go 1.27.2 and the 4badfe9 pseudo-version. Canonical tests run against a different SDK than packaging.
- **Leftover temp files are possible.** If the process dies between writing and renaming, a `.skgo-recorder-*` file is left behind. Cancellation is only checked between files.

## Things that check out
- **Conflict handling:** every file is checked before the first write and rechecked just before its own write. Symlinks are refused and duplicate destinations rejected. Writes go through a temp file plus rename. Identical files are left alone, and the list of changed files is returned on cancellation or I/O error.
- **Staging:** it reads only the explicit allowlist, rejects bad or duplicate paths and non-regular files, and replaces the payload in one rename. It packages no `.build`, `go.mod` or generic host files, and the test guards this.
- **Preserved behavior:** apart from the changes below, the controller, Swift package, transcript store and stream code, and `native-status` match 4bb9373 exactly. The differences are:
  - log prefixes and the package name changed
  - "Relaunch Voice Notes" became "Relaunch the application"
  - the probe identity is parameterized
  - script paths changed
  - the fixture now looks up remote IDs and codecs by name instead of numbered roots (already noted in feedback)
- **One real UI change:** `TranscriptScreen.svelte` adds a new browser-only notice, "Microphone capture uses the native application’s recording controls…". It's tested, but it's new behavior, not preserved behavior, and should be labelled that way.
- **No duplicated generic SKGo host or bootstrap.**

## Not yet proved
- That the generated names the fixture now calls (`EncodeOpenRequest`, `DecodeAck`, `generated.Remotes()` with `.Module()`/`.ID()`) exist in the SKGo version a consumer will get. Canonical remotes were deliberately not compiled.
- That the hard-coded remote ID `13p9jtx` in the test stays the same for a consumer's Vite root.
- `RecorderScreen.swift` has no direct original. It was carved out of the excluded `VoiceNotesApp.swift`, and I didn't fully reconcile it, so the bridge handlers (`snapshotApplied`, `displayAvailability`, `bridgeReady`) and the recovery and foreground logic are untested against the original.
- The Swift and controller checks in `scripts/check-recorder.sh` (not run here), the consumer's copy of that script, generation, build, launch and any native runtime behavior.

## Adjudication

Payload generation is the selected source-distribution contract, documented before build. The final wrapper will stage before invoking the shared builder; no generated second source tree is committed. A main compile entry now permits ordinary root builds. Native permission/background contributions are explicitly recorded in template/native-requirements.md, pending exact core config. Identity validation and manifest-based executable modes now reject invalid identity and incorrect existing modes. Removed unused nameLiteral. Canonical module captures the original SDK requirement; consumers use their generated framework dependency and the two source-only Go checks have no SKGo runtime dependency. Native view/bridge and freshconsumer evidence remain pending.
