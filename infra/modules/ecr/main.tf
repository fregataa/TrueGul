################################################################################
# ECR Repositories
################################################################################

resource "aws_ecr_repository" "api_server" {
  name                 = "${var.project}-api-server"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name    = "${var.project}-api-server"
    Project = var.project
  }
}

resource "aws_ecr_repository" "ml_server" {
  name                 = "${var.project}-ml-server"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name    = "${var.project}-ml-server"
    Project = var.project
  }
}

################################################################################
# Lifecycle Policies
################################################################################

resource "aws_ecr_lifecycle_policy" "api_server" {
  repository = aws_ecr_repository.api_server.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 10 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

resource "aws_ecr_lifecycle_policy" "ml_server" {
  repository = aws_ecr_repository.ml_server.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Keep last 10 images"
        selection = {
          tagStatus   = "any"
          countType   = "imageCountMoreThan"
          countNumber = 10
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}
