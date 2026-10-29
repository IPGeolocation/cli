package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/IPGeolocation/cli/v2/internal/utils"
)

func testHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return filepath.Join(home, ".ipgeolocation", "config.json")
}

func TestSaveRestrictsExistingPermissions(t *testing.T) {
	path := testHome(t)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"apikey":"old"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Save(Config{ApiKey: "replacement"}); err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		for name, want := range map[string]os.FileMode{path: 0600, filepath.Dir(path): 0700} {
			info, err := os.Stat(name)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != want {
				t.Errorf("%s permissions = %o, want %o", name, info.Mode().Perm(), want)
			}
		}
	}
	cfg, err := Load()
	if err != nil || cfg.ApiKey != "replacement" {
		t.Fatalf("saved key not readable: %+v %v", cfg, err)
	}
}

func TestLoadLegacyKeys(t *testing.T) {
	path := testHome(t)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	key := "legacy-test-key"
	encrypted, err := utils.EncryptString(key)
	if err != nil {
		t.Fatal(err)
	}
	for _, stored := range []string{key, encrypted} {
		data, err := json.Marshal(Config{ApiKey: stored})
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
		cfg, err := Load()
		if err != nil || cfg.ApiKey != key {
			t.Fatalf("legacy key not readable: %+v %v", cfg, err)
		}
	}
}

func TestSaveReportsInvalidDirectory(t *testing.T) {
	path := testHome(t)
	if err := os.WriteFile(filepath.Dir(path), []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Save(Config{ApiKey: "test"}); err == nil {
		t.Fatal("save accepted a non-directory")
	}
	data, err := os.ReadFile(filepath.Dir(path))
	if err != nil || string(data) != "keep" {
		t.Fatalf("existing file changed: %q %v", data, err)
	}
}
