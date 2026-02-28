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

func TestFileSessionStoreWritesMetadataOnlyWith0600Permissions(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	store := NewFileSessionStore(configPath)

	now := time.Date(2026, 2, 26, 12, 0, 0, 0, time.UTC)
	err := store.SaveSession(context.Background(), ports.SessionMetadata{
		Backend:    "os-keychain",
		APIKeyHint: "ab***yz",
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
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

	session, found, err := store.GetSession(context.Background())
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if !found {
		t.Fatalf("expected session to exist")
	}
	if session.Backend != "os-keychain" {
		t.Fatalf("expected backend os-keychain, got %q", session.Backend)
	}
}
