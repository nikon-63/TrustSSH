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

### 6) Create a Cognito user

Create a user in the deployed Cognito User Pool via the AWS Console. Record the user `sub` (Cognito user ID) for mapping.

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
