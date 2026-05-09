output "user_pool_id" {
  description = "ID of the Cognito User Pool."
  value       = aws_cognito_user_pool.this.id
}

output "user_pool_arn" {
  description = "ARN of the Cognito User Pool."
  value       = aws_cognito_user_pool.this.arn
}

output "app_client_id" {
  description = "Client ID for the CLI app client."
  value       = aws_cognito_user_pool_client.cli.id
}

output "domain_prefix" {
  description = "Configured Cognito prefix value. The deployed login domain is managed by the cognito_auth_domain module."
  value       = var.cognito_domain_prefix
}
