# Go Recipe Web API + Next.js Frontend

A production-ready full-stack recipe discovery application with secure authentication, demonstrating modern containerized application development and cloud deployment.

## 🎯 Project Overview

A full-stack web application that helps users discover recipes, manage their favorite dishes, and generate shopping lists based on available ingredients. Built with Go backend, Next.js frontend, PostgreSQL database, and deployed on AWS using Terraform and ECS Fargate.

## 🛠️ Tech Stack

**Backend**
- Go 1.24.0
- PostgreSQL 16 (production) / 15 (local)
- JWT authentication with Argon2id password hashing
- RESTful API with 14 endpoints

**Frontend**
- Next.js 15.5.0 with React 19.1.0
- TypeScript
- Tailwind CSS 4
- Context API for state management

**Infrastructure & DevOps**
- **Local Development**: Docker Compose, Kubernetes (Minikube)
- **Production**: AWS (Terraform, ECS Fargate, RDS, ECR, ALB)
- Docker multi-stage builds
- CloudWatch logging
- AWS Secrets Manager

**Database Schema**
- 5 tables: `ingredients`, `recipes`, `recipe_ingredients`, `users`, `user_liked_recipes`
- Normalized design with proper foreign keys and indexes

## ✨ Key Features

- 🔐 **Secure Authentication**: JWT tokens with Argon2 password hashing
- 🔍 **Smart Recipe Search**: Find recipes by available ingredients with match scoring
- ❤️ **Personal Collections**: Save and manage favorite recipes
- 🛒 **Shopping Lists**: Generate ingredient lists from recipes
- 📱 **Responsive Design**: Mobile-first UI with Tailwind CSS 4
- ☁️ **Cloud-Native**: Containerized deployment on AWS
- 🐳 **Multi-Environment**: Docker Compose for local, Kubernetes for staging, AWS for production

## 🏗️ Architecture

### Production (AWS)

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

### Local (Kubernetes)

```
┌────────────────────────────────────────────┐
│  Kubernetes Cluster (Minikube)             │
│                                            │
│  Frontend Service (3 replicas)            │
│  └─ Next.js + React + TypeScript          │
│                                            │
│  Backend Service (3 replicas)             │
│  └─ Go REST API                            │
│                                            │
│  PostgreSQL Service (1 replica)           │
│  └─ Database + PersistentVolume (5Gi)     │
└────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **For Docker Compose**: Docker Desktop
- **For Kubernetes**: Docker Desktop + Minikube + kubectl
- **For AWS**: AWS CLI configured + Terraform installed
- Go 1.22+ and Node.js 18+ (for local development without containers)

### Option 1: Docker Compose (Fastest)

```bash
# Clone repository
git clone https://github.com/ngthecoder/go_web_api.git
cd go_web_api

# Start all services
docker-compose up --build

# Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8000
# Database: localhost:5432
```

### Option 2: Kubernetes with Minikube (Learning K8s)

```bash
# Start Minikube
minikube start --cpus=4 --memory=8192

# Build images in Minikube's Docker daemon
eval $(minikube docker-env)
cd backend && docker build -t recipe-backend:v1 .
cd ../frontend && docker build -t recipe-frontend:v1 .
cd ..

# Deploy to Kubernetes
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/frontend.yaml

# Wait for pods to be ready
kubectl wait --for=condition=ready pod --all -n recipe-app --timeout=120s

# Access application via port-forward
kubectl config set-context --current --namespace=recipe-app
kubectl port-forward svc/frontend 3000:3000 &
kubectl port-forward svc/backend 8000:8000 &

# Open browser
open http://localhost:3000
```

### Option 3: Local Development (No Containers)

```bash
# Backend
cd backend
cat > .env << EOF
DATABASE_URL=postgresql://recipeadmin:password@localhost:5432/recipes?sslmode=disable
JWT_SECRET=local-dev-secret-at-least-32-characters-long
PORT=8000
ALLOWED_ORIGINS=http://localhost:3000
ENVIRONMENT=development
EOF
go run main.go  # Runs on http://localhost:8000

# Frontend (in new terminal)
cd frontend
cat > .env.local << EOF
NEXT_PUBLIC_API_URL=http://localhost:8000
EOF
npm install
npm run dev  # Runs on http://localhost:3000
```

## ☁️ AWS Production Deployment

### Prerequisites

- **AWS Account** with CLI configured (`aws configure`)
- **Terraform** 1.0+ installed
- **Docker** installed
- **Git** for version control

### Environment Configuration

Customize these for your environment in `terraform/variables.tf`:

- **AWS Region**: Default is `us-east-1`
- **Project Name**: Default is `recipe-app`
- **Resource Sizing**: Adjust CPU/memory based on your needs
- **Replica Count**: Default 1 backend + 1 frontend (increase for production)

### Step 1: Clone Repository

```bash
git clone https://github.com/yourusername/go_web_api.git
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
terraform output backend_ecr_repository_url
terraform output frontend_ecr_repository_url
```

**What gets created**: VPC, subnets, security groups, RDS PostgreSQL, ECR repositories, ECS cluster, Application Load Balancers, CloudWatch logs, Secrets Manager (~$40-50/month)

### Step 3: Build and Push Backend

```bash
cd ../backend/

