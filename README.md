# Go Recipe Web API + Next.js Frontend

A production-ready full-stack recipe discovery application with secure authentication, built to demonstrate modern containerized application development and deployment skills.

## 🎯 Project Overview

A full-stack web application that helps users discover recipes, manage their favorite dishes, and generate shopping lists based on available ingredients. Built with Go backend, Next.js frontend, and deployed using Docker and Kubernetes.

**Live Demo**: [Coming Soon]

## ✨ Key Features

- 🔐 **Secure Authentication**: JWT tokens with Argon2 password hashing
- 🔍 **Smart Recipe Search**: Find recipes by available ingredients with match scoring
- ❤️ **Personal Collections**: Save and manage favorite recipes
- 🛒 **Shopping Lists**: Generate shopping lists based on selected recipes
- 📱 **Responsive Design**: Mobile-first UI with Tailwind CSS 4
- 🐳 **Container-Ready**: Docker and Kubernetes deployment configurations
- ☸️ **Production-Grade**: Kubernetes manifests with ConfigMaps, Secrets, and persistent storage

## 🛠️ Tech Stack

**Backend**
- Go 1.24.0
- PostgreSQL 15 (production) / SQLite3 (development)
- JWT authentication with Argon2id password hashing
- RESTful API with 14 endpoints

**Frontend**
- Next.js 15.5.0 with React 19.1.0
- TypeScript
- Tailwind CSS 4
- Context API for state management

**Infrastructure**
- Docker with multi-stage builds
- Kubernetes with Minikube (local) / AWS EKS (production-ready)
- PostgreSQL with PersistentVolumes
- Load balancing with 3 frontend and 3 backend replicas

## 🚀 Quick Start

### Prerequisites

- **For Docker Compose**: Docker Desktop
- **For Kubernetes**: Docker Desktop + Minikube + kubectl
- Go 1.22+ and Node.js 18+ (for local development)

### Option 1: Docker Compose (Recommended for Quick Start)

```bash
# Clone repository
git clone https://github.com/ngthecoder/go_web_api.git
cd go_web_api

# Start all services
docker-compose up --build

# Access application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8000
```

### Option 2: Kubernetes with Minikube

```bash
# Start Minikube
minikube start --cpus=4 --memory=8192

# Build images in Minikube
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

### Option 3: Local Development

```bash
# Backend
cd backend
echo "JWT_SECRET=your-secret-key-here" > .env
go run main.go  # Runs on http://localhost:8000

# Frontend (in new terminal)
cd frontend
npm install
npm run dev  # Runs on http://localhost:3000
```

## 📦 Architecture

### System Overview

```
┌─────────────────────────────────────────────────────┐
│  Kubernetes Cluster                                 │
│                                                     │
│  Frontend Service (3 replicas)                     │
│  └─ Next.js + React + TypeScript                   │
│                                                     │
│  Backend Service (3 replicas)                      │
│  └─ Go REST API                                    │
│                                                     │
│  PostgreSQL Service (1 replica)                    │
│  └─ Database + PersistentVolume (5Gi)             │
└─────────────────────────────────────────────────────┘
```

### Database Schema

**5 Tables**: `ingredients`, `recipes`, `recipe_ingredients` (junction), `users`, `user_liked_recipes`

**Key Relationships**:
- Many-to-many: Recipes ↔ Ingredients (via junction table)
- One-to-many: Users → Liked Recipes

### API Endpoints

| Endpoint | Method | Description | Auth |
|----------|--------|-------------|------|
| `/auth/register` | POST | User registration | No |
| `/auth/login` | POST | User login | No |
| `/recipes` | GET | Browse recipes with filters | Optional |
| `/recipes/{id}` | GET | Recipe details | Optional |
| `/recipes/find-by-ingredients` | GET | Find recipes by ingredients | Optional |
| `/recipes/shopping-list/{id}` | GET | Generate shopping list | No |
| `/ingredients` | GET | Browse ingredients | No |
| `/ingredients/{id}` | GET | Ingredient details | No |
| `/user/profile` | GET | User profile | Yes |
| `/user/liked-recipes` | GET | User's liked recipes | Yes |
| `/user/liked-recipes/add` | POST | Add liked recipe | Yes |
| `/user/liked-recipes/{id}` | DELETE | Remove liked recipe | Yes |
| `/categories` | GET | Category statistics | No |
| `/stats` | GET | Overall statistics | No |

## 🐳 Docker Configuration

### Multi-Stage Builds

**Backend** (Alpine-based, CGO-enabled for PostgreSQL)
- Stage 1: Build with Go 1.24 + gcc
- Stage 2: Runtime with Alpine + PostgreSQL libs
- Image size: ~50MB

**Frontend** (Node.js standalone output)
- Stage 1: Install dependencies
- Stage 2: Build Next.js
- Stage 3: Runtime with minimal Node.js
- Image size: ~150MB

### Docker Compose Services

```yaml
services:
  postgres:    # PostgreSQL 15 with volume persistence
  backend:     # Go API (port 8000)
  frontend:    # Next.js UI (port 3000)
