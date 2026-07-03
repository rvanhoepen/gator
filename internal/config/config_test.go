package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	path := filepath.Join(tempDir, configFileName)

	err := os.WriteFile(path, []byte(`{
		"db_url": "postgres://example",
		"current_user_name": "dev"
	}`), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Read()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DBURL != "postgres://example" {
		t.Errorf("expected DBURL %q, got %q", "postgres://example", cfg.DBURL)
	}

	if cfg.CurrentUserName != "dev" {
		t.Errorf("expected CurrentUserName %q, got %q", "dev", cfg.CurrentUserName)
	}
}

func TestSetUserWritesConfig(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	cfg := Config{
		DBURL:           "postgres://example",
		CurrentUserName: "old_user",
	}

	err := cfg.SetUser("new_user")
	if err != nil {
		t.Fatal(err)
	}

	readCfg, err := Read()

	if readCfg.DBURL != "postgres://example" {
		t.Errorf("expected DBURL to stay %q, got %q", "postgres://example", readCfg.DBURL)
	}

	if readCfg.CurrentUserName != "new_user" {
		t.Errorf("expected CurrentUserName to be %q, got %q", "new_user", readCfg.CurrentUserName)
	}
}
