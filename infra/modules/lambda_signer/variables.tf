variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "function_name" {
  description = "Name of the signer Lambda function."
  type        = string
}

variable "source_dir" {
  description = "Path to the Lambda source directory."
  type        = string
}

variable "role_arn" {
  description = "IAM role ARN used by the Lambda function."
  type        = string
}

variable "user_mapping_table_name" {
  description = "Name of the DynamoDB user mapping table."
  type        = string
}

variable "audit_events_table_name" {
  description = "Name of the DynamoDB audit events table."
  type        = string
}

variable "ca_private_key_parameter_name" {
  description = "Name of the SSM parameter containing the CA private key."
  type        = string
}

variable "default_duration_seconds" {
  description = "Default certificate duration."
  type        = number
}

variable "max_duration_seconds" {
  description = "Maximum certificate duration."
  type        = number
}
