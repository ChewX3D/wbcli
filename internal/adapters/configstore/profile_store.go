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
	configSchemaVersionV1  = 1
	timestampLayoutRFC3339 = time.RFC3339Nano
)

type storedConfig struct {
	SchemaVersion int                      `json:"schema_version"`
	ActiveProfile string                   `json:"active_profile,omitempty"`
	Profiles      map[string]storedProfile `json:"profiles,omitempty"`
}

type storedProfile struct {
	Backend    string `json:"backend,omitempty"`
	APIKeyHint string `json:"api_key_hint,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	LastUsedAt string `json:"last_used_at,omitempty"`
}

// FileProfileStore stores profile metadata in local config file.
type FileProfileStore struct {
	path string
	mu   sync.Mutex
}

// NewDefaultProfileStore constructs profile store at ~/.wbcli/config.yaml.
func NewDefaultProfileStore() (*FileProfileStore, error) {
	configPath, err := defaultConfigPath()
	if err != nil {
		return nil, err
	}

	return NewFileProfileStore(configPath), nil
}

// NewFileProfileStore constructs profile store at custom path.
func NewFileProfileStore(path string) *FileProfileStore {
	return &FileProfileStore{path: path}
}

func defaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home directory: %w", err)
	}

	return filepath.Join(homeDir, configDirName, defaultConfigFileName), nil
}

// UpsertProfile creates or updates non-secret profile metadata.
func (store *FileProfileStore) UpsertProfile(_ context.Context, profile ports.ProfileMetadata) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if profile.Name == "" {
		return errors.New("profile name is required")
	}

	config, err := store.loadConfig()
	if err != nil {
		return err
	}
	if config.Profiles == nil {
		config.Profiles = map[string]storedProfile{}
	}

	config.Profiles[profile.Name] = profileToStored(profile)

	return store.saveConfig(config)
}

// GetProfile returns metadata for a single profile.
func (store *FileProfileStore) GetProfile(_ context.Context, profileName string) (ports.ProfileMetadata, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return ports.ProfileMetadata{}, false, err
	}

	profile, ok := config.Profiles[profileName]
	if !ok {
		return ports.ProfileMetadata{}, false, nil
	}

	metadata, err := storedToProfile(profileName, profile)
	if err != nil {
		return ports.ProfileMetadata{}, false, err
	}

	return metadata, true, nil
}

// ListProfiles returns all profile metadata rows.
func (store *FileProfileStore) ListProfiles(_ context.Context) ([]ports.ProfileMetadata, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return nil, err
	}

	profiles := make([]ports.ProfileMetadata, 0, len(config.Profiles))
	for profileName, profile := range config.Profiles {
		metadata, conversionErr := storedToProfile(profileName, profile)
		if conversionErr != nil {
			return nil, conversionErr
		}
		profiles = append(profiles, metadata)
	}

	return profiles, nil
}

// DeleteProfile removes a metadata profile row.
func (store *FileProfileStore) DeleteProfile(_ context.Context, profileName string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return err
	}

	if config.Profiles == nil {
		return nil
	}

	delete(config.Profiles, profileName)

	return store.saveConfig(config)
}

// SetActiveProfile sets current selected profile name.
func (store *FileProfileStore) SetActiveProfile(_ context.Context, profileName string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return err
	}
	config.ActiveProfile = profileName

	return store.saveConfig(config)
}

// GetActiveProfile gets current selected profile name.
func (store *FileProfileStore) GetActiveProfile(_ context.Context) (string, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return "", err
	}

	return config.ActiveProfile, nil
}

// ClearActiveProfile clears selected active profile.
func (store *FileProfileStore) ClearActiveProfile(_ context.Context) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	config, err := store.loadConfig()
	if err != nil {
		return err
	}
	config.ActiveProfile = ""

	return store.saveConfig(config)
}

func (store *FileProfileStore) loadConfig() (storedConfig, error) {
	fileData, err := os.ReadFile(store.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return storedConfig{
				SchemaVersion: configSchemaVersionV1,
				Profiles:      map[string]storedProfile{},
			}, nil
		}

		return storedConfig{}, fmt.Errorf("read profile config: %w", err)
	}

	if len(fileData) == 0 {
		return storedConfig{
			SchemaVersion: configSchemaVersionV1,
			Profiles:      map[string]storedProfile{},
		}, nil
	}

	var config storedConfig
	if err := json.Unmarshal(fileData, &config); err != nil {
		return storedConfig{}, fmt.Errorf("decode profile config: %w", err)
	}
	if config.SchemaVersion == 0 {
		config.SchemaVersion = configSchemaVersionV1
	}
	if config.Profiles == nil {
		config.Profiles = map[string]storedProfile{}
	}

	return config, nil
}

func (store *FileProfileStore) saveConfig(config storedConfig) error {
	if config.SchemaVersion == 0 {
		config.SchemaVersion = configSchemaVersionV1
	}
	if config.Profiles == nil {
		config.Profiles = map[string]storedProfile{}
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
		return fmt.Errorf("encode profile config: %w", err)
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

func profileToStored(profile ports.ProfileMetadata) storedProfile {
	stored := storedProfile{
		Backend:    profile.Backend,
		APIKeyHint: profile.APIKeyHint,
		CreatedAt:  profile.CreatedAt.UTC().Format(timestampLayoutRFC3339),
		UpdatedAt:  profile.UpdatedAt.UTC().Format(timestampLayoutRFC3339),
	}
	if profile.LastUsedAt != nil {
		stored.LastUsedAt = profile.LastUsedAt.UTC().Format(timestampLayoutRFC3339)
	}

	return stored
}

func storedToProfile(profileName string, profile storedProfile) (ports.ProfileMetadata, error) {
	createdAt, err := parseTimestamp(profile.CreatedAt)
	if err != nil {
		return ports.ProfileMetadata{}, fmt.Errorf("decode created_at for profile %q: %w", profileName, err)
	}
	updatedAt, err := parseTimestamp(profile.UpdatedAt)
	if err != nil {
		return ports.ProfileMetadata{}, fmt.Errorf("decode updated_at for profile %q: %w", profileName, err)
	}
	lastUsedAt, err := parseTimestampPointer(profile.LastUsedAt)
	if err != nil {
		return ports.ProfileMetadata{}, fmt.Errorf("decode last_used_at for profile %q: %w", profileName, err)
	}

	return ports.ProfileMetadata{
		Name:       profileName,
		Backend:    profile.Backend,
		APIKeyHint: profile.APIKeyHint,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		LastUsedAt: lastUsedAt,
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

func parseTimestampPointer(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(timestampLayoutRFC3339, value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
