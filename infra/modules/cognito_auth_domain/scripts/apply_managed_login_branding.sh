#!/usr/bin/env bash

set -euo pipefail

region="$1"
user_pool_id="$2"
client_id="$3"
settings_file="$4"
logo_file="$5"

branding_id="$(
  aws cognito-idp describe-managed-login-branding-by-client \
    --region "$region" \
    --user-pool-id "$user_pool_id" \
    --client-id "$client_id" \
    --query 'ManagedLoginBranding.ManagedLoginBrandingId' \
    --output text 2>/dev/null || true
)"

if [[ -z "$branding_id" || "$branding_id" == "None" ]]; then
  aws cognito-idp create-managed-login-branding \
    --region "$region" \
    --user-pool-id "$user_pool_id" \
    --client-id "$client_id" \
    --use-cognito-provided-values \
    >/dev/null

  branding_id="$(
    aws cognito-idp describe-managed-login-branding-by-client \
      --region "$region" \
      --user-pool-id "$user_pool_id" \
      --client-id "$client_id" \
      --query 'ManagedLoginBranding.ManagedLoginBrandingId' \
      --output text
  )"
fi

request_file="$(mktemp)"
trap 'rm -f "$request_file"' EXIT

jq -n \
  --arg user_pool_id "$user_pool_id" \
  --arg branding_id "$branding_id" \
  --argjson settings "$(jq -c . "$settings_file")" \
  --rawfile logo_base64 <(base64 < "$logo_file" | tr -d '\n') \
  '{
    UserPoolId: $user_pool_id,
    ManagedLoginBrandingId: $branding_id,
    Settings: $settings,
    Assets: [
      {
        Category: "FORM_LOGO",
        ColorMode: "LIGHT",
        Extension: "PNG",
        Bytes: $logo_base64
      }
    ]
  }' > "$request_file"

echo "Updating Cognito Managed Login branding with LIGHT form logo..."
aws cognito-idp update-managed-login-branding \
  --region "$region" \
  --cli-input-json "file://$request_file" \
  >/dev/null
