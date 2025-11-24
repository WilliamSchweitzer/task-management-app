# Create separate database for Kong on same RDS instance
resource "null_resource" "create_kong_database" {
  depends_on = [aws_db_instance.taskmanagement]

  provisioner "local-exec" {
    command = <<-EOT
      PGPASSWORD=${var.db_password} psql -h ${aws_db_instance.taskmanagement.address} -U postgres -d postgres -c "CREATE DATABASE kong;"
      PGPASSWORD=${var.db_password} psql -h ${aws_db_instance.taskmanagement.address} -U postgres -d postgres -c "CREATE USER kong WITH PASSWORD '${var.kong_db_password}';"
      PGPASSWORD=${var.db_password} psql -h ${aws_db_instance.taskmanagement.address} -U postgres -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE kong TO kong;"
    EOT
  }
}