# Build Docker image
docker build -t recipe-app-backend .

# Get your ECR URL from terraform output
BACKEND_ECR=$(cd ../terraform && terraform output -raw backend_ecr_repository_url)

# Tag for ECR
docker tag recipe-app-backend:latest $BACKEND_ECR:latest

# Login to ECR
AWS_REGION=us-east-1  # Change if you used different region
aws ecr get-login-password --region $AWS_REGION | \
  docker login --username AWS --password-stdin $BACKEND_ECR

# Push to ECR
docker push $BACKEND_ECR:latest
```

### Step 4: Build and Push Frontend

**Important**: Frontend needs the backend URL at build time.

```bash
cd ../frontend/

# Get backend ALB URL from terraform
BACKEND_URL=$(cd ../terraform && terraform output -raw backend_url)
FRONTEND_ECR=$(cd ../terraform && terraform output -raw frontend_ecr_repository_url)

# Build with backend URL
docker build \
  --build-arg NEXT_PUBLIC_API_URL=$BACKEND_URL \
  -t recipe-app-frontend .

# Tag for ECR
docker tag recipe-app-frontend:latest $FRONTEND_ECR:latest

# Push to ECR (already logged in from step 3)
docker push $FRONTEND_ECR:latest
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

# Get your application URL
cd terraform/
terraform output frontend_url
```

**Access your app**: Open the frontend URL in your browser

## 🔄 Update Workflows

### Update Backend Code

```bash
cd backend/
# Make your code changes

docker build -t recipe-app-backend .
BACKEND_ECR=$(cd ../terraform && terraform output -raw backend_ecr_repository_url)
docker tag recipe-app-backend:latest $BACKEND_ECR:latest
docker push $BACKEND_ECR:latest

aws ecs update-service \
  --cluster recipe-app-cluster \
  --service recipe-app-backend-service \
  --force-new-deployment \
  --region us-east-1
```

### Update Frontend Code

```bash
cd frontend/
# Make your code changes

BACKEND_URL=$(cd ../terraform && terraform output -raw backend_url)
FRONTEND_ECR=$(cd ../terraform && terraform output -raw frontend_ecr_repository_url)

docker build --build-arg NEXT_PUBLIC_API_URL=$BACKEND_URL -t recipe-app-frontend .
docker tag recipe-app-frontend:latest $FRONTEND_ECR:latest
docker push $FRONTEND_ECR:latest

aws ecs update-service \
  --cluster recipe-app-cluster \
  --service recipe-app-frontend-service \
  --force-new-deployment \
  --region us-east-1
```

### Update Infrastructure

```bash
cd terraform/
# Edit .tf files

terraform plan    # Review changes
terraform apply   # Apply changes
```

## 🛑 Teardown

### Stop Services (saves ~80% cost, keeps data)

```bash
aws ecs update-service --cluster recipe-app-cluster --service recipe-app-backend-service --desired-count 0 --region us-east-1
aws ecs update-service --cluster recipe-app-cluster --service recipe-app-frontend-service --desired-count 0 --region us-east-1
```

### Restart Services

```bash
aws ecs update-service --cluster recipe-app-cluster --service recipe-app-backend-service --desired-count 1 --region us-east-1
aws ecs update-service --cluster recipe-app-cluster --service recipe-app-frontend-service --desired-count 1 --region us-east-1
```

### Destroy Everything

```bash
cd terraform/
terraform destroy
# Type 'yes' when prompted
```

**Warning**: This deletes all data including the database and Docker images in ECR.

## 🐳 Docker Configuration

### Multi-Stage Builds

**Backend** (Alpine-based, CGO-enabled for PostgreSQL)
```dockerfile
# Stage 1: Build with Go 1.24 + gcc
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev
# ... build process

# Stage 2: Runtime with Alpine + PostgreSQL libs
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
# ... copy binary and run
```
**Image size**: ~50MB

**Frontend** (Node.js standalone output)
```dockerfile
# Stage 1: Install dependencies
# Stage 2: Build Next.js with backend URL
# Stage 3: Runtime with minimal Node.js
```
**Image size**: ~150MB

### Docker Compose Services

```yaml
services:
  postgres:    # PostgreSQL 15 with volume persistence
  backend:     # Go API (port 8000)
  frontend:    # Next.js UI (port 3000)
