package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/tylergannon/skgo/nativeapp"
)

// Preserve application-owned configuration, including fields from other
// add-ons. Conflicting recorder values are errors, never silent replacements.
func recorderConfig(data []byte) (string, []byte, error) {
	var c map[string]any
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := d.Decode(&c); err != nil {
		return "", nil, fmt.Errorf("native configuration: %w", err)
	}
	var trailing any
	if err := d.Decode(&trailing); err != io.EOF {
		return "", nil, fmt.Errorf("native configuration must contain exactly one JSON object")
	}
	if c == nil || c["version"] != json.Number("1") {
		return "", nil, fmt.Errorf("voice-notes requires native configuration version 1")
	}
	id, ok := c["bundleId"].(string)
	if !ok || id == "" {
		return "", nil, fmt.Errorf("native configuration needs bundleId")
	}
	wanted := map[string]any{
		"packages":  []any{map[string]any{"name": "RecorderTranscription", "path": "native/Transcription", "products": []any{"TranscriptSource", "AppleTranscription", "ScriptedTranscription"}}},
		"resources": []any{"native/Resources"},
		"swiftRemotes": []any{
			"src/lib/transcript.remote.ts#abortSession", "src/lib/transcript.remote.ts#acceptTranscript", "src/lib/transcript.remote.ts#activateSession", "src/lib/transcript.remote.ts#closeSession", "src/lib/transcript.remote.ts#openSession",
		},
		"deploymentTargets": map[string]any{"macos": "27.0", "iphone": "26.0", "simulator": "27.0"},
		"info": map[string]any{
			"NSMicrophoneUsageDescription":        "Use the microphone to transcribe speech into live notes.",
			"NSSpeechRecognitionUsageDescription": "Use Apple's on-device speech recognition to transcribe your notes.",
			"VNDefaultSource":                     "real",
		},
		"iosInfo": map[string]any{"UIBackgroundModes": []any{"audio"}},
		"presets": map[string]any{
			"real":             map[string]any{"bundleSuffix": "real", "info": map[string]any{"VNDefaultSource": "real"}},
			"synthetic":        syntheticPreset("synthetic"),
			"synthetic-faults": syntheticPreset("synthetic.faults"),
			"synthetic-burst": map[string]any{
				"bundleSuffix": "synthetic.burst", "goTags": []any{"synthetic"},
				"info":      map[string]any{"VNDefaultSource": "synthetic", "VNTranscriptFixture": "burst-transcript"},
				"resources": []any{"native/Tests/Fixtures/burst-transcript.json", "native/Resources/Licenses"},
			},
		},
	}
	if err := mergeConfig(c, wanted, "native/skgo-native.json"); err != nil {
		return "", nil, err
	}
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", nil, err
	}
	var known nativeapp.Config
	if err := json.Unmarshal(out, &known); err != nil {
		return "", nil, fmt.Errorf("native configuration: %w", err)
	}
	if err := known.Validate(); err != nil {
		return "", nil, err
	}
	return id, append(out, '\n'), nil
}

func syntheticPreset(suffix string) map[string]any {
	return map[string]any{"bundleSuffix": suffix, "goTags": []any{"synthetic"}, "info": map[string]any{"VNDefaultSource": "synthetic"}}
}

func mergeConfig(dst, src map[string]any, at string) error {
	for key, want := range src {
		path := at + "." + key
		have, exists := dst[key]
		if !exists {
			dst[key] = want
			continue
		}
		switch w := want.(type) {
		case map[string]any:
			h, ok := have.(map[string]any)
			if !ok {
				return fmt.Errorf("%s conflicts with recorder configuration", path)
			}
			if err := mergeConfig(h, w, path); err != nil {
				return err
			}
		case []any:
			// A preset selects one exact build mode and resource override. Combining
			// competing tags or fixtures would quietly change that mode's meaning.
			if strings.Contains(at, ".presets.") && (key == "goTags" || key == "resources") {
				if !reflect.DeepEqual(have, want) {
					return fmt.Errorf("%s conflicts with recorder configuration", path)
				}
				continue
			}
			h, ok := have.([]any)
			if !ok {
				return fmt.Errorf("%s must be an array", path)
			}
			for _, item := range w {
				found := false
				for _, prior := range h {
					if reflect.DeepEqual(prior, item) {
						found = true
						break
					}
					if key == "packages" {
						p, _ := prior.(map[string]any)
						n, _ := item.(map[string]any)
						if p["name"] == n["name"] || p["path"] == n["path"] {
							return fmt.Errorf("%s has a conflicting recorder package", path)
						}
					}
				}
				if !found {
					h = append(h, item)
				}
			}
			dst[key] = h
		default:
			if !reflect.DeepEqual(have, want) {
				return fmt.Errorf("%s conflicts with recorder configuration", path)
			}
		}
	}
	return nil
}
