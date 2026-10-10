// transcript-fixture feeds synthetic text through the running application's generated remotes.
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	devalue "github.com/tylergannon/devalue/v5"
	generated "github.com/tylergannon/skgo-template-voice-recorder/recorder/internal/skgo"
	"github.com/tylergannon/skgo-template-voice-recorder/recorder/internal/transcript"
)

type event struct {
	AtMs    int64  `json:"atMs"`
	StartMs int64  `json:"startMs"`
	EndMs   int64  `json:"endMs"`
	Text    string `json:"text"`
	IsFinal bool   `json:"isFinal"`
}

var defaultFixture = []event{
	{0, 0, 400, "This was", false},
	{250, 500, 800, " a test", false},
	{500, 0, 450, "This is", false},
	{750, 0, 450, "This is", true},
	{1000, 500, 900, " synthetic text.", false},
	{1250, 500, 900, " synthetic text.", true},
}

func validateOrigin(origin string) error {
	u, err := url.Parse(origin)
	if err != nil {
		return err
	}
	ip := net.ParseIP(u.Hostname())
	if u.Scheme != "http" || ip == nil || !ip.IsLoopback() || u.Port() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("--origin must be an explicit HTTP loopback IP origin with port, such as http://127.0.0.1:12345")
	}
	return nil
}

type client struct {
	origin string
	http   *http.Client
}

