output "endpoint" {
  description = "Valkey primary endpoint"
  value       = aws_elasticache_replication_group.main.primary_endpoint_address
}

output "port" {
  description = "Valkey port"
  value       = aws_elasticache_replication_group.main.port
}

output "redis_url" {
  description = "Valkey connection URL (redis:// protocol)"
  value       = "redis://${aws_elasticache_replication_group.main.primary_endpoint_address}:${aws_elasticache_replication_group.main.port}"
}
