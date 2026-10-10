# SKGo extraction feedback

## Open: native application extension seam

Version: templateapi SDK commit `4badfe9` (prerelease).
Reference: merged application-template contract and framework peer handoff.
Expected: recorder replaces the application view, contributes its Settings and
lifecycle callbacks, and adds package/resource/plist inputs without copying a
generic host or build pipeline.
Observed: SDK is available; native-app extension bytes/config are not yet
delivered. This is an implementation prerequisite, not an observed SDK defect.

## Confirmed: fixture codec tied to generated root ordering

Version: Voice Notes source `4bb9373`, SKGo `v0.26.1`.
Reference: original `app/cmd/transcript-fixture/main.go`.
Observed: feeder calls generated `EncodeRoot0`, `EncodeRoot2` and similar numbered
functions. A fresh consumer's different type graph may reorder these identifiers.
Expected: the reusable feeder encodes named transcript contracts without assuming
the original graph order. Resolving this coupling belongs to extraction; no
fresh-consumer failure has been observed yet.
