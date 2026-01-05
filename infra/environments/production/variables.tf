variable "project" {
  description = "Project name"
  type        = string
  default     = "truegul"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-northeast-2"
}

variable "vpc_cidr" {
  description = "VPC CIDR block"
  type        = string
  default     = "10.1.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["ap-northeast-2a", "ap-northeast-2c"]
}

variable "enable_nat_gateway" {
  description = "Enable NAT Gateway"
  type        = bool
  default     = true
}

# RDS
variable "rds_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.small"
}

variable "rds_multi_az" {
  description = "Enable Multi-AZ for RDS"
  type        = bool
  default     = false
}

variable "database_name" {
  description = "Database name"
  type        = string
  default     = "truegul"
}

variable "database_username" {
  description = "Database master username"
  type        = string
  default     = "truegul"
}

variable "database_password" {
  description = "Database master password"
  type        = string
  sensitive   = true
}

# Redis
variable "redis_node_type" {
  description = "ElastiCache node type"
  type        = string
  default     = "cache.t3.micro"
}

# ECS - API Server
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
  default     = 2
}

variable "api_server_max_count" {
  description = "API Server maximum task count"
  type        = number
  default     = 5
}

# ECS - ML Server
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
