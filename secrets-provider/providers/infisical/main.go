// A module that adheres the Infisical module to the SecretsProvider interface

package main

import (
	"context"
	"dagger/infisical/internal/dagger"
)

type ProviderInfisical struct {
	Client *dagger.Secret
	Secret *dagger.Secret
}

// Configure the Infisical provider with a client and secret
func (m *ProviderInfisical) Configure(ctx context.Context, client, token *dagger.Secret) (*ProviderInfisical, error) {
	m.Client = client
	m.Secret = token
	return m, nil
}

// TODO: infisical module does not provide a function to write secrets
func (m *ProviderInfisical) PutSecret(ctx context.Context, path, key string, secret *dagger.Secret) error {
	return nil
}

// Get a secret from Infisical
func (m *ProviderInfisical) GetSecret(ctx context.Context, environment, path, key string) (*dagger.Secret, error) {
	// TODO: Implement
	project := "6545b31d52162a370f2141fe"

	secret := dag.Infisical().
		WithUniversalAuth(m.Client, m.Secret).
		GetSecretByName(
			key,
			project,
			environment,
			"/"+path)
	return secret, nil
}
