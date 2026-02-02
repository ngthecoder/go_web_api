# CLAUDE.md - AI Assistant Guide for go_web_api

## Project Overview

A full-stack recipe discovery application with a **Go backend** and **Next.js frontend**. Users can browse recipes, search by ingredients, manage favorites, and generate shopping lists. The application features JWT authentication with Argon2id password hashing.

## Tech Stack

### Backend
- **Language**: Go 1.24
- **Database**: PostgreSQL 16 (production) / 15 (local)
- **ORM**: Raw SQL with `database/sql` and `lib/pq` driver
- **Testing**: Uses SQLite in-memory databases via `mattn/go-sqlite3`
- **Authentication**: Custom JWT implementation (HMAC-SHA256, 24h expiry)
- **Password Hashing**: Argon2id

### Frontend
- **Framework**: Next.js 15.5.0 with React 19.1.0
- **Language**: TypeScript
- **Styling**: Tailwind CSS 4
- **State Management**: React Context API
- **Build Tool**: Turbopack

### Infrastructure
- **Local**: Docker Compose, Kubernetes (Minikube)
- **Production**: AWS (Terraform, ECS Fargate, RDS, ECR, ALB)
- **Container Registry**: AWS ECR
- **Secrets**: AWS Secrets Manager

## Directory Structure

```
go_web_api/
├── backend/
│   ├── main.go                    # Entry point, routing, middleware
│   ├── internal/
│   │   ├── auth/                  # JWT auth, password hashing
│   │   │   ├── handlers.go        # HTTP handlers for auth endpoints
│   │   │   ├── services.go        # Business logic (register, login, JWT)
│   │   │   ├── models.go          # Request/response structs
│   │   │   ├── password.go        # Argon2id implementation
│   │   │   └── services_test.go   # Auth tests
│   │   ├── database/
│   │   │   ├── database.go        # DB init, table creation, indexes
│   │   │   └── seed.go            # Sample data seeding
│   │   ├── errors/
│   │   │   └── http_errors.go     # HTTPError type and helpers
│   │   ├── recipes/               # Recipe CRUD and search
│   │   ├── ingredients/           # Ingredient management
│   │   ├── users/                 # User profile management
│   │   └── stats/                 # Statistics endpoints
│   ├── Dockerfile                 # Multi-stage Alpine build
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── app/                       # Next.js App Router pages
│   │   ├── page.tsx               # Home page
│   │   ├── layout.tsx             # Root layout with AuthProvider
│   │   ├── recipes/[id]/          # Recipe detail page
│   │   ├── ingredients/           # Ingredients pages
│   │   ├── login/                 # Login page
│   │   ├── register/              # Registration page
│   │   ├── profile/               # User profile pages
│   │   ├── liked-recipes/         # Saved recipes
│   │   └── find-recipes/          # Recipe search by ingredients
│   ├── components/                # Reusable React components
│   ├── contexts/
│   │   └── AuthContext.tsx        # Auth state management
│   ├── lib/
│   │   ├── api-config.ts          # API endpoint definitions
│   │   ├── api.ts                 # API helper functions
│   │   ├── auth.ts                # Auth API calls
│   │   └── types.ts               # TypeScript interfaces
│   ├── Dockerfile                 # Multi-stage Node.js build
│   ├── package.json
│   └── tsconfig.json
├── k8s/                           # Kubernetes manifests
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── postgres.yaml
│   ├── backend.yaml
│   └── frontend.yaml
├── terraform/                     # AWS infrastructure
│   ├── main.tf                    # Provider config
│   ├── vpc.tf                     # VPC, subnets
│   ├── rds.tf                     # PostgreSQL RDS
│   ├── ecs.tf                     # ECS cluster, services
│   ├── alb.tf                     # Load balancers
│   ├── ecr.tf                     # Container registry
│   ├── secrets.tf                 # Secrets Manager
│   ├── security.tf                # Security groups
│   ├── cloudwatch.tf              # Logging
│   ├── variables.tf               # Config variables
│   └── outputs.tf                 # Output values
├── docker-compose.yml             # Local development
└── README.md
```

## Development Commands

### Backend (Go)
```bash
cd backend

# Run locally (requires .env file)
go run main.go

# Run tests
go test ./... -v

# Run tests with coverage
go test ./... -v -cover

# Build binary
go build -o main .
```

### Frontend (Next.js)
```bash
cd frontend

# Install dependencies
npm install

# Development server (with Turbopack)
npm run dev

# Production build
npm run build

# Start production server
npm start

# Lint
npm run lint
```

### Docker Compose (Full Stack Local)
```bash
# Start all services
docker-compose up --build

# Start in background
docker-compose up -d --build

# Stop services
docker-compose down

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Kubernetes (Minikube)
```bash
# Start minikube
minikube start --cpus=4 --memory=8192

# Build images in minikube's docker
eval $(minikube docker-env)
cd backend && docker build -t recipe-backend:v1 .
cd ../frontend && docker build -t recipe-frontend:v1 .

# Deploy
kubectl apply -f k8s/

