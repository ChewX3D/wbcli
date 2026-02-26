package auth

import "errors"

var (
	// ErrCredentialAlreadyExists indicates existing credentials for a profile.
	ErrCredentialAlreadyExists = errors.New("credentials already exist")
	// ErrProfileNotFound indicates requested profile metadata was not found.
	ErrProfileNotFound = errors.New("profile not found")
	// ErrNoActiveProfile indicates that active profile is not selected.
	ErrNoActiveProfile = errors.New("no active profile selected")
	// ErrAuthTestNotImplemented indicates deferred auth test flow.
	ErrAuthTestNotImplemented = errors.New("auth test is not implemented yet")
)
