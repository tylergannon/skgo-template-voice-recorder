# Released v0.28.3 fresh-package updater qualification

decision: Coordinator authorized one time-sensitive ordinary released update on one fresh exact copy of the untouched edited v0.27 baseline, with publication/attempt times and observed age-window behavior. No input patches or development overrides; stop on failure. Native build and native UI remain distinct; no locked UI retry.

friction: Preflight incorrectly assumed the npm adapter version equalled the Go release and stopped before invoking update when 0.28.3 was absent. Registry latest is 0.28.2; capture actual selection and timestamps instead. Same exact copy, no updater attempt had run and no input changed.

outcome: The sole ordinary update exited 0. Actual adapter 0.28.2 was 3606.747 seconds old and pnpm logged its exact release-age exception; preserved workspace diff shows the accepted addition. Held baseline 120 files and four authored receipts unchanged. Production web/Go and macOS synthetic compilation completed.

friction: Ordinary just test refused ambient pnpm 12.10.1 against pinned 12.9.1. Retained refusal; same recipe via vendor vp env exec --package-manager pnpm@12.9.1 passed all 11 tests. No manifest weakening or updater retry. Initial assumed bin/skgo path absent; resolved actual executable via go tool -n skgo and recorded v0.28.3 buildinfo.

proof: Go/Svelte checks passed. External feeder accepted six revisions; normal browser CUA observed final synthetic text and screenshot. No microphone/native UI attempted. Owned web server stopped and tab closed. Report and receipts prepared for independent review.
