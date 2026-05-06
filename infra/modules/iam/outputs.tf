output "signer_role_arn" {
  description = "ARN of the signer Lambda IAM role."
  value       = aws_iam_role.signer.arn
}

output "signer_role_name" {
  description = "Name of the signer Lambda IAM role."
  value       = aws_iam_role.signer.name
}
