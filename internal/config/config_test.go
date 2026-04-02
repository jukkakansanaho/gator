package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// writeFixture creates a .gatorconfig.json in dir with the given content.
func writeFixture(t *testing.T, dir string, cfg Config) {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("writeFixture: marshal: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, configFileName), data, 0600); err != nil {
		t.Fatalf("writeFixture: write: %v", err)
	}
}

// useHome redirects HOME to dir for the duration of the test.
func useHome(t *testing.T, dir string) {
	t.Helper()
	original := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	t.Cleanup(func() { os.Setenv("HOME", original) })
}

func TestGetConfigFilePath(t *testing.T) {
	tmp := t.TempDir()
	useHome(t, tmp)

	got, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := filepath.Join(tmp, configFileName)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRead_ValidFile(t *testing.T) {
	tmp := t.TempDir()
	useHome(t, tmp)

	want := Config{DbURL: "postgres://test", CurrentUserName: "alice"}
	writeFixture(t, tmp, want)

	got, err := Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestRead_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	useHome(t, tmp)

	_, err := Read()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRead_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	useHome(t, tmp)

	if err := os.WriteFile(filepath.Join(tmp, configFileName), []byte("not json"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := Read()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSetUser(t *testing.T) {
	tmp := t.TempDir()
	useHome(t, tmp)

	writeFixture(t, tmp, Config{DbURL: "postgres://test"})

	cfg, err := Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if err := cfg.SetUser("bob"); err != nil {
		t.Fatalf("SetUser: %v", err)
	}

	persisted, err := Read()
	if err != nil {
		t.Fatalf("Read after SetUser: %v", err)
	}
	if persisted.CurrentUserName != "bob" {
		t.Errorf("CurrentUserName = %q, want %q", persisted.CurrentUserName, "bob")
	}
	if persisted.DbURL != "postgres://test" {
		t.Errorf("DbURL = %q, want %q", persisted.DbURL, "postgres://test")
	}
}

func TestConfigJSONMarshal(t *testing.T) {
	cfg := Config{DbURL: "postgres://example", CurrentUserName: "jukka"}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	want := `{"db_url":"postgres://example","current_user_name":"jukka"}`
	if got := string(data); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestConfigJSONMarshal_EmptyFields(t *testing.T) {
	cfg := Config{}

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	want := `{"db_url":"","current_user_name":""}`
	if got := string(data); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestWrite(t *testing.T) {
	tmp := t.TempDir()
	useHome(t, tmp)

	cfg := Config{DbURL: "postgres://write-test", CurrentUserName: "carol"}
	if err := write(cfg); err != nil {
		t.Fatalf("write: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, configFileName))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var got Config
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got != cfg {
		t.Errorf("got %+v, want %+v", got, cfg)
	}
}
