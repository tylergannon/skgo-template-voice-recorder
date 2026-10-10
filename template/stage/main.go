// stage materializes the reviewed canonical sources for Go embedding.
package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type entry struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Executable  bool   `json:"executable,omitempty"`
}

func main() {
	if err := stage(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func stage() error {
	// go generate executes in template/, whereas direct invocation is from root.
	root := "."
	if _, err := os.Stat("allowlist.json"); err == nil {
		root = ".."
	}
	data, err := os.ReadFile(filepath.Join(root, "template", "allowlist.json"))
	if err != nil {
		return err
	}
	var entries []entry
	if err = json.Unmarshal(data, &entries); err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("empty recorder allowlist")
	}
	tmp, err := os.MkdirTemp(filepath.Join(root, "template"), "payload-stage-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	seen := map[string]bool{}
	for _, e := range entries {
		if !fs.ValidPath(e.Source) || e.Source == "." || !fs.ValidPath(e.Destination) || e.Destination == "." || seen[e.Destination] {
			return fmt.Errorf("invalid or duplicate allowlist entry: %+v", e)
		}
		seen[e.Destination] = true
		source := filepath.Join(root, filepath.FromSlash(e.Source))
		info, err := os.Lstat(source)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("nonregular source %s", e.Source)
		}
		content, err := os.ReadFile(source)
		if err != nil {
			return err
		}
		dest := filepath.Join(tmp, filepath.FromSlash(e.Destination+".txt"))
		if err = os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		mode := os.FileMode(0644)
		if e.Executable {
			mode = 0755
		}
		if err = os.WriteFile(dest, content, mode); err != nil {
			return err
		}
	}
	final := filepath.Join(root, "template", "payload")
	if err = os.RemoveAll(final); err != nil {
		return err
	}
	return os.Rename(tmp, final)
}
