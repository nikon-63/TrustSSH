variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "api_name" {
  description = "Name of the HTTP API."
  type        = string
}

variable "route_path" {
  description = "Route path for certificate issuance."
  type        = string
}

variable "cognito_user_pool_id" {
  description = "Cognito User Pool ID used as the JWT issuer."
  type        = string
}

variable "cognito_client_id" {
  description = "Cognito app client ID used as the JWT audience."
  type        = string
}

variable "signer_function_name" {
  description = "Name of the signer Lambda function."
  type        = string
}

variable "signer_invoke_arn" {
  description = "Invoke ARN of the signer Lambda function."
  type        = string
}
