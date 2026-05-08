resource "aws_s3_bucket" "config" {
  bucket        = var.bucket_name
  force_destroy = true

  tags = {
    Project = var.project_name
  }
}

resource "aws_s3_bucket_ownership_controls" "config" {
  bucket = aws_s3_bucket.config.id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_s3_bucket_public_access_block" "config" {
  bucket = aws_s3_bucket.config.id

  block_public_acls       = true
  block_public_policy     = false
  ignore_public_acls      = true
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "config" {
  bucket = aws_s3_bucket.config.id
  policy = data.aws_iam_policy_document.config.json

  depends_on = [aws_s3_bucket_public_access_block.config]
}

data "aws_iam_policy_document" "config" {
  statement {
    sid = "AllowPublicReadStaticContent"
    principals {
      type        = "*"
      identifiers = ["*"]
    }
    actions = ["s3:GetObject"]
    resources = [
      "${aws_s3_bucket.config.arn}/config.json",
      "${aws_s3_bucket.config.arn}/public_key.txt",
    ]
  }
}

resource "aws_s3_object" "config" {
  bucket       = aws_s3_bucket.config.id
  key          = "config.json"
  content      = var.config_json
  content_type = "application/json"
  etag         = md5(var.config_json)

  tags = {
    Project = var.project_name
  }
}

resource "aws_s3_object" "public_key" {
  bucket       = aws_s3_bucket.config.id
  key          = "public_key.txt"
  content      = trimspace(var.public_key)
  content_type = "text/plain"
  etag         = md5(trimspace(var.public_key))

  tags = {
    Project = var.project_name
  }
}

resource "aws_apigatewayv2_integration" "config" {
  api_id             = var.api_id
  integration_type   = "HTTP_PROXY"
  integration_method = "GET"
  integration_uri    = "https://${aws_s3_bucket.config.bucket_regional_domain_name}/${aws_s3_object.config.key}"
}

resource "aws_apigatewayv2_route" "config" {
  api_id             = var.api_id
  route_key          = "GET /config.json"
  target             = "integrations/${aws_apigatewayv2_integration.config.id}"
  authorization_type = "NONE"
}

resource "aws_apigatewayv2_integration" "public_key" {
  api_id             = var.api_id
  integration_type   = "HTTP_PROXY"
  integration_method = "GET"
  integration_uri    = "https://${aws_s3_bucket.config.bucket_regional_domain_name}/${aws_s3_object.public_key.key}"
}

resource "aws_apigatewayv2_route" "public_key" {
  api_id             = var.api_id
  route_key          = "GET /public_key.txt"
  target             = "integrations/${aws_apigatewayv2_integration.public_key.id}"
  authorization_type = "NONE"
}
