package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	args := []string{"-out", dir, "-width", "32", "-height", "32", "-grid", "5"}

	if err := run(append(args, "-count", "2")); err != nil {
		t.Fatal(err)
	}
	if err := run(append(args, "-seed", "alice", "-prefix", "avatar")); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"identicon_001.png", "identicon_002.png", "avatar.png"} {
		file, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := png.Decode(file); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		file.Close()
	}
}
