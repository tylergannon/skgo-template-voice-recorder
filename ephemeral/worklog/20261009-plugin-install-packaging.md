# Recorder source-install packaging

decision: The coordinator routed the user's source installer UX to shared SKGo. Core owns fetch, host-matched build, installation and discovery; recorder owns only its package convention and payload preparation. Preserve canonical application behavior and released27 evidence. Archive distribution is deferred.

finding: Clean178135e source has no embedded payload, but go generate ./template succeeds and runs its existing directive, staging the explicit47-entry allowlist. A conventional package path plus ordinary Go generation may avoid a custom preparation manifest or installer hook. Coordinate schema before editing rather than invent duplicate installer contracts.

finding: Canonical recorder/ is a nested Go module. x/mod v0.41 zip.CheckFiles omits nested module files; root Go-module zip retrieval cannot provide the reviewed47-file payload inputs. Sent this constraint to the installer owner before choosing a fetch convention. Preserve the full repository source tree instead of changing recorder architecture to fit an incomplete fetch.
