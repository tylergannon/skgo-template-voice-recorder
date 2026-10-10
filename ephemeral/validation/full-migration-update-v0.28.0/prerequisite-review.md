# Reviewed complete prerequisite

Every supplied changed-file hash matches the copied original/post bytes. Complete frontend comparison has no additional or missing files: only package.json, pnpm-workspace.yaml, pnpm-lock.yaml and the Welcome example change. Full unified diff, exact original/post files and framework receipt/log are retained.

The package test script selects vp test. Its already async Welcome test awaits render before DOM assertions. Vitest/browser-playwright versions5.0.1 move into the catalog, with a Vitest override; pnpm aligns12.10.1→12.9.1 including platform lockfile entries/integrities. The override narrows a lockfile Vitest peer selection to5.0.1. Existing Storybook peer rule remains.

Explicitly broader policy: peerDependencyRules.allowAny adds vitest, and allowedVersions.vitest becomes '*'. These suppress Vitest peer-version constraints rather than just select a package version. Coordinator authorizes them solely as this disposable experiment's prerequisite; canonical recorder and installed application policy adoption is not authorized or demonstrated. No guards are weakened.

Named authored receipts and full native configuration bytes are preserved, and only four reviewed baseline files differ before the updater. This can demonstrate conditional update after explicit full vp migration; it cannot establish unattended usability on the original27 input.
