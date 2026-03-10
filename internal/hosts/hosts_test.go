package hosts

import (
	"os"
	"strings"
	"testing"
)

func TestBlockTemplate(t *testing.T) {
	result := BlockTemplate([]string{"twitter.com", "reddit.com"})

	if !strings.Contains(result, "### Pomodoro CLI - Begin Blocked sites ###") {
		t.Error("missing begin marker")
	}
	if !strings.Contains(result, "### Pomodoro CLI - End Blocked sites ###") {
		t.Error("missing end marker")
	}
	if !strings.Contains(result, "127.0.0.1 twitter.com") {
		t.Error("missing twitter.com entry")
	}
	if !strings.Contains(result, "127.0.0.1 www.twitter.com") {
		t.Error("missing www.twitter.com entry")
	}
	if !strings.Contains(result, "127.0.0.1 reddit.com") {
		t.Error("missing reddit.com entry")
	}
}

func tempFile(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "hosts-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return f
}

func readFile(t *testing.T, f *os.File) string {
	t.Helper()
	if _, err := f.Seek(0, 0); err != nil {
		t.Fatalf("seek failed: %v", err)
	}
	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	return string(data)
}

func TestBlock(t *testing.T) {
	original := "127.0.0.1 localhost\n"
	f := tempFile(t, original)
	defer f.Close()

	tmpl := BlockTemplate([]string{"example.com"})
	if err := Block(tmpl, f); err != nil {
		t.Fatalf("Block returned error: %v", err)
	}

	content := readFile(t, f)
	if !strings.Contains(content, original) {
		t.Error("original content was removed")
	}
	if !strings.Contains(content, tmpl) {
		t.Error("block template was not appended")
	}
}

func TestBlock_idempotent(t *testing.T) {
	f := tempFile(t, "127.0.0.1 localhost\n")
	defer f.Close()

	tmpl := BlockTemplate([]string{"example.com"})
	if err := Block(tmpl, f); err != nil {
		t.Fatalf("first Block error: %v", err)
	}
	if err := Block(tmpl, f); err != nil {
		t.Fatalf("second Block error: %v", err)
	}

	content := readFile(t, f)
	if count := strings.Count(content, "### Pomodoro CLI - Begin Blocked sites ###"); count != 1 {
		t.Errorf("expected block marker once, got %d", count)
	}
}

func TestUnblock(t *testing.T) {
	original := "127.0.0.1 localhost\n"
	f := tempFile(t, original)
	defer f.Close()

	tmpl := BlockTemplate([]string{"example.com"})
	if err := Block(tmpl, f); err != nil {
		t.Fatalf("Block error: %v", err)
	}
	if err := Unblock(tmpl, f); err != nil {
		t.Fatalf("Unblock error: %v", err)
	}

	content := readFile(t, f)
	if content != original {
		t.Errorf("expected %q, got %q", original, content)
	}
}

func TestUnblock_noop(t *testing.T) {
	original := "127.0.0.1 localhost\n"
	f := tempFile(t, original)
	defer f.Close()

	tmpl := BlockTemplate([]string{"example.com"})
	if err := Unblock(tmpl, f); err != nil {
		t.Fatalf("Unblock error: %v", err)
	}

	content := readFile(t, f)
	if content != original {
		t.Errorf("expected file unchanged, got %q", content)
	}
}
