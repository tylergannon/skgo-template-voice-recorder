# Fresh recorder consumer — released foundation qualification

Recorder application/plugin source: `a4813333e0f0a6ea117fa40f7e67f0ecbe442f9e`.
Host and generated dependency: released SKGo `v0.27.0`, source commit
`e64e3509973e3aa4bad0c2ab7e46dbaff9251dd1`, without local module replacement.
Fresh recipe consumer is `/private/tmp/skgo-recorder-consumer.GthvK4/released`,
named FieldNotes, module `example.com/field-notes`, base bundle
`com.example.fieldnotes`. It is separate from both this template and original
Voice Notes. [Released metadata](released-build-metadata.json) records the copied
source hashes and executable identities. [Development report](report.md) retains
prior build-specific observations and focused checks.

## Released commands and composition

The actual shared helper built and loaded the plugin against a separately installed
v0.27.0 standalone CLI and the actual v0.27.0 consumer Go-tool host.
[Standalone helper](../plugin-released-build.log),
[project helper](../plugin-released-project.log).
The released standalone `new --template voice-notes` created the ordinary project,
applied builtin native-app and then the recorder plugin.
[Creation](../released-new.log).
A separate released ordinary `new`, `add native-app`, `add voice-notes` composition
matches 55 reviewed recorder/foundation/config/module/identity/server files exactly.
[Comparison with exact paths](released-manual-equivalence.json),
[manual commands](../released-manual.log).
The v0.27.0 project-tool host reapplied voice-notes with no changed files.
[No-op output](../released-project-reapply.log).
Root race checks passed against released dependencies.
[Checks](../plugin-released-checks.log).

## Running released products

- macOS 27 synthetic product built and opened with normal computer controls.
  Record produced text in the Svelte view through Swift and Go. Stop retained it.
  [Recording](released-macos-recording.png), [Stopped](released-macos-stopped.png),
  [build](../released-macos.log).
- iPhone 16/iOS 27 simulator synthetic product built, installed, launched and
  displayed streamed fixture text. Stop ended the visible recording controls.
  [Recording](released-simulator-recording.png), [build](../released-simulator.log).
- Real iPhone product signed and installed as `com.example.fieldnotes.real` on the
  connected iPhone 16/iOS 26.6.2. CLI launch was denied while the phone was locked;
  the user opened the installed release-based update and answered **“FieldNotes
  opens.”** This is human-reported physical launch, not an agent screenshot.
  [Build](../released-iphone.log), [install](released-iphone-install.json),
  [locked CLI launch](released-iphone-launch.json).

## Plugin and template checkout unavailable

The template worktree was temporarily renamed away from its original path,
including all its plugin directories. The unrelated original Voice Notes checkout
was left intact. The pinned development toolchain was moved to an independent
temporary directory so compilation did not borrow tooling from the missing source
path. A bounded guard restored both paths and stopped its server afterward.

During that absence, the released consumer rebuilt its standalone Go server and
its macOS synthetic native application. The copied native product ran from
`released-runtime/FieldNotes.app`; Record visibly delivered text. The independently
rebuilt Go binary ran from that runtime directory, and the consumer's external
feeder submitted six provisional/final updates. Normal browser computer use
observed **“This is synthetic text.”** The native app was stopped and quit before
restoration. [Build log](released-detached-build.log),
[native process identity](released-detached-process.txt),
[native screenshot](released-detached-native.png),
[browser screenshot](released-detached-browser.png),
[feeder](released-detached-feeder.log),
[server](released-detached-server.log),
[one-off guard](released-detached-check.py).
These establish build and runtime independence for the generated web/macOS app.
The physical phone runs a self-contained installed product; no plugin is shipped.
The guard's `ADDR` environment variable was ignored: the server log established
the actual default origin at port 8080. The feeder and CUA observations were
performed separately against that origin while the guard kept the checkout absent.

## Separate real speech observation

The previous fresh FieldNotes signed product used the published development module
at `2bde811` and the same recorder source `a481333`. The user performed Record,
spoke before backgrounding, on Home and while locked, reopened it and stopped.
Their reply was **“Yes, all three appeared.”** This demonstrates the short physical
microphone/Apple recognition/background session for that identified development
product. [Product hash and exact question/reply](iphone-user-check.json),
[successful install](iphone-install.json). It is not attributed to the later tagged
release product and does not establish routing-change or long-duration behavior.
The user did not separately answer the meter part of the question.

The development macOS synthetic-burst product also visibly reached its final
fixture sentence while Recording and retained it after Stop.
[Burst recording](burst-recording.png), [Stopped](burst-stopped.png).
Its visible sentence is not a claim that every intermediate revision was rendered.

## Limits

No new real desktop microphone assessment was performed; synthetic Mac/simulator
and browser results cannot establish Apple's recognition. Existing original-app
release assessments remain separate. Comprehensive routing changes, startup Stop,
foreground recovery fault scenarios and long recording reassessment are outside
this fresh-template qualification. SKGo compatible-update work is separately owned
and not included in v0.27.0 qualification; the next update-capable release needs its
own consumer check. The released native SDK still stages under `sdk`; the peer's
subsequent `.sdk` correction was not silently applied here.

## Independent release review and clean source build

Claude Opus 5.5 reviewed the requirements, implementation, release reports,
metadata, logs and screenshots read-only. It found only nitpicks, with no material
requirement gap or misattributed observation.
[Complete review](../../reviews/release-qualification-review.md).
The retained guard/default-port explanation and host/toolchain rebuild guidance
address its documentation notes. Root `ignore ./build` excludes disposable tooling
from Go discovery; [focused root race output](../plugin-released-focused-checks.log)
contains only owned packages. This does not alter consumer application sources.
A clean GitHub clone of checkpoint `736953c` staged its own allowlisted payload and
built/loaded the plugin with the released standalone host, with no preexisting
payload or original-app dependency. [Clean clone/helper log](../final-clean-clone.log).
