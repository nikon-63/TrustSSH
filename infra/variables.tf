variable "aws_region" {
  description = "AWS region for TrustSSH resources."
  type        = string
  default     = "eu-west-2"
}

variable "project_name" {
  description = "Short project name used for tagging and resource names."
  type        = string
  default     = "trustssh"
}

variable "cognito_user_pool_name" {
  description = "Name of the Cognito User Pool."
  type        = string
  default     = "trustssh-users"
}

variable "cognito_app_client_name" {
  description = "Name of the Cognito public/native app client used by the CLI."
  type        = string
  default     = "trustssh-cli"
}

variable "cognito_domain_prefix" {
  description = "Cognito Hosted UI domain prefix. This must be globally unique within the AWS region."
  type        = string
  default     = "trustssh-cli-dev"
}

variable "callback_url" {
  description = "Loopback callback URL used by the CLI OAuth flow."
  type        = string
  default     = "http://localhost:8765/callback"
}

variable "logout_url" {
  description = "Loopback logout URL for the CLI app client."
  type        = string
  default     = "http://localhost:8765/logout"
}

variable "cognito_access_token_validity_minutes" {
  description = "Cognito access token lifetime in minutes."
  type        = number
  default     = 60
}

variable "cognito_id_token_validity_minutes" {
  description = "Cognito ID token lifetime in minutes."
  type        = number
  default     = 60
}

variable "cognito_refresh_token_validity_days" {
  description = "Cognito refresh token lifetime in days."
  type        = number
  default     = 30
}

variable "hosted_zone_id" {
  description = "Existing Route 53 hosted zone ID."
  type        = string
}

variable "base_domain" {
  description = "Base domain for TrustSSH endpoints."
  type        = string
}

variable "api_subdomain" {
  description = "Subdomain label for the TrustSSH API endpoint."
  type        = string
  default     = "trustssh"
}

variable "dynamodb_user_mapping_table_name" {
  description = "Name of the DynamoDB table that maps Cognito users to allowed SSH principals."
  type        = string
  default     = "trustssh-user-mappings"
}

variable "dynamodb_audit_events_table_name" {
  description = "Name of the DynamoDB table that will store certificate issuance audit events."
  type        = string
  default     = "trustssh-audit-events"
}

variable "dynamodb_deletion_protection" {
  description = "Whether DynamoDB deletion protection is enabled for TrustSSH tables."
  type        = bool
  default     = false
}

variable "dynamodb_point_in_time_recovery" {
  description = "Whether DynamoDB point-in-time recovery is enabled for TrustSSH tables."
  type        = bool
  default     = false
}

variable "ca_private_key_parameter_name" {
  description = "Name of the SSM SecureString parameter that will hold the OpenSSH CA private key."
  type        = string
  default     = "/trustssh/ca/private-key"
}

variable "ca_private_key_value" {
  description = "OpenSSH CA private key material to store in SSM SecureString. Set this in terraform.tfvars only."
  type        = string
  sensitive   = true
}

variable "ca_public_key_parameter_name" {
  description = "Name of the SSM String parameter that stores the OpenSSH CA public key for reference."
  type        = string
  default     = "/trustssh/ca/public-key"
}

variable "ca_public_key_value" {
  description = "OpenSSH CA public key material to store in SSM String."
  type        = string
}

variable "signer_lambda_role_name" {
  description = "Name of the IAM role used by the signer Lambda."
  type        = string
  default     = "trustssh-signer-role"
}

variable "signer_lambda_function_name" {
  description = "Name of the signer Lambda function."
  type        = string
  default     = "trustssh-signer"
}

variable "api_gateway_name" {
  description = "Name of the TrustSSH HTTP API."
  type        = string
  default     = "trustssh-api"
}

variable "issue_certificate_route_path" {
  description = "HTTP API route path for requesting a certificate."
  type        = string
  default     = "/issue-cert"
}

variable "default_certificate_duration_seconds" {
  description = "Default certificate lifetime requested by the CLI."
  type        = number
  default     = 1800
}

variable "max_certificate_duration_seconds" {
  description = "Global maximum certificate lifetime for the MVP."
  type        = number
  default     = 14400
}
