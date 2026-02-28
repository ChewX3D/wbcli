package configstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ChewX3D/wbcli/internal/app/ports"
)

const (
	configDirName          = ".wbcli"
	defaultConfigFileName  = "config.yaml"
	configSchemaVersionV2  = 2
	timestampLayoutRFC3339 = time.RFC3339Nano
)

type storedConfig struct {
	SchemaVersion int            `json:"schema_version"`
	Session       *storedSession `json:"session,omitempty"`
}

type storedSession struct {
	Backend    string `json:"backend,omitempty"`
	APIKeyHint string `json:"api_key_hint,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

// FileSessionStore stores auth session metadata in local config file.
type FileSessionStore struct {
	path string
	mu   sync.Mutex
}

// NewDefaultSessionStore constructs session store at ~/.wbcli/config.yaml.
func NewDefaultSessionStore() (*FileSessionStore, error) {
	configPath, err := defaultConfigPath()
	if err != nil {
		return nil, err
	}

	return NewFileSessionStore(configPath), nil
}

// NewFileSessionStore constructs session store at custom path.
func NewFileSessionStore(path string) *FileSessionStore {
	return &FileSessionStore{path: path}
}

func defaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home directory: %w", err)
	}

	return filepath.Join(homeDir, configDirName, defaultConfigFileName), nil
}

// SaveSession creates or updates non-secret auth session metadata.
func (store *FileSessionStore) SaveSession(_ context.Context, session ports.SessionMetadata) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return err
	}

	config.Session = &storedSession{
		Backend:    session.Backend,
		APIKeyHint: session.APIKeyHint,
		CreatedAt:  session.CreatedAt.UTC().Format(timestampLayoutRFC3339),
		UpdatedAt:  session.UpdatedAt.UTC().Format(timestampLayoutRFC3339),
	}

	return store.saveConfig(config)
}

// GetSession returns current auth session metadata.
func (store *FileSessionStore) GetSession(_ context.Context) (ports.SessionMetadata, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return ports.SessionMetadata{}, false, err
	}
	if config.Session == nil {
		return ports.SessionMetadata{}, false, nil
	}

	metadata, err := storedToSession(*config.Session)
	if err != nil {
		return ports.SessionMetadata{}, false, err
	}

	return metadata, true, nil
}

// ClearSession clears auth session metadata.
func (store *FileSessionStore) ClearSession(_ context.Context) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return err
	}
	config.Session = nil

	return store.saveConfig(config)
}

func (store *FileSessionStore) loadConfig() (storedConfig, error) {
	fileData, err := os.ReadFile(store.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return storedConfig{SchemaVersion: configSchemaVersionV2}, nil
		}

		return storedConfig{}, fmt.Errorf("read session config: %w", err)
	}

	if len(fileData) == 0 {
		return storedConfig{SchemaVersion: configSchemaVersionV2}, nil
	}

	var config storedConfig
	if err := json.Unmarshal(fileData, &config); err != nil {
		return storedConfig{}, fmt.Errorf("decode session config: %w", err)
	}
	if config.SchemaVersion == 0 {
		config.SchemaVersion = configSchemaVersionV2
	}

	return config, nil
}

func (store *FileSessionStore) saveConfig(config storedConfig) error {
	if config.SchemaVersion == 0 {
		config.SchemaVersion = configSchemaVersionV2
	}

	configDir := filepath.Dir(store.path)
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	tempFile, err := os.CreateTemp(configDir, "config-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp config file: %w", err)
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath)

	if err := tempFile.Chmod(0o600); err != nil {
		tempFile.Close()
		return fmt.Errorf("set temp config mode: %w", err)
	}

	encoded, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		tempFile.Close()
		return fmt.Errorf("encode session config: %w", err)
	}
	encoded = append(encoded, '\n')

	if _, err := tempFile.Write(encoded); err != nil {
		tempFile.Close()
		return fmt.Errorf("write temp config: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp config: %w", err)
	}

	if err := os.Rename(tempFilePath, store.path); err != nil {
		return fmt.Errorf("replace config file: %w", err)
	}
	if err := os.Chmod(store.path, 0o600); err != nil {
		return fmt.Errorf("set config mode: %w", err)
	}

	return nil
}

func storedToSession(session storedSession) (ports.SessionMetadata, error) {
	createdAt, err := parseTimestamp(session.CreatedAt)
	if err != nil {
		return ports.SessionMetadata{}, fmt.Errorf("decode created_at: %w", err)
	}
	updatedAt, err := parseTimestamp(session.UpdatedAt)
	if err != nil {
		return ports.SessionMetadata{}, fmt.Errorf("decode updated_at: %w", err)
	}

	return ports.SessionMetadata{
		Backend:    session.Backend,
		APIKeyHint: session.APIKeyHint,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func parseTimestamp(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	parsed, err := time.Parse(timestampLayoutRFC3339, value)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}
