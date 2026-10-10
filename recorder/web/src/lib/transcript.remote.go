package transcriptremote

import (
	"context"
	"errors"
	"github.com/tylergannon/skgo"
	"github.com/tylergannon/skgo-template-voice-recorder/recorder/internal/skgo/params"
	"github.com/tylergannon/skgo-template-voice-recorder/recorder/internal/transcript"
)

func openSession(ctx context.Context, event params.RequestEvent, r transcript.OpenRequest) (transcript.Ack, error) {
	return transcript.Default.Open(r)
}
func acceptTranscript(ctx context.Context, event params.RequestEvent, r transcript.Update) (transcript.Ack, error) {
	if err := transcript.BeforeAccept(ctx); err != nil {
		return transcript.Ack{}, err
	}
	return transcript.Default.Accept(r)
}
func closeSession(ctx context.Context, event params.RequestEvent, r transcript.CloseRequest) (transcript.Ack, error) {
	return transcript.Default.Close(r)
}
func watchTranscript(ctx context.Context, yield func(transcript.Snapshot) error) error {
	ch, cancel, now := transcript.Default.Watch()
	defer cancel()
	failure, stopFailure := transcript.StreamFailure()
	defer stopFailure()
	if err := transcript.BeforeYield(ctx); err != nil {
		return err
	}
	if err := yield(now); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-failure:
			return errors.New("Synthetic test disconnected the transcript stream")
		case next := <-ch:
			if err := transcript.BeforeYield(ctx); err != nil {
				return err
			}
			if err := yield(next); err != nil {
				return err
			}
		}
	}
}

var _ = skgo.Command(openSession)
var _ = skgo.Command(acceptTranscript)
var _ = skgo.Command(closeSession)
var _ = skgo.LiveQuery(watchTranscript)

func activateSession(ctx context.Context, event params.RequestEvent, r transcript.SessionRequest) (transcript.Ack, error) {
	return transcript.Default.Activate(r)
}
func abortSession(ctx context.Context, event params.RequestEvent, r transcript.SessionRequest) (transcript.Ack, error) {
	return transcript.Default.Abort(r)
}

var _ = skgo.Command(activateSession)
var _ = skgo.Command(abortSession)
