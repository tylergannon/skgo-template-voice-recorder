# SKGo voice recorder template

Reusable recorder sources and a development-time SKGo plugin, extracted from the
working Voice Notes application at `4bb9373840ca36a88c39b318da7dfb432c8196c2`.
The original Voice Notes checkout remains the working application baseline.

`recorder/` is the canonical recorder source. `template/` packages a reviewed
allowlist from that tree; generated payloads are build products, not a second
application maintained by hand. The plugin exports the `voice-notes` add-on and
the `voice-notes` recipe: SKGo's built-in `native-app`, then `voice-notes`.

Implementation uses SKGo's actual `templateapi`, `nativeapp` application slot
and shared source-build helper at development commit `2bde811`. The recorder
exports its plugin factory and Apply implementation; see [source packaging](template/README.md)
for host-specific builds. A fresh FieldNotes consumer now builds and runs in
browser mode, on macOS and in the iPhone simulator with synthetic transcription.
Manual composition matches the template's sources and configuration. The real
iPhone code compiles; signing/install, fresh real speech and the shared framework
release remain pending. See [retained validation](ephemeral/validation/fresh-consumer/report.md)
for evidence and limits.

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

## Intended consumer commands

After the shared SDK, native foundation and plugin build recipe are released:

```sh
CI=1 skgo new --name FieldNotes --module example.com/field-notes \
  --template voice-notes --set bundle-id=com.example.fieldnotes field-notes
```

The equivalent manual composition is ordinary `skgo new`, followed by
`skgo add --set bundle-id=com.example.fieldnotes native-app` and
`skgo add voice-notes`. A generated application owns its copied sources and does
not load the plugin at runtime. Signing inputs belong to the consumer developer.
