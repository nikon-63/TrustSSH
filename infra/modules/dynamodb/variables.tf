variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "user_mapping_table_name" {
  description = "Name of the DynamoDB table that maps Cognito subjects to SSH principals."
  type        = string
}

variable "audit_events_table_name" {
  description = "Name of the DynamoDB table that stores certificate issuance audit events."
  type        = string
}

variable "deletion_protection" {
  description = "Whether deletion protection is enabled for the DynamoDB tables."
  type        = bool
}

variable "point_in_time_recovery" {
  description = "Whether point-in-time recovery is enabled for the DynamoDB tables."
  type        = bool
}
