variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "add_function_name" {
  description = "Name of the add users Lambda function."
  type        = string
}

variable "remove_function_name" {
  description = "Name of the remove users Lambda function."
  type        = string
}

variable "source_dir" {
  description = "Base path to the Lambda users deployment package source."
  type        = string
}

variable "role_arn" {
  description = "ARN of the IAM role for the Lambda functions."
  type        = string
}

variable "cognito_user_pool_id" {
  description = "ID of the Cognito user pool."
  type        = string
}

variable "user_mapping_table_name" {
  description = "Name of the DynamoDB user mapping table."
  type        = string
}
