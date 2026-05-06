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
  default     = ""
}

variable "base_domain" {
  description = "Base domain for TrustSSH endpoints."
  type        = string
  default     = ""
}

variable "api_subdomain" {
  description = "Subdomain label for the TrustSSH API endpoint."
  type        = string
  default     = "trustssh"
}
