# TrustSSH

![TrustSSH Logo](images/logo.png)

SSH login broker for issuing short-lived OpenSSH user certificates using AWS Cognito and Lambda.

The Go CLI provides `trustssh configure`, `trustssh login`, and `trustssh logout` commands that work with the deployed AWS infrastructure to authenticate users, and sign short-lived SSH certificates for use with OpenSSH.

TrustSSH does not upload the user's SSH private key. The CLI signs nothing locally and sends only the user's SSH public key to the signing API.

Demostation of the CLI login flow:

```bash
# Before login, SSH access is denied:
user@MacBook TrustSSH % ssh -i ~/.trustssh/id_ed25519 ubuntu@demo.com
ubuntu@demo.com: Permission denied (publickey).
# Login with the CLI, which opens the browser for Cognito auth:
user@MacBook TrustSSH % ./trustssh login                                   
Opening browser for Cognito login...
Requesting 30 minute certificate...
Authenticated as demo@example.com
Using SSH key: ~/.trustssh/id_ed25519.pub
Allowed SSH principals: ubuntu, demo
Certificate saved: ~/.trustssh/id_ed25519-cert.pub
Certificate valid until: 2026-05-09T14:38:28Z
Tokens saved: /Users/user/.trustssh/tokens.json
You can now use normal SSH commands.
# After login, SSH access is granted with the short-lived certificate:
user@MacBook TrustSSH % ssh -i ~/.trustssh/id_ed25519 ubuntu@192.168.100.31
Welcome to Ubuntu 24.04 LTS (GNU/Linux 6.17.2-1-pve x86_64)
ubuntu@ssh-demo:~$ 
```

## Current CLI Commands

```bash
trustssh configure <base-url>
trustssh passkeys add
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
lambda/           Python Lambda Functions
docs/             Deployment and flow documentation
helpers/          Helper scripts and shared code
```
