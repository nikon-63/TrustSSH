output "domain_name" {
  description = "Custom API domain name."
  value       = aws_apigatewayv2_domain_name.api.domain_name
}

output "api_base_url" {
  description = "Custom API base URL."
  value       = "https://${aws_apigatewayv2_domain_name.api.domain_name}"
}

output "certificate_arn" {
  description = "ARN of the ACM certificate used by the API custom domain."
  value       = aws_acm_certificate_validation.api.certificate_arn
}
