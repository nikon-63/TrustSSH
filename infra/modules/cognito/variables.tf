variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "user_pool_name" {
  description = "Name of the Cognito User Pool."
  type        = string
}

variable "app_client_name" {
  description = "Name of the Cognito app client."
  type        = string
}

variable "cognito_domain_prefix" {
  description = "Hosted UI domain prefix. Must be unique within the AWS region."
  type        = string
}

variable "callback_urls" {
  description = "Allowed OAuth callback URLs."
  type        = list(string)
}

variable "logout_urls" {
  description = "Allowed OAuth logout URLs."
  type        = list(string)
}

variable "access_token_validity" {
  description = "Access token lifetime in minutes."
  type        = number
}

variable "id_token_validity" {
  description = "ID token lifetime in minutes."
  type        = number
}

variable "refresh_token_validity" {
  description = "Refresh token lifetime in days."
  type        = number
}
