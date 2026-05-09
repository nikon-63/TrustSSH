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

output "remove_function_name" {
  description = "Name of the remove users Lambda function."
  value       = aws_lambda_function.remove.function_name
}

output "remove_function_arn" {
  description = "ARN of the remove users Lambda function."
  value       = aws_lambda_function.remove.arn
}

output "remove_invoke_arn" {
  description = "Invoke ARN of the remove users Lambda function."
  value       = aws_lambda_function.remove.invoke_arn
}

