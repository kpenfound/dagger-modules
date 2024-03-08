// Interact with HashiCorp Vault
//
// A utility module for working with secrets in HashiCorp Vault

package main

import (
	"context"
	"time"

	vcg "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type Vault struct {
	ApproleID     string
	ApproleSecret string
	Address       string
}

// Authenticate to Vault with an approle secret
func (v *Vault) Auth(approleID, approleSecret, address string) *Vault {
	v.ApproleID = approleID
	v.ApproleSecret = approleSecret
	v.Address = address

	return v
}

// Get a secret from Vault
func (v *Vault) GetSecret(ctx context.Context, secret, key string) (*Secret, error) {
	client, err := getClient(ctx, v)
	if err != nil {
		return nil, err
	}
	s, err := client.Secrets.KvV2Read(ctx, secret, vcg.WithMountPath("secret"))
	if err != nil {
		return nil, err
	}
	dagSecret := dag.SetSecret(key, s.Data.Data[key].(string))
	return dagSecret, nil
}

// Put a secret in Vault
func (v *Vault) PutSecret(ctx context.Context, secret, key, value string) (*Vault, error) {
	client, err := getClient(ctx, v)
	if err != nil {
		return nil, err
	}
	_, err = client.Secrets.KvV2Write(ctx, secret, schema.KvV2WriteRequest{
		Data: map[string]any{
			key: value,
		}},
		vcg.WithMountPath("secret"),
	)
	return v, err
}


func getClient(ctx context.Context, v *Vault) (*vcg.Client, error) {
	client, err := vcg.New(
		vcg.WithAddress(v.Address),
		vcg.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, err
	}
	resp, err := client.Auth.AppRoleLogin(
		ctx,
		schema.AppRoleLoginRequest{
			RoleId:   v.ApproleID,
			SecretId: v.ApproleSecret,
		},
	)
	if err != nil {
		return nil, err
	}

	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		return nil, err
	}
	return client, nil
}
