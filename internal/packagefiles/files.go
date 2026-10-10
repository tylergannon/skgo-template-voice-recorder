// Package packagefiles installs a prepared file plan without overwriting edits.
package packagefiles

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type File struct {
	Path string
	Data []byte
	Mode fs.FileMode
	// Replace permits replacement of exactly one known scaffold file.
	Replace []byte
}

// Apply checks every conflict before its first write, and reports completed
// files on cancellation or I/O failure. Identical files are left alone.
func Apply(ctx context.Context, root string, files []File) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	files = append([]File(nil), files...)
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	for i, f := range files {
		if i > 0 && f.Path == files[i-1].Path {
			return nil, fmt.Errorf("duplicate destination %s", f.Path)
		}
		if err := check(root, f); err != nil {
			return nil, err
		}
	}
	var changed []string
	for _, f := range files {
		if err := ctx.Err(); err != nil {
			return changed, err
		}
		// Recheck immediately before writing to avoid replacing edits made since
		// preflight. Applications should not run competing add-ons concurrently.
		if err := check(root, f); err != nil {
			return changed, err
		}
		path := filepath.Join(root, filepath.FromSlash(f.Path))
		old, readErr := os.ReadFile(path)
		if readErr == nil && bytes.Equal(old, f.Data) {
			continue
		}
		parent := filepath.Dir(path)
		if err := os.MkdirAll(parent, 0755); err != nil {
			return changed, err
		}
		tmp, err := os.CreateTemp(parent, ".skgo-recorder-*")
		if err != nil {
			return changed, err
		}
		tmpPath := tmp.Name()
		if _, err = tmp.Write(f.Data); err == nil {
			err = tmp.Chmod(f.Mode)
		}
		closeErr := tmp.Close()
		if err == nil {
			err = closeErr
		}
		if err == nil {
			err = os.Rename(tmpPath, path)
		}
		if err != nil {
			_ = os.Remove(tmpPath)
			return changed, err
		}
		changed = append(changed, f.Path)
	}
	return changed, nil
}

func check(root string, f File) error {
	if !fs.ValidPath(f.Path) || f.Path == "." || !filepath.IsLocal(filepath.FromSlash(f.Path)) {
		return fmt.Errorf("invalid destination %q", f.Path)
	}
	path := filepath.Join(root, filepath.FromSlash(f.Path))
	for part := path; part != root; part = filepath.Dir(part) {
		info, err := os.Lstat(part)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s uses a symlink; refusing installation", f.Path)
		}
	}
	old, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%s: %w", f.Path, err)
	}
	if bytes.Equal(old, f.Data) {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if info.Mode().Perm() != f.Mode.Perm() {
			return fmt.Errorf("%s has mode %04o; restore expected mode %04o before reapplying", f.Path, info.Mode().Perm(), f.Mode.Perm())
		}
		return nil
	}
	if f.Replace != nil && bytes.Equal(old, f.Replace) {
		return nil
	}
	return fmt.Errorf("%s conflicts with existing content; preserve your edits and reconcile this file before retrying", f.Path)
}
