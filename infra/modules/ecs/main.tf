################################################################################
# ECS Cluster
################################################################################

resource "aws_ecs_cluster" "main" {
  name = "${var.project}-${var.environment}"

  setting {
    name  = "containerInsights"
    value = var.environment == "production" ? "enabled" : "disabled"
  }

  tags = {
    Name        = "${var.project}-${var.environment}"
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name = aws_ecs_cluster.main.name

  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

################################################################################
# CloudWatch Log Groups
################################################################################

resource "aws_cloudwatch_log_group" "api_server" {
  name              = "/ecs/${var.project}-${var.environment}/api-server"
  retention_in_days = var.log_retention_days

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_cloudwatch_log_group" "ml_server" {
  name              = "/ecs/${var.project}-${var.environment}/ml-server"
  retention_in_days = var.log_retention_days

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

################################################################################
# IAM Roles
################################################################################

# ECS Task Execution Role
resource "aws_iam_role" "ecs_task_execution" {
  name = "${var.project}-${var.environment}-ecs-task-execution"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "ecs_task_execution_secrets" {
  name = "${var.project}-${var.environment}-secrets-access"
  role = aws_iam_role.ecs_task_execution.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = "arn:aws:secretsmanager:${var.aws_region}:*:secret:${var.project}/${var.environment}/*"
      }
    ]
  })
}

# ECS Task Role
resource "aws_iam_role" "ecs_task" {
  name = "${var.project}-${var.environment}-ecs-task"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

################################################################################
# Application Load Balancer
################################################################################

resource "aws_lb" "main" {
  name               = "${var.project}-${var.environment}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.alb_security_group_id]
  subnets            = var.public_subnet_ids

  enable_deletion_protection = var.environment == "production"

  tags = {
    Name        = "${var.project}-${var.environment}-alb"
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_lb_target_group" "api_server" {
  name        = "${var.project}-${var.environment}-api"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/health"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 3
  }

  tags = {
    Name        = "${var.project}-${var.environment}-api"
    Environment = var.environment
    Project     = var.project
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api_server.arn
  }
}

################################################################################
# API Server Task Definition
################################################################################

resource "aws_ecs_task_definition" "api_server" {
  family                   = "${var.project}-${var.environment}-api-server"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.api_server_cpu
  memory                   = var.api_server_memory
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name  = "api-server"
      image = "${var.api_server_image}:${var.api_server_image_tag}"

      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "PORT", value = "8080" },
        { name = "ENVIRONMENT", value = var.environment },
        { name = "DATABASE_URL", value = var.database_url },
        { name = "REDIS_URL", value = var.redis_url },
        { name = "ML_SERVER_URL", value = "http://localhost:8000" },
        { name = "CALLBACK_BASE_URL", value = "http://localhost:8080" },
        { name = "CORS_ORIGINS", value = var.cors_origins }
      ]

      secrets = [
        {
          name      = "JWT_SECRET"
          valueFrom = "arn:aws:secretsmanager:${var.aws_region}:${var.aws_account_id}:secret:${var.project}/${var.environment}/jwt"
        },
        {
          name      = "ML_CALLBACK_SECRET"
          valueFrom = "arn:aws:secretsmanager:${var.aws_region}:${var.aws_account_id}:secret:${var.project}/${var.environment}/callback"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.api_server.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }

      healthCheck = {
        command     = ["CMD-SHELL", "wget -q --spider http://localhost:8080/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 60
      }

      essential = true
    }
  ])

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

################################################################################
# ML Server Task Definition
################################################################################

resource "aws_ecs_task_definition" "ml_server" {
  family                   = "${var.project}-${var.environment}-ml-server"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.ml_server_cpu
  memory                   = var.ml_server_memory
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name  = "ml-server"
      image = "${var.ml_server_image}:${var.ml_server_image_tag}"

      portMappings = [
        {
          containerPort = 8000
          hostPort      = 8000
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "PORT", value = "8000" },
        { name = "ENVIRONMENT", value = var.environment },
        { name = "REDIS_URL", value = var.redis_url },
        { name = "CALLBACK_BASE_URL", value = "http://${aws_lb.main.dns_name}" }
      ]

      secrets = [
        {
          name      = "ML_CALLBACK_SECRET"
          valueFrom = "arn:aws:secretsmanager:${var.aws_region}:${var.aws_account_id}:secret:${var.project}/${var.environment}/callback"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ml_server.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }

      healthCheck = {
        command     = ["CMD-SHELL", "wget -q --spider http://localhost:8000/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 120
      }

      essential = true
    }
  ])

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

################################################################################
# ECS Services
################################################################################

resource "aws_ecs_service" "api_server" {
  name            = "api-server"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.api_server.arn
  desired_count   = var.api_server_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.ecs_security_group_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api_server.arn
    container_name   = "api-server"
    container_port   = 8080
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  tags = {
    Environment = var.environment
    Project     = var.project
  }

  depends_on = [aws_lb_listener.http]
}

resource "aws_ecs_service" "ml_server" {
  name            = "ml-server"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.ml_server.arn
  desired_count   = var.ml_server_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.ecs_security_group_id]
    assign_public_ip = false
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  tags = {
    Environment = var.environment
    Project     = var.project
  }
}

################################################################################
# Auto Scaling
################################################################################

resource "aws_appautoscaling_target" "api_server" {
  max_capacity       = var.api_server_max_count
  min_capacity       = var.api_server_desired_count
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.api_server.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "api_server_cpu" {
  name               = "${var.project}-${var.environment}-api-cpu"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.api_server.resource_id
  scalable_dimension = aws_appautoscaling_target.api_server.scalable_dimension
  service_namespace  = aws_appautoscaling_target.api_server.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value       = 70.0
    scale_in_cooldown  = 300
    scale_out_cooldown = 60
  }
}
