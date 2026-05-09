variable "project_name" {
  description = "Short project name used for tags."
  type        = string
}

variable "ca_private_key_parameter_name" {
  description = "Name of the SSM SecureString parameter that will hold the OpenSSH CA private key."
  type        = string
}

variable "ca_private_key_value" {
  description = "OpenSSH CA private key material to store in SSM SecureString."
  type        = string
  sensitive   = true
}

variable "ca_public_key_parameter_name" {
  description = "Name of the SSM String parameter that stores the OpenSSH CA public key."
  type        = string
}

variable "ca_public_key_value" {
  description = "OpenSSH CA public key material to store in SSM String."
  type        = string
}

variable "version_parameter_name" {
  description = "Name of the SSM String parameter that stores the TrustSSH version."
  type        = string
}

variable "version_value" {
  description = "TrustSSH version read from the root VERSION file."
  type        = string
}
