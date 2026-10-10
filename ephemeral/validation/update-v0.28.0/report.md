# Released updater qualification: safe refusal

Outcome: the installed released SKGo v0.28.0 CLI safely refused to migrate an edited v0.27.0 recorder consumer because staged VitePlus migration would change authored package.json configuration. This does **not** demonstrate successful stack migration or a rebuilt updated native app. Original released v0.27.0 evidence remains intact.

Source consumer: the qualified released-manual FieldNotes app, module example.com/field-notes, bundle com.example.fieldnotes. Disposable copy: /private/tmp/skgo-recorder-consumer.GthvK4/edited-update-consumer. README and ApplicationView carry a user-authored receipt, recorder-update-20261009. The application template was not reapplied.

Updater installed independently with go install github.com/tylergannon/skgo/cmd/skgo@v0.28.0. Released source:3f0a9c8018d7b375cf84ccc2d0efa0340d69d8ca. Release gate: https://github.com/tylergannon/skgo/actions/runs/38019124856. Starting v0.27.0 has no update command; using a later CLI on an older app is not proof of upgrading an older update-capable CLI.

## Exact invocation

```sh
env VP_HOME=/private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp \
  GOWORK=off CI=1 GOBIN=/private/tmp/skgo-recorder-consumer.GthvK4/update-host \
  /private/tmp/skgo-recorder-consumer.GthvK4/update-host/skgo update \
  --root /private/tmp/skgo-recorder-consumer.GthvK4/edited-update-consumer \
  --vp /private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp/bin/vp \
  --native-platform macos --native-preset synthetic
```

Exit1. update.log records staged VitePlus1.0 migration, dependency install with managed pnpm12.9.1, and Kit3 queue patch verification. Then: `VitePlus migration would change authored configuration in package.json; application configuration is untouched`. Native build is not reached, so the forwarded macOS/synthetic choices are invocation evidence only.

The native global vp executable was copied into a disposable VP_HOME, then downgraded1.1→1.0 with both HOME and VP_HOME isolated for that subprocess, from a directory outside the consumer. The user's ordinary global vp was not replaced. Startup files contain no isolated-home references and their captured hashes remain unchanged across update. Keep the isolated home alive for diagnosis.

The initial attempt supplied npm-packaged vp instead of the required global executable; first-npm-vp-attempt.log is a separate caller setup error, not a framework defect. It failed before app mutation. One corrected attempt reached the authored-config guard; no repeated retry or manual weakening followed.

## Verified preservation

See before.json and after.json for authored SHA256 receipts: README.md, native/Sources/ApplicationView.swift, native/skgo-native.json and web/src/lib/TranscriptScreen.svelte. All match. go.mod/go.sum and web/package.json/pnpm-lock.yaml are byte-identical to the released-manual source. The app remains on SKGo v0.27.0. Real shell startup file hashes are unchanged. The preservation checks concern these named files, not an assertion that every build/cache path is unchanged.

No phone replacement, user speech test, or new runtime claim was made. Framework peer and coordinator received the exact safe-refusal result. Successful update, native build and runtime preservation remain unproved until the configuration migration issue is resolved.
