package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/tylergannon/skgo/nativeapp"
	"github.com/tylergannon/skgo/templateapi"
)

const applyTestWelcome = "<h1>Welcome to SvelteKit</h1>\n<p>Visit <a href=\"https://svelte.dev/docs/kit\">svelte.dev/docs/kit</a> to read the documentation</p>\n"

func applyTestWrite(t *testing.T, root, path string, data []byte) {
	t.Helper()
	file := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func applyTestProject(t *testing.T, config map[string]any) templateapi.Project {
	t.Helper()
	root := t.TempDir()
	if config == nil {
		config = map[string]any{"version": 1, "bundleId": "com.example.fieldnotes"}
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	applyTestWrite(t, root, "native/skgo-native.json", data)
	applyTestWrite(t, root, "native/Sources/ApplicationView.swift", []byte(nativeapp.DefaultApplicationView))
	applyTestWrite(t, root, "native/Sources/NativeHost.swift", []byte("// Generic host owned by SKGo; application add-ons preserve it.\n"))
	applyTestWrite(t, root, "web/src/routes/+page.svelte", []byte(applyTestWelcome))
	return templateapi.Project{Root: root, Web: filepath.Join(root, "web"), Name: "FieldNotes", Module: "example.com/field-notes"}
}

func applyTestTree(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path != root {
				relative, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				files[filepath.ToSlash(relative)+"/"] = "directory"
			}
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(relative)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func applyTestRequest(project templateapi.Project) templateapi.Request {
	return templateapi.Request{Addon: "voice-notes", Project: project}
}

func applyTestConfig(t *testing.T, project templateapi.Project) nativeapp.Config {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(project.Root, "native/skgo-native.json"))
	if err != nil {
		t.Fatal(err)
	}
	var config nativeapp.Config
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	return config
}

func TestApplyInstallsRecorderWithConsumerIdentity(t *testing.T) {
	project := applyTestProject(t, nil)
	before := applyTestTree(t, project.Root)
	result, err := SKGoPluginV1().Apply(context.Background(), applyTestRequest(project))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changed) == 0 {
		t.Fatal("successful installation reported no changes")
	}
	for _, path := range result.Changed {
		if !fs.ValidPath(path) || path == "." {
			t.Errorf("changed path must be project-relative: %q", path)
		}
	}
	files := applyTestTree(t, project.Root)
	if files["native/Sources/NativeHost.swift"] != before["native/Sources/NativeHost.swift"] {
		t.Fatal("recorder altered the generic host")
	}
	for path, want := range map[string]string{
		"web/src/routes/+page.svelte":           "<title>FieldNotes</title>",
		"web/src/lib/transcript.remote.go":      `"example.com/field-notes/internal/transcript"`,
		"native/Transcription/ProbeInfo.plist":  "com.example.fieldnotes.transcription-probe",
		"native/Sources/ApplicationView.swift":  "RecorderScreen(controller:",
		"native/Sources/RecorderControls.swift": "VStack(spacing: 2)",
	} {
		if !strings.Contains(files[path], want) {
			t.Errorf("%s missing %q", path, want)
		}
	}
	if !strings.Contains(files["native/Transcription/ProbeInfo.plist"], "FieldNotes Transcription Probe") {
		t.Fatal("probe display name not derived from consumer")
	}
	config := applyTestConfig(t, project)
	if config.Version != 1 || config.BundleID != "com.example.fieldnotes" || config.Info["VNDefaultSource"] != "real" {
		t.Fatalf("consumer identity/default source lost: %+v", config)
	}
	for _, key := range []string{"NSMicrophoneUsageDescription", "NSSpeechRecognitionUsageDescription"} {
		if value, ok := config.Info[key].(string); !ok || strings.TrimSpace(value) == "" {
			t.Errorf("missing meaningful permission description %s", key)
		}
	}
	if !applyTestContains(config.IOSInfo["UIBackgroundModes"], "audio") {
		t.Fatal("iOS background audio capability missing")
	}
	foundPackage := false
	for _, pkg := range config.Packages {
		if pkg.Path == "native/Transcription" {
			foundPackage = true
			for _, product := range []string{"TranscriptSource", "AppleTranscription", "ScriptedTranscription"} {
				if !applyTestContains(pkg.Products, product) {
					t.Errorf("transcription product %s not linked", product)
				}
			}
		}
	}
	if !foundPackage {
		t.Fatal("independent transcription package not linked")
	}
	for name, suffix := range map[string]string{"real": "real", "synthetic": "synthetic", "synthetic-burst": "synthetic.burst", "synthetic-faults": "synthetic.faults"} {
		preset, ok := config.Presets[name]
		if !ok || preset.BundleSuffix != suffix {
			t.Errorf("missing preset identity %s: %+v", name, preset)
		}
	}
	for _, name := range []string{"synthetic", "synthetic-faults", "synthetic-burst"} {
		preset := config.Presets[name]
		if preset.Info["VNDefaultSource"] != "synthetic" || !reflect.DeepEqual(preset.GoTags, []string{"synthetic"}) {
			t.Errorf("%s does not select the synthetic source and hooks: %+v", name, preset)
		}
	}
	if !applyTestContains(config.Resources, "native/Resources") || !applyTestContains(config.Presets["synthetic-burst"].Resources, "native/Tests/Fixtures/burst-transcript.json") {
		t.Fatal("synthetic fixture resources missing")
	}
	if config.Presets["synthetic-burst"].Info["VNTranscriptFixture"] != "burst-transcript" {
		t.Fatal("burst preset does not select its named fixture")
	}
	second, err := SKGoPluginV1().Apply(context.Background(), applyTestRequest(project))
	if err != nil || len(second.Changed) != 0 {
		t.Fatalf("same configuration reapply: changes=%v error=%v", second.Changed, err)
	}
	if after := applyTestTree(t, project.Root); !reflect.DeepEqual(after, files) {
		t.Fatal("no-op reapply changed destination files")
	}
}

func applyTestContains(value any, want string) bool {
	switch values := value.(type) {
	case []string:
		for _, item := range values {
			if item == want {
				return true
			}
		}
	case []any:
		for _, item := range values {
			if item == want {
				return true
			}
		}
	}
	return false
}

func TestApplyRejectsEditedApplicationBeforeWriting(t *testing.T) {
	for _, path := range []string{"web/src/routes/+page.svelte", "native/Sources/ApplicationView.swift"} {
		t.Run(path, func(t *testing.T) {
			project := applyTestProject(t, nil)
			applyTestWrite(t, project.Root, path, []byte("User-owned custom application content.\n"))
			before := applyTestTree(t, project.Root)
			result, err := SKGoPluginV1().Apply(context.Background(), applyTestRequest(project))
			if err == nil || len(result.Changed) != 0 {
				t.Fatalf("conflict not rejected before writes: result=%+v error=%v", result, err)
			}
			if !reflect.DeepEqual(before, applyTestTree(t, project.Root)) {
				t.Fatal("conflict changed project files")
			}
		})
	}
}

func TestApplyRequiresNativeFoundationAndHonorsCancellation(t *testing.T) {
	for _, missing := range []string{"native/skgo-native.json", "native/Sources/ApplicationView.swift", "native/Sources/NativeHost.swift"} {
		t.Run(missing, func(t *testing.T) {
			project := applyTestProject(t, nil)
			if err := os.Remove(filepath.Join(project.Root, filepath.FromSlash(missing))); err != nil {
				t.Fatal(err)
			}
			before := applyTestTree(t, project.Root)
			result, err := SKGoPluginV1().Apply(context.Background(), applyTestRequest(project))
			if err == nil || len(result.Changed) != 0 || !reflect.DeepEqual(before, applyTestTree(t, project.Root)) {
				t.Fatalf("missing foundation was not rejected without writes: result=%+v error=%v", result, err)
			}
		})
	}
	t.Run("cancelled", func(t *testing.T) {
		project := applyTestProject(t, nil)
		before := applyTestTree(t, project.Root)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := SKGoPluginV1().Apply(ctx, applyTestRequest(project))
		if !errors.Is(err, context.Canceled) || len(result.Changed) != 0 || !reflect.DeepEqual(before, applyTestTree(t, project.Root)) {
			t.Fatalf("cancelled application wrote files or lost cancellation: result=%+v error=%v", result, err)
		}
	})
}

func TestApplyPreservesUnrelatedConfigAndSource(t *testing.T) {
	project := applyTestProject(t, map[string]any{
		"version": 1, "bundleId": "com.example.fieldnotes",
		"sources": []string{"native/Extra.swift"}, "resources": []string{"native/Extra.txt"},
		"packages":           []map[string]any{{"name": "ExtraNative", "path": "native/Extra", "products": []string{"Extra"}}},
		"info":               map[string]any{"CustomSetting": "keep"},
		"iosInfo":            map[string]any{"UIBackgroundModes": []string{"fetch"}},
		"macosInfo":          map[string]any{"LSUIElement": true},
		"presets":            map[string]any{"custom": map[string]any{"bundleSuffix": "custom", "goTags": []string{"custom"}, "resources": []string{}}},
		"unrelatedExtension": map[string]any{"keep": true},
	})
	applyTestWrite(t, project.Root, "native/Extra.swift", []byte("// Keep my unrelated native extension.\n"))
	result, err := SKGoPluginV1().Apply(context.Background(), applyTestRequest(project))
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changed) == 0 {
		t.Fatal("recorder not installed")
	}
	config := applyTestConfig(t, project)
	if !applyTestContains(config.Sources, "native/Extra.swift") || !applyTestContains(config.Resources, "native/Extra.txt") || config.Info["CustomSetting"] != "keep" || config.MacOSInfo["LSUIElement"] != true {
		t.Fatalf("unrelated fields lost: %+v", config)
	}
	if !applyTestContains(config.IOSInfo["UIBackgroundModes"], "fetch") || !applyTestContains(config.IOSInfo["UIBackgroundModes"], "audio") {
		t.Fatal("background audio replaced unrelated background modes")
	}
	if custom, ok := config.Presets["custom"]; !ok || custom.BundleSuffix != "custom" || !applyTestContains(custom.GoTags, "custom") {
		t.Fatal("unrelated preset lost")
	}
	foundExtra := false
	for _, pkg := range config.Packages {
		if pkg.Name == "ExtraNative" && pkg.Path == "native/Extra" && reflect.DeepEqual(pkg.Products, []string{"Extra"}) {
			foundExtra = true
		}
	}
	if !foundExtra {
		t.Fatal("unrelated package lost")
	}
	files := applyTestTree(t, project.Root)
	if files["native/Extra.swift"] != "// Keep my unrelated native extension.\n" {
		t.Fatal("unrelated handwritten source altered")
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(files["native/skgo-native.json"]), &raw); err != nil {
		t.Fatal(err)
	}
	var extension map[string]bool
	if err := json.Unmarshal(raw["unrelatedExtension"], &extension); err != nil || !extension["keep"] {
		t.Fatal("unknown config extension lost")
	}
}

