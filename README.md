# Go Recipe Web API + Next.js
A web application using Go backend API and Next.js frontend with user authentication

## Project Purpose
- Learn Web API development using Go language
- Build a cross-reference system for ingredients and recipes
- Gain team development experience
- Implement secure user authentication and authorization

## Tech Stack
- **Backend**: Go 1.22.4 + SQLite3
- **Frontend**: Next.js 15.5.0 + TypeScript + Tailwind CSS
- **Database**: SQLite3 (for local development)
- **Authentication**: JWT tokens with Argon2 password hashing
- **Deployment**: Docker + AWS ECS with Terraform

## Prerequisites
Before starting, ensure the following are installed:
1. Go 1.22 or higher
2. Node.js 18 or higher
3. Git
4. Text editor like VSCode
5. Docker Desktop (for production, optional)

## Getting Started
1. Clone the repository
```bash
git clone git@github.com:ngthecoder/go_web_api.git
cd go_web_api
```

2. Start the backend
```bash
cd backend
go run main.go
```
- Runs on http://localhost:8000

3. Start the frontend
```bash
cd frontend
npm install
npm run dev
```
- Runs on http://localhost:3000

## System Architecture

### Database Design

#### Table Structure

**ingredients (Ingredients Table)**
- id (INTEGER)
- name (TEXT)
- category (TEXT)
- calories_per_100g (INTEGER)
- description (TEXT)
- PRIMARY KEY: id

**recipes (Recipes Table)**
- id (INTEGER)
- name (TEXT)
- category (TEXT)
- prep_time_minutes (INTEGER)
- cook_time_minutes (INTEGER)
- servings (INTEGER)
- difficulty (TEXT)
- instructions (TEXT)
- description (TEXT)
- PRIMARY KEY: id

**recipe_ingredients (Recipe-Ingredient Junction Table)**
- recipe_id (INTEGER)
- ingredient_id (INTEGER)
- quantity (REAL)
- unit (TEXT)
- notes (TEXT)
- PRIMARY KEY: recipe_id, ingredient_id
- FOREIGN KEY: recipe_id, ingredient_id

**users (User Authentication Table)**
- id (TEXT) - UUID primary key
- username (TEXT) - unique username
- email (TEXT) - unique email address
- password_hash (TEXT) - Argon2 hashed password
- created_at (DATETIME)
- updated_at (DATETIME)

**user_liked_recipes (User Preferences Table)**
- user_id (TEXT)
- recipe_id (INTEGER)
- created_at (DATETIME)
- PRIMARY KEY: user_id, recipe_id
- FOREIGN KEY: user_id, recipe_id

#### Why the recipe_ingredients Table is Necessary

##### Bad Example
recipes table:
```
| id | name        | ingredients_list          |
|----|-------------|---------------------------|
| 1  | Tomato Rice | "tomato, onion, rice"     |
| 2  | Tomato Pasta| "tomato, onion, garlic"   |
```

