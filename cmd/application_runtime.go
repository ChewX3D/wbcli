package cmd

import (
	"fmt"
	"sync"

	appcontainer "github.com/ChewX3D/wbcli/internal/app/application"
)

type applicationProvider func() (*appcontainer.Application, error)

var applicationFactory = appcontainer.NewDefault

func newApplicationProvider() applicationProvider {
	var (
		once      sync.Once
		cached    *appcontainer.Application
		cachedErr error
	)

	return func() (*appcontainer.Application, error) {
		once.Do(func() {
			cached, cachedErr = applicationFactory()
			if cachedErr != nil {
				cachedErr = fmt.Errorf("init application: %w", cachedErr)
			}
		})

		return cached, cachedErr
	}
}

// SetApplicationFactoryForTest overrides runtime application factory.
func SetApplicationFactoryForTest(factory func() (*appcontainer.Application, error)) func() {
	previousFactory := applicationFactory
	applicationFactory = factory

	return func() {
		applicationFactory = previousFactory
	}
}