# Port forward
kubectl port-forward svc/frontend 3000:3000 -n recipe-app &
kubectl port-forward svc/backend 8000:8000 -n recipe-app &
```

## Environment Variables

### Backend (.env)
```bash
DATABASE_URL=postgresql://user:pass@host:5432/recipes?sslmode=disable
# OR individual components:
POSTGRES_USER=recipeuser
POSTGRES_PASSWORD=recipepassword
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=recipes

JWT_SECRET=your-secret-at-least-32-characters-long
PORT=8000
ALLOWED_ORIGINS=http://localhost:3000
```

### Frontend (.env.local)
```bash
NEXT_PUBLIC_API_URL=http://localhost:8000
```

## Code Conventions

### Backend (Go)

1. **Package Structure**: Each domain (auth, recipes, ingredients, users, stats) has:
   - `handlers.go` - HTTP handlers
   - `services.go` - Business logic
   - `models.go` - Data structures

2. **Handler Pattern**:
```go
type SomeHandler struct {
    service *SomeService
}

func NewSomeHandler(s *SomeService) *SomeHandler {
    return &SomeHandler{service: s}
}

func (h *SomeHandler) HandleSomething(w http.ResponseWriter, r *http.Request) {
    // Method check
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Business logic via service
    result, err := h.service.DoSomething()
    if err != nil {
        errors.WriteHTTPError(w, err)
        return
    }

    // JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

3. **Error Handling**: Use `internal/errors` package:
```go
errors.NewBadRequestError("Invalid input")
errors.NewNotFoundError("Recipe not found")
errors.NewInternalServerError("DB error", err)
errors.WriteHTTPError(w, err)
```

4. **Middleware Stack**: Handlers are wrapped: `loggingMiddleware(enableCORS(allowedOrigins, handler))`

5. **Auth Context**: User ID is passed via context:
```go
userID := ""
if userIDValue := r.Context().Value("user_id"); userIDValue != nil {
    userID = userIDValue.(string)
}
```

6. **SQL Queries**: Use parameterized queries with `$1`, `$2` placeholders (PostgreSQL style)

### Frontend (TypeScript/React)

1. **API Configuration**: All endpoints defined in `lib/api-config.ts`

2. **Auth Pattern**: Use `useAuth()` hook from `contexts/AuthContext.tsx`:
```tsx
const { user, token, isAuthenticated, login, logout } = useAuth();
```

3. **API Calls**: Pass Bearer token in Authorization header:
```typescript
const response = await fetch(endpoint, {
    headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    }
});
```

4. **Page Components**: Use Next.js App Router with `'use client'` directive for client components

5. **Types**: Define interfaces in `lib/types.ts`

## Database Schema

5 tables with proper foreign keys:

- **ingredients**: id, name, category, calories_per_100g, description
- **recipes**: id, name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description
- **recipe_ingredients**: recipe_id, ingredient_id, quantity, unit, notes (composite PK)
- **users**: id (UUID), username, email, password_hash, created_at, updated_at
- **user_liked_recipes**: user_id, recipe_id, created_at (composite PK, CASCADE delete)

## API Endpoints

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
| `/api/health` | GET | No | Health check |

## Testing

### Backend Tests
- Tests use **SQLite in-memory** databases for isolation
- Each test file has a `setupTestDB()` function that creates schema and seed data
- Test files are co-located with source: `services_test.go` next to `services.go`

```bash
cd backend
go test ./... -v -cover
```

### Test Pattern
```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    // Create tables and seed data
    return db
}

func TestSomething(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    service := NewSomeService(db)
    // Test assertions
}
```

## Security Notes

- **Passwords**: Hashed with Argon2id (time=3, memory=64MB, 16-byte salt)
- **JWT**: HMAC-SHA256 signed, 24-hour expiration, custom implementation
- **CORS**: Configured via `ALLOWED_ORIGINS` environment variable
- **SQL Injection**: All queries use parameterized statements
- **Secrets**: Production uses AWS Secrets Manager

## Common Tasks

### Adding a New API Endpoint

1. Create/update model structs in `internal/{domain}/models.go`
2. Add business logic in `internal/{domain}/services.go`
3. Create handler in `internal/{domain}/handlers.go`
4. Register route in `main.go` with middleware

### Adding a New Frontend Page

1. Create page in `app/{route}/page.tsx`
2. Add any new types to `lib/types.ts`
3. Add API endpoint to `lib/api-config.ts` if needed
4. Create API helper in `lib/api.ts` if needed

### Updating Database Schema

1. Modify table creation in `internal/database/database.go`
2. Update seed data in `internal/database/seed.go` if needed
3. Update relevant models and queries in domain packages

## Deployment

### Local: Docker Compose
```bash
docker-compose up --build
# Frontend: http://localhost:3000
# Backend: http://localhost:8000
```

### Production: AWS via Terraform
```bash
cd terraform
terraform init
terraform plan
terraform apply

# Build and push images
# Update ECS services
```

See README.md for detailed AWS deployment instructions.
