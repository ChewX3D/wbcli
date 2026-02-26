package auth

import "time"

// SystemClock implements ports.Clock using wall-clock time.
type SystemClock struct{}

// Now returns current time.
func (SystemClock) Now() time.Time {
	return time.Now()
}
