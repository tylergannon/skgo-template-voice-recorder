package packagefiles

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestConflictPreventsEarlierNewFiles(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "z.txt"), []byte("user edit"), 0644)
	changed, err := Apply(context.Background(), root, []File{{Path: "a/new.txt", Data: []byte("new"), Mode: 0644}, {Path: "z.txt", Data: []byte("generated"), Mode: 0644}})
	if err == nil || len(changed) != 0 {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if _, err := os.Stat(filepath.Join(root, "a")); !os.IsNotExist(err) {
		t.Fatalf("preflight mutated project: %v", err)
	}
}

func TestReapplyPreservesUnrelatedContentAndRejectsEditedOutput(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "page"), []byte("known scaffold"), 0644)
	os.WriteFile(filepath.Join(root, "user"), []byte("keep me"), 0644)
	plan := []File{{Path: "page", Data: []byte("recorder page"), Replace: []byte("known scaffold"), Mode: 0644}, {Path: "scripts/check.sh", Data: []byte("#!/bin/sh\n"), Mode: 0755}}
	changed, err := Apply(context.Background(), root, plan)
	if err != nil || len(changed) != 2 {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	changed, err = Apply(context.Background(), root, plan)
	if err != nil || len(changed) != 0 {
		t.Fatalf("reapply=%v err=%v", changed, err)
	}
	keep, _ := os.ReadFile(filepath.Join(root, "user"))
	if string(keep) != "keep me" {
		t.Fatal("lost unrelated content")
	}
	os.WriteFile(filepath.Join(root, "page"), []byte("my recorder changes"), 0644)
	if _, err := Apply(context.Background(), root, plan); err == nil {
		t.Fatal("overwrote edited output")
	}
}

func TestRejectsSymlinkEscapeAndCanceledRequest(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "native")); err != nil {
		t.Fatal(err)
	}
	plan := []File{{Path: "native/file", Data: []byte("no"), Mode: 0644}}
	if _, err := Apply(context.Background(), root, plan); err == nil {
		t.Fatal("followed symlink")
	}
	if _, err := os.Stat(filepath.Join(outside, "file")); !os.IsNotExist(err) {
		t.Fatal("wrote outside project")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Apply(ctx, root, []File{{Path: "new", Data: []byte("no"), Mode: 0644}}); err == nil {
		t.Fatal("ignored cancellation")
	}
}

func TestReapplyReportsLostExecutableMode(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "check.sh")
	data := []byte("#!/bin/sh\n")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(context.Background(), root, []File{{Path: "check.sh", Data: data, Mode: 0755}}); err == nil {
		t.Fatal("nonexecutable script treated as usable")
	}
}
