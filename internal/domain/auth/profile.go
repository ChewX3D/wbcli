package auth

import (
	"errors"
	"regexp"
	"strings"
)

var profileNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,64}$`)

var (
	// ErrProfileNameRequired indicates that the profile value is empty.
	ErrProfileNameRequired = errors.New("profile is required")
	// ErrProfileNameInvalid indicates that the profile value violates naming rules.
	ErrProfileNameInvalid = errors.New("profile format is invalid")
)

// ValidateProfileName validates the profile naming contract.
func ValidateProfileName(profile string) error {
	if profile == "" {
		return ErrProfileNameRequired
	}
	if strings.TrimSpace(profile) != profile {
		return ErrProfileNameInvalid
	}
	if !profileNamePattern.MatchString(profile) {
		return ErrProfileNameInvalid
	}

	return nil
}
