# TrueGul Infrastructure

Terraform configuration for deploying TrueGul to AWS.

## Architecture

```
                              User Browser
                              /          \
                   (1) Load  /            \ (2) API calls
                       HTML/JS             \
                           /                \
               ┌──────────▼──┐         ┌─────▼─────┐
               │  Frontend   │         │    ALB    │
               │  (Vercel)   │         └─────┬─────┘
               └─────────────┘               │
                                        ┌────▼────┐
                              callback  │   API   │
                           ┌────────────│ Server  │
                           │            │ (ECS)   │
                     ┌─────▼─────┐      └─┬────┬──┘
                     │    ML     │        │    │
                     │  Server   │        │    │
                     │  (ECS)    │        │    │
                     └─────┬─────┘        │    │
                           │              │    │
                     ┌─────▼─────┐ ┌──────▼──┐ │
                     │   Redis   │ │   RDS   │ │
                     │(ElastiCache)│(Postgres)│◄┘
                     └───────────┘ └─────────┘
```

**Data Flow:**
1. Browser loads frontend from Vercel
2. Browser makes API calls directly to ALB
3. API Server enqueues task to Redis
4. ML Server consumes task from Redis, processes it
5. ML Server sends result via HTTP callback to API Server
6. API Server stores result in RDS

## Prerequisites

1. **AWS CLI** configured with appropriate credentials
2. **Terraform** >= 1.0
3. **S3 Bucket** for Terraform state: `truegul-terraform-state`
4. **DynamoDB Table** for state locking: `truegul-terraform-locks`

### Create State Backend (One-time setup)

```bash
# Create S3 bucket for state
aws s3 mb s3://truegul-terraform-state --region ap-northeast-2

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket truegul-terraform-state \
  --versioning-configuration Status=Enabled

# Create DynamoDB table for locking
aws dynamodb create-table \
  --table-name truegul-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region ap-northeast-2
```

## Directory Structure

```
infra/
├── modules/
│   ├── vpc/          # VPC, Subnets, Security Groups
│   ├── ecr/          # ECR Repositories
│   ├── rds/          # RDS PostgreSQL
│   ├── elasticache/  # ElastiCache Redis
│   └── ecs/          # ECS Cluster, Services, ALB
├── environments/
│   └── production/   # Production environment
└── README.md
```

## Deployment Steps

### 1. Deploy Infrastructure

```bash
cd infra/environments/production

# Copy and configure variables
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your database password

# Initialize Terraform
terraform init

# Review plan
terraform plan

# Apply
terraform apply
```

### 2. Create Secrets in AWS Secrets Manager

```bash
# JWT Secret
aws secretsmanager create-secret \
  --name truegul/production/jwt \
  --secret-string "<REPLACE_WITH_JWT_SECRET>"

# Callback Secret
aws secretsmanager create-secret \
  --name truegul/production/callback \
  --secret-string "<REPLACE_WITH_CALLBACK_SECRET>"
```

### 3. Build and Push Docker Images

Images are automatically built and pushed via GitHub Actions on merge to `main`.

For manual push:
```bash
# Login to ECR
aws ecr get-login-password --region ap-northeast-2 | \
  docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-northeast-2.amazonaws.com

# Build and push API server
cd api-server
docker build -t truegul-api-server .
docker tag truegul-api-server:latest <account-id>.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-api-server:latest
docker push <account-id>.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-api-server:latest

# Build and push ML server
cd ../ml-server
docker build -t truegul-ml-server .
docker tag truegul-ml-server:latest <account-id>.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-ml-server:latest
docker push <account-id>.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-ml-server:latest
```

### 4. Run Database Migrations

```bash
# Using ECS Exec (requires AWS CLI v2)
aws ecs execute-command \
  --cluster truegul-production \
  --task <task-id> \
  --container api-server \
  --interactive \
  --command "/bin/sh"

# Inside container, migrations run automatically on startup
```

## GitHub Actions Secrets Required

| Secret | Description |
|--------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key with ECR and ECS permissions |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |

## Outputs

After deployment, Terraform outputs:

- `alb_dns_name`: ALB DNS name for API access
- `api_server_repository_url`: ECR repository URL for API server
- `ml_server_repository_url`: ECR repository URL for ML server
- `rds_endpoint`: RDS endpoint
- `redis_endpoint`: Redis endpoint

## Cost Optimization

### NAT Gateway Alternatives
NAT Gateway costs ~$32/month. Alternatives:
- **VPC Endpoints** for AWS services (ECR, S3, CloudWatch)
- **NAT Instance** (less reliable, but cheaper)

### RDS
- Start with **db.t3.micro** for low traffic
- Enable **Multi-AZ** for high availability

### ElastiCache
- Single node sufficient for low traffic
- Consider upgrading as traffic grows

## Troubleshooting

### ECS Service Not Starting
1. Check CloudWatch logs: `/ecs/truegul-production/api-server`
2. Verify Secrets Manager permissions
3. Check security group rules

### Database Connection Failed
1. Verify RDS security group allows ECS security group
2. Check DATABASE_URL format
3. Ensure RDS is in same VPC as ECS

### Redis Connection Failed
1. Verify ElastiCache security group
2. Check REDIS_URL format
3. Ensure Redis is accessible from private subnets
