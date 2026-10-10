# Reuse the working recorder

Create an independently named SKGo application from the recorder that already
works, without rebuilding its behavior or copying machine-specific setup.

This repository owns the reusable recorder source and the `voice-notes` plugin.
Shared SKGo owns project creation, `templateapi`, discovery, persisted identity,
the generic `native-app` shell/Go host/WebView/build foundation and compatible
updates. The selected central contract is
[application templates](https://github.com/tylergannon/skgo-project/blob/4f04eaf5027a893cdbffb7b3e4d17b4c5555d969/ephemeral/plans/application-templates.md).
The user subsequently authorized this new repository as canonical reusable
source; the original Voice Notes application remains separate and intact.

The source covers one continuous microphone transcription session: Record,
Pause/Resume, Stop, timestamped provisional/final text, native device selection,
compact timer/level meter, foreground display recovery and active iPhone
background recording. A scripted source and external synthetic transcript feeder
exercise delivery independently of Apple's audio adapter. No saving, transcript
editing, projects, sharing, archive or simultaneous sources are added.

Success requires a fresh directory, module, display name and bundle identity.
Template creation and manual add-on composition must yield equivalent configured
behavior. The resulting ordinary web server, macOS app and iPhone app must build
and launch without the source checkout or installed plugin. Native Swift commands
must reach Go and appear in the live Svelte view. Browser mode explains native
capture and accepts the documented synthetic feeder; no browser microphone is
added. The standalone transcription package must still run its independent checks.

Synthetic delivery, real microphone/Apple recognition, simulator behavior and
physical-device behavior remain separately identified observations. Existing
Voice Notes observations cannot prove a newly generated consumer. Mission
Control stays unproved until it holds that consumer's identified evidence.
