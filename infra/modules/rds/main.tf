################################################################################
# DB Subnet Group
################################################################################

resource "aws_db_subnet_group" "main" {
  name       = "${var.project}-${var.environment}-db-subnet"
  subnet_ids = var.subnet_ids

  tags = {
    Name        = "${var.project}-${var.environment}-db-subnet"
    Environment = var.environment
    Project     = var.project
  }
}

################################################################################
# RDS Instance
################################################################################

resource "aws_db_instance" "main" {
  identifier = "${var.project}-${var.environment}-postgres"

  engine               = "postgres"
  engine_version       = var.engine_version
  instance_class       = var.instance_class
  allocated_storage    = var.allocated_storage
  storage_type         = "gp3"
  storage_encrypted    = true

  db_name  = var.database_name
  username = var.database_username
  password = var.database_password

  vpc_security_group_ids = var.security_group_ids
  db_subnet_group_name   = aws_db_subnet_group.main.name

  multi_az               = var.multi_az
  publicly_accessible    = false
  deletion_protection    = var.deletion_protection
  skip_final_snapshot    = var.environment != "production"
  final_snapshot_identifier = var.environment == "production" ? "${var.project}-${var.environment}-final-snapshot" : null

  backup_retention_period = var.backup_retention_period
  backup_window          = "03:00-04:00"
  maintenance_window     = "Mon:04:00-Mon:05:00"

  performance_insights_enabled = var.environment == "production"

  tags = {
    Name        = "${var.project}-${var.environment}-postgres"
    Environment = var.environment
    Project     = var.project
  }
}
