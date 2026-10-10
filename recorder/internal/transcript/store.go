package transcript

import (
	"encoding/json"
	"errors"
	"sync"
)

const MaxPayloadBytes = 65536
const MaxCanceledSessions = 1024
const MaxSessionBytes = 8 * 1024 * 1024
const MaxInteger = 9007199254740991

type Chunk struct {
	StartMs int64  `json:"startMs"`
	EndMs   int64  `json:"endMs"`
	Text    string `json:"text"`
}
type OpenRequest struct {
	SessionID string `json:"sessionID"`
	Source    string `json:"source"`
}
type Update struct {
	SessionID   string  `json:"sessionID"`
	Sequence    int64   `json:"sequence"`
	Finalized   []Chunk `json:"finalized"`
	Provisional []Chunk `json:"provisional"`
}
type CloseRequest struct {
	SessionID string `json:"sessionID"`
	Sequence  int64  `json:"sequence"`
}
type Ack struct {
	SessionID string `json:"sessionID"`
	Sequence  int64  `json:"sequence"`
	Closed    bool   `json:"closed"`
}
type Snapshot struct {
	SessionID string `json:"sessionID"`
	Sequence  int64  `json:"sequence"`
	Text      string `json:"text"`
	Source    string `json:"source"`
	Closed    bool   `json:"closed"`
}
type SessionRequest struct {
	SessionID string `json:"sessionID"`
}
type savedState struct {
	snapshot Snapshot
	final    string
	boundary int64
	last     string
}
type Store struct {
	visible  Snapshot
	pending  bool
	previous savedState
	aborted  map[string]Ack
	mu       sync.Mutex
	snapshot Snapshot
	final    string
	boundary int64
	last     string
	watchers map[chan Snapshot]struct{}
}

var Default = New()

