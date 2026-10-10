package main

import (
	"bytes"
	"context"
	"errors"
	app "github.com/tylergannon/skgo-template-voice-recorder/recorder"
	"github.com/tylergannon/skgo-template-voice-recorder/recorder/internal/transcript"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

func application(t *testing.T) client {
	t.Helper()
	transcript.Default = transcript.New()
	server := httptest.NewUnstartedServer(nil)
	origin := "http://" + server.Listener.Addr().String()
	handler, _, err := app.NewHandler(os.DirFS("../../web/build"), "", origin)
	if err != nil {
		t.Fatal(err)
	}
	server.Config.Handler = handler
	server.Start()
	t.Cleanup(server.Close)
	return client{origin, server.Client()}
}
func snapshot() transcript.Snapshot { _, cancel, s := transcript.Default.Watch(); cancel(); return s }
func TestFixtureThroughBuiltApplication(t *testing.T) {
	c := application(t)
	events := []event{{0, 0, 400, "I has", false}, {0, 500, 700, " a cat", false}, {0, 0, 450, "I have", false}, {0, 0, 450, "I have", true}, {0, 500, 700, " a dog.", false}, {0, 500, 700, " a dog.", true}}
	var output bytes.Buffer
	if err := play(context.Background(), c, events, &output); err != nil {
		t.Fatal(err)
	}
	s := snapshot()
	if s.Text != "I have a dog." || s.Source != "synthetic" || s.Sequence != 6 || !s.Closed {
		t.Fatalf("snapshot: %+v", s)
	}
	if !strings.Contains(output.String(), "acknowledged final text: I have a dog.") {
		t.Fatal(output.String())
	}
}
func TestCompetingSessionPreserved(t *testing.T) {
	c := application(t)
	ctx := context.Background()
	if _, err := c.call(ctx, "openSession", transcript.OpenRequest{SessionID: "microphone-owner", Source: "microphone"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.call(ctx, "activateSession", transcript.SessionRequest{SessionID: "microphone-owner"}); err != nil {
		t.Fatal(err)
	}
	u := transcript.Update{SessionID: "microphone-owner", Sequence: 1, Finalized: []transcript.Chunk{{StartMs: 0, EndMs: 500, Text: "Keep this."}}, Provisional: []transcript.Chunk{}}
	if _, err := c.call(ctx, "acceptTranscript", u); err != nil {
		t.Fatal(err)
	}
	if _, err := c.call(ctx, "acceptTranscript", u); err != nil {
		t.Fatalf("identical request retry: %v", err)
	}
	before := snapshot()
	if err := play(ctx, c, []event{{0, 0, 500, "Injected", true}}, &bytes.Buffer{}); err == nil {
		t.Fatal("competing playback accepted")
	}
	if after := snapshot(); after != before {
		t.Fatalf("competing session changed: before=%+v after=%+v", before, after)
	}
}
func TestCancellationClosesOwnedSession(t *testing.T) {
	c := application(t)
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	if err := play(ctx, c, []event{{0, 0, 500, "Partial", false}, {1000, 0, 500, "Final", true}}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected cancellation")
	}
	if s := snapshot(); !s.Closed {
		t.Fatalf("owned session remains open: %+v", s)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestLostOpenResponseAbortsPendingSession(t *testing.T) {
	c := application(t)
	ctx := context.Background()
	if err := play(ctx, c, []event{{0, 0, 100, "Previous text.", true}}, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	before := snapshot()
	transport := c.http.Transport
	lost := false
	c.http.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		res, err := transport.RoundTrip(r)
		if err == nil && !lost && strings.HasSuffix(r.URL.Path, "/openSession") {
			lost = true
			res.Body.Close()
			return nil, errors.New("injected lost open response")
		}
		return res, err
	})
	if err := play(ctx, c, []event{{0, 0, 100, "Discard this.", true}}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "injected lost open response") {
		t.Fatalf("expected lost open response: %v", err)
	}
	if after := snapshot(); after != before {
		t.Fatalf("pending session changed prior display: before=%+v after=%+v", before, after)
	}
	if err := play(ctx, c, []event{{0, 0, 100, "Next session.", true}}, &bytes.Buffer{}); err != nil {
		t.Fatalf("pending session blocks next recording: %v", err)
	}
}

func TestFixtureSignalHelperProcess(t *testing.T) {
	if os.Getenv("VN_FIXTURE_SIGNAL_HELPER") != "1" {
		return
	}
	os.Args = []string{"transcript-fixture", "--origin", os.Getenv("VN_FIXTURE_SIGNAL_ORIGIN"), "--fixture", os.Getenv("VN_FIXTURE_SIGNAL_PATH")}
	main()
}

func TestSIGTERMCleansUpActiveSession(t *testing.T) {
	c := application(t)
	updates, unsubscribe, _ := transcript.Default.Watch()
	defer unsubscribe()
	fixturePath := filepath.Join(t.TempDir(), "signal-fixture.json")
	fixture := `[{"atMs":0,"startMs":0,"endMs":400,"text":"Waiting for cancellation.","isFinal":false},{"atMs":20000,"startMs":0,"endMs":400,"text":"Finished without cancellation.","isFinal":true}]`
	if err := os.WriteFile(fixturePath, []byte(fixture), 0600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestFixtureSignalHelperProcess$")
	cmd.Env = append(os.Environ(), "VN_FIXTURE_SIGNAL_HELPER=1", "VN_FIXTURE_SIGNAL_ORIGIN="+c.origin, "VN_FIXTURE_SIGNAL_PATH="+fixturePath)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	for {
		select {
		case s := <-updates:
			if s.Text != "Waiting for cancellation." {
				continue
			}
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				t.Fatal(err)
			}
		case err := <-done:
			t.Fatalf("fixture exited before signal: %v; %s", err, output.String())
		case <-deadline.C:
			t.Fatal("fixture did not submit its first transcript")
		}
		break
	}
	select {
	case err := <-done:
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 1 {
			t.Fatalf("expected graceful canceled exit: %v; %s", err, output.String())
		}
	case <-deadline.C:
		t.Fatal("fixture did not finish SIGTERM cleanup")
	}
	if s := snapshot(); !s.Closed || s.Text != "Waiting for cancellation." {
		t.Fatalf("SIGTERM left session open or changed visible text: %+v", s)
	}
	if err := play(context.Background(), c, []event{{0, 0, 100, "Next session.", true}}, &bytes.Buffer{}); err != nil {
		t.Fatalf("SIGTERM cleanup blocks next recording: %v", err)
	}
}
func TestLoopbackOrigin(t *testing.T) {
	for _, origin := range []string{"https://127.0.0.1:1", "http://localhost:1", "http://example.com:1", "http://127.0.0.1:1/path", "http://user@127.0.0.1:1", "http://127.0.0.1:1?x=y", "http://127.0.0.1"} {
		if validateOrigin(origin) == nil {
			t.Errorf("accepted %s", origin)
		}
	}
	for _, origin := range []string{"http://127.0.0.1:1234", "http://[::1]:1234"} {
		if err := validateOrigin(origin); err != nil {
			t.Errorf("%s: %v", origin, err)
		}
	}
}

func TestTranscriptRemoteRegistration(t *testing.T) {
	for _, name := range []string{"openSession", "activateSession", "acceptTranscript", "closeSession", "abortSession"} {
		id, err := transcriptRemoteID(name)
		if err != nil || id != "13p9jtx/"+name {
			t.Fatalf("%s ID = %q, error = %v", name, id, err)
		}
	}
	if _, err := transcriptRemoteID("missingSession"); err == nil {
		t.Fatal("missing generated remote accepted")
	}
}
