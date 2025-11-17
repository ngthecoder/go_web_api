# Go Recipe Web API + Next.js Frontend

A production-ready full-stack recipe discovery application with secure authentication, deployed to AWS using Infrastructure as Code.

## Project Overview

A full-stack web application that helps users discover recipes, manage their favorite dishes, and generate shopping lists. Built with Go backend, Next.js frontend, PostgreSQL database, and deployed on AWS using Terraform and ECS Fargate.

**Tech Stack**: Go 1.24 • Next.js 15 • PostgreSQL 16 • Docker • Terraform • AWS ECS

## Key Features

- **Secure Authentication**: JWT tokens with Argon2 password hashing
- **Smart Recipe Search**: Find recipes by available ingredients
- **Personal Collections**: Save and manage favorite recipes
- **Shopping Lists**: Generate ingredient lists from recipes
- **Responsive Design**: Mobile-first UI with Tailwind CSS 4
- **Cloud-Native**: Containerized deployment on AWS

## Architecture

```
Internet
    ↓
┌─────────────────────────────────────────────┐
│  AWS VPC (10.0.0.0/16)                      │
│                                             │
│  ┌─────────────────┐  ┌─────────────────┐  │
│  │ Frontend ALB    │  │ Backend ALB     │  │
│  │ (Port 80)       │  │ (Port 80)       │  │
│  └────────┬────────┘  └────────┬────────┘  │
│           ↓                    ↓            │
│  ┌─────────────────┐  ┌─────────────────┐  │
│  │ Frontend ECS    │  │ Backend ECS     │  │
│  │ (Next.js)       │  │ (Go API)        │  │
│  │ Public Subnet   │  │ Public Subnet   │  │
│  └─────────────────┘  └────────┬────────┘  │
│                                 ↓            │
│                        ┌─────────────────┐  │
│                        │ RDS PostgreSQL  │  │
│                        │ Private Subnet  │  │
│                        └─────────────────┘  │
└─────────────────────────────────────────────┘

Storage:
- ECR: Docker image registry
- Secrets Manager: Passwords, JWT secret
- CloudWatch: Application logs
```

## Deployment Guide

### Prerequisites

- **AWS Account** with CLI configured (`aws configure`)
- **Terraform** 1.0+ installed
- **Docker** installed
- **Git** for version control

### Environment Setup

Your setup will vary, so configure these for your environment:

1. **AWS Region**: Default is `us-east-1` (change in `terraform/variables.tf`)
2. **Project Name**: Default is `recipe-app` (change in `terraform/variables.tf`)
3. **Resource Sizing**: Adjust CPU/memory in `terraform/variables.tf` based on your needs

### Step 1: Clone Repository

```bash
git clone https://github.com/ngthecoder/go_web_api.git
cd go_web_api
```

### Step 2: Deploy Infrastructure

```bash
cd terraform/

# Initialize Terraform
terraform init

# Review what will be created
terraform plan

# Create AWS infrastructure
terraform apply
# Type 'yes' when prompted

# Save important outputs
terraform output backend_alb_dns     # Copy this URL
terraform output frontend_alb_dns    # Copy this URL
```

**What gets created**: VPC, subnets, security groups, RDS database, ECR repositories, ECS cluster, load balancers, secrets (~$40-50/month)

### Step 3: Build and Push Backend

```bash
cd ../backend/

# Build Docker image
docker build -t recipe-app-backend .

# Tag for ECR (replace with your account ID from terraform output)
docker tag recipe-app-backend:latest <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-backend:latest

# Login to ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com

# Push to ECR
docker push <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-backend:latest
```

### Step 4: Build and Push Frontend

**Important**: Frontend needs the backend URL at build time.

```bash
cd ../frontend/

# Build with backend URL (use the URL from Step 2)
docker build \
  --build-arg NEXT_PUBLIC_API_URL=http://<BACKEND_ALB_DNS> \
  -t recipe-app-frontend .

# Tag for ECR
docker tag recipe-app-frontend:latest <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-frontend:latest

# Push to ECR
docker push <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-frontend:latest
```

### Step 5: Deploy to ECS

```bash
# Force ECS to pull and run the new images
aws ecs update-service \
  --cluster recipe-app-cluster \
  --service recipe-app-backend-service \
  --force-new-deployment \
  --region us-east-1

aws ecs update-service \
  --cluster recipe-app-cluster \
  --service recipe-app-frontend-service \
  --force-new-deployment \
  --region us-east-1
```

### Step 6: Verify Deployment

```bash
# Check backend logs
aws logs tail /ecs/recipe-app-backend --follow --region us-east-1
# Look for: "Database connection established successfully"

# Check frontend logs
aws logs tail /ecs/recipe-app-frontend --follow --region us-east-1
```

**Access your app**: Open `http://<FRONTEND_ALB_DNS>` in your browser (from terraform output)

