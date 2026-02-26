package secretstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ChewX3D/wbcli/internal/app/ports"
	domainauth "github.com/ChewX3D/wbcli/internal/domain/auth"
)

const serviceName = "wbcli"

// OSKeychainStore stores credentials in platform keychain tools.
type OSKeychainStore struct{}

// NewOSKeychainStore constructs OS keychain-backed credential store.
func NewOSKeychainStore() *OSKeychainStore {
	return &OSKeychainStore{}
}

// BackendName returns backend identifier.
func (store *OSKeychainStore) BackendName() string {
	return "os-keychain"
}

// Save writes profile credential to OS keychain.
func (store *OSKeychainStore) Save(ctx context.Context, profile string, credential domainauth.Credential) error {
	payload, err := marshalCredential(credential)
	if err != nil {
		return fmt.Errorf("marshal credential payload: %w", err)
	}

	switch runtime.GOOS {
	case "darwin":
		return saveDarwin(ctx, profile, payload)
	case "linux":
		return saveLinux(ctx, profile, payload)
	default:
		return ports.ErrSecretStoreUnavailable
	}
}

// Load reads profile credential from OS keychain.
func (store *OSKeychainStore) Load(ctx context.Context, profile string) (domainauth.Credential, error) {
	var (
		payload string
		err     error
	)

	switch runtime.GOOS {
	case "darwin":
		payload, err = loadDarwin(ctx, profile)
	case "linux":
		payload, err = loadLinux(ctx, profile)
	default:
		return domainauth.Credential{}, ports.ErrSecretStoreUnavailable
	}
	if err != nil {
		return domainauth.Credential{}, err
	}

	credential, err := unmarshalCredential(payload)
	if err != nil {
		return domainauth.Credential{}, fmt.Errorf("unmarshal credential payload: %w", err)
	}

	return credential, nil
}

// Exists checks whether profile credentials exist.
func (store *OSKeychainStore) Exists(ctx context.Context, profile string) (bool, error) {
	_, err := store.Load(ctx, profile)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, ports.ErrCredentialNotFound) {
		return false, nil
	}

	return false, err
}

// Delete removes profile credential from OS keychain.
func (store *OSKeychainStore) Delete(ctx context.Context, profile string) error {
	switch runtime.GOOS {
	case "darwin":
		return deleteDarwin(ctx, profile)
	case "linux":
		return deleteLinux(ctx, profile)
	default:
		return ports.ErrSecretStoreUnavailable
	}
}

type credentialPayload struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

func marshalCredential(credential domainauth.Credential) (string, error) {
	encoded, err := json.Marshal(credentialPayload{
		APIKey:    credential.APIKey,
		APISecret: string(credential.APISecret),
	})
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func unmarshalCredential(payload string) (domainauth.Credential, error) {
	var decoded credentialPayload
	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		return domainauth.Credential{}, err
	}

	return domainauth.Credential{
		APIKey:    decoded.APIKey,
		APISecret: []byte(decoded.APISecret),
	}, nil
}

func saveDarwin(ctx context.Context, profile string, payload string) error {
	if _, err := exec.LookPath("security"); err != nil {
		return ports.ErrSecretStoreUnavailable
	}

	output, err := runCommand(ctx, "", "security",
		"add-generic-password",
		"-a", profile,
		"-s", serviceName,
		"-w", payload,
		"-U",
	)
	if err != nil {
		return mapCommandError(err, output)
	}

	return nil
}

func loadDarwin(ctx context.Context, profile string) (string, error) {
	if _, err := exec.LookPath("security"); err != nil {
		return "", ports.ErrSecretStoreUnavailable
	}

	output, err := runCommand(ctx, "", "security",
		"find-generic-password",
		"-a", profile,
		"-s", serviceName,
		"-w",
	)
	if err != nil {
		return "", mapCommandError(err, output)
	}

	return strings.TrimSpace(string(output)), nil
}

func deleteDarwin(ctx context.Context, profile string) error {
	if _, err := exec.LookPath("security"); err != nil {
		return ports.ErrSecretStoreUnavailable
	}

	output, err := runCommand(ctx, "", "security",
		"delete-generic-password",
		"-a", profile,
		"-s", serviceName,
	)
	if err != nil {
		mappedErr := mapCommandError(err, output)
		if errors.Is(mappedErr, ports.ErrCredentialNotFound) {
			return ports.ErrCredentialNotFound
		}

		return mappedErr
	}

	return nil
}

func saveLinux(ctx context.Context, profile string, payload string) error {
	if _, err := exec.LookPath("secret-tool"); err != nil {
		return ports.ErrSecretStoreUnavailable
	}

	output, err := runCommand(ctx, payload, "secret-tool",
		"store",
		"--label=wbcli "+profile,
		"service", serviceName,
		"profile", profile,
	)
	if err != nil {
		return mapCommandError(err, output)
	}

	return nil
}

func loadLinux(ctx context.Context, profile string) (string, error) {
	if _, err := exec.LookPath("secret-tool"); err != nil {
		return "", ports.ErrSecretStoreUnavailable
	}

	output, err := runCommand(ctx, "", "secret-tool",
		"lookup",
		"service", serviceName,
		"profile", profile,
	)
	if err != nil {
		return "", mapCommandError(err, output)
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return "", ports.ErrCredentialNotFound
	}

	return trimmed, nil
}

func deleteLinux(ctx context.Context, profile string) error {
	if _, err := exec.LookPath("secret-tool"); err != nil {
		return ports.ErrSecretStoreUnavailable
	}

	output, err := runCommand(ctx, "", "secret-tool",
		"clear",
		"service", serviceName,
		"profile", profile,
	)
	if err != nil {
		mappedErr := mapCommandError(err, output)
		if errors.Is(mappedErr, ports.ErrCredentialNotFound) {
			return ports.ErrCredentialNotFound
		}

		return mappedErr
	}

	return nil
}

func runCommand(ctx context.Context, stdin string, binary string, arguments ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, binary, arguments...)
	if stdin != "" {
		command.Stdin = strings.NewReader(stdin)
	}

	output, err := command.CombinedOutput()
	if err != nil {
		return output, err
	}

	return output, nil
}

func mapCommandError(err error, output []byte) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, exec.ErrNotFound) {
		return ports.ErrSecretStoreUnavailable
	}

	var commandErr *exec.Error
	if errors.As(err, &commandErr) && errors.Is(commandErr.Err, exec.ErrNotFound) {
		return ports.ErrSecretStoreUnavailable
	}

	lowerOutput := strings.ToLower(string(output))
	if strings.Contains(lowerOutput, "could not be found") ||
		strings.Contains(lowerOutput, "not found") ||
		strings.Contains(lowerOutput, "no such secret") ||
		strings.Contains(lowerOutput, "item not found") {
		return ports.ErrCredentialNotFound
	}
	if strings.Contains(lowerOutput, "permission denied") ||
		strings.Contains(lowerOutput, "interaction is not allowed") {
		return ports.ErrSecretStorePermissionDenied
	}

	return fmt.Errorf("os-keychain command failed: %w", err)
}
