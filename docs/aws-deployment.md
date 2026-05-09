# AWS Deployment

This guide covers deploying the AWS side of TrustSSH.

TODO: add a diagram of the AWS architecture and how the components interact.

## Overview

This document contains the minimal, repeatable steps to build, authenticate, and deploy the Lambda-based signer and supporting Terraform resources for TrustSSH.

## Prerequisites

- AWS CLI installed and configured
- `terraform` installed
- Access to the target AWS account/hosted zone

## Steps

### 1) Build Lambda dependencies

Run the build script for the Lambda signer dependencies:

```bash
./lambda/signer/build_requirements.sh
```

### 2) Authenticate with AWS

Example flow for obtaining the `process` profile used for deployment:

```bash
aws login --profile signin
aws sts get-caller-identity --profile signin
aws configure set region eu-west-2 --profile process
aws configure set credential_process \
  "aws configure export-credentials --profile signin --format process" \
  --profile process
export AWS_PROFILE=process
export AWS_SDK_LOAD_CONFIG=1
export AWS_EC2_METADATA_DISABLED=true
```

### 3) Create an SSH CA key pair (local, for upload to SSM)

Generate an Ed25519 key pair to act as the OpenSSH CA:

```bash
ssh-keygen -t ed25519 -f /tmp/trustssh_ca -N "" -C "trustssh-ca"
```

### 4) Fill in Terraform variables

Edit `terraform.tfvars` (or create one from `terraform.tfvars.example`) with your environment values.

### 5) Deploy Terraform

Run the standard Terraform deployment from the `infra/` directory:

```bash
terraform init
terraform validate
terraform plan
terraform apply
```

Terraform creates two custom HTTPS domains:

- `trustssh.demo.com` for the TrustSSH API
- `auth.trustssh.demo.com` for Cognito managed login and passkeys

The Cognito auth domain uses an ACM certificate in `us-east-1`, because Cognito custom domains are backed by CloudFront. Terraform also creates the Route 53 validation and alias records.

Terraform also creates the default Cognito managed-login branding style for the CLI app client. Without that style, Cognito managed login can show `Login pages unavailable`.

### 6) Create a Cognito user

The Cognito pool is configured for passkey-capable sign-in:

- email one-time password
- WebAuthn/passkey
- password

Cognito currently requires `PASSWORD` to be present when `WEB_AUTHN` is enabled as a first auth factor. TrustSSH uses the custom auth domain as the WebAuthn relying party ID:

```text
auth.trustssh.demo.com
```

Create users as an administrator with an email address. Public self sign-up is disabled. After the user exists, use the managed login page to complete email verification and register a passkey. Record the user `sub` (Cognito user ID) for mapping.

### 7) Seed the DynamoDB user mapping table

Write a mapping item for the Cognito `sub` that lists allowed SSH principals. Example DynamoDB item (JSON format for `aws dynamodb put-item --item`):

```json
{
  "cognito_sub": {
    "S": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  },
  "email": {
    "S": "example@example.com"
  },
  "enabled": {
    "BOOL": true
  },
  "ssh_principals": {
    "L": [
      {
        "S": "ubuntu"
      }
    ]
  },
  "max_duration_seconds": {
    "N": "1800"
  },
  "created_at": {
    "S": "2026-05-07T18:00:00Z"
  },
  "updated_at": {
    "S": "2026-05-07T18:00:00Z"
  }
}
```
