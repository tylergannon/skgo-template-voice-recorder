import { page } from 'vite-plus/test/browser';
import { describe, expect, it, vi } from 'vite-plus/test';
import { render } from 'vitest-browser-svelte';
import { createSubscriber } from 'svelte/reactivity';
import { watchTranscript } from './transcript.remote.js';
import TranscriptScreen from './TranscriptScreen.svelte';
vi.mock('./transcript.remote.js', () => ({ watchTranscript: vi.fn() }));

describe('rendered transcript delivery', () => {
  it('hydrates cached native status, re-sends an early receipt, and waits for delayed final text', async () => {
    let value = { sessionID: 'session', sequence: 0, text: 'First words', source: 'synthetic', closed: false };
    let notify = () => {};
    const subscribe = createSubscriber((update) => { notify = update; });
    const stream = {
      get current() { subscribe(); return value; },
      ready: true, connected: true, done: false, error: undefined,
      reconnect: vi.fn().mockResolvedValue(undefined),
    };
    vi.mocked(watchTranscript).mockReturnValue(stream as unknown as ReturnType<typeof watchTranscript>);
    const receipts: { sessionID: string; sequence: number }[] = [];
    window.webkit = { messageHandlers: {
      bridgeReady: { postMessage: vi.fn() },
      snapshotApplied: { postMessage: (receipt) => {
        // Receipt must follow rendered text, not merely assigning a snapshot.
        expect(document.querySelector('.transcript')?.textContent).toBe(value.text);
        receipts.push(receipt);
      } },
      displayAvailability: { postMessage: vi.fn() },
    } };
    window.voiceNotesLatestStatus = { sessionID: '', sequence: 0, state: 'idle', source: 'synthetic' };
    const screen = await render(TranscriptScreen);
    await expect.element(page.getByText('Synthetic transcript')).toBeInTheDocument();
    expect(document.querySelector('.capture-notice')).toBeNull();
    await expect.element(page.getByText('First words')).toBeInTheDocument();
    await vi.waitFor(() => expect(receipts.length).toBeGreaterThan(0));
    const before = receipts.length;
    window.voiceNotesNativeStatus?.({ sessionID: 'session', sequence: 0, state: 'recording', source: 'synthetic' });
    await vi.waitFor(() => expect(receipts.length).toBeGreaterThan(before));
    window.voiceNotesNativeStatus?.({ sessionID: 'session', sequence: 1, state: 'paused', source: 'synthetic' });
    await expect.element(page.getByText('Pausing…', { exact: true })).toBeInTheDocument();
    value = { ...value, sequence: 1, text: 'First words and final words' };
    notify();
    await expect.element(page.getByText('Paused', { exact: true })).toBeInTheDocument();
    await expect.element(page.getByText(value.text)).toBeInTheDocument();
    await vi.waitFor(() => expect(receipts.at(-1)?.sequence).toBe(1));
    await screen.unmount();
    delete window.webkit;
    delete window.voiceNotesLatestStatus;
  });
});

describe('stream failures preserve rendered text', () => {
  it('automatically reconnects, then expires repeated failure without losing accepted words', async () => {
    let connected = true;
    let value = { sessionID: 'session', sequence: 1, text: 'Accepted words remain', source: 'synthetic', closed: false };
    let notify = () => {};
    const subscribe = createSubscriber((update) => { notify = update; });
    const reconnect = vi.fn().mockResolvedValue(undefined);
    const stream = {
      get current() { subscribe(); return value; },
      get connected() { return connected; },
      ready: true, done: false, error: undefined, reconnect,
    };
    vi.mocked(watchTranscript).mockReturnValue(stream as unknown as ReturnType<typeof watchTranscript>);
    const availability = vi.fn(({ available, error }: { available: boolean; error?: string }) => {
      window.voiceNotesNativeStatus?.({ sessionID: 'session', sequence: value.sequence, state: 'recording', source: 'synthetic', error: available ? undefined : error });
    });
    window.webkit = { messageHandlers: { displayAvailability: { postMessage: availability } } };
    const screen = await render(TranscriptScreen);
    await expect.element(page.getByText(value.text)).toBeInTheDocument();
    connected = false;
    await vi.waitFor(() => expect(reconnect).toHaveBeenCalled());
    // A fresh snapshot on recovery must clear the error and acknowledge availability.
    connected = true;
    value = { ...value, sequence: 2, text: 'Accepted words remain after recovery' };
    notify();
    await expect.element(page.getByText(value.text)).toBeInTheDocument();
    await vi.waitFor(() => expect(document.querySelector('[role="alert"]')).toBeNull());
    expect(availability).toHaveBeenCalledWith({ available: true, error: undefined });
    connected = false;
    await vi.waitFor(() => expect(document.querySelector('[role="alert"]')?.textContent).toContain('final text is unconfirmed'), { timeout: 9500 });
    expect(document.querySelectorAll('[role="alert"]').length).toBe(1);
    window.voiceNotesNativeStatus?.({ sessionID: 'session', sequence: value.sequence, state: 'error', source: 'synthetic', error: 'The speech source failed.' });
    await vi.waitFor(() => expect(document.querySelectorAll('[role="alert"]').length).toBe(2));
    const attempts = reconnect.mock.calls.length;
    await new Promise((resolve) => setTimeout(resolve, 900));
    expect(reconnect.mock.calls.length).toBe(attempts);
    await expect.element(page.getByText(value.text)).toBeInTheDocument();
    expect(availability.mock.lastCall?.[0].available).toBe(false);
    await screen.unmount();
    delete window.webkit;
  }, 15000);
});


