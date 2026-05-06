output "user_mapping_table_name" {
  description = "Name of the user mapping table."
  value       = aws_dynamodb_table.user_mappings.name
}

output "user_mapping_table_arn" {
  description = "ARN of the user mapping table."
  value       = aws_dynamodb_table.user_mappings.arn
}

output "audit_events_table_name" {
  description = "Name of the audit events table."
  value       = aws_dynamodb_table.audit_events.name
}

output "audit_events_table_arn" {
  description = "ARN of the audit events table."
  value       = aws_dynamodb_table.audit_events.arn
}
