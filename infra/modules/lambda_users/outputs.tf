output "add_function_name" {
  description = "Name of the add users Lambda function."
  value       = aws_lambda_function.add.function_name
}

output "add_function_arn" {
  description = "ARN of the add users Lambda function."
  value       = aws_lambda_function.add.arn
}

output "add_invoke_arn" {
  description = "Invoke ARN of the add users Lambda function."
  value       = aws_lambda_function.add.invoke_arn
}
