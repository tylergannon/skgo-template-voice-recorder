package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/tylergannon/skgo/templateapi"
)

func TestPackagedSourceUsesFreshConsumerIdentity(t *testing.T) {
	files, err := recorderFiles(templateapi.Project{Module: "example.com/field-notes", Name: "FieldNotes"}, "com.example.fieldnotes")
	if err != nil {
		t.Fatal(err)
	}
	byPath := map[string]string{}
	for _, f := range files {
		byPath[f.Path] = string(f.Data)
		if bytes.Contains(f.Data, []byte(canonicalModule)) && strings.HasSuffix(f.Path, ".go") {
			t.Errorf("source import leaked into %s", f.Path)
		}
		if strings.Contains(f.Path, ".build") || strings.Contains(f.Path, "ephemeral") || strings.HasSuffix(f.Path, ".so") {
			t.Errorf("build artifact packaged: %s", f.Path)
		}
	}
	if !strings.Contains(byPath["web/src/lib/transcript.remote.go"], `"example.com/field-notes/internal/transcript"`) {
		t.Fatal("destination Go identity missing")
	}
	if !strings.Contains(byPath["web/src/routes/+page.svelte"], "<title>FieldNotes</title>") {
		t.Fatal("destination display identity missing")
	}
	plist := byPath["native/Transcription/ProbeInfo.plist"]
	if !strings.Contains(plist, "com.example.fieldnotes.transcription-probe") || !strings.Contains(plist, "FieldNotes Transcription Probe") {
		t.Fatal("probe identity not derived")
	}
	if !strings.Contains(byPath["native/Sources/RecorderControls.swift"], "VStack(spacing: 2)") {
		t.Fatal("stacked meter not retained")
	}
	if !strings.Contains(byPath["scripts/check-recorder.sh"], `cd "$repo_root"`) || strings.Contains(byPath["scripts/check-recorder.sh"], `cd "$repo_root/recorder"`) {
		t.Fatal("consumer check script uses source layout")
	}
	for _, generic := range []string{"native/Sources/NativeHost.swift", "native/host/main.go", "native/Generated.swift", "scripts/build-native.sh", "go.mod"} {
		if _, ok := byPath[generic]; ok {
			t.Errorf("duplicated generic foundation %s", generic)
		}
	}
}

func TestPayloadMatchesReviewedAllowlist(t *testing.T) {
	data, err := os.ReadFile("allowlist.json")
	if err != nil {
		t.Fatal(err)
	}
	var entries []struct{ Source, Destination string }
	if err = json.Unmarshal(data, &entries); err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		source, err := os.ReadFile("../" + entry.Source)
		if err != nil {
			t.Fatal(err)
		}
		staged, err := payload.ReadFile("payload/" + entry.Destination + ".txt")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(source, staged) {
			t.Errorf("stale payload: %s; rerun stage", entry.Source)
		}
	}
}

func TestRecipeUsesCoreNativeFoundation(t *testing.T) {
	d := descriptor("actual-host-version")
	if d.SkgoVersion != "actual-host-version" || len(d.Addons) != 1 || d.Addons[0].Name != "voice-notes" {
		t.Fatal("descriptor ownership mismatch")
	}
	steps := d.Templates[0].Steps
	if len(steps) != 2 || steps[0].Addon != "native-app" || steps[1].Addon != "voice-notes" || steps[0].Bindings["bundle-id"] != "bundle-id" {
		t.Fatal("recipe diverged from manual composition")
	}
}

func TestIdentityValidationRejectsMarkupAndInvalidBundle(t *testing.T) {
	for _, p := range []templateapi.Project{{Module: "example.com/notes", Name: "{injection}"}, {Module: "not a module", Name: "Notes"}} {
		if _, err := recorderFiles(p, "com.example.notes"); err == nil {
			t.Fatalf("accepted invalid identity %+v", p)
		}
	}
	if _, err := recorderFiles(templateapi.Project{Module: "example.com/notes", Name: "Notes"}, "../escape"); err == nil {
		t.Fatal("accepted invalid bundle")
	}
}
