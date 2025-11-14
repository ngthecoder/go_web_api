output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "database_endpoint" {
  description = "RDS instance endpoint for database connections"
  value       = aws_db_instance.main.endpoint
}

output "database_name" {
  description = "Name of the database"
  value       = aws_db_instance.main.db_name
}

output "database_username" {
  description = "Master username for the database"
  value       = aws_db_instance.main.username
}

output "database_port" {
  description = "Port the database is listening on"
  value       = aws_db_instance.main.port
}

output "database_password_secret" {
  description = "Database password (sensitive)"
  value       = random_password.db_password.result
  sensitive   = true
}

output "public_subnet_ids" {
  description = "IDs of public subnets (for ECS tasks)"
  value       = aws_subnet.public[*].id
}

output "private_subnet_ids" {
  description = "IDs of private subnets (for RDS)"
  value       = aws_subnet.private[*].id
}

output "alb_security_group_id" {
  description = "Security group ID for Application Load Balancer"
  value       = aws_security_group.alb.id
}

output "backend_security_group_id" {
  description = "Security group ID for backend ECS tasks"
  value       = aws_security_group.backend.id
}

output "frontend_security_group_id" {
  description = "Security group ID for frontend ECS tasks"
  value       = aws_security_group.frontend.id
}

output "rds_security_group_id" {
  description = "Security group ID for RDS instance"
  value       = aws_security_group.rds.id
}

output "aws_region" {
  description = "AWS region where resources are deployed"
  value       = var.aws_region
}

output "backend_url" {
  description = "URL to access the backend API"
  value       = "http://${aws_lb.backend.dns_name}"
}

output "frontend_url" {
  description = "URL to access the frontend application"
  value       = "http://${aws_lb.frontend.dns_name}"
}

output "backend_alb_dns" {
  description = "Backend ALB DNS name"
  value       = aws_lb.backend.dns_name
}

output "frontend_alb_dns" {
  description = "Frontend ALB DNS name"
  value       = aws_lb.frontend.dns_name
}

output "backend_ecr_repository_url" {
  description = "ECR repository URL for backend Docker images"
  value       = aws_ecr_repository.backend.repository_url
}

output "frontend_ecr_repository_url" {
  description = "ECR repository URL for frontend Docker images"
  value       = aws_ecr_repository.frontend.repository_url
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "database_config_secret_arn" {
  description = "ARN of the database config secret in Secrets Manager"
  value       = aws_secretsmanager_secret.database_config.arn
}

output "jwt_secret_arn" {
  description = "ARN of the JWT secret in Secrets Manager"
  value       = aws_secretsmanager_secret.jwt_secret.arn
}
