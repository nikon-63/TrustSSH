data "archive_file" "package" {
  type        = "zip"
  source_dir  = var.source_dir
  output_path = "${path.module}/${var.function_name}.zip"
}

resource "aws_lambda_function" "this" {
  function_name    = var.function_name
  role             = var.role_arn
  handler          = "app.handler"
  runtime          = "python3.12"
  filename         = data.archive_file.package.output_path
  source_code_hash = data.archive_file.package.output_base64sha256
  timeout          = 15
  memory_size      = 128

  environment {
    variables = {
      USER_MAPPING_TABLE            = var.user_mapping_table_name
      AUDIT_EVENTS_TABLE            = var.audit_events_table_name
      CA_PRIVATE_KEY_PARAMETER_NAME = var.ca_private_key_parameter_name
      DEFAULT_DURATION_SECONDS      = tostring(var.default_duration_seconds)
      MAX_DURATION_SECONDS          = tostring(var.max_duration_seconds)
    }
  }

  tags = {
    Project = var.project_name
  }
}
