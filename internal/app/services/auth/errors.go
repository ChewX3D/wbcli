package auth

import "errors"

var (
	// ErrNotLoggedIn indicates missing logged-in auth session.
	ErrNotLoggedIn = errors.New("not logged in")
)
