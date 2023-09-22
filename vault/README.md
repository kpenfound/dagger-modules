# Vault

Dagger module for interacting with Vault

## Examples

`echo "{vault { auth(approleID:\"$MY_APPROLE_ROLE_ID\", approleSecret:\"$MY_APPROLE_SECRET_ID\", address:\"$VAULT_ADDR\") { putSecret(secret:\"myapp\", key:\"foo\", value:\"bar\") { getSecret(secret:\"myapp\", key:\"foo\")}}}}" | dagger query`
