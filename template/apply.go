package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tylergannon/skgo-template-voice-recorder/internal/packagefiles"
	"github.com/tylergannon/skgo/nativeapp"
	"github.com/tylergannon/skgo/templateapi"
)

// The core source-build helper stamps the invoking host's version.
var skgoVersion = "development"

type recorderPlugin struct{}

func SKGoPluginV1() templateapi.Plugin                  { return recorderPlugin{} }
func (recorderPlugin) Describe() templateapi.Descriptor { return descriptor(skgoVersion) }

// Only the ordinary minimal SvelteKit page is replaceable. Edited pages must
// be reconciled by their owner, just like an edited native application slot.
const defaultPage = `<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
`

func (recorderPlugin) Apply(ctx context.Context, req templateapi.Request) (result templateapi.Result, err error) {
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if req.Addon != "voice-notes" {
		return result, fmt.Errorf("unknown recorder add-on %q", req.Addon)
	}
	if len(req.Options) != 0 {
		return result, fmt.Errorf("voice-notes takes no options; set bundle-id on native-app")
	}
	if !filepath.IsAbs(req.Project.Root) || filepath.Clean(req.Project.Web) != filepath.Join(req.Project.Root, "web") {
		return result, fmt.Errorf("voice-notes requires the ordinary SKGo project with web/ at its root")
	}
	oldConfig, err := os.ReadFile(filepath.Join(req.Project.Root, "native/skgo-native.json"))
	if err != nil {
		return result, fmt.Errorf("apply native-app before voice-notes: %w", err)
	}
	for _, name := range []string{"native/Sources/NativeHost.swift", "native/Sources/ApplicationView.swift"} {
		if _, err := os.Stat(filepath.Join(req.Project.Root, name)); err != nil {
			return result, fmt.Errorf("native-app foundation is incomplete (%s): %w", name, err)
		}
	}
	bundleID, config, err := recorderConfig(oldConfig)
	if err != nil {
		return result, err
	}
	files, err := recorderFiles(req.Project, bundleID)
	if err != nil {
		return result, err
	}
	for i := range files {
		switch files[i].Path {
		case "web/src/routes/+page.svelte":
			files[i].Replace = []byte(defaultPage)
		case "native/Sources/ApplicationView.swift":
			files[i].Replace = []byte(nativeapp.DefaultApplicationView)
		}
	}
	files = append(files, packagefiles.File{Path: "native/skgo-native.json", Data: config, Mode: 0644, Replace: oldConfig})
	result.Changed, err = packagefiles.Apply(ctx, req.Project.Root, files)
	if err != nil {
		return result, err
	}
	result.Next = []string{
		"Build the native recorder with go tool skgo native build --platform macos --preset real (or iphone / simulator).",
		"For microphone-free playback, select --preset synthetic; synthetic does not establish microphone capture or Apple recognition.",
		"Browser mode serves transcripts without capture. Start the web server and use scripts/feed-transcript.sh --help for external synthetic delivery.",
	}
	return result, nil
}