## Update Workflow

### Update Backend Code

```bash
cd backend/
# Make your code changes

docker build -t recipe-app-backend .
docker tag recipe-app-backend:latest <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-backend:latest
docker push <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-backend:latest

aws ecs update-service --cluster recipe-app-cluster --service recipe-app-backend-service --force-new-deployment --region us-east-1
```

### Update Frontend Code

```bash
cd frontend/
# Make your code changes

docker build --build-arg NEXT_PUBLIC_API_URL=http://<BACKEND_ALB_DNS> -t recipe-app-frontend .
docker tag recipe-app-frontend:latest <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-frontend:latest
docker push <YOUR_ACCOUNT_ID>.dkr.ecr.us-east-1.amazonaws.com/recipe-app-frontend:latest

aws ecs update-service --cluster recipe-app-cluster --service recipe-app-frontend-service --force-new-deployment --region us-east-1
```

### Update Infrastructure

```bash
cd terraform/
# Edit .tf files

terraform plan    # Review changes
terraform apply   # Apply changes
```

## Teardown

### Stop Services (saves ~80% cost, keeps data)

```bash
aws ecs update-service --cluster recipe-app-cluster --service recipe-app-backend-service --desired-count 0 --region us-east-1
aws ecs update-service --cluster recipe-app-cluster --service recipe-app-frontend-service --desired-count 0 --region us-east-1
```

### Destroy Everything

```bash
cd terraform/
terraform destroy
# Type 'yes' when prompted
```

**Warning**: This deletes all data including the database. Images in ECR are also deleted.

## Security Features

- **Network Isolation**: Database in private subnet, no internet access
- **Security Groups**: Strict firewall rules (ALB → ECS → RDS only)
- **SSL/TLS**: RDS requires SSL connections
- **Secrets Management**: AWS Secrets Manager for passwords and keys
- **Password Hashing**: Argon2id (time=3, memory=64MB, 16-byte salt)
- **JWT Tokens**: HMAC-SHA256, 24-hour expiration
- **CORS Protection**: Configured allowed origins

## API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/auth/register` | POST | No | Create account |
| `/api/auth/login` | POST | No | Login |
| `/api/recipes` | GET | Optional | Browse recipes |
| `/api/recipes/{id}` | GET | Optional | Recipe details |
| `/api/recipes/find-by-ingredients` | GET | Optional | Find by ingredients |
| `/api/user/profile` | GET | Yes | User profile |
| `/api/user/liked-recipes` | GET | Yes | Saved recipes |
| `/api/ingredients` | GET | No | Browse ingredients |
| `/api/stats` | GET | No | Statistics |

## Project Structure

```
go_web_api/
├── backend/
│   ├── internal/
│   │   ├── auth/           # JWT & authentication
│   │   ├── database/       # DB connection & migrations
│   │   ├── recipes/        # Recipe logic
│   │   ├── ingredients/    # Ingredient logic
│   │   └── users/          # User management
│   ├── main.go
│   └── Dockerfile
├── frontend/
│   ├── app/
│   │   ├── recipes/        # Recipe pages
│   │   ├── profile/        # User profile
│   │   └── page.tsx        # Home
│   ├── components/         # React components
│   ├── contexts/           # Auth context
│   └── Dockerfile
└── terraform/
    ├── main.tf             # Provider config
    ├── vpc.tf              # Network setup
    ├── rds.tf              # Database
    ├── ecs.tf              # Container orchestration
    ├── alb.tf              # Load balancers
    ├── ecr.tf              # Image registry
    ├── secrets.tf          # Secrets Manager
    ├── security.tf         # Security groups
    ├── cloudwatch.tf       # Logging
    ├── variables.tf        # Configuration
    └── outputs.tf          # Important values
```

### Deployment Flow Summary

1. **Terraform** creates AWS infrastructure (one-time)
2. **Docker** builds application images (whenever code changes)
3. **ECR** stores Docker images (AWS Docker registry)
4. **ECS** pulls images and runs containers
5. **ALB** routes traffic to containers
6. **RDS** provides database service
7. **Secrets Manager** injects credentials at runtime

## Testing

```bash
cd backend
go test ./... -v -cover
```

**Coverage**: ~40% focusing on authentication security and core business logic

## What I Learned

- **Infrastructure as Code**: Terraform for reproducible AWS deployments
- **Container Orchestration**: ECS Fargate for serverless container management
- **Networking**: VPCs, subnets, security groups, load balancers
- **Security**: Secrets management, SSL/TLS, least-privilege IAM roles
- **DevOps**: Docker multi-stage builds, ECR, CloudWatch logging
- **Full-Stack**: Go REST API, Next.js SSR, PostgreSQL database design
- **Debugging**: AWS CloudWatch logs, ECS task troubleshooting
