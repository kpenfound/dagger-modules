package main

import (
	"context"
	"fmt"
	"time"

	vcg "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type Vault struct {
	ApproleID     string
	ApproleSecret string
	Address       string
}

func (v *Vault) Auth(approleID, approleSecret, address string) *Vault {
	v.ApproleID = approleID
	v.ApproleSecret = approleSecret
	v.Address = address

	return v
}

func (v *Vault) GetSecret(ctx context.Context, secret, key string) (string, error) {
	client, err := getClient(ctx, v)
	if err != nil {
		return "", err
	}
	s, err := client.Secrets.KvV2Read(ctx, secret, vcg.WithMountPath("secret"))
	if err != nil {
		return "", err
	}
	return fmt.Sprint(s.Data.Data[key]), nil
}

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

func (c *Container) WithVaultSecret(ctx context.Context, approleID, approleSecret, address, secret, key, name string) (*Container, error) {
	v := &Vault{
		ApproleID:     approleID,
		ApproleSecret: approleSecret,
		Address:       address,
	}

	s, err := v.GetSecret(ctx, secret, key)
	if err != nil {
		return nil, err
	}

	dagSecret := dag.SetSecret(name, s)
	return c.WithMountedSecret(name, dagSecret), nil
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
