output "ca_private_key_parameter_name" {
  description = "Name of the CA private key SSM parameter."
  value       = aws_ssm_parameter.ca_private_key.name
}

output "ca_private_key_parameter_arn" {
  description = "ARN of the CA private key SSM parameter."
  value       = aws_ssm_parameter.ca_private_key.arn
}

output "ca_public_key_parameter_name" {
  description = "Name of the CA public key SSM parameter."
  value       = aws_ssm_parameter.ca_public_key.name
}

output "ca_public_key_parameter_arn" {
  description = "ARN of the CA public key SSM parameter."
  value       = aws_ssm_parameter.ca_public_key.arn
}
