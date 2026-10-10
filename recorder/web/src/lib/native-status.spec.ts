import { describe, expect, it } from 'vite-plus/test';
import { DisplayRecovery, visibleRecorderState, displayCaughtUp, type NativeStatus } from './native-status';
const paused: NativeStatus = { sessionID: 'new', sequence: 4, state: 'paused', source: 'synthetic' };
describe('display completion', () => {
  it('requires its own session and final sequence for Paused and Stopped', () => {
    expect(visibleRecorderState(paused, { sessionID: 'old', sequence: 99 })).toBe('pausing');
    expect(visibleRecorderState(paused, { sessionID: 'new', sequence: 3 })).toBe('pausing');
    expect(visibleRecorderState(paused, { sessionID: 'new', sequence: 4 })).toBe('paused');
    expect(visibleRecorderState({ ...paused, state: 'stopped' }, { sessionID: 'new', sequence: 3 })).toBe('stopping');
    expect(visibleRecorderState({ ...paused, state: 'stopped' }, { sessionID: 'new', sequence: 5 })).toBe('stopped');
  });
});
describe('automatic display recovery', () => {
  it('retries without manual action and clears when a fresh display catches up', () => {
    const recovery = new DisplayRecovery();
    expect(recovery.check(0, false, true)).toBe(false);
    expect(recovery.check(249, false, true)).toBe(false);
    expect(recovery.check(250, false, true)).toBe(true);
    expect(recovery.check(999, false, true)).toBe(false);
    expect(recovery.check(1000, false, true)).toBe(true);
    expect(recovery.check(2000, true, false)).toBe(false);
    expect(recovery.error).toBe('');
    expect(recovery.expired).toBe(false);
  });
  it('expires at eight seconds and never retries or hides error after expiry', () => {
    const recovery = new DisplayRecovery();
    recovery.check(100, false, true);
    expect(recovery.check(8099, false, true)).toBe(true);
    expect(recovery.check(8100, false, true)).toBe(false);
    expect(recovery.expired).toBe(true);
    expect(recovery.error).toContain('final text is unconfirmed');
    expect(recovery.check(8200, true, true)).toBe(false);
    expect(recovery.error).toContain('relaunch');
  });
  it('does not extend a connected stream deadline while final display remains missing', () => {
    const recovery = new DisplayRecovery();
    recovery.check(0, false, false);
    expect(recovery.check(7000, false, false)).toBe(false);
    expect(recovery.check(8000, false, false)).toBe(false);
    expect(recovery.expired).toBe(true);
  });
});

// The same app can display an independently fed remote session without starting audio.
describe('external synthetic session', () => {
  it('reflects remote capture/close while native controls remain idle', () => {
    const idle: NativeStatus = { sessionID: '', sequence: 0, state: 'idle', source: 'microphone' };
    const incoming = { sessionID: 'external', sequence: 1, source: 'synthetic', closed: false };
    expect(visibleRecorderState(idle, incoming)).toBe('recording');
    const prior = { ...idle, sessionID: 'old-native', state: 'stopped' as const, sequence: 40 };
    expect(displayCaughtUp(prior, incoming)).toBe(true);
    expect(displayCaughtUp(prior, { ...incoming, closed: true })).toBe(true);
    const recovery = new DisplayRecovery();
    recovery.check(0, displayCaughtUp(prior, incoming), false);
    recovery.check(9000, displayCaughtUp(prior, incoming), false);
    expect(recovery.expired).toBe(false);
    expect(visibleRecorderState(idle, { ...incoming, closed: true })).toBe('stopped');
    const preparing = { ...prior, state: 'preparing' as const };
    expect(visibleRecorderState(preparing, incoming)).toBe('preparing');
    const setupRecovery = new DisplayRecovery();
    setupRecovery.check(0, displayCaughtUp(preparing, incoming), false);
    setupRecovery.check(9000, displayCaughtUp(preparing, incoming), false);
    expect(setupRecovery.expired).toBe(false);
    // Preparation must not excuse an actually disconnected subscription.
    setupRecovery.check(10000, false, true);
    setupRecovery.check(18000, false, true);
    expect(setupRecovery.expired).toBe(true);
  });
});
