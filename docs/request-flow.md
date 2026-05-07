# Request Flow

This document shows how data moves through TrustSSH during `trustssh login`, certificate issuance, and normal SSH login.

TrustSSH does not send the user's SSH private key to AWS. The CLI signs nothing locally and sends only the SSH public key plus a Cognito bearer token to the API.

## End-to-End View

```mermaid
flowchart LR
    User[User] --> CLI[TrustSSH CLI]
    CLI --> Browser[Browser]
    Browser --> Cognito[Cognito Hosted UI]
    Cognito --> CLI
    CLI --> API[API Gateway HTTP API]
    API --> Lambda[Signer Lambda]
    Lambda --> Mappings[(DynamoDB user mappings)]
    Lambda --> SSM[SSM SecureString CA private key]
    Lambda --> Audit[(DynamoDB audit events)]
    Lambda --> API
    API --> CLI
    CLI --> LocalFiles[Local TrustSSH files]
    User --> SSH[OpenSSH client]
    SSH --> Server[Linux SSH server]
```

## Login and Token Flow

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant CLI as TrustSSH CLI
    participant Browser
    participant Cognito as Cognito Hosted UI
    participant TokenEndpoint as Cognito token endpoint
    participant LocalFiles as ~/.trustssh

    User->>CLI: trustssh login
    CLI->>CLI: Load ~/.trustssh/config.json
    CLI->>CLI: Generate PKCE verifier, challenge, and OAuth state
    CLI->>CLI: Start localhost callback server
    CLI->>Browser: Open /oauth2/authorize URL
    Browser->>Cognito: GET /oauth2/authorize
    Cognito-->>Browser: Login and MFA challenge if configured
    User->>Cognito: Authenticate
    Cognito-->>Browser: Redirect to http://localhost:8765/callback?code=...&state=...
    Browser->>CLI: GET /callback with auth code and state
    CLI->>CLI: Verify returned state
    CLI->>TokenEndpoint: POST /oauth2/token with code and PKCE verifier
    TokenEndpoint-->>CLI: access_token, id_token, refresh_token
    CLI->>LocalFiles: Save tokens.json with 0600 permissions
