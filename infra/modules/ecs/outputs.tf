output "cluster_id" {
  description = "ECS cluster ID"
  value       = aws_ecs_cluster.main.id
}

output "cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "alb_dns_name" {
  description = "ALB DNS name"
  value       = aws_lb.main.dns_name
}

output "alb_zone_id" {
  description = "ALB zone ID"
  value       = aws_lb.main.zone_id
}

output "api_server_service_name" {
  description = "API Server ECS service name"
  value       = aws_ecs_service.api_server.name
}

output "ml_server_service_name" {
  description = "ML Server ECS service name"
  value       = aws_ecs_service.ml_server.name
}

output "api_server_task_definition_arn" {
  description = "API Server task definition ARN"
  value       = aws_ecs_task_definition.api_server.arn
}

output "ml_server_task_definition_arn" {
  description = "ML Server task definition ARN"
  value       = aws_ecs_task_definition.ml_server.arn
}

output "api_log_group_name" {
  description = "API Server CloudWatch log group name"
  value       = aws_cloudwatch_log_group.api_server.name
}

output "ml_log_group_name" {
  description = "ML Server CloudWatch log group name"
  value       = aws_cloudwatch_log_group.ml_server.name
}
