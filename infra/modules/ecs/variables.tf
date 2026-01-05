variable "project" {
  description = "Project name"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "aws_region" {
  description = "AWS region"
  type        = string
}

variable "aws_account_id" {
  description = "AWS account ID"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs"
  type        = list(string)
}

variable "alb_security_group_id" {
  description = "ALB security group ID"
  type        = string
}

variable "ecs_security_group_id" {
  description = "ECS security group ID"
  type        = string
}

variable "database_url" {
  description = "Database connection URL"
  type        = string
  sensitive   = true
}

variable "redis_url" {
  description = "Redis connection URL"
  type        = string
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 14
}

# API Server
variable "api_server_image" {
  description = "API Server Docker image URL"
  type        = string
}

variable "api_server_image_tag" {
  description = "API Server Docker image tag"
  type        = string
  default     = "latest"
}

variable "api_server_cpu" {
  description = "API Server CPU units"
  type        = number
  default     = 256
}

variable "api_server_memory" {
  description = "API Server memory in MB"
  type        = number
  default     = 512
}

variable "api_server_desired_count" {
  description = "API Server desired task count"
  type        = number
  default     = 1
}

variable "api_server_max_count" {
  description = "API Server maximum task count"
  type        = number
  default     = 3
}

# ML Server
variable "ml_server_image" {
  description = "ML Server Docker image URL"
  type        = string
}

variable "ml_server_image_tag" {
  description = "ML Server Docker image tag"
  type        = string
  default     = "latest"
}

variable "ml_server_cpu" {
  description = "ML Server CPU units"
  type        = number
  default     = 512
}

variable "ml_server_memory" {
  description = "ML Server memory in MB"
  type        = number
  default     = 1024
}

variable "ml_server_desired_count" {
  description = "ML Server desired task count"
  type        = number
  default     = 1
}
