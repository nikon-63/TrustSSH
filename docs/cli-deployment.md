# CLI Deployment

This guide covers building, installing, and configuring the TrustSSH CLI.

## Build the Binary

!NOTE: Dev build ```go build -o trustssh .``` does not include the version string, so the CLI will report version as `dev`.

From the CLI directory:

```bash
VERSION=$(tr -d '\n' < ../VERSION)
go build -ldflags "-X github.com/nikon-63/TrustSSH/cli/cmd.Version=${VERSION}" -o trustssh .
```

Run it directly:

```bash
./trustssh configure https://trustssh.demo.com
./trustssh passkeys add
./trustssh login
./trustssh logout
```

## Install on PATH

Install the binary to `~/bin` and ensure that is on your `PATH`:

```bash
mkdir -p ~/bin
VERSION=$(tr -d '\n' < ../VERSION)
go build -ldflags "-X github.com/nikon-63/TrustSSH/cli/cmd.Version=${VERSION}" -o ~/bin/trustssh .
```

Edit `~/.zshrc`:

```bash
export PATH="$HOME/bin:$PATH"
```

## CLI Config File

The CLI reads:

```text
~/.trustssh/config.json
```

Create the directory:

```bash
mkdir -p ~/.trustssh
chmod 700 ~/.trustssh
```

Recommended setup is to fetch the generated config from the deployed TrustSSH API:

```bash
trustssh configure https://trustssh.demo.com
```

This downloads:

```text
https://trustssh.demo.com/config.json
```

and saves it to:

```text
~/.trustssh/config.json
```

with `0600` permissions.

You can also create `~/.trustssh/config.json` manually from Terraform outputs:

```json
{
  "region": "eu-west-2",
  "cognito_domain": "https://trustssh-cli.auth.eu-west-2.amazoncognito.com",
  "client_id": "the-cognito-client-id",
  "redirect_uri": "http://localhost:8765/callback",
  "api_base_url": "https://trustssh.demo.com",
  "default_duration_seconds": 1800
}
```

Set permissions:

```bash
chmod 600 ~/.trustssh/config.json
```

## Getting Config Values

From the Terraform directory:

```bash
terraform output -raw aws_region
terraform output -raw cognito_domain
terraform output -raw cognito_client_id
terraform output -raw callback_url
terraform output -raw api_base_url
terraform output -raw api_gateway_default_endpoint
terraform output -raw cli_config_url
```

Map them into `config.json`:

| Config field | Terraform output |
| --- | --- |
| `region` | `aws_region` |
| `cognito_domain` | `cognito_domain` |
| `client_id` | `cognito_client_id` |
| `redirect_uri` | `callback_url` |
| `api_base_url` | `api_base_url` |
| `default_duration_seconds` | `1800` |

`cli_config_url` is the static config document fetched by `trustssh configure`.

## Login Flow

Run:

```bash
trustssh login
```

The CLI will:

```text
1. Open Cognito managed login in your browser.
2. Receive the localhost callback on http://localhost:8765/callback.
3. Exchange the auth code for Cognito tokens using PKCE.
4. Save tokens to ~/.trustssh/tokens.json.
5. Create or reuse ~/.trustssh/id_ed25519.
6. Send ~/.trustssh/id_ed25519.pub to the TrustSSH API.
7. Save the returned certificate to ~/.trustssh/id_ed25519-cert.pub.
```

## Passkey Enrollment

To add a passkey, run:

```bash
trustssh passkeys add
```

This opens the Cognito managed login passkey enrollment page using the
client ID and redirect URI from ~/.trustssh/config.json.

Expected local files:

```text
~/.trustssh/config.json
~/.trustssh/tokens.json
~/.trustssh/id_ed25519
~/.trustssh/id_ed25519.pub
~/.trustssh/id_ed25519-cert.pub
```

Expected permissions:

```text
~/.trustssh                 0700
config.json                 0600
tokens.json                 0600
id_ed25519                  0600
id_ed25519.pub              0644
id_ed25519-cert.pub         0644
```

## Logout

Run:

```bash
trustssh logout
```

Logout removes:

```text
~/.trustssh/tokens.json
~/.trustssh/id_ed25519-cert.pub
```

Logout does not remove:

```text
~/.trustssh/id_ed25519
~/.trustssh/id_ed25519.pub
```

The private/public SSH key pair stays on the local machine for reuse.
