package ports

import (
	"context"
	"errors"
	"time"

	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

var (
	// ErrCredentialNotFound indicates missing credentials in secure store.
	ErrCredentialNotFound = errors.New("credential not found")
	// ErrSecretStoreUnavailable indicates unavailable secure backend.
	ErrSecretStoreUnavailable = errors.New("secret store unavailable")
	// ErrSecretStorePermissionDenied indicates denied backend access.
	ErrSecretStorePermissionDenied = errors.New("secret store permission denied")
)

// SessionMetadata holds non-secret auth session information.
type SessionMetadata struct {
	Backend    string
	APIKeyHint string
	HedgeMode  *bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CredentialStore persists and retrieves secret credentials as a single slot.
type CredentialStore interface {
	BackendName() string
	Save(ctx context.Context, credential domainauth.Credential) error
	Load(ctx context.Context) (domainauth.Credential, error)
	Exists(ctx context.Context) (bool, error)
	Delete(ctx context.Context) error
}

// SessionStore persists and retrieves non-secret session metadata.
type SessionStore interface {
	SaveSession(ctx context.Context, session SessionMetadata) error
	GetSession(ctx context.Context) (SessionMetadata, bool, error)
	ClearSession(ctx context.Context) error
}

// Clock supplies deterministic time in services.
type Clock interface {
	Now() time.Time
}

// CredentialVerificationResult contains non-secret verification metadata.
type CredentialVerificationResult struct {
	Endpoint  string
	HedgeMode bool
}

// CredentialVerifier verifies credential validity and required permission against the exchange.
type CredentialVerifier interface {
	Verify(ctx context.Context, credential domainauth.Credential) (CredentialVerificationResult, error)
}
