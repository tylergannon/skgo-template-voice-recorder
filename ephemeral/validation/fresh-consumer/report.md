# Fresh recorder template consumer — development evidence

Consumer: FieldNotes, module `example.com/field-notes`, base bundle
`com.example.fieldnotes`. Generated outside all source checkouts at
`/private/tmp/skgo-recorder-consumer.GthvK4/field-notes`. Framework dependency is
the actual published development module for `2bde811`, with no local replacement.
Recorder source fingerprints and executable identities are in
[build-metadata.json](build-metadata.json). The original Voice Notes application
at4bb9373 remains intact. These are fresh-consumer observations, not a release
reassessment of that original application.

## Demonstrated

- The shared source-build helper staged/built/loaded the actual factory for a
  separately Go-installed standalone CLI and for this consumer's actual `go tool
  skgo` host. Help lists voice-notes template/add-on and its option. Reapplying
  voice-notes through the project host reported no changed files.
- `skgo new --template voice-notes` applied builtin native-app followed by the
  real recorder plugin to an ordinary minimal TypeScript scaffold. A separately
  generated ordinary project followed by `add native-app`, `add voice-notes`
  matches55 reviewed source/config/foundation/module/identity files byte for byte.
  [manual-equivalence.json](manual-equivalence.json) retains the comparison result.
- macOS synthetic app compiled and launched through normal computer controls.
  Record changed native controls and synthetic words appeared in the Svelte
  view. Pause kept text and froze the timer at00:46 while the simulated meter
  cleared; Resume returned to Recording; Stop hid timer/meter and kept final text.
  Command-Comma opened FieldNotes Settings showing microphone-off synthetic mode.
  [Recording](macos-recording.png), [Paused](macos-paused.png),
  [Stopped](macos-stopped.png), [Settings](macos-settings.png).
- iPhone16/iOS27 simulator synthetic app compiled, installed and launched. Normal
  Device Hub computer use observed Record, streamed words, Pause/Resume and Stop
  retaining words. These are simulator/synthetic observations.
  [Recording](simulator-recording.png), [Paused](simulator-paused.png),
  [Stopped](simulator-stopped.png).
- Standalone production Go binary served browser mode at127.0.0.1:8188, explaining
  native microphone controls. Its external synthetic feeder submitted six
  provisional/final updates using fresh generated named codecs/remote registrations;
  the browser visibly ended with "This is synthetic text."
  [Browser](browser-synthetic.png), [feeder log](../field-notes-browser-feeder.log).
- macOS synthetic-burst compiled with the named burst fixture packaged unchanged,
  VNTranscriptFixture=burst-transcript and its independently suffixed bundle.
  [Resource check](burst-resources.json). Burst runtime is not yet observed.
- Real iPhone26 code compiled unsigned. Generated plist contains microphone and
  speech usage descriptions and background audio; this does not prove recording.

## Focused checks

Root Go race tests cover Apply/no-op/prerequisites, conflicts before writes,
unrelated configuration preservation, identity, filesystem safety and staging.
Fresh consumer Go race tests include the real built-app feeder handler, retry,
cancellation and competing-session cases. Svelte check reports zero errors or
warnings; four frontend test files / eleven tests pass. Existing independent Swift
audio/scripted package and recorder/input/controller checks pass without microphone
capture. Opus5.5 independent static integration review found no material defect;
it explicitly did not witness new hardware behavior. Logs and review are retained
under the sibling validation/reviews directories.

## Subsequent qualification

The initial signing issue was resolved by selecting the already signed-in Personal
Team in the generated Xcode project; Xcode created a profile without account
reauthentication. Signed installation and a user-observed real background speech
session subsequently passed for this development build. Detached native/browser
observations also passed. The tagged v0.27.0 commands, products, composition and
source/plugin-unavailable rebuild/runtime checks are separately recorded in
[release-report.md](release-report.md). Preserve this report and its metadata as
the earlier development snapshot; do not attribute its hardware result to the
later release-based product.
