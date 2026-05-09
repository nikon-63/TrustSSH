data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "signer" {
  name               = var.signer_role_name
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json

  tags = {
    Project = var.project_name
  }
}

resource "aws_iam_role_policy_attachment" "basic_execution" {
  role       = aws_iam_role.signer.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "signer" {
  statement {
    sid = "ReadUserMappings"
    actions = [
      "dynamodb:GetItem",
      "dynamodb:Query",
    ]
    resources = [var.user_mapping_table_arn]
  }

  statement {
    sid = "WriteAuditEvents"
    actions = [
      "dynamodb:PutItem",
    ]
    resources = [var.audit_events_table_arn]
  }

  statement {
    sid = "ReadCAPrivateKeyParameter"
    actions = [
      "ssm:GetParameter",
    ]
    resources = [var.ca_private_key_parameter_arn]
  }
}

resource "aws_iam_policy" "signer" {
  name        = "${var.signer_role_name}-policy"
  description = "Least-privilege policy for the TrustSSH signer Lambda."
  policy      = data.aws_iam_policy_document.signer.json

  tags = {
    Project = var.project_name
  }
}

resource "aws_iam_role_policy_attachment" "signer" {
  role       = aws_iam_role.signer.name
  policy_arn = aws_iam_policy.signer.arn
}

resource "aws_iam_role" "users" {
  name               = var.users_role_name
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json

  tags = {
    Project = var.project_name
  }
}

resource "aws_iam_role_policy_attachment" "users_basic_execution" {
  role       = aws_iam_role.users.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "users" {
  statement {
    sid = "ManageCognitoUsers"
    actions = [
      "cognito-idp:AdminCreateUser",
    ]
    resources = [var.cognito_user_pool_arn]
  }

  statement {
    sid = "WriteUserMappings"
    actions = [
      "dynamodb:PutItem",
    ]
    resources = [var.user_mapping_table_arn]
  }
}

resource "aws_iam_policy" "users" {
  name        = "${var.users_role_name}-policy"
  description = "Least-privilege policy for the TrustSSH users Lambda."
  policy      = data.aws_iam_policy_document.users.json

  tags = {
    Project = var.project_name
  }
}

resource "aws_iam_role_policy_attachment" "users" {
  role       = aws_iam_role.users.name
  policy_arn = aws_iam_policy.users.arn
}

