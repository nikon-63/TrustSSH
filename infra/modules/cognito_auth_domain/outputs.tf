output "auth_domain" {
  description = "Custom Cognito auth domain name."
  value       = var.domain_name
}

output "auth_domain_url" {
  description = "Custom Cognito auth domain URL."
  value       = "https://${var.domain_name}"
}

output "cloudfront_distribution" {
  description = "CloudFront distribution backing the custom Cognito auth domain."
  value       = aws_cognito_user_pool_domain.auth.cloudfront_distribution
}

output "hosted_ui_login_url" {
  description = "Hosted UI login URL for the CLI app client."
  value       = "https://${var.domain_name}/oauth2/authorize?client_id=${var.client_id}&response_type=code&scope=openid+email+profile&redirect_uri=${urlencode(var.callback_urls[0])}"
}

output "domain_name" {
  description = "Custom Cognito auth domain name."
  value       = aws_cognito_user_pool_domain.auth.domain
}

output "hosted_ui_domain" {
  description = "Managed login base URL."
  value       = "https://${aws_cognito_user_pool_domain.auth.domain}"
}

output "certificate_arn" {
  description = "ARN of the ACM certificate used by the Cognito custom auth domain."
  value       = aws_acm_certificate_validation.auth.certificate_arn
}

output "webauthn_relying_party_id" {
  description = "WebAuthn relying party ID for passkeys."
  value       = var.domain_name
}
