package ports

import (
	"context"
	"errors"
	"time"

	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

var (
	// ErrCredentialNotFound indicates missing credentials for a profile.
	ErrCredentialNotFound = errors.New("credential not found")
	// ErrSecretStoreUnavailable indicates unavailable secure backend.
	ErrSecretStoreUnavailable = errors.New("secret store unavailable")
	// ErrSecretStorePermissionDenied indicates denied backend access.
	ErrSecretStorePermissionDenied = errors.New("secret store permission denied")
)

// ProfileMetadata holds non-secret profile information.
type ProfileMetadata struct {
	Name       string
	Backend    string
	APIKeyHint string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastUsedAt *time.Time
}

// CredentialStore persists and retrieves secret credentials.
type CredentialStore interface {
	BackendName() string
	Save(ctx context.Context, profile string, credential domainauth.Credential) error
	Load(ctx context.Context, profile string) (domainauth.Credential, error)
	Exists(ctx context.Context, profile string) (bool, error)
	Delete(ctx context.Context, profile string) error
}

// ProfileStore persists and retrieves non-secret profile metadata.
type ProfileStore interface {
	UpsertProfile(ctx context.Context, profile ProfileMetadata) error
	GetProfile(ctx context.Context, profile string) (ProfileMetadata, bool, error)
	ListProfiles(ctx context.Context) ([]ProfileMetadata, error)
	DeleteProfile(ctx context.Context, profile string) error
	SetActiveProfile(ctx context.Context, profile string) error
	GetActiveProfile(ctx context.Context) (string, error)
	ClearActiveProfile(ctx context.Context) error
}

// Clock supplies deterministic time in services.
type Clock interface {
	Now() time.Time
}

// AuthProbe verifies a stored credential against WhiteBIT.
type AuthProbe interface {
	Probe(ctx context.Context, credential domainauth.Credential) error
}
