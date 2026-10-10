export type RecorderState = 'idle' | 'preparing' | 'recording' | 'pausing' | 'paused' | 'resuming' | 'stopping' | 'stopped' | 'error';
export interface NativeStatus {
  sessionID: string;
  sequence: number;
  state: RecorderState;
  source: string;
  error?: string;
}
export interface SnapshotReceipt { sessionID: string; sequence: number; closed?: boolean; source?: string }
declare global {
  interface Window {
    voiceNotesLatestStatus?: NativeStatus;
    voiceNotesNativeStatus?: (status: NativeStatus) => void;
    webkit?: { messageHandlers?: {
      bridgeReady?: { postMessage: (ready: Record<string, never>) => void };
      snapshotApplied?: { postMessage: (receipt: SnapshotReceipt) => void };
      displayAvailability?: { postMessage: (status: { available: boolean; error?: string }) => void };
    } };
  }
}
export function reportDisplay(available: boolean, error?: string) {
  window.webkit?.messageHandlers?.displayAvailability?.postMessage({ available, error });
}
export function reportSnapshot(receipt: SnapshotReceipt) {
  window.webkit?.messageHandlers?.snapshotApplied?.postMessage(receipt);
}

export function isExternalSnapshot(status: NativeStatus, snapshot: SnapshotReceipt): boolean {
  return !!snapshot.sessionID && snapshot.sessionID !== status.sessionID && snapshot.source === 'synthetic' && ['idle', 'stopped', 'error'].includes(status.state);
}

export function displayCaughtUp(status: NativeStatus, snapshot: SnapshotReceipt): boolean {
  // Preparation preserves the previous transcript until activation. No new
  // native session has been published yet; stream connectivity is checked separately.
  return status.state === 'preparing' || !status.sessionID || isExternalSnapshot(status, snapshot) || (status.sessionID === snapshot.sessionID && snapshot.sequence >= status.sequence);
}

export function visibleRecorderState(status: NativeStatus, snapshot: SnapshotReceipt): RecorderState {
  if (isExternalSnapshot(status, snapshot)) {
    return snapshot.closed ? 'stopped' : 'recording';
  }
  const behind = status.sessionID && (status.sessionID !== snapshot.sessionID || snapshot.sequence < status.sequence);
  return behind && status.state === 'paused' ? 'pausing' : behind && status.state === 'stopped' ? 'stopping' : status.state;
}

export class DisplayRecovery {
  private started: number | undefined;
  private retryAt = 0;
  expired = false;
  error = '';
  check(now: number, healthy: boolean, reconnectNeeded: boolean): boolean {
    if (this.expired) return false;
    if (healthy) { this.started = undefined; this.error = ''; return false; }
    if (this.started === undefined) { this.started = now; this.retryAt = now + 250; }
    if (now - this.started >= 8000) {
      this.expired = true;
      this.error = 'Transcript display unavailable. The final text is unconfirmed. Stop recording, then relaunch the application.';
      return false;
    }
    this.error = 'Transcript delivery interrupted. Reconnecting…';
    if (reconnectNeeded && now >= this.retryAt) { this.retryAt = now + 750; return true; }
    return false;
  }
}
