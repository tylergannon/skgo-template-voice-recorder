<script lang="ts">
  import { onMount, tick, untrack } from 'svelte';
  import { watchTranscript } from './transcript.remote.js';
  import { reportDisplay, reportSnapshot, type NativeStatus, visibleRecorderState, displayCaughtUp, DisplayRecovery } from './native-status';

  let stream = $state<ReturnType<typeof watchTranscript>>();
  let snapshot = $state({ sessionID: '', sequence: 0, text: '', source: '', closed: false });
  let status = $state<NativeStatus>({ sessionID: '', sequence: 0, state: 'idle', source: '' });
  let deliveryError = $state('');
  let mounted = $state(false);
  let nativeCapture = $state(true);
  const recovery = new DisplayRecovery();
  let reconnecting = false;
  const labels = { idle: 'Ready', preparing: 'Preparing…', recording: 'Recording', pausing: 'Pausing…', paused: 'Paused', resuming: 'Resuming…', stopping: 'Stopping…', stopped: 'Stopped', error: 'Recording ended' };
  const displayState = $derived(visibleRecorderState(status, snapshot));

  $effect(() => {
    const current = stream?.current;
    if (!mounted || !current || recovery.expired) return;
    if (untrack(() => current.sessionID === snapshot.sessionID && current.sequence < snapshot.sequence)) return;
    snapshot = current;
    const receipt = { sessionID: current.sessionID, sequence: current.sequence, closed: current.closed };
    void tick().then(() => {
      if (!mounted || receipt.sessionID !== snapshot.sessionID || receipt.sequence !== snapshot.sequence) return;
      reportSnapshot(receipt);
      const caughtUp = displayCaughtUp(status, snapshot);
      reportDisplay(caughtUp, caughtUp ? undefined : 'Waiting for the final transcript…');
    });
  });

  onMount(() => {
    mounted = true;
    nativeCapture = !!window.webkit?.messageHandlers?.bridgeReady;
    stream = watchTranscript();
    reportDisplay(false, 'Connecting transcript display…');
    window.voiceNotesNativeStatus = (next) => {
      if (status.sessionID && next.sessionID !== status.sessionID && next.state !== 'preparing' && status.state !== 'preparing' && next.sessionID !== snapshot.sessionID) return;
      status = next;
      void tick().then(() => {
        if (mounted && snapshot.sessionID === next.sessionID) reportSnapshot({ sessionID: snapshot.sessionID, sequence: snapshot.sequence, closed: snapshot.closed });
      });
    };
    if (window.voiceNotesLatestStatus) window.voiceNotesNativeStatus(window.voiceNotesLatestStatus);
    window.webkit?.messageHandlers?.bridgeReady?.postMessage({});
    const timer = setInterval(() => {
      if (!stream || recovery.expired) return;
      const behind = !displayCaughtUp(status, snapshot);
      const disconnected = !!stream.error || stream.done || !stream.connected;
      const retry = recovery.check(Date.now(), !disconnected && !behind && stream.ready, disconnected && !reconnecting);
      deliveryError = recovery.error;
      if (deliveryError) reportDisplay(false, deliveryError);
      if (retry) {
        reconnecting = true;
        void stream.reconnect().catch(() => {}).finally(() => { reconnecting = false; });
      }
    }, 200);
    return () => {
      mounted = false;
      clearInterval(timer);
      delete window.voiceNotesNativeStatus;
    };
  });
</script>

<main>
  {#if !nativeCapture}<p class="capture-notice">Microphone capture uses the native application’s recording controls. This browser displays incoming transcripts, including synthetic playback.</p>{/if}
  <p class="status" aria-live="polite">{labels[displayState]}</p>
  {#if snapshot.source === 'synthetic' || status.source === 'synthetic'}<p class="source">Synthetic transcript</p>{/if}
  {#if status.error && status.error !== deliveryError}<p class="error" role="alert">{status.error}</p>{/if}
  {#if deliveryError}<p class="error" role="alert">{deliveryError}</p>{/if}
  <section aria-label="Transcript">
    {#if snapshot.text}<p class="transcript">{snapshot.text}</p>
    {:else}<p class="placeholder">Your words will appear here as you speak.</p>{/if}
  </section>
</main>

<style>
  :global(html) { color-scheme: light dark; background: light-dark(#faf9f6, #191b1e); }
  :global(body) { margin: 0; font-family: -apple-system, BlinkMacSystemFont, sans-serif; color: light-dark(#25292e, #f0f1f3); }
  main { max-width: 52rem; margin: 0 auto; padding: clamp(1.5rem, 5vw, 3.5rem); }
  .status { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip-path: inset(50%); white-space: nowrap; }
  .capture-notice { margin: 0 0 1rem; font-size: 0.9rem; line-height: 1.5; color: light-dark(#7a8087, #959ca7); }
  .source { margin: 0 0 0.75rem; color: light-dark(#855614, #e4b770); font-size: 0.75rem; }
  .transcript { font-size: clamp(1.25rem, 3vw, 1.7rem); line-height: 1.65; white-space: pre-wrap; overflow-wrap: anywhere; margin: 0; }
  .placeholder { font-size: 1.2rem; line-height: 1.6; color: light-dark(#7a8087, #959ca7); margin: 0; }
  .error { padding: 0.85rem 1rem; background: light-dark(#fff0ec, #3b2525); color: light-dark(#9b3028, #ffb9af); border-radius: 0.45rem; line-height: 1.5; }
</style>