```

## ☸️ Kubernetes Features

### Resources

- **ConfigMaps**: Application configuration (ports, URLs, database connection)
- **Secrets**: Sensitive data (passwords, JWT secret)
- **Deployments**: 
  - Frontend (3 replicas)
  - Backend (3 replicas)
  - PostgreSQL (1 replica with PersistentVolume)
- **Services**:
  - Frontend: NodePort (30300)
  - Backend: NodePort (30800)
  - PostgreSQL: ClusterIP (internal only)

### Key Features

- **Persistent Storage**: 5Gi PersistentVolumeClaim for database
- **Load Balancing**: Service-level load balancing across replicas
- **Self-Healing**: Automatic pod restart on failure
- **Rolling Updates**: Zero-downtime deployments
- **Namespace Isolation**: All resources in `recipe-app` namespace

## 🔐 Security Features

**Network Security**
- Database in private subnet with no internet access
- Security groups with strict firewall rules (ALB → ECS → RDS only)
- SSL/TLS required for RDS connections

**Application Security**
- Password Hashing: Argon2id (time=3, memory=64MB, 16-byte salt)
- JWT Tokens: HMAC-SHA256, 24-hour expiration
- CORS Protection: Configured allowed origins
- Input Validation: All endpoints validate parameters
- SQL Injection Prevention: Parameterized queries

**Infrastructure Security**
- AWS Secrets Manager for credentials
- IAM roles with least-privilege permissions
- No hardcoded secrets in code or containers

## 📊 API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/auth/register` | POST | No | User registration |
| `/api/auth/login` | POST | No | User login |
| `/api/recipes` | GET | Optional | Browse recipes with filters |
| `/api/recipes/{id}` | GET | Optional | Recipe details |
| `/api/recipes/find-by-ingredients` | GET | Optional | Find recipes by ingredients |
| `/api/recipes/shopping-list/{id}` | GET | No | Generate shopping list |
| `/api/ingredients` | GET | No | Browse ingredients |
| `/api/ingredients/{id}` | GET | No | Ingredient details |
| `/api/user/profile` | GET | Yes | User profile |
| `/api/user/liked-recipes` | GET | Yes | User's liked recipes |
| `/api/user/liked-recipes/add` | POST | Yes | Add liked recipe |
| `/api/user/liked-recipes/{id}` | DELETE | Yes | Remove liked recipe |
| `/api/user/profile/update` | PUT | Yes | Update profile |
| `/api/user/password` | PUT | Yes | Change password |
| `/api/user/account` | DELETE | Yes | Delete account |
| `/api/categories` | GET | No | Category statistics |
| `/api/stats` | GET | No | Overall statistics |

## 🧪 Testing

```bash
cd backend
go test ./... -v -cover

# Current coverage: ~40%
# Focus: Authentication security, core business logic
```

## 📁 Project Structure

```
go_web_api/
├── backend/
│   ├── internal/
│   │   ├── auth/           # JWT & authentication
│   │   ├── database/       # DB connection & migrations
│   │   ├── recipes/        # Recipe business logic
│   │   ├── ingredients/    # Ingredient management
│   │   ├── users/          # User profile management
│   │   └── stats/          # Statistics
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── app/
│   │   ├── recipes/        # Recipe pages
│   │   ├── ingredients/    # Ingredient pages
│   │   ├── profile/        # User profile
│   │   └── page.tsx        # Home page
│   ├── components/         # Reusable React components
│   ├── contexts/           # Auth context
│   ├── Dockerfile
│   └── package.json
├── k8s/                    # Kubernetes manifests
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── postgres.yaml
│   ├── backend.yaml
│   └── frontend.yaml
├── terraform/              # AWS infrastructure
│   ├── main.tf             # Provider configuration
│   ├── vpc.tf              # VPC, subnets, networking
│   ├── rds.tf              # PostgreSQL database
│   ├── ecs.tf              # ECS cluster and services
│   ├── alb.tf              # Application Load Balancers
│   ├── ecr.tf              # Container registry
│   ├── secrets.tf          # Secrets Manager
│   ├── security.tf         # Security groups
│   ├── cloudwatch.tf       # Logging
│   ├── variables.tf        # Configuration variables
│   └── outputs.tf          # Output values
├── docker-compose.yml      # Local development
└── README.md
```

## 🎓 What I Learned

**Infrastructure as Code**
- Terraform for reproducible AWS deployments
- Managing state and dependencies between resources
- Cost optimization with proper resource sizing

**Container Orchestration**
- Docker multi-stage builds for optimization
- Kubernetes deployments with multiple replicas
- ECS Fargate for serverless container management

**Cloud Networking**
- VPCs, subnets (public vs private), security groups
- Application Load Balancers and target groups
- Network isolation and firewall rules

**Security Best Practices**
- Secrets management with AWS Secrets Manager
- SSL/TLS for database connections
- Least-privilege IAM roles
- Password hashing and JWT authentication

**Full-Stack Development**
- RESTful API design with Go
- Server-side rendering with Next.js
- Database schema design and migrations
- CORS and cross-origin authentication

## 🚧 Future Enhancements

- HTTPS with ACM certificates
- Custom domain with Route 53
- CI/CD pipeline with GitHub Actions
- Monitoring with Prometheus + Grafana
- Rate limiting and API throttling
- Email verification and password reset
- Recipe image uploads to S3
