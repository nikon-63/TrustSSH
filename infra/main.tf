terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
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

module "dynamodb" {
  source = "./modules/dynamodb"

  project_name            = var.project_name
  user_mapping_table_name = var.dynamodb_user_mapping_table_name
  audit_events_table_name = var.dynamodb_audit_events_table_name
  deletion_protection     = var.dynamodb_deletion_protection
  point_in_time_recovery  = var.dynamodb_point_in_time_recovery
}