```

The authorization request contains:

```text
response_type=code
client_id=<cognito client id>
redirect_uri=http://localhost:8765/callback
scope=openid email profile
code_challenge=<PKCE challenge>
code_challenge_method=S256
state=<random state>
```

The token exchange request is form-encoded:

```text
grant_type=authorization_code
client_id=<cognito client id>
code=<authorization code>
redirect_uri=http://localhost:8765/callback
code_verifier=<PKCE verifier>
```

Token response:

```json
{
  "access_token": "...",
  "id_token": "...",
  "refresh_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

## Local SSH Key Preparation

```mermaid
flowchart TD
    Start([After Cognito auth succeeds]) --> EnsureDir[Ensure ~/.trustssh exists]
    EnsureDir --> DirPerms[Set directory permissions to 0700]
    DirPerms --> HasPrivate{Does id_ed25519 exist?}
    HasPrivate -->|No| HasOrphanPublic{Does id_ed25519.pub exist?}
    HasOrphanPublic -->|Yes| Stop[Stop with error]
    HasOrphanPublic -->|No| Generate[Run ssh-keygen -t ed25519]
    HasPrivate -->|Yes| PrivatePerms[Set private key permissions to 0600]
    Generate --> PrivatePerms
    PrivatePerms --> HasPublic{Does id_ed25519.pub exist?}
    HasPublic -->|No| DerivePublic[Run ssh-keygen -y]
    HasPublic -->|Yes| PublicPerms[Set public key permissions to 0644]
    DerivePublic --> PublicPerms
    PublicPerms --> Done([Public key ready])
```

Local files:

```text
~/.trustssh/config.json
~/.trustssh/tokens.json
~/.trustssh/id_ed25519
~/.trustssh/id_ed25519.pub
~/.trustssh/id_ed25519-cert.pub
```

Permissions:

```text
~/.trustssh                 0700
id_ed25519                  0600
id_ed25519.pub              0644
id_ed25519-cert.pub         0644
tokens.json                 0600
```

## Certificate Issuance Flow

```mermaid
sequenceDiagram
    autonumber
    participant CLI as TrustSSH CLI
    participant API as API Gateway HTTP API
    participant Authorizer as Cognito JWT authorizer
    participant Lambda as Signer Lambda
    participant Mappings as DynamoDB user mappings
    participant SSM as SSM Parameter Store
    participant Audit as DynamoDB audit events
    participant LocalFiles as ~/.trustssh

    CLI->>CLI: Read ~/.trustssh/id_ed25519.pub
    CLI->>API: POST /issue-cert with Authorization bearer token
    API->>Authorizer: Validate JWT issuer and audience
    Authorizer-->>API: Cognito claims including sub and email
    API->>Lambda: Invoke with claims and request body
    Lambda->>Lambda: Validate request body and SSH public key
    Lambda->>Mappings: GetItem by cognito_sub
    Mappings-->>Lambda: enabled, ssh_principals, max_duration_seconds
    Lambda->>Lambda: Enforce enabled mapping and duration cap
    Lambda->>SSM: GetParameter WithDecryption for CA private key
    SSM-->>Lambda: OpenSSH CA private key
    Lambda->>Lambda: Build and sign OpenSSH user certificate
    Lambda->>Audit: PutItem issuance audit event
    Lambda-->>API: certificate, valid_until, principals, serial
    API-->>CLI: 200 OK JSON response
    CLI->>LocalFiles: Save id_ed25519-cert.pub
```

CLI request:

```http
POST /issue-cert
Authorization: Bearer <cognito access token>
Content-Type: application/json
```

```json
{
  "public_key": "ssh-ed25519 AAAA... trustssh",
  "requested_duration_seconds": 1800
}
```

API Gateway checks:

```text
issuer   = https://cognito-idp.<region>.amazonaws.com/<user-pool-id>
audience = <cognito app client id>
```

Lambda uses these JWT claims:

```json
{
  "sub": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "email": "user@example.com"
}
```

Mapping lookup key:

```json
{
  "cognito_sub": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}
```

Expected mapping item:

```json
{
  "cognito_sub": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "email": "user@example.com",
  "enabled": true,
  "ssh_principals": ["ubuntu"],
  "max_duration_seconds": 1800
}
```

Successful response:

```json
{
  "certificate": "ssh-ed25519-cert-v01@openssh.com AAAA...",
  "valid_until": "2026-05-07T18:30:00Z",
  "principals": ["ubuntu"],
  "serial": 12345
}
```

The certificate is saved to:

```text
~/.trustssh/id_ed25519-cert.pub
```

## Lambda Decision Logic

```mermaid
flowchart TD
    Start([Lambda invoked]) --> Claims{Cognito sub present?}
    Claims -->|No| Unauth[Return 401 unauthenticated]
    Claims -->|Yes| Body{Valid JSON body?}
    Body -->|No| BadRequest[Return 400 invalid_request]
    Body -->|Yes| Key{Valid ssh-ed25519 public key?}
    Key -->|No| BadRequest
    Key -->|Yes| Mapping[Read DynamoDB mapping by cognito_sub]
    Mapping --> Enabled{Mapping enabled with principals?}
    Enabled -->|No| Deny[Write denied audit event and return 403]
    Enabled -->|Yes| Duration{Duration within mapping and global caps?}
    Duration -->|No| BadRequest
    Duration -->|Yes| CA[Load CA private key from SSM]
    CA --> Sign[Sign OpenSSH user certificate]
    Sign --> Audit[Write issued audit event]
    Audit --> Return[Return certificate JSON]
```

## Audit Event

The Lambda writes one audit event for successful certificate issuance. It also writes denied events when the Cognito user has no enabled mapping.

Audit table keys:

```text
cognito_sub
issued_at_serial
```

Issued event:

```json
{
  "cognito_sub": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "issued_at_serial": "2026-05-07T18:00:00Z#12345",
  "email": "user@example.com",
  "principals": ["ubuntu"],
  "serial": 12345,
  "valid_until": "2026-05-07T18:30:00Z",
  "outcome": "issued"
}
```

Denied event:

```json
{
  "cognito_sub": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "issued_at_serial": "2026-05-07T18:00:00Z#0",
  "email": "user@example.com",
  "principals": [],
  "serial": 0,
  "valid_until": "",
  "outcome": "denied"
}
```

## Normal SSH Login

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant SSH as OpenSSH client
    participant LocalFiles as ~/.trustssh
    participant Server as Linux SSH server

    User->>SSH: ssh ubuntu@server
    SSH->>LocalFiles: Read id_ed25519
    SSH->>LocalFiles: Read id_ed25519-cert.pub
    SSH->>Server: Present public key certificate
    Server->>Server: Verify cert signed by TrustedUserCAKeys
    Server->>Server: Check cert validity window
    Server->>Server: Check ubuntu principal in AuthorizedPrincipalsFile
    Server-->>SSH: Accept or reject login
```

The CLI does not run SSH. OpenSSH automatically pairs:

```text
~/.trustssh/id_ed25519
~/.trustssh/id_ed25519-cert.pub
```

## Data Handling Summary

| Data | Source | Destination | Notes |
| --- | --- | --- | --- |
| PKCE verifier | CLI | Cognito token endpoint | Never sent in the browser authorization request |
| Authorization code | Cognito | CLI localhost callback | Exchanged once for tokens |
| Access token | Cognito | CLI, API Gateway | Used as bearer token for `/issue-cert` |
| ID token | Cognito | CLI local token file | Decoded locally for display |
| SSH private key | CLI local machine | Local filesystem only | Never sent to AWS |
| SSH public key | CLI local machine | Signer Lambda | Sent in `/issue-cert` request |
| CA private key | `terraform.tfvars` then SSM | Signer Lambda | Stored as SSM SecureString and present in Terraform state |
| CA public key | `terraform.tfvars` then SSM | Server setup/reference | Lambda does not need this value |
| SSH certificate | Signer Lambda | CLI local filesystem | Short-lived OpenSSH user cert |
| Audit event | Signer Lambda | DynamoDB audit table | No tokens or private keys |
