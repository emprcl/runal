package js

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestIsSketchEvent(t *testing.T) {
	r := runtime{filename: filepath.Join("sketches", "sketch.js")}

	tests := []struct {
		name  string
		event fsnotify.Event
		want  bool
	}{
		{
			name:  "write to the sketch",
			event: fsnotify.Event{Name: filepath.Join("sketches", "sketch.js"), Op: fsnotify.Write},
			want:  true,
		},
		{
			// Editors that save atomically rename a temp file over the
			// original, which never produces a Write for the sketch.
			name:  "create from an atomic save",
			event: fsnotify.Event{Name: filepath.Join("sketches", "sketch.js"), Op: fsnotify.Create},
			want:  true,
		},
		{
			name:  "rename of the sketch",
			event: fsnotify.Event{Name: filepath.Join("sketches", "sketch.js"), Op: fsnotify.Rename},
			want:  true,
		},
		{
			// A suffix match would wrongly reload on this one.
			name:  "different file sharing a suffix",
			event: fsnotify.Event{Name: filepath.Join("sketches", "my-sketch.js"), Op: fsnotify.Write},
			want:  false,
		},
		{
			name:  "unrelated file in the same directory",
			event: fsnotify.Event{Name: filepath.Join("sketches", "other.js"), Op: fsnotify.Write},
			want:  false,
		},
		{
			name:  "chmod is not an edit",
			event: fsnotify.Event{Name: filepath.Join("sketches", "sketch.js"), Op: fsnotify.Chmod},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.shouldReload(tt.event); got != tt.want {
				t.Errorf("shouldReload(%v) = %v, want %v", tt.event, got, tt.want)
			}
		})
	}
}

func TestSameFile(t *testing.T) {
	if !sameFile("sketch.js", "./sketch.js") {
		t.Error("sameFile should resolve equivalent relative paths")
	}
	if !sameFile("a/../sketch.js", "sketch.js") {
		t.Error("sameFile should clean paths before comparing")
	}
	if sameFile("sketch.js", "my-sketch.js") {
		t.Error("sameFile must not match on suffix")
	}
}

func TestParseJSRequiresSetupAndDraw(t *testing.T) {
	tests := []struct {
		name    string
		script  string
		wantErr bool
	}{
		{
			name:   "valid sketch",
			script: `function setup(c) {}; function draw(c) {}`,
		},
		{
			name:    "missing draw",
			script:  `function setup(c) {}`,
			wantErr: true,
		},
		{
			name:    "missing setup",
			script:  `function draw(c) {}`,
			wantErr: true,
		},
		{
			name:    "syntax error",
			script:  `function setup( {`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, err := parseJS(tt.script)
			if tt.wantErr && err == nil {
				t.Error("expected an error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseJSOptionalCallbacks(t *testing.T) {
	_, _, _, cb, err := parseJS(`function setup(c) {}; function draw(c) {}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Absent callbacks must stay nil so the canvas can skip the event.
	if cb.onKey != nil || cb.onMouseMove != nil {
		t.Error("expected absent callbacks to be nil")
	}

	_, _, _, cb, err = parseJS(`function setup(c) {}; function draw(c) {}; function onKey(c, e) {}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cb.onKey == nil {
		t.Error("expected onKey to be picked up")
	}
}

// A sketch reached through a symlink must still match the events
// fsnotify reports, which name the real path.
func TestSameFileFollowsSymlinks(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "sketch.js")
	if err := os.WriteFile(real, []byte("// sketch"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.js")
	if err := os.Symlink(real, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if !sameFile(link, real) {
		t.Error("a symlinked sketch should match its target")
	}
	if !sameFile(real, real) {
		t.Error("a path should match itself")
	}
	if sameFile(real, filepath.Join(dir, "other.js")) {
		t.Error("distinct sketches should not match")
	}
}

// A file that does not exist yet (an atomic save seen mid-rename) must
// still compare by absolute path rather than dropping the reload.
func TestSameFileOnMissingPath(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "sketch.js")
	if !sameFile(missing, missing) {
		t.Error("a missing path should still match itself")
	}
}
