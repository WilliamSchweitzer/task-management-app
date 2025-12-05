# ACM Certificate for API subdomain
resource "aws_acm_certificate" "api" {
  domain_name       = "api.task-management.wschweitzer.com"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name        = "task-management-api-cert"
  }
}

# Output validation records to add in Netlify DNS
output "acm_validation_record" {
  description = "Add this CNAME record in Netlify DNS to validate the certificate"
  value = {
    for dvo in aws_acm_certificate.api.domain_validation_options : dvo.domain_name => {
      record_name  = dvo.resource_record_name
      record_type  = dvo.resource_record_type
      record_value = dvo.resource_record_value
    }
  }
}

output "nlb_https_dns_name" {
  description = "Add CNAME in Netlify: api.task-management -> this value"
  value       = aws_lb.main.dns_name
}