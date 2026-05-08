variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "bucket_name" {
  description = "Name of the S3 bucket that stores static public content."
  type        = string
}

variable "api_id" {
  description = "ID of the API Gateway HTTP API."
  type        = string
}

variable "config_json" {
  description = "Rendered TrustSSH CLI config JSON."
  type        = string
}

variable "public_key" {
  description = "OpenSSH CA public key to publish as static text."
  type        = string
}
