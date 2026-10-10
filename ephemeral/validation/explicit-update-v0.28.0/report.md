# Explicit prerequisite migration: second safe refusal

Outcome: one released SKGo v0.28.0 updater attempt after an explicitly reviewed generated-baseline migration still refused authored pnpm-workspace.yaml changes. App remains v0.27.0; native build was not reached. Successful update after a prerequisite migration is **not demonstrated**. No guard was bypassed and no repeated update ran.

The coordinator authorized this bounded attempt on a NEW disposable edited consumer. Original safe-refusal evidence remains at [c1cc76d](https://github.com/tylergannon/skgo-template-voice-recorder/blob/c1cc76d/ephemeral/validation/update-v0.28.0/report.md). Original edited-update-consumer, canonical recorder sources and installed apps were not changed. No phone request or new native runtime claim was made.

## Prerequisite reviewed and applied

NEW consumer: /private/tmp/skgo-recorder-consumer.GthvK4/explicit-update-consumer, copied from the qualified released-manual27 FieldNotes product. It carries a distinct explicit-recorder-update-20261009 receipt in README and native ApplicationView. before.json records its full source inventory and named authored hashes.

Framework's retained vp1.0 diagnostic and [issue297 comment](https://github.com/tylergannon/skgo/issues/297#issuecomment-6093551831) establish two generated baseline rewrites. explicit-prerequisite.patch records exactly those substitutions: scripts.test:unit vitest→vp test; the already async Welcome example test awaits render before its DOM assertions. The exact migrated example matches the framework's retained post-migration bytes. prerequisite-review.md explains their semantics. No dependency/workspace/lock changes were copied from the probe. pre-update.json confirms exactly two changed files, all authored receipt hashes unchanged and Go still27.

## Exact released attempt

Released updater source3f0a9c8018d7b375cf84ccc2d0efa0340d69d8ca, stable v0.28.0 independently verified before invocation (released-framework.json). Exact executable path and SHA are in before.json. Native isolated vp reports1.0.0; updater log reports managed pnpm12.9.1 and selected adapter0.22.1/sv0.17.1.

```sh
env VP_HOME=/private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp \
  GOWORK=off CI=1 GOBIN=/private/tmp/skgo-recorder-consumer.GthvK4/update-host \
  /private/tmp/skgo-recorder-consumer.GthvK4/update-host/skgo update \
  --root /private/tmp/skgo-recorder-consumer.GthvK4/explicit-update-consumer \
  --vp /private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp/bin/vp \
  --native-platform macos --native-preset synthetic
```

Exit1. Staged migration, install and Kit queue patch verification passed. Then: `VitePlus migration would change authored configuration in pnpm-workspace.yaml; application configuration is untouched`. The installed CLI is28, but project update is incomplete. Native macOS/synthetic selection was forwarded; generation/build did not reach it.

## Workspace diagnosis and preservation

peer-workspace-diagnostic.patch and peer-workspace-after-migration.yaml retain the **prior framework probe**, not the removed stage from this actual updater attempt. That probe adds Vitest/browser catalog pins, a Vitest override, and peerDependencyRules allowAny/allowedVersions entries. From released frontend.go authoredConfig, catalog and override additions are managed/stripped, while these peer rules survive its comparison. The peer-rule additions are therefore an inferred cause of the new guard rejection; the actual log directly establishes only the workspace rejection. Those additions were not applied to the new consumer. Further baseline changes require explicit review/authorization before another run.

After.json confirms the entire new consumer source inventory is byte-identical to pre-update, all named authored receipts and shell startup hashes remain unchanged, and the original refusal consumer's named receipts still match retained c1cc76d hashes. Old released27 and development observations remain separate. No app rebuild/runtime was rechecked because the updater stopped before app mutation/build.

Framework peer and coordinator received the new bounded result. Hold the consumers/toolchain for diagnosis; do not describe either guarded attempt as unattended or successful stack migration.
