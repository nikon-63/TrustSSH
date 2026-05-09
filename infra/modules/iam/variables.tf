variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "signer_role_name" {
  description = "Name of the IAM role used by the signer Lambda."
  type        = string
}

variable "users_role_name" {
  description = "Name of the IAM role used by the users Lambda."
  type        = string
}

variable "cognito_user_pool_arn" {
  description = "ARN of the Cognito User Pool."
  type        = string
}

variable "user_mapping_table_arn" {
  description = "ARN of the DynamoDB user mapping table."
  type        = string
}

variable "audit_events_table_arn" {
  description = "ARN of the DynamoDB audit events table."
  type        = string
}

variable "ca_private_key_parameter_arn" {
  description = "ARN of the SSM parameter containing the CA private key."
  type        = string
}
