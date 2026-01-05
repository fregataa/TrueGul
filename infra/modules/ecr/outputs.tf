output "api_server_repository_url" {
  description = "API Server ECR repository URL"
  value       = aws_ecr_repository.api_server.repository_url
}

output "ml_server_repository_url" {
  description = "ML Server ECR repository URL"
  value       = aws_ecr_repository.ml_server.repository_url
}

output "api_server_repository_arn" {
  description = "API Server ECR repository ARN"
  value       = aws_ecr_repository.api_server.arn
}

output "ml_server_repository_arn" {
  description = "ML Server ECR repository ARN"
  value       = aws_ecr_repository.ml_server.arn
}
