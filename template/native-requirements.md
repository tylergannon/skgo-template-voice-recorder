# Recorder contributions to native-app

The generic shell and host remain SKGo-owned. The recorder application slot must
retain one controller/input owner for the life of its host, including background
and foreground transitions, and contribute a macOS Settings scene. It calls
controller.shutdown during app shutdown and setAppBackgrounded during lifecycle
changes. Host recovery retains its handler, transcript state and loopback origin.

Recorder-specific build inputs must add the independent Transcription package and
its TranscriptSource, AppleTranscription and ScriptedTranscription products,
fixture and license resources, and the following application permission/lifecycle
settings before any fresh consumer can be considered runnable:

- NSMicrophoneUsageDescription: explain live microphone transcription.
- NSSpeechRecognitionUsageDescription: explain Apple's on-device transcription.
- UIBackgroundModes includes audio on iOS, so active capture may continue while
  hidden or locked.

Preserve real/synthetic and optional synthetic variant bundle IDs, VNDefaultSource
metadata and VN_SOURCE launch selection. No developer team, device or credentials
are packaged. Exact config edits await the shared native-app handoff; copying the
original generic Info.plist/build pipeline is intentionally outside this source
allowlist. Current package/controller checks do not prove those consumer settings.
