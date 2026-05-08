output "bucket_name" {
  description = "Name of the config bucket."
  value       = aws_s3_bucket.config.bucket
}

output "route_key" {
  description = "API Gateway route key for the static config JSON."
  value       = aws_apigatewayv2_route.config.route_key
}

output "public_key_route_key" {
  description = "API Gateway route key for the static CA public key."
  value       = aws_apigatewayv2_route.public_key.route_key
}
