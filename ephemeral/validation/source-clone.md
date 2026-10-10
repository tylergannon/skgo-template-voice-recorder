# Clean source-clone packaging check

Cloned the private remote branch `codex/recorder-template` into ignored
`build/source-clone`, with no original Voice Notes checkout access or import
overlays. Actual clone revision: `8a4e398113dd78a0bf5ca79c3ef77820fdc07e63`.

- `git clone --branch codex/recorder-template --single-branch https://github.com/tylergannon/skgo-template-voice-recorder.git /Users/tyler/.codex/worktrees/recorder-template/skgo-template-voice-recorder/build/source-clone`: exit `0`.
- `git rev-parse HEAD`: exit `0`.
- `go run ./template/stage`: exit `0`.
- `go test ./...`: exit `0`.
- `go build ./...`: exit `0`.
- `go test -C recorder ./internal/transcript`: exit `0`.

This demonstrates source-distribution packaging from a fresh repository clone,
canonical Go transcript checks, and root package tests/build. It does not establish
a host-compatible usable plugin, generated application, native runtime, real
microphone capture, or physical-device behavior. No microphone was launched.

Full command output and exit statuses: [source-clone.log](source-clone.log).