func TestApplyRejectsContradictoryRecorderConfigBeforeWriting(t *testing.T) {
	for name, patch := range map[string]map[string]any{
		"source metadata":  {"info": map[string]any{"VNDefaultSource": "synthetic"}},
		"synthetic preset": {"presets": map[string]any{"synthetic": map[string]any{"bundleSuffix": "different", "resources": []string{}}}},
		"synthetic tags":   {"presets": map[string]any{"synthetic": map[string]any{"bundleSuffix": "synthetic", "goTags": []string{"real"}}}},
		"burst fixture":    {"presets": map[string]any{"synthetic-burst": map[string]any{"bundleSuffix": "synthetic.burst", "info": map[string]any{"VNTranscriptFixture": "different-fixture"}}}},
		"burst resources":  {"presets": map[string]any{"synthetic-burst": map[string]any{"bundleSuffix": "synthetic.burst", "resources": []string{"native/UserFixture.json"}}}},
		"package":          {"packages": []map[string]any{{"name": "RecorderTranscription", "path": "native/Transcription", "products": []string{"IncompatibleProduct"}}}},
	} {
		t.Run(name, func(t *testing.T) {
			config := map[string]any{"version": 1, "bundleId": "com.example.fieldnotes"}
			for key, value := range patch {
				config[key] = value
			}
			project := applyTestProject(t, config)
			before := applyTestTree(t, project.Root)
			result, err := SKGoPluginV1().Apply(context.Background(), applyTestRequest(project))
			if err == nil || len(result.Changed) != 0 || !reflect.DeepEqual(before, applyTestTree(t, project.Root)) {
				t.Fatalf("contradictory config was not rejected before writes: result=%+v error=%v", result, err)
			}
		})
	}
}
