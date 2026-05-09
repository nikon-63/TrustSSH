data "archive_file" "add_package" {
  type        = "zip"
  source_dir  = "${var.source_dir}/add"
  output_path = "${path.module}/add_${var.add_function_name}.zip"
  excludes    = ["test_*.py", "__pycache__/*"]
}

resource "aws_lambda_function" "add" {
  function_name    = var.add_function_name
  role             = var.role_arn
  handler          = "app.handler"
  runtime          = "python3.12"
  filename         = data.archive_file.add_package.output_path
  source_code_hash = data.archive_file.add_package.output_base64sha256
  timeout          = 15
  memory_size      = 128

  environment {
    variables = {
      USER_POOL_ID         = var.cognito_user_pool_id
      USER_MAPPING_TABLE   = var.user_mapping_table_name
    }
  }

  tags = {
    Project = var.project_name
  }

  depends_on = [aws_cloudwatch_log_group.add]
}

resource "aws_cloudwatch_log_group" "add" {
  name              = "/aws/lambda/${var.add_function_name}"
  retention_in_days = 30

  tags = {
    Project = var.project_name
  }
}

data "archive_file" "remove_package" {
  type        = "zip"
  source_dir  = "${var.source_dir}/remove"
  output_path = "${path.module}/remove_${var.remove_function_name}.zip"
  excludes    = ["test_*.py", "__pycache__/*"]
}

resource "aws_lambda_function" "remove" {
  function_name    = var.remove_function_name
  role             = var.role_arn
  handler          = "app.handler"
  runtime          = "python3.12"
  filename         = data.archive_file.remove_package.output_path
  source_code_hash = data.archive_file.remove_package.output_base64sha256
  timeout          = 15
  memory_size      = 128

  environment {
    variables = {
      USER_POOL_ID         = var.cognito_user_pool_id
      USER_MAPPING_TABLE   = var.user_mapping_table_name
    }
  }

  tags = {
    Project = var.project_name
  }

  depends_on = [aws_cloudwatch_log_group.remove]
}

resource "aws_cloudwatch_log_group" "remove" {
  name              = "/aws/lambda/${var.remove_function_name}"
  retention_in_days = 30

  tags = {
    Project = var.project_name
  }
}

