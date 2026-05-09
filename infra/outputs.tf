output "aws_region" {
  description = "AWS region used by the deployment."
  value       = var.aws_region
}

output "callback_url" {
  description = "OAuth callback URL for the TrustSSH CLI."
  value       = var.callback_url
}

output "cognito_user_pool_id" {
  description = "ID of the TrustSSH Cognito User Pool."
  value       = module.cognito.user_pool_id
}

output "cognito_user_pool_arn" {
  description = "ARN of the TrustSSH Cognito User Pool."
  value       = module.cognito.user_pool_arn
}

output "cognito_client_id" {
  description = "Client ID for the TrustSSH CLI Cognito app client."
  value       = module.cognito.user_pool_client_id
}

output "cognito_domain" {
  description = "Managed login base URL for Cognito authentication."
  value       = module.cognito_auth_domain.auth_domain_url
}

output "api_base_url" {
  description = "Custom API base URL."
  value       = module.route53.api_base_url
}

output "api_gateway_default_endpoint" {
  description = "Default execute-api endpoint for the HTTP API."
  value       = module.api_gateway.api_endpoint
}

output "dynamodb_user_mapping_table" {
  description = "DynamoDB table for Cognito-user-to-SSH-principal mappings."
  value       = module.dynamodb.user_mapping_table_name
}

output "dynamodb_user_mapping_table_arn" {
  description = "ARN of the DynamoDB user mapping table."
  value       = module.dynamodb.user_mapping_table_arn
}

output "dynamodb_audit_table" {
  description = "DynamoDB table for certificate issuance audit events."
  value       = module.dynamodb.audit_events_table_name
}

output "dynamodb_audit_table_arn" {
  description = "ARN of the DynamoDB audit events table."
  value       = module.dynamodb.audit_events_table_arn
}

output "issue_certificate_url" {
  description = "URL for the certificate issuance route."
  value       = "${module.route53.api_base_url}${var.issue_certificate_route_path}"
}

output "cli_config_url" {
  description = "URL for the static TrustSSH CLI configuration document."
  value       = "${module.route53.api_base_url}/config.json"
}

output "static_content_bucket" {
  description = "S3 bucket storing static public TrustSSH content."
  value       = module.static_config.bucket_name
}

output "ca_public_key_url" {
  description = "URL for the public OpenSSH CA public key."
  value       = "${module.route53.api_base_url}/public_key.txt"
}

output "signer_lambda_function_name" {
  description = "Name of the signer Lambda function."
  value       = module.lambda_signer.function_name
}

output "ca_private_key_parameter_name" {
  description = "SSM parameter name for the OpenSSH CA private key."
  value       = module.ssm.ca_private_key_parameter_name
}

output "ca_public_key_parameter_name" {
  description = "SSM parameter name for the OpenSSH CA public key."
  value       = module.ssm.ca_public_key_parameter_name
}

output "version" {
  description = "TrustSSH version read from the root VERSION file."
  value       = data.external.version.result.version
}

output "version_parameter_name" {
  description = "SSM parameter name for the TrustSSH version."
  value       = module.ssm.version_parameter_name
}

output "api_custom_domain_name" {
  description = "Custom domain name for the TrustSSH API."
  value       = module.route53.domain_name
}

output "cognito_custom_domain_name" {
  description = "Custom domain name for Cognito managed login."
  value       = module.cognito_auth_domain.auth_domain
}

output "webauthn_relying_party_id" {
  description = "WebAuthn relying party ID used for Cognito passkeys."
  value       = module.cognito_auth_domain.webauthn_relying_party_id
}
