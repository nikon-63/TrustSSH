variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "aws_region" {
  description = "AWS region where the Cognito User Pool is deployed."
  type        = string
}

variable "hosted_zone_id" {
  description = "Route 53 hosted zone ID for the Cognito custom auth domain."
  type        = string
}

variable "domain_name" {
  description = "Fully qualified Cognito custom auth domain name."
  type        = string
}

variable "user_pool_id" {
  description = "Cognito User Pool ID."
  type        = string
}

variable "client_id" {
  description = "Cognito app client ID that should receive a managed login branding style."
  type        = string
}

variable "callback_urls" {
  description = "Allowed OAuth callback URLs for hosted UI links."
  type        = list(string)
}

variable "webauthn_user_verification" {
  description = "WebAuthn user verification setting."
  type        = string
  default     = "required"

  validation {
    condition     = contains(["required", "preferred"], var.webauthn_user_verification)
    error_message = "webauthn_user_verification must be either required or preferred."
  }
}