describe('external sessions between native recordings', () => {
  it('acknowledges a new external session and a native restart before its live snapshot arrives', async () => {
    let value = { sessionID: 'old-native', sequence: 4, text: 'Old words', source: 'microphone', closed: true };
    let notify = () => {};
    const subscribe = createSubscriber((update) => { notify = update; });
    const stream = { get current() { subscribe(); return value; }, ready: true, connected: true, done: false, error: undefined, reconnect: vi.fn().mockResolvedValue(undefined) };
    vi.mocked(watchTranscript).mockReturnValue(stream as unknown as ReturnType<typeof watchTranscript>);
    const availability = vi.fn();
    window.webkit = { messageHandlers: { displayAvailability: { postMessage: availability } } };
    window.voiceNotesLatestStatus = { sessionID: 'old-native', sequence: 4, state: 'stopped', source: 'microphone' };
    const screen = await render(TranscriptScreen);
    await expect.element(page.getByText('Old words')).toBeInTheDocument();
    value = { sessionID: 'external', sequence: 1, text: 'Corrected synthetic words', source: 'synthetic', closed: false }; notify();
    await expect.element(page.getByText(value.text)).toBeInTheDocument();
    await expect.element(page.getByText('Recording', { exact: true })).toBeInTheDocument();
    await vi.waitFor(() => expect(availability.mock.lastCall?.[0].available).toBe(true));
    value = { ...value, closed: true }; notify();
    await expect.element(page.getByText('Stopped', { exact: true })).toBeInTheDocument();
    let now = Date.now();
    const clock = vi.spyOn(Date, 'now').mockImplementation(() => now);
    window.voiceNotesNativeStatus?.({ sessionID: 'old-native', sequence: 4, state: 'preparing', source: 'microphone' });
    await expect.element(page.getByText('Preparing…', { exact: true })).toBeInTheDocument();
    await new Promise((resolve) => setTimeout(resolve, 250));
    now += 9000;
    await new Promise((resolve) => setTimeout(resolve, 250));
    expect(document.querySelector('[role="alert"]')).toBeNull();
    expect(availability.mock.lastCall?.[0].available).toBe(true);
    clock.mockRestore();
    window.voiceNotesNativeStatus?.({ sessionID: 'new-native', sequence: 0, state: 'recording', source: 'microphone' });
    value = { sessionID: 'new-native', sequence: 1, text: 'New native words', source: 'microphone', closed: false }; notify();
    await expect.element(page.getByText(value.text)).toBeInTheDocument();
    await expect.element(page.getByText('Recording', { exact: true })).toBeInTheDocument();
    await vi.waitFor(() => expect(availability.mock.lastCall?.[0].available).toBe(true));
    await screen.unmount(); delete window.webkit; delete window.voiceNotesLatestStatus;
  });
});


describe('browser-only transcript display', () => {
  it('explains native capture and still displays externally submitted synthetic text', async () => {
    delete window.webkit;
    delete window.voiceNotesLatestStatus;
    const stream = { current: { sessionID: 'external-browser', sequence: 1, text: 'Synthetic browser words', source: 'synthetic', closed: false }, ready: true, connected: true, done: false, error: undefined, reconnect: vi.fn().mockResolvedValue(undefined) };
    vi.mocked(watchTranscript).mockReturnValue(stream as unknown as ReturnType<typeof watchTranscript>);
    const screen = await render(TranscriptScreen);
    await expect.element(page.getByText('Synthetic browser words')).toBeInTheDocument();
    await expect.element(page.getByText('Microphone capture uses the native application’s recording controls. This browser displays incoming transcripts, including synthetic playback.')).toBeInTheDocument();
    expect(document.querySelector('button')).toBeNull();
    await screen.unmount();
  });
});
