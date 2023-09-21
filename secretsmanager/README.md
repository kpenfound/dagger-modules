# Secretsmanager

Dagger module for interacting with AWS Secrets Manager

## Examples

`echo "{secretsmanager { auth(key:\"$AWS_ACCESS_KEY_ID\", secret:\"$AWS_SECRET_ACCESS_KEY\") { putSecret(name:\"foo\", value:\"bar\") { getSecret(name:\"foo\")}}}}" | dagger query`