func (c client) call(ctx context.Context, name string, value any) (transcript.Ack, error) {
	var tree any
	var err error
	switch v := value.(type) {
	case transcript.OpenRequest:
		tree, err = generated.EncodeOpenRequest(v)
	case transcript.SessionRequest:
		tree, err = generated.EncodeSessionRequest(v)
	case transcript.Update:
		tree, err = generated.EncodeUpdate(v)
	case transcript.CloseRequest:
		tree, err = generated.EncodeCloseRequest(v)
	default:
		return transcript.Ack{}, errors.New("unsupported command argument")
	}
	if err != nil {
		return transcript.Ack{}, err
	}
	wire, err := devalue.Stringify(tree)
	if err != nil {
		return transcript.Ack{}, err
	}
	body, err := json.Marshal(struct {
		Payload   string   `json:"payload"`
		Refreshes []string `json:"refreshes"`
	}{base64.RawURLEncoding.EncodeToString([]byte(wire)), []string{}})
	if err != nil {
		return transcript.Ack{}, err
	}
	id, err := transcriptRemoteID(name)
	if err != nil {
		return transcript.Ack{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.origin+"/_app/remote/"+id, bytes.NewReader(body))
	if err != nil {
		return transcript.Ack{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", c.origin)
	res, err := c.http.Do(req)
	if err != nil {
		return transcript.Ack{}, err
	}
	defer res.Body.Close()
	var envelope struct {
		Type  string `json:"type"`
		Data  string `json:"data"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&envelope); err != nil {
		return transcript.Ack{}, err
	}
	if res.StatusCode != 200 || envelope.Type != "result" {
		return transcript.Ack{}, fmt.Errorf("%s rejected (%d): %s", name, res.StatusCode, envelope.Error.Message)
	}
	parsed, err := devalue.Parse(envelope.Data, nil)
	if err != nil {
		return transcript.Ack{}, err
	}
	obj, ok := parsed.(*devalue.Object)
	if !ok {
		return transcript.Ack{}, errors.New("invalid command result")
	}
	result, ok := obj.Get("_")
	if !ok {
		return transcript.Ack{}, errors.New("missing command result")
	}
	return generated.DecodeAck(result)
}

// Resolve the protocol ID from the actual generated registration. SvelteKit
// hashes the Vite-root-relative module path, independently of Go module/name.
func transcriptRemoteID(name string) (string, error) {
	for _, remote := range generated.Remotes() {
		if remote.Module() == "src/lib/transcript.remote.ts" && remote.Name() == name {
			return remote.ID(), nil
		}
	}
	return "", fmt.Errorf("generated transcript remote %q is unavailable", name)
}

func play(ctx context.Context, c client, events []event, out io.Writer) (err error) {
	token := make([]byte, 16)
	if _, err = rand.Read(token); err != nil {
		return err
	}
	id := "fixture-" + hex.EncodeToString(token)
	closed := false
	defer func() {
		// Open may have reached the store even when its response was lost. Abort
		// only this random ID; the store also remembers absent IDs to reject a
		// delayed open, without changing another session's display.
		if !closed {
			cleanup, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			ack, cleanupErr := c.call(cleanup, "abortSession", transcript.SessionRequest{SessionID: id})
			if cleanupErr == nil && (ack.SessionID != id || !ack.Closed) {
				cleanupErr = errors.New("unexpected abort acknowledgement")
			}
			if cleanupErr != nil {
				err = errors.Join(err, fmt.Errorf("session cleanup: %w", cleanupErr))
			}
		}
	}()
	ack, err := c.call(ctx, "openSession", transcript.OpenRequest{SessionID: id, Source: "synthetic"})
	if err != nil {
		return err
	}
	if ack.SessionID != id || ack.Sequence != 0 || ack.Closed {
		return errors.New("unexpected open acknowledgement")
	}
	ack, err = c.call(ctx, "activateSession", transcript.SessionRequest{SessionID: id})
	if err != nil {
		return err
	}
	if ack.SessionID != id || ack.Sequence != 0 || ack.Closed {
		return errors.New("unexpected activation acknowledgement")
	}
	fmt.Fprintln(out, "Synthetic session active:", id)
	start := time.Now()
	boundary := int64(0)
	tail := []transcript.Chunk{}
	sequence := int64(0)
	finalText := ""
	for _, e := range events {
		timer := time.NewTimer(time.Until(start.Add(time.Duration(e.AtMs) * time.Millisecond)))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
		if e.StartMs < boundary || e.EndMs < e.StartMs || e.EndMs > transcript.MaxInteger {
			return errors.New("invalid fixture audio range")
		}
		chunk := transcript.Chunk{StartMs: e.StartMs, EndMs: e.EndMs, Text: e.Text}
		kept := tail[:0]
		for _, p := range tail {
			if !(p.EndMs > chunk.StartMs && p.StartMs < chunk.EndMs || p.StartMs == chunk.StartMs) {
				kept = append(kept, p)
			}
		}
		tail = kept
		finals := []transcript.Chunk{}
		if e.IsFinal {
			boundary = e.EndMs
			kept = tail[:0]
			for _, p := range tail {
				if p.StartMs >= boundary {
					kept = append(kept, p)
				}
			}
			tail = kept
			finals = append(finals, chunk)
			finalText += e.Text
		} else {
			tail = append(tail, chunk)
			sort.Slice(tail, func(i, j int) bool { return tail[i].StartMs < tail[j].StartMs })
		}
		sequence++
		ack, err = c.call(ctx, "acceptTranscript", transcript.Update{SessionID: id, Sequence: sequence, Finalized: finals, Provisional: tail})
		if err != nil {
			return err
		}
		if ack.SessionID != id || ack.Sequence != sequence || ack.Closed {
			return errors.New("unexpected transcript acknowledgement")
		}
		fmt.Fprintf(out, "Accepted %d (final=%t)\n", sequence, e.IsFinal)
	}
	if len(tail) != 0 {
		return errors.New("fixture ends with provisional text; add final events")
	}
	ack, err = c.call(ctx, "closeSession", transcript.CloseRequest{SessionID: id, Sequence: sequence})
	if err != nil {
		return err
	}
	if ack.SessionID != id || ack.Sequence != sequence || !ack.Closed {
		return errors.New("unexpected close acknowledgement")
	}
	closed = true
	fmt.Fprintf(out, "Synthetic session closed; acknowledged final text: %s\n", finalText)
	return nil
}

func run(ctx context.Context, args []string, out io.Writer) error {
	flags := flag.NewFlagSet("transcript-fixture", flag.ContinueOnError)
	flags.SetOutput(out)
	origin := flags.String("origin", "", "running app's explicit loopback origin")
	path := flags.String("fixture", "", "JSON array of atMs/startMs/endMs/text/isFinal events")
	timeout := flags.Duration("timeout", 30*time.Second, "overall limit (maximum 30s, including cleanup)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional arguments")
	}
	if err := validateOrigin(*origin); err != nil {
		return err
	}
	if *timeout <= 2*time.Second || *timeout > 30*time.Second {
		return errors.New("--timeout must be greater than 2s and at most 30s")
	}
	events := defaultFixture
	if *path != "" {
		f, err := os.Open(*path)
		if err != nil {
			return err
		}
		defer f.Close()
		decoder := json.NewDecoder(io.LimitReader(f, transcript.MaxSessionBytes+1))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&events); err != nil {
			return err
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			return errors.New("fixture must contain exactly one JSON array")
		}
	}
	if len(events) == 0 {
		return errors.New("fixture is empty")
	}
	previous := int64(0)
	for _, e := range events {
		if e.AtMs < previous || e.AtMs < 0 || e.AtMs > int64((*timeout-2*time.Second)/time.Millisecond) {
			return errors.New("fixture timing exceeds playback budget or is out of order")
		}
		previous = e.AtMs
	}
	playback, cancel := context.WithTimeout(ctx, *timeout-2*time.Second)
	defer cancel()
	return play(playback, client{*origin, &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("remote redirects are not allowed") }}}, events, out)
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
