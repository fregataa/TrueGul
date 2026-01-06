################################################################################
# ElastiCache Subnet Group
################################################################################

resource "aws_elasticache_subnet_group" "main" {
  name       = "${var.project}-${var.environment}-cache-subnet"
  subnet_ids = var.subnet_ids

  tags = {
    Name        = "${var.project}-${var.environment}-cache-subnet"
    Environment = var.environment
    Project     = var.project
  }
}

################################################################################
# ElastiCache Valkey Cluster
################################################################################

resource "aws_elasticache_cluster" "main" {
  cluster_id           = "${var.project}-${var.environment}-valkey"
  engine               = "valkey"
  engine_version       = var.engine_version
  node_type            = var.node_type
  num_cache_nodes      = 1
  parameter_group_name = "default.valkey8"
  port                 = 6379

  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = var.security_group_ids

  snapshot_retention_limit = var.environment == "production" ? 7 : 0
  snapshot_window          = "02:00-03:00"
  maintenance_window       = "sun:03:00-sun:04:00"

  tags = {
    Name        = "${var.project}-${var.environment}-valkey"
    Environment = var.environment
    Project     = var.project
  }
}
