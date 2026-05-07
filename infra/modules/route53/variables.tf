variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "hosted_zone_id" {
  description = "Route 53 hosted zone ID for the API custom domain."
  type        = string
}

variable "domain_name" {
  description = "Fully qualified API custom domain name."
  type        = string
}

variable "api_id" {
  description = "ID of the API Gateway HTTP API."
  type        = string
}

variable "stage_name" {
  description = "API Gateway stage name to map to the custom domain."
  type        = string
}
