# After explicit full vp migration: update still refused

Outcome: **safe refusal after explicit full vp migration**, not successful update. One released SKGo v0.28.0 updater run refused another change to src/lib/vitest-examples/Welcome.svelte.spec.ts. It stopped before app Go update, generation or native build. Source files and authored receipts remain byte-identical to the reviewed migrated baseline.

This was the coordinator-authorized bounded experiment on NEW /private/tmp/skgo-recorder-consumer.GthvK4/full-migration-update-consumer. Original consumers, canonical recorder and installed apps remain separate. No guards were bypassed, no extra edit/migration/update was attempted after the refusal, and no phone request was made. Successful conditional migration and normal original-input usability remain unproved.

## Complete explicit prerequisite

Independently compared every frontend source file with the framework's retained pinned-vp probe: exactly four changed paths, no added/missing files. Complete diff including all lockfile hunks was inspected. Every provided before/post SHA matched copied bytes. Retained originals and exact reviewed post files are under original-web/ and reviewed-migrated-web/; complete-prerequisite.diff and framework-complete-migration.diff retain the full changes, with framework-changed-files.json and its prior probe log. That diagnostic is a separate prior staged probe, not this run's removed stage.

Changes: generated test:unit script vitest→vp test; the already async Welcome example awaits render; Vitest/browser-playwright5.0.1 catalog selection and Vitest override; pnpm12.10.1→12.9.1 manifest/lock pins. **Broadened peer prerequisite:** peerDependencyRules.allowAny adds vitest and allowedVersions.vitest becomes '*', relaxing peer-version constraints. These exact rules were accepted solely in this disposable experiment. No canonical or installed-app policy adoption is authorized or demonstrated. Existing Storybook compatibility policy remains.

A distinct authored receipt, full-migration-recorder-update-20261009, is in README and native ApplicationView. Before/pre-update inventories show only the four reviewed prerequisite files changed. Named README, native source, transcript screen and full native config hashes are preserved. See prerequisite-review.md. No later guessed patch was applied.

## One released updater invocation

Stable released framework v0.28.0 independently verified before invocation, source3f0a9c8018d7b375cf84ccc2d0efa0340d69d8ca (released-framework.json). before.json identifies the actual independently installed CLI and executable SHA. Native isolated vp1.0.0 and managed pnpm12.9.1 are identified by the update log.

```sh
env VP_HOME=/private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp \
  GOWORK=off CI=1 GOBIN=/private/tmp/skgo-recorder-consumer.GthvK4/update-host \
  /private/tmp/skgo-recorder-consumer.GthvK4/update-host/skgo update \
  --root /private/tmp/skgo-recorder-consumer.GthvK4/full-migration-update-consumer \
  --vp /private/tmp/skgo-recorder-consumer.GthvK4/update-isolated-home/vp/bin/vp \
  --native-platform macos --native-preset synthetic
```

Exit1. Staged migration/install and Kit queue patch verification passed, then: `VitePlus migration would change application source src/lib/vitest-examples/Welcome.svelte.spec.ts; migration is incomplete and the original file is untouched`. CLI28 is installed but app remains27. Forwarded macOS/synthetic choices are invocation evidence; no native build or updated runtime exists.

The fully migrated baseline is not stable under the updater's migration pipeline. The actual updater stage is removed on exit, so this report does not assert the specific subsequent Welcome rewrite or whether its codemod/formatter caused it. Framework peer received the failure for diagnosis; the experiment stopped at this changed-source condition as instructed.

## Preservation and earlier refusals

After.json records the complete post-run source inventory, identical to pre-update, unchanged authored hashes, full native config, existing Storybook rule and shell startup files. Original safe-refusal consumer named hashes still match retained evidence. App/native behavior was not rerun because the update stopped before app mutation/build. Prior released27 and development native observations remain attached to their identified products.

- [Original27 refusal on package configuration](https://github.com/tylergannon/skgo-template-voice-recorder/blob/c1cc76d/ephemeral/validation/update-v0.28.0/report.md).
- [Two-file explicit prerequisite refusal on workspace policy](https://github.com/tylergannon/skgo-template-voice-recorder/blob/6ce835b/ephemeral/validation/explicit-update-v0.28.0/report.md).

This third bounded result neither supersedes those preservation witnesses nor closes skgo297. Any canonical baseline correction or further experiment is separate work, not part of this authorized one-run test.
