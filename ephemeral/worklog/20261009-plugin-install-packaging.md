# Recorder source-install packaging

decision: The coordinator routed the user's source installer UX to shared SKGo. Core owns fetch, host-matched build, installation and discovery; recorder owns only its package convention and payload preparation. Preserve canonical application behavior and released27 evidence. Archive distribution is deferred.

finding: Clean178135e source has no embedded payload, but go generate ./template succeeds and runs its existing directive, staging the explicit47-entry allowlist. A conventional package path plus ordinary Go generation may avoid a custom preparation manifest or installer hook. Coordinate schema before editing rather than invent duplicate installer contracts.

finding: Canonical recorder/ is a nested Go module. x/mod v0.41 zip.CheckFiles omits nested module files; root Go-module zip retrieval cannot provide the reviewed47-file payload inputs. Sent this constraint to the installer owner before choosing a fetch convention. Preserve the full repository source tree instead of changing recorder architecture to fit an incomplete fetch.

decision: Agreed proposed core contract: module-root skgo-plugin.json has only package=./template; core fetches immutable full Git repository from Go resolver origin, runs ordinary Go generation in isolated writable source with actual-host Go/GOWORK=off, then uses its existing actual-host builder and managed discovery. No custom preparation field. Added exactly that metadata and pending-release docs; qualified27 procedure and prior runtime reports remain intact.

validation: Metadata parses as exactly one package field. Existing allowlist remains47 entries. go generate ./template and root go test ./... pass, including payload/source byte equivalence. This verifies packaging preparation only; core installer is not yet implemented/released and no automatic install/discovery or new native behavior is claimed. Independent packaging review requested before PR handoff.

review: Independent packaging review inspected the current changes against repository instructions and the proposed source-install contract; no findings. Full report retained. Core contract review and actual installer implementation remain separate prerequisites for source-install runtime proof.

correction: Shared contract review found plugin root's go1.27.2 minimum would reject otherwise admissible released SKGo hosts built with Go1.27.1. Lowered only the plugin root directive to1.27.1. Canonical recorder/go.mod and retained release/native evidence are unchanged. With ignored payload temporarily moved away, GOTOOLCHAIN=go1.27.1 GOWORK=off Go generation and all root packaging tests pass; exact version and commands retained in go1271-tests.log. No dependency/API blocker found.

doc_bug: Manual helper directions must remain historical release qualification, not a parallel normal install flow. Label actual27 qualification without inventing28 runtime proof. Both27/28 predate managed installer. Docs explain that manual and managed copies on the same host conflict, SKGO_PLUGIN_DIRS does not hide managed plugins, and planned current-host remove MODULE is the remedy when keeping a manual copy. Commands remain unreleased; PR stays draft for second core contract review.

review: Independent packaging follow-up review against the revised shared contract reports no findings. It verified exact Go1.27.1 preparation/test evidence and scope of root directive/doc changes, while leaving actual installer/ABI/runtime qualification unproved. Complete report retained; packaging PR remains draft as requested pending second core contract review.
