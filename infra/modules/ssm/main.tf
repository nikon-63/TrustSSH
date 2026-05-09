resource "aws_ssm_parameter" "ca_private_key" {
  name        = var.ca_private_key_parameter_name
  description = "OpenSSH CA private key for TrustSSH certificate signing."
  type        = "SecureString"
  value       = var.ca_private_key_value

  tags = {
    Project = var.project_name
    Purpose = "trustssh-ca-private-key"
  }
}

resource "aws_ssm_parameter" "ca_public_key" {
  name        = var.ca_public_key_parameter_name
  description = "OpenSSH CA public key for TrustSSH server trust configuration."
  type        = "String"
  value       = var.ca_public_key_value

  tags = {
    Project = var.project_name
    Purpose = "trustssh-ca-public-key"
  }
}

resource "aws_ssm_parameter" "version" {
  name        = var.version_parameter_name
  description = "TrustSSH application version."
  type        = "String"
  value       = var.version_value

  tags = {
    Project = var.project_name
    Purpose = "trustssh-version"
  }
}
