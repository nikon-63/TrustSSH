terraform {
  required_providers {
    aws = {
      source                = "hashicorp/aws"
      configuration_aliases = [aws.us_east_1]
    }
  }
}

resource "aws_acm_certificate" "auth" {
  provider = aws.us_east_1

  domain_name       = var.domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Project = var.project_name
  }
}

resource "aws_route53_record" "certificate_validation" {
  for_each = {
    for option in aws_acm_certificate.auth.domain_validation_options : option.domain_name => {
      name   = option.resource_record_name
      record = option.resource_record_value
      type   = option.resource_record_type
    }
  }

  zone_id = var.hosted_zone_id
  name    = each.value.name
  type    = each.value.type
  ttl     = 60
  records = [each.value.record]
}

resource "aws_acm_certificate_validation" "auth" {
  provider = aws.us_east_1

  certificate_arn         = aws_acm_certificate.auth.arn
  validation_record_fqdns = [for record in aws_route53_record.certificate_validation : record.fqdn]
}

resource "aws_cognito_user_pool_domain" "auth" {
  domain                = var.domain_name
  certificate_arn       = aws_acm_certificate_validation.auth.certificate_arn
  managed_login_version = 2
  user_pool_id          = var.user_pool_id
}

resource "aws_route53_record" "auth" {
  zone_id = var.hosted_zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    name                   = aws_cognito_user_pool_domain.auth.cloudfront_distribution
    zone_id                = aws_cognito_user_pool_domain.auth.cloudfront_distribution_zone_id
    evaluate_target_health = false
  }
}

resource "terraform_data" "managed_login_branding" {
  triggers_replace = {
    client_id    = var.client_id
    user_pool_id = var.user_pool_id
  }

  provisioner "local-exec" {
    interpreter = ["/bin/bash", "-c"]
    command     = <<EOT
set -e
if aws cognito-idp describe-managed-login-branding-by-client --region ${var.aws_region} --user-pool-id ${var.user_pool_id} --client-id ${var.client_id} >/dev/null 2>&1; then
  exit 0
fi
aws cognito-idp create-managed-login-branding --region ${var.aws_region} --user-pool-id ${var.user_pool_id} --client-id ${var.client_id} --use-cognito-provided-values
EOT
  }

  depends_on = [
    aws_cognito_user_pool_domain.auth,
  ]
}

resource "terraform_data" "webauthn_config" {
  triggers_replace = {
    relying_party_id  = var.domain_name
    user_pool_id      = var.user_pool_id
    user_verification = var.webauthn_user_verification
  }

  provisioner "local-exec" {
    interpreter = ["/bin/bash", "-c"]
    command     = <<EOT
set -e
for attempt in {1..30}; do
  if aws cognito-idp set-user-pool-mfa-config --region ${var.aws_region} --user-pool-id ${var.user_pool_id} --web-authn-configuration RelyingPartyId=${var.domain_name},UserVerification=${var.webauthn_user_verification}; then
    exit 0
  fi
  sleep 20
done
aws cognito-idp set-user-pool-mfa-config --region ${var.aws_region} --user-pool-id ${var.user_pool_id} --web-authn-configuration RelyingPartyId=${var.domain_name},UserVerification=${var.webauthn_user_verification}
EOT
  }

  depends_on = [
    aws_cognito_user_pool_domain.auth,
    aws_route53_record.auth,
    terraform_data.managed_login_branding,
  ]
}