Problems:
- Difficult to search (e.g., can't search for recipes using only tomatoes)
- Cannot store quantity information

##### Good Example
ingredients table:
```
| id | name   |
|----|--------|
| 1  | Tomato |
| 2  | Onion  |
```

recipes table:
```
| id | name         |
|----|--------------|
| 1  | Tomato Rice  |
| 2  | Tomato Pasta |
```

recipe_ingredients table:
```
| recipe_id | ingredient_id | quantity | unit  | notes     |
|-----------|---------------|----------|-------|-----------|
| 1         | 1             | 2        | pieces| diced     |
| 1         | 2             | 1        | piece | minced    |
| 2         | 1             | 3        | pieces| sliced    |
| 2         | 2             | 1        | piece | minced    |
```

Benefits:
- Efficient searching
- Detailed information storage
- Flexibility

#### Database Relationship Diagram
```
ingredients (Ingredients)    recipe_ingredients (Junction)    recipes (Recipes)
┌──────────────┐            ┌─────────────────────┐          ┌─────────────────┐
│ id (PK)      │◄───────────┤ ingredient_id (FK)  │          │ id (PK)         │
│ name         │            │ recipe_id (FK)      ├─────────►│ name            │
│ category     │            │ quantity            │          │ category        │
│ calories     │            │ unit                │          │ prep_time       │
│ description  │            │ notes               │          │ cook_time       │
└──────────────┘            └─────────────────────┘          │ servings        │
                                                             │ difficulty      │
                                                             │ instructions    │
                                                             │ description     │
                                                             └─────────────────┘
                                                                      △
                                                                      │
users (Users)              user_liked_recipes (User Preferences)      │
┌──────────────┐          ┌─────────────────────┐                     │
│ id (PK)      │◄─────────┤ user_id (FK)        │                     │
│ username     │          │ recipe_id (FK)      ├─────────────────────┘
│ email        │          │ created_at          │
│ password_hash│          └─────────────────────┘
│ created_at   │
│ updated_at   │
└──────────────┘
```

### API Endpoint Design

#### Endpoint List
| HTTP Method | Endpoint | Description | Required Parameters | Optional Parameters | Auth Required |
|-------------|----------|-------------|--------------------|--------------------|---------------|
| POST | `/auth/register` | User registration | username, email, password | - | No |
| POST | `/auth/login` | User login | email, password | - | No |
| GET | `/ingredients` | Get ingredient list, search & filter | - | `search`, `category`, `sort`, `order`, `page`, `limit` | No |
| GET | `/ingredients/{id}` | Get ingredient details + related recipes | `id` | - | No |
| GET | `/recipes` | Get recipe list, search & filter | - | `search`, `category`, `max_time`, `difficulty`, `sort`, `order`, `page`, `limit` | No |
| GET | `/recipes/find-by-ingredients` | Find recipes by available ingredients | `ingredients` | `match_type`, `page`, `limit` | No |
| GET | `/recipes/{id}` | Get recipe details + required ingredients | `id` | - | No |
| GET | `/recipes/shopping-list/{id}` | Generate shopping list for recipe | `id` | `have_ingredients` | No |
| GET | `/categories` | Get category statistics | - | - | No |
| GET | `/stats` | Get overall statistics | - | - | No |
| GET | `/user/profile` | Get user profile | - | - | Yes |
| GET | `/user/liked-recipes` | Get user's liked recipes | - | - | Yes |
| POST | `/user/liked-recipes/add` | Add recipe to liked list | recipe_id | - | Yes |
| DELETE | `/user/liked-recipes/{id}` | Remove recipe from liked list | `id` | - | Yes |

#### Authentication Endpoints

**1: POST /api/auth/register**
```bash
POST /api/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**2: POST /api/auth/login**
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Recipe Search Endpoints

**3: GET /api/recipes/find-by-ingredients**
```bash
GET /api/recipes/find-by-ingredients?ingredients=2,26&match_type=partial&limit=1
```
```json
[
  {
    "id": 5,
    "name": "Omurice",
    "category": "lunch",
    "prep_time_minutes": 15,
    "cook_time_minutes": 20,
    "servings": 2,
    "difficulty": "medium",
    "instructions": "1. Stir-fry onions\n2. Add rice to make fried rice\n3. Make thin omelet with eggs\n4. Wrap fried rice in omelet",
    "description": "Everyone's favorite omurice",
    "matched_ingredients_count": 2,
    "total_ingredients_count": 6,
    "match_score": 0.33333334
  }
]
```

**4: GET /api/recipes/shopping-list/{id}**
```bash
GET /api/recipes/shopping-list/5?have_ingredients=1,4
```
```json
{
  "recipe_id": 5,
  "shopping_list": [
    {
      "ingredient_id": 6,
      "name": "Potato",
      "quantity": 3,
      "unit": "pieces",
      "notes": "diced"
    }
  ]
}
```

**5: GET /api/categories**
```json
{
  "ingredient_categories": [
    {"category": "vegetables", "count": 25},
    {"category": "protein", "count": 10},
    {"category": "grains", "count": 8}
  ],
  "recipe_categories": [
    {"category": "breakfast", "count": 4},
    {"category": "lunch", "count": 6},
    {"category": "dinner", "count": 12}
  ]
}
```

**6: GET /api/stats**
```json
{
  "total_ingredients": 89,
  "total_recipes": 27,
  "avg_prep_time": 12.592593,
  "avg_cook_time": 17.777779,
  "difficulty_distribution": {
    "easy": 12,
    "hard": 2,
    "medium": 13
  }
}
```

#### User Profile Endpoints

**7: GET /api/user/profile**
```bash
GET /api/user/profile
Authorization: Bearer {token}
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**8: GET /api/user/liked-recipes**
```bash
GET /api/user/liked-recipes
Authorization: Bearer {token}
```

Response:
```json
{
  "liked_recipes": [
    {
      "id": 1,
      "name": "Scrambled Eggs",
      "category": "Breakfast",
      "prep_time_minutes": 5,
      "cook_time_minutes": 5,
      "servings": 2,
      "difficulty": "easy",
      "instructions": "1. Crack eggs into bowl...",
      "description": "Fluffy and creamy scrambled eggs"
    }
  ],
  "total": 1
}
```

**9: POST /api/user/liked-recipes/add**
```bash
POST /api/user/liked-recipes/add
Authorization: Bearer {token}
Content-Type: application/json

{
  "recipe_id": 1
}
```

Response:
```json
{
  "message": "Recipe added to liked list"
}
```

**10: DELETE /api/user/liked-recipes/{id}**
```bash
DELETE /api/user/liked-recipes/1
Authorization: Bearer {token}
```

Response:
```json
{
  "message": "Recipe removed from liked list"
}
```

## Frontend Design (Next.js)

### Page Details

#### Main Pages
- **Home (/)**: Recipe browsing with search, filtering, and pagination
- **Ingredients (/ingredients)**: Ingredient browsing with search and filtering
- **Find Recipes (/find-recipes)**: Recipe search by available ingredients
- **Register (/register)**: User registration page
- **Login (/login)**: User login page
- **Profile (/profile)**: User profile management (protected route)

#### Ingredient Detail Page (/ingredients/[id])
- Display ingredient details (calories, category, description)
- List of recipes using this ingredient
- Click recipe cards to navigate to recipe details

#### Recipe Detail Page (/recipes/[id])
- Recipe details (cooking time, servings, difficulty, instructions)
- List of required ingredients (with quantities and units)
- Click ingredients to navigate to ingredient details
- Shopping list generation functionality

### User Authentication Flow

#### Registration Flow
1. User accesses `/register` page
2. Fills out registration form (username, email, password)
3. Frontend sends POST request to `/api/auth/register`
4. Backend validates input and creates user account
5. JWT token returned and stored in localStorage
6. User redirected to profile page

#### Login Flow
1. User accesses `/login` page
2. Enters email and password
3. Frontend sends POST request to `/api/auth/login`
4. Backend validates credentials
5. JWT token returned and stored in localStorage
6. User can access protected features

### User Experience Flows

#### Flow 1: Recipe Search by Available Ingredients
1. Navigate to "Find Recipes" page (/find-recipes)
2. Search and select available ingredients from checkbox list
3. Click "Search Recipes" button
4. View matched recipes with match scores
5. Click on recipe to view full details

#### Flow 2: User Registration and Profile
1. Click "Register" in navigation
2. Fill out registration form
3. Automatically logged in after successful registration
4. Access profile page to view account information
5. Logout functionality available

#### Flow 3: User Profile and Liked Recipes
1. User logs into their account
2. Navigates to profile page to view account information
3. Browses through recipe catalog
4. Adds interesting recipes to their liked list
5. Views their liked recipes collection
6. Removes recipes they're no longer interested in

### Project Structure
```
frontend/
├── app/
│   ├── page.tsx                 # Home page (recipe browsing)
│   ├── ingredients/
│   │   ├── page.tsx             # Ingredients list
│   │   └── [id]/page.tsx        # Ingredient detail
│   ├── recipes/
│   │   └── [id]/page.tsx        # Recipe detail
│   ├── find-recipes/
│   │   └── page.tsx             # Recipe search by ingredients
│   ├── register/
│   │   └── page.tsx             # User registration
│   ├── login/
│   │   └── page.tsx             # User login
│   ├── profile/
│   │   └── page.tsx             # User profile (protected)
│   ├── layout.tsx               # Root layout with navigation
│   └── globals.css              # Global styles (Tailwind CSS)
├── components/
│   └── Navigation.tsx           # Navigation component
├── lib/
│   ├── auth.ts                  # Authentication functions
│   └── types.ts                 # TypeScript interfaces
```

## Security Features

### Password Security
- **Argon2 Hashing**: Passwords are hashed using Argon2id algorithm
- **Salt Generation**: Random 16-byte salt generated for each password
- **Memory-hard Function**: Argon2 parameters: time=3, memory=64MB, threads=2

### JWT Token Security
- **HMAC-SHA256 Signing**: Tokens signed with secret key
- **24-Hour Expiration**: Tokens expire after 24 hours
- **Stateless Authentication**: No server-side session storage required

### API Security
- **CORS Protection**: Configured for frontend origin only
- **Input Validation**: All endpoints validate input parameters
- **Error Handling**: Secure error messages without sensitive information

## Features

### Core Features
1. **Recipe Management**: Browse, search, and filter recipes by various criteria
2. **Ingredient Management**: Browse ingredients with detailed nutritional information
3. **Recipe Discovery**: Find recipes based on available ingredients with match scoring
4. **Shopping Lists**: Generate shopping lists for recipes based on available ingredients
5. **User Authentication**: Secure user registration and login system
6. **User Profiles**: Personal user accounts with profile management
7. **Liked Recipes**: Users can save and manage their favorite recipes

### User Stories
1. **Browse Recipes**: Users can search and filter recipes by category, difficulty, cooking time
2. **Find Recipes by Ingredients**: Users can select available ingredients and find matching recipes
3. **View Details**: Users can view detailed recipe instructions and ingredient requirements
4. **User Registration**: New users can create accounts securely
5. **User Login**: Existing users can log in to access their profiles
6. **User Profile Management**: Users can view their profile information
7. **Liked Recipes**: Users can add and remove recipes from their liked list
8. **Generate Shopping Lists**: Users can create shopping lists for recipes they want to make
9. **Bidirectional Navigation**: Navigate between ingredients and recipes seamlessly

### Use Cases

**Scenario 1: New User Registration**
1. User visits the application
2. Clicks "Register" in navigation
3. Fills out username, email, and password
4. System creates account and logs user in automatically
5. User can now access profile and future features

**Scenario 2: Recipe Discovery Workflow**
1. User selects "Find Recipes" from navigation
2. Searches and selects available ingredients (e.g., tomato, onion, rice)
3. System shows matching recipes with match scores
4. User clicks on recipe to view full instructions
5. User can generate shopping list for missing ingredients

**Scenario 3: User Profile and Liked Recipes**
1. User logs into their account
2. Navigates to profile page to view account information
3. Browses through recipe catalog
4. Adds interesting recipes to their liked list
5. Views their liked recipes collection
6. Removes recipes they're no longer interested in

## Team Development Workflow

### Backend Tasks
- Database schema design and implementation
- User authentication system (registration, login, JWT)
- Password security with Argon2 hashing
- Basic recipe and ingredient CRUD endpoints
- Recipe search by ingredients functionality
- Shopping list generation
- Search, filtering, and pagination
- Statistics and category endpoints
- CORS and security middleware
- User-specific features (liked recipes, user preferences)
- API documentation server
- Rate limiting and advanced security features
- Custom error handling with HTTP status codes
- Clean service layer architecture
- Input validation and sanitization

### Frontend Tasks (Next.js + TypeScript)
- Next.js project setup with TypeScript
- User authentication UI (register, login, profile)
- Recipe and ingredient browsing pages
- Recipe search by ingredients interface
- Responsive design with Tailwind CSS
- Navigation component with auth state
- Local storage for authentication tokens
- Protected routes and auth guards
- User preferences and liked recipes
- Advanced search and filtering UI
- Shopping list management

### Security & Performance Tasks
- JWT token authentication
- Password hashing with Argon2
- Input validation and sanitization
- CORS configuration
- Rate limiting implementation
- SQL injection prevention
- XSS protection
- Performance optimization and caching

## Docker Deployment

### Docker Configuration

#### Backend Dockerfile
```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["./main"]
```

#### Frontend Dockerfile
```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000

CMD ["node", "server.js"]
```

#### Docker Compose
```yaml
version: '3.8'

services:
  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8000:8000"
    environment:
      - JWT_SECRET=${JWT_SECRET}
    volumes:
      - ./data:/root/data
    networks:
      - recipe-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://backend:8000
    depends_on:
      - backend
    networks:
      - recipe-network

networks:
  recipe-network:
    driver: bridge

volumes:
  recipe-data:
```

### Local Docker Development
```bash
# Build and run with Docker Compose
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## AWS Deployment Strategy

### Infrastructure Architecture

#### AWS Services Used
- **ECS Fargate**: Container orchestration for scalable deployment
- **Application Load Balancer**: Traffic distribution and SSL termination
- **RDS PostgreSQL**: Production database (migrated from SQLite)
- **ECR**: Container image registry
- **VPC**: Network isolation and security
- **Route 53**: DNS management
- **Certificate Manager**: SSL certificates
- **CloudWatch**: Monitoring and logging

#### Deployment Architecture Diagram
```
Internet → Route 53 → ALB → ECS Fargate Cluster
                      │
                      ├── Frontend Service (3000)
                      └── Backend Service (8000)
                              │
                              └── RDS PostgreSQL
```

### Terraform Configuration

#### Directory Structure
```
terraform/
├── main.tf                 # Main infrastructure
├── variables.tf            # Input variables
├── outputs.tf              # Output values
├── vpc.tf                  # VPC and networking
├── ecs.tf                  # ECS cluster and services
├── rds.tf                  # Database configuration
├── alb.tf                  # Load balancer
├── ecr.tf                  # Container registry
├── iam.tf                  # IAM roles and policies
├── security-groups.tf      # Security group rules
└── terraform.tfvars       # Variable values (not committed)
```

#### Key Infrastructure Components

**VPC and Networking**
- Multi-AZ VPC with public and private subnets
- Internet Gateway and NAT Gateways
- Route tables and security groups

**ECS Fargate Cluster**
- Auto-scaling container services
- Service discovery and load balancing
- Blue-green deployment capability

**RDS PostgreSQL**
- Multi-AZ deployment for high availability
- Automated backups and monitoring
- Security group isolation

**Application Load Balancer**
- SSL termination with ACM certificates
- Path-based routing to frontend/backend
- Health checks and auto-scaling triggers

### Database Migration

#### SQLite to PostgreSQL Migration
```sql
-- PostgreSQL schema creation
CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    calories_per_100g INTEGER NOT NULL,
    description TEXT
);

CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    prep_time_minutes INTEGER NOT NULL,
    cook_time_minutes INTEGER NOT NULL,
    servings INTEGER NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    instructions TEXT NOT NULL,
    description TEXT
);

CREATE TABLE recipe_ingredients (
    recipe_id INTEGER REFERENCES recipes(id),
    ingredient_id INTEGER REFERENCES ingredients(id),
    quantity DECIMAL NOT NULL,
    unit VARCHAR(50) NOT NULL,
    notes TEXT,
    PRIMARY KEY (recipe_id, ingredient_id)
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_liked_recipes (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, recipe_id)
);
```

### Deployment Pipeline

#### CI/CD Workflow
1. **Code Push**: Developer pushes to main branch
2. **Build**: GitHub Actions builds Docker images
3. **Test**: Run automated tests and security scans
4. **Push**: Upload images to ECR
5. **Deploy**: Terraform applies infrastructure changes
6. **Update**: ECS updates services with new images
7. **Verify**: Health checks confirm successful deployment

#### Environment Configuration
```bash
# Production environment variables
JWT_SECRET=${RANDOM_SECRET_FROM_AWS_SECRETS_MANAGER}
DATABASE_URL=${RDS_CONNECTION_STRING}
REDIS_URL=${ELASTICACHE_ENDPOINT}
ALLOWED_ORIGINS=${FRONTEND_DOMAIN}
```

### Monitoring and Operations

#### CloudWatch Metrics
- Container CPU and memory utilization
- Database connection counts and query performance
- Application response times and error rates
- Custom business metrics (user registrations, recipe views)

#### Alerting
- High error rate alerts
- Database performance degradation
- Container health check failures
- SSL certificate expiration warnings

### Cost Optimization

#### Resource Sizing
- **ECS Tasks**: Start with 0.25 vCPU, 512 MB memory
- **RDS Instance**: db.t3.micro for development, db.t3.small for production
- **Auto Scaling**: Scale based on CPU utilization and request count

#### Cost Monitoring
- AWS Cost Explorer for spend analysis
- Resource tagging for cost allocation
- Scheduled scaling for predictable traffic patterns

## Troubleshooting

### Common Issues

**1. Database file not found**
```bash
# Solution: Start server in backend directory
cd backend
go run main.go
```

**2. CORS Error**
```bash
# Issue: Cannot call API from frontend  
# Solution: Check CORS settings in main.go
# Ensure frontend origin is whitelisted
```

**3. JWT Token Issues**
```bash
# Check JWT_SECRET environment variable
# Ensure token is properly stored in localStorage
# Verify token format in Authorization header
```

**4. Docker Build Issues**
```bash
# Clear Docker cache
docker system prune

