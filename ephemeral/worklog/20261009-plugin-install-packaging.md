# Recorder source-install packaging

decision: The coordinator routed the user's source installer UX to shared SKGo. Core owns fetch, host-matched build, installation and discovery; recorder owns only its package convention and payload preparation. Preserve canonical application behavior and released27 evidence. Archive distribution is deferred.

finding: Clean178135e source has no embedded payload, but go generate ./template succeeds and runs its existing directive, staging the explicit47-entry allowlist. A conventional package path plus ordinary Go generation may avoid a custom preparation manifest or installer hook. Coordinate schema before editing rather than invent duplicate installer contracts.

finding: Canonical recorder/ is a nested Go module. x/mod v0.41 zip.CheckFiles omits nested module files; root Go-module zip retrieval cannot provide the reviewed47-file payload inputs. Sent this constraint to the installer owner before choosing a fetch convention. Preserve the full repository source tree instead of changing recorder architecture to fit an incomplete fetch.

decision: Agreed proposed core contract: module-root skgo-plugin.json has only package=./template; core fetches immutable full Git repository from Go resolver origin, runs ordinary Go generation in isolated writable source with actual-host Go/GOWORK=off, then uses its existing actual-host builder and managed discovery. No custom preparation field. Added exactly that metadata and pending-release docs; qualified27 procedure and prior runtime reports remain intact.

validation: Metadata parses as exactly one package field. Existing allowlist remains47 entries. go generate ./template and root go test ./... pass, including payload/source byte equivalence. This verifies packaging preparation only; core installer is not yet implemented/released and no automatic install/discovery or new native behavior is claimed. Independent packaging review requested before PR handoff.

review: Independent packaging review inspected the current changes against repository instructions and the proposed source-install contract; no findings. Full report retained. Core contract review and actual installer implementation remain separate prerequisites for source-install runtime proof.
