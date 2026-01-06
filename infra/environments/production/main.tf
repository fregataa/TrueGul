terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket       = "truegul-terraform-state"
    key          = "production/terraform.tfstate"
    region       = "ap-northeast-2"
    encrypt      = true
    use_lockfile = true
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}

################################################################################
# Data Sources
################################################################################

data "aws_caller_identity" "current" {}

################################################################################
# ECR
################################################################################

module "ecr" {
  source = "../../modules/ecr"

  project = var.project
}

################################################################################
# VPC
################################################################################

module "vpc" {
  source = "../../modules/vpc"

  project            = var.project
  environment        = var.environment
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
  enable_nat_gateway = var.enable_nat_gateway
}

################################################################################
# RDS
################################################################################

module "rds" {
  source = "../../modules/rds"

  project            = var.project
  environment        = var.environment
  subnet_ids         = module.vpc.private_subnet_ids
  security_group_ids = [module.vpc.rds_security_group_id]

  instance_class      = var.rds_instance_class
  database_name       = var.database_name
  database_username   = var.database_username
  database_password   = var.database_password
  deletion_protection = true
  multi_az            = var.rds_multi_az
}

################################################################################
# ElastiCache
################################################################################

module "elasticache" {
  source = "../../modules/elasticache"

  project            = var.project
  environment        = var.environment
  subnet_ids         = module.vpc.private_subnet_ids
  security_group_ids = [module.vpc.redis_security_group_id]

  node_type = var.redis_node_type
}

################################################################################
# ECS
################################################################################

module "ecs" {
  source = "../../modules/ecs"

  project        = var.project
  environment    = var.environment
  aws_region     = var.aws_region
  aws_account_id = data.aws_caller_identity.current.account_id

  vpc_id             = module.vpc.vpc_id
  public_subnet_ids  = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids

  alb_security_group_id = module.vpc.alb_security_group_id
  ecs_security_group_id = module.vpc.ecs_security_group_id

  database_url = module.rds.database_url
  redis_url    = module.elasticache.redis_url

  api_server_image = module.ecr.api_server_repository_url
  ml_server_image  = module.ecr.ml_server_repository_url

  api_server_cpu           = var.api_server_cpu
  api_server_memory        = var.api_server_memory
  api_server_desired_count = var.api_server_desired_count
  api_server_max_count     = var.api_server_max_count

  ml_server_cpu           = var.ml_server_cpu
  ml_server_memory        = var.ml_server_memory
  ml_server_desired_count = var.ml_server_desired_count

  depends_on = [module.rds, module.elasticache]
}

################################################################################
# Outputs
################################################################################

output "alb_dns_name" {
  description = "ALB DNS name for API access"
  value       = module.ecs.alb_dns_name
}

output "api_server_repository_url" {
  description = "API Server ECR repository URL"
  value       = module.ecr.api_server_repository_url
}

output "ml_server_repository_url" {
  description = "ML Server ECR repository URL"
  value       = module.ecr.ml_server_repository_url
}

output "rds_endpoint" {
  description = "RDS endpoint"
  value       = module.rds.endpoint
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = module.elasticache.endpoint
}
