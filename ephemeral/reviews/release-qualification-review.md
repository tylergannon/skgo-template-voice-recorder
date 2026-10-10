I've reviewed the recorder template and its released-foundation evidence, read-only; I didn't edit any files or record audio. Outcome: **only nitpicks remain**. Below is the complete review for you to save as `ephemeral/reviews/release-qualification-review.md`.

---

# Adversarial review: released-foundation qualification for the recorder template

**Target:** `skgo-template-voice-recorder` at source checkpoint `a4813333e0f0a6ea117fa40f7e67f0ecbe442f9e`, plus uncommitted changes:
- the dependency move to SKGo `v0.27.0` (`go.mod`, `go.sum`);
- `README.md` "Create a consumer";
- `ephemeral/skgo-feedback.md`, the worklog and `fresh-consumer/report.md`;
- the qualification artifact `ephemeral/validation/fresh-consumer/release-report.md`.

**Scope note:** The scope came from AGENTS.md, VISION.md and README.md. The caller's framing, that the physical observations are human reports, was treated as context only. It did not narrow what was reviewed.

## Evidence inspected

- **Requirements:** AGENTS.md, VISION.md (whole file), the README diff, the go.mod diff, the skgo-feedback diff, `.gitignore`, `scripts/`, and the git-tracked file list (allowlisted `recorder/`, `internal/packagefiles`).
- **Release report and metadata:**
  - `release-report.md` (whole file).
  - `released-build-metadata.json`: recorder commit matches HEAD, `replaces: []`, SKGo v0.27.0/e64e350, three product bundle IDs and hashes, iPhone/simulator `backgroundModes: [audio]`, macOS none.
  - `released-manual-equivalence.json`: 55 paths compared.
- **Logs:** `released-project-reapply.log` (ends `applied voice-notes`; the report's no-change claim is based on the caller's comparison), `plugin-released-checks.log`, `released-iphone-launch.json` (CoreDevice 10002 / FBS `RequestDenied`, matching the report's "locked CLI launch").
- **Detached run:**
  - `released-detached-check.py`: renames the source checkout and toolchain away, rebuilds the server and macOS synthetic app, copies the app to runtime, and restores everything in `finally`.
  - `released-detached-build.log`, `released-detached-process.txt` (the copied runtime app was running, and a `test ! -e <checkout>` guard ran before the feeder).
  - `released-detached-feeder.log`: six updates accepted; final "This is synthetic text."
  - `released-detached-server.log`.
- **Screenshots:** `released-detached-native.png` (FieldNotes recording at 00:31 with fixture text, meter and timer visible) and `released-macos-stopped.png` (text kept after Stop, Record control shown).
- **Leak check:** `git grep` across all non-ephemeral tracked files for `DEVELOPMENT_TEAM`, `/Users/tyler` and the user's email found nothing. The released metadata contains no home-directory paths.

## Assessment against the requirements

- **Fresh name, module and bundle:** Met (FieldNotes, `example.com/field-notes`, `com.example.fieldnotes[.synthetic|.real]`).
- **Template creation vs. manual composition:** The 55 compared files are identical.
- **Built and launched without the source checkout or plugin:**
  - Web and macOS: proven by the detached run.
  - iPhone: the installed product is self-contained, so the requirement doesn't apply in the same way. The report says so explicitly.
- **Swift → Go → Svelte path:** Shown on macOS and in the simulator.
- **Browser feeder:** Exercised.
- **Evidence kept separate:**
  - Synthetic, simulator and physical observations are kept apart.
  - Physical speech is correctly attributed to the earlier `2bde811` development product, not to the released product.
  - The released iPhone's launch is labelled as the user's own report ("FieldNotes opens").
  - Original Voice Notes evidence isn't used as proof.
  - Mission Control isn't claimed.
- **Uses the real shared foundation:** Yes. Released v0.27.0 with no local module replacement.
- **No local SDK or competing bootstrap:** The README creation path uses the real `skgo new --template voice-notes`.
- **Licenses:** Both licenses are kept and pass through the allowlist.
- **Known gap logged:** The SDK staging gap is recorded in skgo-feedback.md and disclosed under Limits.

No material requirement gaps, incorrect claims or misattributed evidence were found.

## Findings

1. **Nitpick: the detached-run script's port setting had no effect, so the script and its record disagree.**
   - `released-detached-check.py:25` starts the rebuilt server with `ADDR='127.0.0.1:8188'`.
   - But `released-detached-server.log` shows `listening on http://127.0.0.1:8080`, and `released-detached-process.txt` shows the feeder targeting `:8080`.
   - The generated server evidently ignores `ADDR`. The proof still holds: the detached binary wrote that log and owned 8080, and the feeder's success matches it.
   - The impact is that a reader can't tell from the script which port was meant. The retained script also doesn't record that the feeder and browser check were run by hand.
   - Fix: correct or annotate the script, or note that the default port was used.

2. **Nitpick: root `go test ./...` picks up the ignored toolchain.**
   - `plugin-released-checks.log:1` lists `build/toolchain/node_modules/flatted/golang/pkg/flatted` as a package of this module.
   - Because `build/` is in `.gitignore`, git ignores it, but Go doesn't. Any Go files under the local toolchain become part of the repository's test and vet set.
   - It's harmless today, but it's a fragile dependency on the local machine's state.
   - Fix: `./internal/... ./template/...`, or a `go.mod` boundary or `_`-prefixed directory for the toolchain.

3. **Nitpick: the README's creation steps are only partly qualified.**
   - `README.md` "Create a consumer" requires `scripts/build-plugin.sh` and a plugin host in a matching build.
   - The release report covers this through the shared helper logs, but the README doesn't say the plugin `.so` has to be rebuilt whenever the host or Go toolchain changes.
   - It also doesn't link to the v0.27.0 `sdk` staging limitation recorded in the report.
   - The impact is limited to how clear the README is for consumers.

**Outcome: only nitpicks remain**
