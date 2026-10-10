I found no material defects in the scoped changes. I didn't modify anything or run capture. My only checks were reading the source, the tests and the SKGo `2bde811` code in the module cache. I didn't re-run the Go tests or native builds.

**What I checked and found consistent**
- **`template/apply.go`:**
  - Apply runs only when `native/skgo-native.json`, `NativeHost.swift` and `ApplicationView.swift` already exist, which is what the `native-app` → `voice-notes` template order guarantees.
  - The application view is replaced only if it still matches the core default exactly (`nativeapp.DefaultApplicationView`).
  - The page and config also use the narrow replace-if-unedited rule in `internal/packagefiles`.
  - Running Apply a second time changes nothing: the merged config re-serializes identically and the installed view already matches.
  - Re-running `native-app` after `voice-notes` is safe, because `nativeapp/plugin.go:112` skips both app-owned files once configured.
- **`template/config.go`:**
  - It keeps unknown fields and fields from other add-ons, and an unchanged value is a no-op.
  - Real conflicts return an error instead of overwriting: a different scalar, a package with the same name or path, or a mismatched preset `goTags`/`resources`.
  - The result is re-checked with the real `nativeapp.Config.Validate`.
  - Presets that omit `resources` keep the base resources, which is correct under the core's nil-inherits rule (`config.go:42`). The burst preset's explicit resource list is deliberate.
- **`template/recorder.go`:**
  - It renames Go imports to the consumer's module and changes only the probe plist identity and the project-name marker.
  - It fails on any leftover `__SKGO_` marker or any payload file not in the allowlist.
  - The bridge and product names stay fixed, so a consumer with a different name won't break the Swift/JS bridge contract.
  - The new `ApplicationView.swift` is in the allowlist.
- **Plugin factory and version stamp:** The exported `SKGoPluginV1` and the `skgoVersion` variable match the core build helper, and `scripts/build-plugin.sh` passes arguments through to the real `skgo-plugin-build` after staging.
- **`ApplicationView.swift`:** Its type signatures match exactly what the core `Application.swift` constructs. It creates one input and one controller per host, which survive resume. It connects `onShutdown` to shutdown, `onBackgroundChange` to `setAppBackgrounded`, and forwards availability loss only. Ignoring availability gain is correct because the revision bump triggers a reload and WebKit then reports the display as available.
- **`RecordingInput` fixture selection:**
  - `VN_TRANSCRIPT_FIXTURE` still takes precedence.
  - The default name is still `synthetic-transcript`.
  - Invalid names give a configuration error and never fall back to the microphone.
  - A missing resource still points to a URL that doesn't exist, as before.
  - The tests confirm none of these paths lists microphones.

**Minor (not material)**
- **`RecorderScreen.swift:72–78` / `:109` (macOS):**
  - **Trigger:** a new web view coordinator is created after a recovery has already happened (host revision above 0), for example when a new window is opened.
  - **Effect:** the coordinator starts at revision 0, so its first update treats this as a recovery. The page reloads right after the initial load and briefly shows "reconnecting".
  - **Scope:** this is harmless and unlikely, since closing the last window quits the app. The iOS version already behaved this way.
- **Defaults I couldn't confirm locally:** `defaultPage` (`apply.go:22`) is only correct if the scaffold's generated SvelteKit page matches it byte for byte. The shared SKGo module doesn't include that page, so if they differ, a fresh consumer fails with a page conflict. Your fresh scaffold run will settle this.

**Claims not yet backed by running evidence (not defects)**
- None of the following has been shown on a fresh consumer:
  - That its native build compiles with the merged packages, products and presets.
  - That the macOS forwarding actually reloads the transcript after a real host recovery.
  - That background and shutdown hooks behave correctly in a running app.
  - That the burst preset's bundle actually contains `burst-transcript.json`.
- The existing package/controller tests and the original app's baseline don't prove any of these.
