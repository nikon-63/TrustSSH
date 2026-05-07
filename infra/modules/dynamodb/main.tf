resource "aws_dynamodb_table" "user_mappings" {
  name         = var.user_mapping_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "cognito_sub"

  deletion_protection_enabled = var.deletion_protection

  attribute {
    name = "cognito_sub"
    type = "S"
  }

  point_in_time_recovery {
    enabled = var.point_in_time_recovery
  }

  server_side_encryption {
    enabled = true
  }

  tags = {
    Project = var.project_name
    Purpose = "trustssh-user-mappings"
  }
}

resource "aws_dynamodb_table" "audit_events" {
  name         = var.audit_events_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "cognito_sub"
  range_key    = "issued_at_serial"

  deletion_protection_enabled = var.deletion_protection

  attribute {
    name = "cognito_sub"
    type = "S"
  }

  attribute {
    name = "issued_at_serial"
    type = "S"
  }

  point_in_time_recovery {
    enabled = var.point_in_time_recovery
  }

  server_side_encryption {
    enabled = true
  }

  tags = {
    Project = var.project_name
    Purpose = "trustssh-audit-events"
  }
}
