variable "aws_region" {
  description = "AWS region for all resources"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "recipe-app"
}

variable "environment" {
  description = "Environment name (production, staging, dev)"
  type        = string
  default     = "production"
}

variable "db_name" {
  description = "PostgreSQL database name"
  type        = string
  default     = "recipes"
}

variable "db_username" {
  description = "PostgreSQL master username"
  type        = string
  default     = "recipeadmin"
}

variable "backend_cpu" {
  description = "CPU units for backend container (256 = 0.25 vCPU)"
  type        = number
  default     = 256
}

variable "backend_memory" {
  description = "Memory for backend container in MB"
  type        = number
  default     = 512
}

variable "frontend_cpu" {
  description = "CPU units for frontend container (256 = 0.25 vCPU)"
  type        = number
  default     = 256
}

variable "frontend_memory" {
  description = "Memory for frontend container in MB"
  type        = number
  default     = 512
}

variable "backend_replicas" {
  description = "Number of backend container instances to run"
  type        = number
  default     = 1
}

variable "frontend_replicas" {
  description = "Number of frontend container instances to run"
  type        = number
  default     = 1
}
