package auth

import "errors"

var (
	// ErrNotLoggedIn indicates missing logged-in auth session.
	ErrNotLoggedIn = errors.New("not logged in")
	// ErrAuthTestNotImplemented indicates deferred auth test flow.
	ErrAuthTestNotImplemented = errors.New("auth test is not implemented yet")
)
