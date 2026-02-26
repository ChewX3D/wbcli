package configstore

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
)

func TestFileProfileStoreWritesMetadataOnlyWith0600Permissions(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	store := NewFileProfileStore(configPath)

	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)
	err := store.UpsertProfile(context.Background(), ports.ProfileMetadata{
		Name:       "prod",
		Backend:    "os-keychain",
		APIKeyHint: "ab***yz",
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := store.SetActiveProfile(context.Background(), "prod"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	fileData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	text := string(fileData)
	if strings.Contains(text, "api_secret") || strings.Contains(text, "secret") {
		t.Fatalf("config must not contain secret values, got: %s", text)
	}

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if fileInfo.Mode().Perm() != 0o600 {
		t.Fatalf("expected config mode 0600, got %o", fileInfo.Mode().Perm())
	}

	activeProfile, err := store.GetActiveProfile(context.Background())
	if err != nil {
		t.Fatalf("get active profile: %v", err)
	}
	if activeProfile != "prod" {
		t.Fatalf("expected active profile prod, got %q", activeProfile)
	}
}
