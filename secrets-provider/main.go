// A SecretsProvider interface

package main

import (
	"context"
	"dagger/secrets-provider/internal/dagger"
)

type SecretsProvider struct {
	P Provider // This needs to be here to access the type
}

// Secret Provider interface
type Provider interface {
	dagger.DaggerObject
	// A function to configure the client ;TODO: add more parameters
	Configure(client dagger.Secret, token dagger.Secret) Provider
	// A function to write a secret
	PutSecret(ctx context.Context, path, key string, secret *dagger.Secret) error
	// A function to read a secret
	GetSecret(environment, path, key string) *dagger.Secret
}
