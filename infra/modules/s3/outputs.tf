output "ml_models_bucket_name" {
  description = "ML models S3 bucket name"
  value       = aws_s3_bucket.ml_models.id
}

output "ml_models_bucket_arn" {
  description = "ML models S3 bucket ARN"
  value       = aws_s3_bucket.ml_models.arn
}
