// A module for Secrets Provider examples in Go

package main

import (
	"context"
	"dagger/go/internal/dagger"
)

type Examples struct{}

// Demo the use of a secrets provider
func (e *Examples) Demo(
	ctx context.Context,
	provider string,
	// +optional
	client *dagger.Secret,
	token *dagger.Secret,
) (string, error) {
	// TODO: HACK
	if client == nil {
		client = dag.SetSecret("client", "")
	}

	// TODO:
	// This part should ideally be configured on the client or engine
	var secretsProvider *dagger.SecretsProviderProvider
	switch provider {
	case "infisical":
		secretsProvider = dag.ProviderInfisical().AsSecretsProviderProvider().Configure(client, token)
	case "op":
		secretsProvider = dag.ProviderOp().AsSecretsProviderProvider().Configure(client, token)
	}

	// Call a function that needs access to a secret provider
	return genericSecretThing(ctx, secretsProvider)
}

// A function that needs access to a secret provider
func genericSecretThing(ctx context.Context, provider *dagger.SecretsProviderProvider) (string, error) {
	// Read a secret

	// Information I probably know about the secret I want
	environment := "dev" // Maybe this is another client config?
	path := "Hackathon"
	key := "SUPER_SECRET"

	secret := provider.GetSecret(environment, path, key)

	return secret.Plaintext(ctx)
}
