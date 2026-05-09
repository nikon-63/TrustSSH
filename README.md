# TrustSSH

SSH login broker for issuing short-lived OpenSSH user certificates using AWS Cognito and Lambda.

The Go CLI provides `trustssh configure`, `trustssh login`, and `trustssh logout` commands that work with the deployed AWS infrastructure to authenticate users, and sign short-lived SSH certificates for use with OpenSSH.

TrustSSH does not upload the user's SSH private key. The CLI signs nothing locally and sends only the user's SSH public key to the signing API.

## Current CLI Commands

```bash
trustssh configure <base-url>
trustssh login
trustssh logout
```

`trustssh configure <base-url>` downloads `<base-url>/config.json` and saves it to `~/.trustssh/config.json`.

`trustssh login`:

1. Opens Cognito managed login in the browser.
2. Receives the localhost OAuth callback.
3. Exchanges the auth code using PKCE.
4. Creates or reuses `~/.trustssh/id_ed25519`.
5. Sends the public key to the signing API.
6. Saves the returned certificate as `~/.trustssh/id_ed25519-cert.pub`.

`trustssh logout` removes local tokens and the short-lived certificate. It does not remove the SSH key pair.

## Documentation

- [AWS deployment](docs/aws-deployment.md)
- [CLI deployment](docs/cli-deployment.md)
- [Request flow](docs/request-flow.md)

## Project Layout

```text
cli/              Go CLI
infra/            Terraform AWS infrastructure
lambda/signer/    Python certificate signer Lambda
docs/             Deployment and flow documentation
```
