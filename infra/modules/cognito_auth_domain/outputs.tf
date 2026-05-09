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
