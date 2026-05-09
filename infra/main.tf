terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.4"
    }
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project = var.project_name
    }
  }
}

provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"

  default_tags {
    tags = {
      Project = var.project_name
    }
  }
}

module "cognito" {
  source = "./modules/cognito"

  project_name           = var.project_name
  user_pool_name         = var.cognito_user_pool_name
  app_client_name        = var.cognito_app_client_name
  cognito_domain_prefix  = var.cognito_domain_prefix
  callback_urls          = [var.callback_url]
  logout_urls            = [var.logout_url]
  access_token_validity  = var.cognito_access_token_validity_minutes
  id_token_validity      = var.cognito_id_token_validity_minutes
  refresh_token_validity = var.cognito_refresh_token_validity_days
}

module "cognito_auth_domain" {
  source = "./modules/cognito_auth_domain"

  providers = {
    aws           = aws
    aws.us_east_1 = aws.us_east_1
  }

  project_name               = var.project_name
  aws_region                 = var.aws_region
  hosted_zone_id             = var.hosted_zone_id
  domain_name                = "${var.auth_subdomain}.${var.base_domain}"
  user_pool_id               = module.cognito.user_pool_id
  client_id                  = module.cognito.user_pool_client_id
  callback_urls              = [var.callback_url]
  webauthn_user_verification = "required"
}

module "dynamodb" {
  source = "./modules/dynamodb"

  project_name            = var.project_name
  user_mapping_table_name = var.dynamodb_user_mapping_table_name
  audit_events_table_name = var.dynamodb_audit_events_table_name
  deletion_protection     = var.dynamodb_deletion_protection
  point_in_time_recovery  = var.dynamodb_point_in_time_recovery
}

module "ssm" {
  source = "./modules/ssm"

  project_name                  = var.project_name
  ca_private_key_parameter_name = var.ca_private_key_parameter_name
  ca_private_key_value          = var.ca_private_key_value
  ca_public_key_parameter_name  = var.ca_public_key_parameter_name
  ca_public_key_value           = var.ca_public_key_value
}

module "iam" {
  source = "./modules/iam"

  project_name                 = var.project_name
  signer_role_name             = var.signer_lambda_role_name
  user_mapping_table_arn       = module.dynamodb.user_mapping_table_arn
  audit_events_table_arn       = module.dynamodb.audit_events_table_arn
  ca_private_key_parameter_arn = module.ssm.ca_private_key_parameter_arn
}

module "lambda_signer" {
  source = "./modules/lambda_signer"

  project_name                  = var.project_name
  function_name                 = var.signer_lambda_function_name
  source_dir                    = "${path.module}/../lambda/signer/build/package"
  role_arn                      = module.iam.signer_role_arn
  user_mapping_table_name       = module.dynamodb.user_mapping_table_name
  audit_events_table_name       = module.dynamodb.audit_events_table_name
  ca_private_key_parameter_name = module.ssm.ca_private_key_parameter_name
  default_duration_seconds      = var.default_certificate_duration_seconds
  max_duration_seconds          = var.max_certificate_duration_seconds
}

module "api_gateway" {
  source = "./modules/api_gateway"

  project_name         = var.project_name
  api_name             = var.api_gateway_name
  route_path           = var.issue_certificate_route_path
  cognito_user_pool_id = module.cognito.user_pool_id
  cognito_client_id    = module.cognito.user_pool_client_id
  signer_function_name = module.lambda_signer.function_name
  signer_invoke_arn    = module.lambda_signer.invoke_arn
}

module "route53" {
  source = "./modules/route53"

  project_name   = var.project_name
  hosted_zone_id = var.hosted_zone_id
  domain_name    = "${var.api_subdomain}.${var.base_domain}"
  api_id         = module.api_gateway.api_id
  stage_name     = module.api_gateway.default_stage_name
}

module "static_config" {
  source = "./modules/static_config"

  project_name = var.project_name
  bucket_name  = var.static_content_bucket_name
  api_id       = module.api_gateway.api_id
  config_json = jsonencode({
    region                   = var.aws_region
    cognito_domain           = module.cognito_auth_domain.auth_domain_url
    client_id                = module.cognito.user_pool_client_id
    redirect_uri             = var.callback_url
    api_base_url             = module.route53.api_base_url
    default_duration_seconds = var.default_certificate_duration_seconds
  })
  public_key = var.ca_public_key_value
}
