package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Secretsmanager struct {
	service *secretsmanager.SecretsManager
}

func (m *Secretsmanager) Auth(key, secret string) (*Secretsmanager, error) {
	config := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	svc := secretsmanager.New(sess)
	m.service = svc
	return m, nil
}

func (m *Secretsmanager) GetSecret(name string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	value, err := m.service.GetSecretValue(input)
	if err != nil {
		return "", err
	}
	return *(value.SecretString), nil
}

func (m *Secretsmanager) PutSecret(name, value string) (*Secretsmanager, error) {
	input := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(name),
		SecretString: aws.String(value),
	}

	_, err := m.service.PutSecretValue(input)
	return m, err
}
