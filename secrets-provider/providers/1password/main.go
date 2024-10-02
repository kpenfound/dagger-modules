// A generated module for 1Password functions

package main

import (
	"context"
	"dagger/1-password/internal/dagger"
)

type ProviderOp struct {
	Token *dagger.Secret
}

// Configure the 1Password provider with a token
func (m *ProviderOp) Configure(ctx context.Context, client, token *dagger.Secret) (*ProviderOp, error) {
	m.Token = token
	return m, nil
}

// TODO: implement
func (m *ProviderOp) PutSecret(ctx context.Context, path, key string, secret *dagger.Secret) error {
	return nil
}

// Get a secret from 1Password
func (m *ProviderOp) GetSecret(ctx context.Context, environment, path, key string) (*dagger.Secret, error) {
	return dag.Onepassword().FindSecret(
		m.Token,
		environment,
		path,
		key,
	), nil
}
