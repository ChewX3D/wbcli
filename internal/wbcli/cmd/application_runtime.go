package cmd

import (
	"fmt"
	"sync"

	appcontainer "github.com/ChewX3D/crypto/internal/app/application"
)

type applicationProvider func() (*appcontainer.Application, error)

func newApplicationProvider(factory func() (*appcontainer.Application, error)) applicationProvider {
	var (
		once      sync.Once
		cached    *appcontainer.Application
		cachedErr error
	)

	return func() (*appcontainer.Application, error) {
		once.Do(func() {
			cached, cachedErr = factory()
			if cachedErr != nil {
				cachedErr = fmt.Errorf("init application: %w", cachedErr)
			}
		})

		return cached, cachedErr
	}
}
