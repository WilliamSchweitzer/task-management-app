# This file configures Terraform to use S3 for state storage
# Run this AFTER creating the S3 bucket and DynamoDB table manually
# or use the bootstrap script in scripts/bootstrap-terraform.sh

terraform {
  backend "s3" {
    bucket         = "task-management-terraform-state"  # Update with your bucket name
    key            = "task-management-app/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "task-management-terraform-locks"  # Update with your table name
  }
}

# Initial setup (before backend is configured):
# 1. Comment out the backend block above
# 2. Run: terraform init
# 3. Run: terraform apply -target=module.bootstrap (if you create a bootstrap module)
# 4. Uncomment the backend block
# 5. Run: terraform init -migrate-state
