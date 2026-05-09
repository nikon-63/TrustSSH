output "signer_role_arn" {
  description = "ARN of the signer Lambda IAM role."
  value       = aws_iam_role.signer.arn
}

output "signer_role_name" {
  description = "Name of the signer Lambda IAM role."
  value       = aws_iam_role.signer.name
}

output "users_role_arn" {
  description = "ARN of the users Lambda IAM role."
  value       = aws_iam_role.users.arn
}

output "users_role_name" {
  description = "Name of the users Lambda IAM role."
  value       = aws_iam_role.users.name
}

