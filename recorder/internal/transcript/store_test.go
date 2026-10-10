package transcript

import (
	"fmt"
	"testing"
)

func mustOpen(t *testing.T, s *Store, id string) {
	t.Helper()
	if _, err := s.Open(OpenRequest{id, "synthetic"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Activate(SessionRequest{id}); err != nil {
		t.Fatal(err)
	}
}
func TestRevisionFinalsAndIdempotence(t *testing.T) {
	s := New()
	mustOpen(t, s, "a")
	ch, cancel, _ := s.Watch()
	defer cancel()
	u := Update{"a", 1, []Chunk{}, []Chunk{{0, 100, "hel"}}}
	if _, e := s.Accept(u); e != nil {
		t.Fatal(e)
	}
	<-ch
	if _, e := s.Accept(u); e != nil {
		t.Fatal(e)
	}
	select {
	case <-ch:
		t.Fatal("duplicate notified")
	default:
	}
	u.Provisional[0].Text = "conflict"
	if _, e := s.Accept(u); e == nil {
		t.Fatal("conflicting duplicate accepted")
	}
	if _, e := s.Accept(Update{"a", 2, []Chunk{{0, 100, "hello "}}, []Chunk{{100, 200, "world"}}}); e != nil {
		t.Fatal(e)
	}
	_, stop, now := s.Watch()
	stop()
	if now.Text != "hello world" {
		t.Fatal(now)
	}
	if _, e := s.Accept(Update{"a", 3, []Chunk{{99, 110, "rewrite"}}, nil}); e == nil {
		t.Fatal("rewrote finalized range")
	}
	if _, e := s.Accept(Update{"a", 3, []Chunk{{100, 200, "world"}}, []Chunk{}}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Close(CloseRequest{"a", 2}); e == nil {
		t.Fatal("close gap accepted")
	}
	if _, e := s.Close(CloseRequest{"a", 3}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Close(CloseRequest{"a", 3}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Accept(Update{"a", 4, nil, nil}); e == nil {
		t.Fatal("write after close")
	}
	mustOpen(t, s, "b")
	if _, e := s.Accept(Update{"a", 4, nil, nil}); e == nil {
		t.Fatal("stale write accepted")
	}
}
func TestInvalidNeverMutatesAndReconnectSnapshot(t *testing.T) {
	s := New()
	mustOpen(t, s, "a")
	for _, u := range []Update{{"a", 2, nil, nil}, {"a", 1, nil, []Chunk{{-1, 2, "bad"}}}, {"a", 1, nil, []Chunk{{0, MaxInteger + 1, "bad"}}}} {
		if _, e := s.Accept(u); e == nil {
			t.Fatal("invalid accepted")
		}
	}
	_, cancel, now := s.Watch()
	cancel()
	if now.Sequence != 0 || now.Text != "" {
		t.Fatal(now)
	}
	u := Update{"a", 1, []Chunk{{0, 100, "text"}}, []Chunk{}}
	if _, e := s.Accept(u); e != nil {
		t.Fatal(e)
	}
	_, cancel, now = s.Watch()
	cancel()
	if now.Sequence != 1 || now.Text != "text" {
		t.Fatal("reconnect missing current", now)
	}
	if _, e := s.Open(OpenRequest{"b", "synthetic"}); e == nil {
		t.Fatal("replaced active session")
	}
}

func TestRepeatOpenAndPayloadBudget(t *testing.T) {
	s := New()
	mustOpen(t, s, "a")
	if ack, e := s.Open(OpenRequest{"a", "synthetic"}); e != nil || ack.Sequence != 0 {
		t.Fatal(ack, e)
	}
	if _, e := s.Open(OpenRequest{"a", "microphone"}); e == nil {
		t.Fatal("conflicting repeated open")
	}
	oversized := make([]byte, MaxPayloadBytes)
	for i := range oversized {
		oversized[i] = 'x'
	}
	if _, e := s.Accept(Update{"a", 1, nil, []Chunk{{0, 1, string(oversized)}}}); e == nil {
		t.Fatal("oversized payload accepted")
	}
	_, cancel, now := s.Watch()
	cancel()
	if now.Sequence != 0 || now.Text != "" {
		t.Fatal("oversize mutated", now)
	}
}

func TestSetupTransactionAndAbortRecovery(t *testing.T) {
	s := New()
	mustOpen(t, s, "a")
	if _, e := s.Accept(Update{"a", 1, []Chunk{{0, 1, "previous words"}}, nil}); e != nil {
		t.Fatal(e)
	}
	s.Close(CloseRequest{"a", 1})
	ch, cancel, _ := s.Watch()
	defer cancel()
	if _, e := s.Open(OpenRequest{"b", "synthetic"}); e != nil {
		t.Fatal(e)
	}
	select {
	case <-ch:
		t.Fatal("pending setup changed display")
	default:
	}
	_, stop, now := s.Watch()
	stop()
	if now.Text != "previous words" || now.SessionID != "a" {
		t.Fatal(now)
	}
	ack, e := s.Abort(SessionRequest{"b"})
	if e != nil || !ack.Closed {
		t.Fatal(ack, e)
	}
	again, e := s.Abort(SessionRequest{"b"})
	if e != nil || again != ack {
		t.Fatal("abort replay", again, e)
	}
	mustOpen(t, s, "c")
	s.Accept(Update{"c", 1, []Chunk{{0, 1, "accepted"}}, nil})
	ack, e = s.Abort(SessionRequest{"c"})
	if e != nil || ack.Sequence != 1 || !ack.Closed {
		t.Fatal(ack, e)
	}
	_, stop, now = s.Watch()
	stop()
	if now.Text != "accepted" || !now.Closed {
		t.Fatal(now)
	}
	mustOpen(t, s, "d")
	_, stop, before := s.Watch()
	stop()
	if _, e := s.Abort(SessionRequest{"b"}); e != nil {
		t.Fatal(e)
	}
	_, stop, after := s.Watch()
	stop()
	if before != after {
		t.Fatal("old abort changed new session")
	}
}

func TestClosedSessionRejectsIdenticalRetry(t *testing.T) {
	s := New()
	mustOpen(t, s, "a")
	u := Update{"a", 1, []Chunk{{0, 1, "final"}}, nil}
	if _, e := s.Accept(u); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Close(CloseRequest{"a", 1}); e != nil {
		t.Fatal(e)
	}
	ch, cancel, before := s.Watch()
	defer cancel()
	if _, e := s.Accept(u); e == nil {
		t.Fatal("identical write was admitted after close")
	}
	select {
	case <-ch:
		t.Fatal("closed retry notified")
	default:
	}
	_, stop, after := s.Watch()
	stop()
	if before != after {
		t.Fatal("closed retry mutated", after)
	}
}

func TestAbsentAbortPreservesDisplayAndRejectsLateOpen(t *testing.T) {
	s := New()
	mustOpen(t, s, "previous")
	s.Accept(Update{"previous", 1, []Chunk{{0, 1, "previous words"}}, nil})
	s.Close(CloseRequest{"previous", 1})
	ch, cancel, before := s.Watch()
	defer cancel()
	ack, e := s.Abort(SessionRequest{"never-delivered"})
	if e != nil || ack != (Ack{"never-delivered", 0, true}) {
		t.Fatal(ack, e)
	}
	select {
	case <-ch:
		t.Fatal("absent abort published")
	default:
	}
	_, stop, after := s.Watch()
	stop()
	if before != after {
		t.Fatal("absent abort cleared prior display", after)
	}
	if _, e := s.Open(OpenRequest{"never-delivered", "synthetic"}); e == nil {
		t.Fatal("late canceled open admitted")
	}
	mustOpen(t, s, "next")
	s.Accept(Update{"next", 1, []Chunk{{0, 1, "current words"}}, nil})
	_, stop, before = s.Watch()
	stop()
	if _, e := s.Abort(SessionRequest{"another-absent"}); e != nil {
		t.Fatal(e)
	}
	_, stop, after = s.Watch()
	stop()
	if before != after || after.Closed {
		t.Fatal("absent abort changed unrelated active session", after)
	}
	if again, e := s.Abort(SessionRequest{"never-delivered"}); e != nil || again != ack {
		t.Fatal("absence receipt forgotten", again, e)
	}
	s.Close(CloseRequest{"next", 1})
	if _, e := s.Open(OpenRequest{"never-delivered", "synthetic"}); e == nil {
		t.Fatal("later abort forgot canceled open")
	}
	mustOpen(t, s, "healthy")
}
func TestCanceledIdentityBudgetNeverEvictsSafety(t *testing.T) {
	s := New()
	for i := 0; i < MaxCanceledSessions; i++ {
		if _, e := s.Abort(SessionRequest{fmt.Sprintf("absent-%d", i)}); e != nil {
			t.Fatal(e)
		}
	}
	if _, e := s.Abort(SessionRequest{"one-too-many"}); e == nil {
		t.Fatal("cancellation budget unbounded")
	}
	if ack, e := s.Abort(SessionRequest{"absent-0"}); e != nil || !ack.Closed {
		t.Fatal("old cancellation evicted", ack, e)
	}
	if _, e := s.Open(OpenRequest{"absent-0", "synthetic"}); e == nil {
		t.Fatal("budget forgot delayed open")
	}
	if len(s.aborted) != MaxCanceledSessions {
		t.Fatal(len(s.aborted))
	}
}
