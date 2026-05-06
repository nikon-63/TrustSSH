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
  description = "API base URL"
  value       = "https://${var.api_subdomain}.${var.base_domain}"
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
