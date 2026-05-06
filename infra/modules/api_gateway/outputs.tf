output "api_id" {
  description = "ID of the HTTP API."
  value       = aws_apigatewayv2_api.this.id
}

output "api_endpoint" {
  description = "Base endpoint for the HTTP API."
  value       = aws_apigatewayv2_api.this.api_endpoint
}

output "issue_certificate_route" {
  description = "Route key for certificate issuance."
  value       = aws_apigatewayv2_route.issue_cert.route_key
}