func New() *Store {
	return &Store{watchers: make(map[chan Snapshot]struct{}), aborted: make(map[string]Ack)}
}
func (s *Store) ack() Ack { return Ack{s.snapshot.SessionID, s.snapshot.Sequence, s.snapshot.Closed} }
func (s *Store) publish() {
	if s.pending {
		return
	}
	s.visible = s.snapshot
	for ch := range s.watchers {
		select {
		case ch <- s.snapshot:
		default:
			select {
			case <-ch:
			default:
			}
			ch <- s.snapshot
		}
	}
}
func (s *Store) Open(r OpenRequest) (Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.SessionID == "" || len(r.SessionID) > 128 || (r.Source != "microphone" && r.Source != "synthetic") {
		return Ack{}, errors.New("invalid session")
	}
	if _, canceled := s.aborted[r.SessionID]; canceled {
		return Ack{}, errors.New("aborted session")
	}
	if r.SessionID == s.snapshot.SessionID {
		if r.Source != s.snapshot.Source {
			return Ack{}, errors.New("conflicting session")
		}
		return s.ack(), nil
	}
	if s.snapshot.SessionID != "" && !s.snapshot.Closed {
		return Ack{}, errors.New("previous session is still open")
	}
	s.previous = savedState{s.snapshot, s.final, s.boundary, s.last}
	s.pending = true
	s.snapshot = Snapshot{SessionID: r.SessionID, Source: r.Source}
	s.final = ""
	s.boundary = 0
	s.last = ""
	s.publish()
	return s.ack(), nil
}
func validate(chunks []Chunk, boundary int64) (int64, string, error) {
	text := ""
	for _, c := range chunks {
		if c.StartMs < boundary || c.EndMs < c.StartMs || c.EndMs > MaxInteger {
			return 0, "", errors.New("invalid transcript audio range")
		}
		boundary = c.EndMs
		text += c.Text
	}
	return boundary, text, nil
}
func (s *Store) Accept(r Update) (Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	payload, err := json.Marshal(r)
	if err != nil || len(payload) > MaxPayloadBytes {
		return Ack{}, errors.New("transcript update exceeds request budget")
	}
	if r.SessionID != s.snapshot.SessionID || r.SessionID == "" {
		return Ack{}, errors.New("stale session")
	}
	if s.snapshot.Closed {
		return Ack{}, errors.New("session is closed")
	}
	if r.Sequence == s.snapshot.Sequence && r.Sequence > 0 {
		if string(payload) == s.last {
			return s.ack(), nil
		}
		return Ack{}, errors.New("conflicting duplicate")
	}
	if r.Sequence != s.snapshot.Sequence+1 || r.Sequence > MaxInteger {
		return Ack{}, errors.New("out of order sequence")
	}
	boundary, final, err := validate(r.Finalized, s.boundary)
	if err != nil {
		return Ack{}, err
	}
	_, tail, err := validate(r.Provisional, boundary)
	if err != nil {
		return Ack{}, err
	}
	if len(s.final)+len(final)+len(tail) > MaxSessionBytes {
		return Ack{}, errors.New("transcript exceeds in-memory budget")
	}
	s.final += final
	s.boundary = boundary
	s.snapshot.Text = s.final + tail
	s.snapshot.Sequence = r.Sequence
	s.last = string(payload)
	s.publish()
	return s.ack(), nil
}
func (s *Store) Close(r CloseRequest) (Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ack, canceled := s.aborted[r.SessionID]; canceled && r.Sequence == ack.Sequence {
		return ack, nil
	}
	if r.SessionID == "" || r.SessionID != s.snapshot.SessionID {
		return Ack{}, errors.New("stale session")
	}
	if r.Sequence != s.snapshot.Sequence {
		return Ack{}, errors.New("close sequence does not match accepted transcript")
	}
	if s.pending {
		if r.Sequence != 0 {
			return Ack{}, errors.New("pending session cannot close after transcript updates")
		}
		return s.abort()
	}
	if !s.snapshot.Closed {
		s.snapshot.Closed = true
		s.publish()
	}
	return s.ack(), nil
}
func (s *Store) Watch() (chan Snapshot, func(), Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan Snapshot, 1)
	s.watchers[ch] = struct{}{}
	return ch, func() { s.mu.Lock(); delete(s.watchers, ch); s.mu.Unlock() }, s.visible
}

func (s *Store) Activate(r SessionRequest) (Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.SessionID == "" || r.SessionID != s.snapshot.SessionID {
		return Ack{}, errors.New("stale session")
	}
	if s.snapshot.Closed {
		return Ack{}, errors.New("session is closed")
	}
	if s.pending {
		s.pending = false
		s.previous = savedState{}
		s.publish()
	}
	return s.ack(), nil
}

func (s *Store) abort() (Ack, error) {
	if len(s.aborted) >= MaxCanceledSessions {
		return Ack{}, errors.New("Session cleanup memory limit reached. Relaunch the application.")
	}
	ack := Ack{s.snapshot.SessionID, s.snapshot.Sequence, true}
	s.aborted[ack.SessionID] = ack
	if s.pending {
		s.snapshot = s.previous.snapshot
		s.final = s.previous.final
		s.boundary = s.previous.boundary
		s.last = s.previous.last
		s.previous = savedState{}
		s.pending = false
	} else {
		s.snapshot.Closed = true
		s.publish()
	}
	return ack, nil
}
func (s *Store) Abort(r SessionRequest) (Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.SessionID == "" || len(r.SessionID) > 128 {
		return Ack{}, errors.New("invalid session")
	}
	if ack, canceled := s.aborted[r.SessionID]; canceled {
		return ack, nil
	}
	if r.SessionID == s.snapshot.SessionID {
		return s.abort()
	}
	if len(s.aborted) >= MaxCanceledSessions {
		return Ack{}, errors.New("Session cleanup memory limit reached. Relaunch the application.")
	}
	ack := Ack{r.SessionID, 0, true}
	s.aborted[r.SessionID] = ack
	return ack, nil
}