# Rebuild without cache
docker-compose build --no-cache

# Check container logs
docker-compose logs backend
```

**5. AWS Deployment Issues**
```bash
# Check ECS service status
aws ecs describe-services --cluster recipe-cluster --services recipe-backend

# View CloudWatch logs
aws logs tail /ecs/recipe-backend --follow

# Verify security group rules
aws ec2 describe-security-groups --group-ids sg-xxxxx
```

## Production Environment

### Local Development with Docker
```bash
# Start all services
docker-compose up --build

# Development with hot reload (requires volume mounts)
docker-compose -f docker-compose.dev.yml up
```

### AWS Production Deployment
```bash
# Initialize Terraform
cd terraform
terraform init

# Plan deployment
terraform plan

# Deploy infrastructure
terraform apply

# Deploy application
./scripts/deploy.sh production
```

## Future Enhancements

### Planned Features
1. **User Preferences**: Save liked recipes and dietary restrictions
2. **Mobile App**: React Native mobile application
3. **Recipe Ratings**: User rating and review system
4. **Social Features**: Share recipes and follow other users
5. **Meal Planning**: Weekly meal planning with shopping lists
6. **Nutrition Tracking**: Detailed nutritional analysis

### Technical Improvements
1. **Caching**: Implement Redis for session management and caching
2. **Search**: Elasticsearch for advanced recipe search
3. **Testing**: Comprehensive unit tests and integration tests
4. **CI/CD**: Automated testing and deployment pipeline
5. **Monitoring**: Application monitoring and logging with DataDog
6. **API Documentation**: Interactive API documentation with Swagger
7. **Performance**: Database query optimization and CDN integration