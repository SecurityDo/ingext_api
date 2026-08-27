package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadProcessorInput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "script.fpl")
	if err := os.WriteFile(path, []byte("function main() {}\n"), 0644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	got, err := readProcessorInput(path, "--script", strings.NewReader(""))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(got) != "function main() {}\n" {
		t.Fatalf("unexpected content %q", got)
	}

	// '@path' is accepted for symmetry with 'processor add --content'.
	if got, err = readProcessorInput("@"+path, "--script", strings.NewReader("")); err != nil {
		t.Fatalf("expected no error, got %v", err)
	} else if string(got) != "function main() {}\n" {
		t.Fatalf("unexpected content %q", got)
	}

	if got, err = readProcessorInput("-", "--event", strings.NewReader(`{"a":1}`)); err != nil {
		t.Fatalf("expected no error, got %v", err)
	} else if string(got) != `{"a":1}` {
		t.Fatalf("unexpected content %q", got)
	}

	if _, err = readProcessorInput("", "--script", strings.NewReader("")); err == nil {
		t.Fatalf("expected an error for an empty path")
	}
	if _, err = readProcessorInput(filepath.Join(dir, "missing.fpl"), "--script", strings.NewReader("")); err == nil {
		t.Fatalf("expected an error for a missing file")
	}
}
