package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"golang.org/x/mod/module"
	"html"
	"io/fs"
	"regexp"
	"strings"

	"github.com/tylergannon/skgo-template-voice-recorder/internal/packagefiles"
	"github.com/tylergannon/skgo/templateapi"
)

//go:generate go run ./stage
//go:embed all:payload allowlist.json
var payload embed.FS

const canonicalModule = "github.com/tylergannon/skgo-template-voice-recorder/recorder"

func descriptor(hostVersion string) templateapi.Descriptor {
	return templateapi.Descriptor{
		ID: "voice-notes", Version: "0.1.0", SkgoVersion: hostVersion,
		Addons: []templateapi.Addon{{Name: "voice-notes", Summary: "Continuous native transcription with controls and a live transcript"}},
		Templates: []templateapi.Template{{Name: "voice-notes", Summary: "A native voice recorder built on the standard SKGo project",
			Options: []templateapi.Option{{Name: "bundle-id", Help: "Base application bundle identifier", Required: true}},
			Steps:   []templateapi.Step{{Addon: "native-app", Bindings: map[string]string{"bundle-id": "bundle-id"}}, {Addon: "voice-notes"}},
		}},
	}
}

func recorderFiles(project templateapi.Project, bundleID string) ([]packagefiles.File, error) {
	if err := module.CheckImportPath(project.Module); err != nil {
		return nil, err
	}
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`).MatchString(project.Name) {
		return nil, fmt.Errorf("invalid destination name %q", project.Name)
	}
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]*(\.[A-Za-z0-9][A-Za-z0-9-]*)+$`).MatchString(bundleID) {
		return nil, fmt.Errorf("invalid bundle ID %q", bundleID)
	}
	manifest, err := payload.ReadFile("allowlist.json")
	if err != nil {
		return nil, err
	}
	var entries []struct {
		Destination string
		Executable  bool
	}
	if err := json.Unmarshal(manifest, &entries); err != nil {
		return nil, err
	}
	modes := map[string]fs.FileMode{}
	for _, e := range entries {
		mode := fs.FileMode(0644)
		if e.Executable {
			mode = 0755
		}
		modes[e.Destination] = mode
	}
	var files []packagefiles.File
	err = fs.WalkDir(payload, "payload", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := payload.ReadFile(path)
		if err != nil {
			return err
		}
		path = strings.TrimSuffix(strings.TrimPrefix(path, "payload/"), ".txt")
		// Only explicit application identity markers and quoted Go import paths
		// are substituted. Recorder bridge/product identifiers remain stable.
		if strings.HasSuffix(path, ".go") {
			data = bytes.ReplaceAll(data, []byte(`"`+canonicalModule+`/`), []byte(`"`+project.Module+`/`))
			data = bytes.ReplaceAll(data, []byte(`"`+canonicalModule+`"`), []byte(`"`+project.Module+`"`))
		}
		if path == "web/src/routes/+page.svelte" {
			data = bytes.ReplaceAll(data, []byte("__SKGO_PROJECT_NAME__"), []byte(html.EscapeString(project.Name)))
		}
		if path == "native/Transcription/ProbeInfo.plist" {
			data = bytes.ReplaceAll(data, []byte("dev.skgo.recorder.transcription-probe"), []byte(html.EscapeString(bundleID+".transcription-probe")))
			data = bytes.ReplaceAll(data, []byte("Recorder Transcription Probe"), []byte(html.EscapeString(project.Name+" Transcription Probe")))
		}
		if path == "scripts/check-recorder.sh" {
			data = bytes.ReplaceAll(data, []byte(`cd "$repo_root/recorder"`), []byte(`cd "$repo_root"`))
		}
		if bytes.Contains(data, []byte("__SKGO_")) {
			return fmt.Errorf("%s contains an unresolved identity marker", path)
		}
		mode, listed := modes[path]
		if !listed {
			return fmt.Errorf("unlisted payload file %s", path)
		}
		files = append(files, packagefiles.File{Path: path, Data: data, Mode: mode})
		return nil
	})
	return files, err
}
