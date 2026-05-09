#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INFRA_DIR="$ROOT_DIR/infra"

export AWS_EC2_METADATA_DISABLED=true
export TF_IN_AUTOMATION=true

export TF_DATA_DIR="$ROOT_DIR/.terraform-test"

echo "Running TrustSSH Terraform checks"
echo

cd "$INFRA_DIR"

echo "Initializing Terraform without backend..."
terraform init \
    -backend=false \
    -input=false \
    -no-color

echo
echo "Checking Terraform formatting..."
if ! terraform fmt \
    -check \
    -recursive \
    -no-color; then

    echo
    echo "Terraform formatting check failed."
    exit 1
fi

echo
echo "Validating Terraform configuration..."
terraform validate \
    -no-color

echo
echo "All Terraform checks passed"