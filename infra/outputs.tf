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
  value       = module.cognito.app_client_id
}

output "cognito_domain" {
  description = "Hosted UI base URL for Cognito authentication."
  value       = module.cognito.hosted_ui_domain
}

output "api_base_url" {
  description = "API Gateway base URL."
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
  value       = "${module.api_gateway.api_endpoint}${var.issue_certificate_route_path}"
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