```

## ☸️ Kubernetes Deployment

### Resources

**ConfigMaps**: Application configuration (ports, URLs, database connection)
**Secrets**: Sensitive data (passwords, JWT secret)
**Deployments**: 
- Frontend (3 replicas)
- Backend (3 replicas)
- PostgreSQL (1 replica with PersistentVolume)

**Services**:
- Frontend: NodePort (30300)
- Backend: NodePort (30800)
- PostgreSQL: ClusterIP (internal only)

### Key Features

- **Persistent Storage**: 5Gi PersistentVolumeClaim for database
- **Load Balancing**: Service-level load balancing across replicas
- **Self-Healing**: Automatic pod restart on failure
- **Rolling Updates**: Zero-downtime deployments
- **Health Checks**: Liveness and readiness probes (coming soon)

## 🔐 Security Features

- **Password Hashing**: Argon2id with 16-byte salt (time=3, memory=64MB)
- **JWT Tokens**: HMAC-SHA256 signed, 24-hour expiration
- **CORS Protection**: Configured allowed origins
- **Input Validation**: All endpoints validate parameters
- **SQL Injection Prevention**: Parameterized queries
- **Secrets Management**: Kubernetes Secrets for sensitive data

## 📊 Testing

```bash
cd backend
go test ./... -v -cover

# Current coverage: ~40%
# Focus: Authentication security, core business logic
```

## 🎓 What I Learned

### Technical Skills

**Containerization**: 
- Docker multi-stage builds for optimization
- docker-compose for local development orchestration
- Image building and optimization techniques

**Kubernetes**:
- Deployed full-stack application with 7 pods
- Configured Services, Deployments, ConfigMaps, Secrets
- Set up PersistentVolumes for data persistence
- Implemented load balancing with multiple replicas
- Used port-forwarding for local development

**Database Migration**:
- Migrated from SQLite to PostgreSQL
- Adapted SQL placeholder syntax (`?` → `$1, $2`)
- Implemented dynamic connection string building

**Full-Stack Development**:
- Built RESTful API with Go
- Created responsive UI with Next.js 15 and React 19
- Implemented JWT authentication flow
- Handled CORS and networking challenges

**Debugging & Problem-Solving**:
- Debugged Docker networking issues (Docker driver limitations)
- Resolved Kubernetes namespace context problems
- Fixed CORS origin mismatches
- Troubleshot pod crashes and connection failures

## 📁 Project Structure

```
go_web_api/
├── backend/
│   ├── internal/
│   │   ├── auth/           # Authentication logic
│   │   ├── recipes/        # Recipe business logic
│   │   ├── ingredients/    # Ingredient management
│   │   └── users/          # User profile management
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── app/
│   │   ├── page.tsx        # Home page
│   │   ├── recipes/        # Recipe pages
│   │   ├── ingredients/    # Ingredient pages
│   │   └── profile/        # User profile
│   ├── components/         # Reusable components
│   ├── contexts/           # AuthContext
│   ├── Dockerfile
│   └── package.json
├── k8s/
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── postgres.yaml
│   ├── backend.yaml
│   └── frontend.yaml
├── docker-compose.yml
└── README.md
```

## 🚧 Future Enhancements

- [ ] Ingress controller for better routing
- [ ] Horizontal Pod Autoscaling (HPA)
- [ ] Health checks (liveness/readiness probes)
- [ ] CI/CD pipeline with GitHub Actions
- [ ] Monitoring with Prometheus + Grafana
- [ ] AWS EKS deployment
- [ ] Recipe ratings and reviews
- [ ] Email verification and password reset
