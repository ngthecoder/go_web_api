# Go Recipe Web API + Next.js
A production-ready full-stack web application using Go backend API and Next.js frontend with user authentication, built as a portfolio project to demonstrate modern web development skills.

## Quick Demo

**Live Demo**: Coming soon

## Key Features

- 🔐 **Secure Authentication**: JWT tokens with Argon2 password hashing
- 🍳 **Recipe Discovery**: Find recipes based on ingredients you have
- ❤️ **Personal Collection**: Save and manage your favorite recipes
- 🛒 **Smart Shopping Lists**: Generate shopping lists from recipes
- 📱 **Responsive Design**: Works seamlessly on mobile and desktop
- ✅ **Tested**: 40% test coverage focusing on critical security features
- 🐳 **Containerized**: Docker setup for consistent deployment

## Project Purpose
- Learn Web API development using Go language
- Build a cross-reference system for ingredients and recipes
- Gain team development experience
- Implement secure user authentication and authorization

## Tech Stack
- **Backend**: Go 1.24.0 + SQLite3
- **Frontend**: Next.js 15.5.0 + React 19.1.0 + TypeScript + Tailwind CSS 4
- **Database**: SQLite3 (for local development)
- **Authentication**: JWT tokens with Argon2 password hashing
- **Deployment**: Docker + AWS Fargate with Terraform

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

2. Set up environment variables
```bash
cd backend
# Create .env file with your JWT secret
echo "JWT_SECRET=your-secret-key-here" > .env
```

3. Start the backend
```bash
cd backend
go run main.go
```
- Runs on http://localhost:8000

4. Start the frontend
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
| GET | `/recipes` | Get recipe list, search & filter | - | `search`, `category`, `max_time`, `difficulty`, `sort`, `order`, `page`, `limit` | Optional |
| GET | `/recipes/find-by-ingredients` | Find recipes by available ingredients | `ingredients` | `match_type`, `page`, `limit` | Optional |
| GET | `/recipes/{id}` | Get recipe details + required ingredients | `id` | - | Optional |
| GET | `/recipes/shopping-list/{id}` | Generate shopping list for recipe | `id` | `have_ingredients` | No |
| GET | `/categories` | Get category statistics | - | - | No |
| GET | `/stats` | Get overall statistics | - | - | No |
| GET | `/user/profile` | Get user profile | - | - | Yes |
| GET | `/user/liked-recipes` | Get user's liked recipes | - | - | Yes |
| POST | `/user/liked-recipes/add` | Add recipe to liked list | recipe_id | - | Yes |
| DELETE | `/user/liked-recipes/{id}` | Remove recipe from liked list | `id` | - | Yes |

**Note**: Endpoints marked "Optional" for auth will return `is_liked: false` for unauthenticated users and `is_liked: true/false` based on user preference when authenticated.

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
Authorization: Bearer {token} (optional)
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
    "match_score": 0.33333334,
    "is_liked": false
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
[
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
]
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
- **Home (/)**: Recipe browsing with search, filtering, pagination, and like functionality
- **Ingredients (/ingredients)**: Ingredient browsing with search and filtering
- **Find Recipes (/find-recipes)**: Recipe search by available ingredients with like functionality
- **Register (/register)**: User registration page with validation
- **Login (/login)**: User login page with password visibility toggle
- **Profile (/profile)**: User profile management (protected route)
- **Liked Recipes (/liked-recipes)**: View and manage liked recipes (protected route)

#### Ingredient Detail Page (/ingredients/[id])
- Display ingredient details (calories, category, description)
- List of recipes using this ingredient
- Click recipe cards to navigate to recipe details

#### Recipe Detail Page (/recipes/[id])
- Recipe details (cooking time, servings, difficulty, instructions)
- List of required ingredients (with quantities and units)
- Click ingredients to navigate to ingredient details
- Shopping list generation functionality
- Like button for authenticated users

### User Authentication Flow

#### Registration Flow
1. User accesses `/register` page
2. Fills out registration form with client-side validation:
   - Username: minimum 3 characters, alphanumeric with underscores
   - Email: valid email format
   - Password: minimum 6 characters with confirmation
3. Frontend sends POST request to `/api/auth/register`
4. Backend validates input and creates user account with Argon2 hashed password
5. JWT token returned and stored in localStorage
6. User redirected to home page or previous location
7. AuthContext updates application state

#### Login Flow
1. User accesses `/login` page
2. Enters email and password with optional password visibility toggle
3. Frontend sends POST request to `/api/auth/login`
4. Backend validates credentials using Argon2 verification
5. JWT token returned and stored in localStorage
6. User redirected to home page or previous location
7. AuthContext updates application state

### User Experience Flows

#### Flow 1: Recipe Search by Available Ingredients
1. Navigate to "Find Recipes" page (/find-recipes)
2. Search and select available ingredients from checkbox list
3. Click "Search Recipes" button
4. View matched recipes with match scores and like status
5. Click heart icon to like/unlike recipes (requires authentication)
6. Click on recipe to view full details

#### Flow 2: User Registration and Profile
1. Click "Register" in navigation
2. Fill out registration form with validation feedback
3. Automatically logged in after successful registration
4. Access profile page to view account information
5. View account statistics (days as member, liked recipes count)
6. Logout functionality available in navigation dropdown

#### Flow 3: User Profile and Liked Recipes
1. User logs into their account
2. Navigates to profile page to view account information
3. Browses through recipe catalog (home page or find recipes)
4. Clicks heart icon on recipes to add to liked list
5. Views liked recipes collection at /liked-recipes
6. Removes recipes by clicking heart icon again
7. Real-time UI updates reflect like status changes

#### Flow 4: Protected Routes and Authentication
1. Unauthenticated user attempts to access protected route
2. System stores intended destination in localStorage
3. User redirected to login page with contextual message
4. After successful login, user redirected to original destination
5. ProtectedRoute component ensures only authenticated access

### Project Structure
```
frontend/
├── app/
│   ├── page.tsx                 # Home page (recipe browsing with likes)
│   ├── ingredients/
│   │   ├── page.tsx             # Ingredients list
│   │   └── [id]/page.tsx        # Ingredient detail
│   ├── recipes/
│   │   └── [id]/page.tsx        # Recipe detail with like button
│   ├── find-recipes/
│   │   └── page.tsx             # Recipe search by ingredients with likes
│   ├── register/
│   │   └── page.tsx             # User registration with validation
│   ├── login/
│   │   └── page.tsx             # User login
│   ├── profile/
│   │   ├── page.tsx             # User profile (protected)
│   │   ├── edit/
│   │   │   └── page.tsx         # Edit profile (placeholder)
│   │   └── change-password/
│   │       └── page.tsx         # Change password (placeholder)
│   ├── liked-recipes/
│   │   └── page.tsx             # Liked recipes (protected)
│   ├── layout.tsx               # Root layout with navigation & AuthProvider
│   └── globals.css              # Global styles (Tailwind CSS 4)
├── components/
│   ├── Navigation.tsx           # Navigation with auth dropdown
│   ├── ProtectedRoute.tsx       # Route protection HOC
│   └── LikeButton.tsx           # Reusable like button component
├── contexts/
│   └── AuthContext.tsx          # Authentication state management
├── lib/
│   ├── auth.ts                  # Authentication API functions
│   └── types.ts                 # TypeScript interfaces
└── package.json                 # Dependencies (React 19.1.0, Next.js 15.5.0)
```

## Security Features

### Password Security
- **Argon2 Hashing**: Passwords are hashed using Argon2id algorithm
- **Salt Generation**: Random 16-byte salt generated for each password
- **Memory-hard Function**: Argon2 parameters: time=3, memory=64MB, threads=2
- **Constant-time Comparison**: Uses subtle.ConstantTimeCompare for password verification

### JWT Token Security
- **HMAC-SHA256 Signing**: Tokens signed with secret key from environment variables
- **24-Hour Expiration**: Tokens expire after 24 hours for security
- **Stateless Authentication**: No server-side session storage required
- **Bearer Token Format**: Tokens prefixed with "Bearer" in Authorization header
- **Client-side Storage**: Tokens stored in localStorage with automatic cleanup on logout

### API Security
- **CORS Protection**: Configured for frontend origin (http://localhost:3000)
- **Input Validation**: All endpoints validate input parameters
- **Error Handling**: Secure error messages without sensitive information using custom HTTPError types
- **Optional Authentication**: Some endpoints support both authenticated and anonymous access
- **Protected Routes**: Client-side ProtectedRoute component guards sensitive pages
- **Automatic Redirects**: Stores intended destination before login redirect

### Frontend Security
- **Form Validation**: Client-side validation before API calls
- **Password Visibility Toggle**: User-controlled password display
- **Token Management**: Automatic token cleanup on logout
- **Protected Routes**: ProtectedRoute HOC prevents unauthorized access
- **Error Boundaries**: Graceful error handling in AuthContext

## 🎯 Current Project Status

### ✅ Completed Features

#### Backend (Go)
- ✅ **Authentication System**
  - User registration with validation
  - Login with JWT tokens
  - Password hashing with Argon2
  - Protected route middleware
  - Optional authentication support
  
- ✅ **Recipe Management**
  - Browse recipes with search, filtering, sorting
  - Recipe details with ingredients
  - Find recipes by available ingredients
  - Match scoring algorithm
  - Shopping list generation
  - Pagination support
  
- ✅ **Ingredient Management**
  - Browse ingredients with filters
  - Ingredient details with related recipes
  - Category-based organization
  
- ✅ **User Features**
  - User profiles
  - Like/unlike recipes
  - View liked recipes collection
  - User statistics
  
- ✅ **Database**
  - SQLite3 for development
  - Normalized schema (5 tables)
  - Proper foreign key relationships
  - Database indexes for performance
  
- ✅ **Testing**
  - 40% test coverage
  - Authentication security tests
  - Recipe business logic tests
  - In-memory test databases
  - Table-driven test patterns

#### Frontend (Next.js)
- ✅ **Pages**
  - Home page with recipe browsing
  - Recipe detail page
  - Ingredient listing and detail pages
  - Find recipes by ingredients
  - User registration and login
  - User profile page
  - Liked recipes page
  
- ✅ **Components**
  - Navigation with auth dropdown
  - LikeButton component
  - ProtectedRoute wrapper
  - Responsive design (mobile-friendly)
  
- ✅ **State Management**
  - AuthContext for global auth state
  - Local storage for token persistence
  - Real-time UI updates for likes
  
- ✅ **Features**
  - Client-side form validation
  - Password visibility toggle
  - Automatic redirect after login
  - Anonymous browsing support
  - Protected routes with auth guards

### 🚧 In Progress / Next Steps

#### Week 1-2: Production Ready (Current Phase)
- [ ] **Docker Setup**
  - Dockerfile for backend
  - Dockerfile for frontend
  - docker-compose.yml for local development
  - Multi-stage builds for optimization
  
- [ ] **PostgreSQL Migration**
  - Migrate from SQLite to PostgreSQL
  - Update connection strings
  - Test all queries work with PostgreSQL
  - Update SQL syntax differences

#### Week 3-5: AWS Production Deployment
- [ ] **Infrastructure as Code**
  - Terraform configuration
  - VPC and networking setup
  - RDS PostgreSQL configuration
  - ECR for container registry
  - ECS Fargate for containers
  - CloudFront for CDN
  
- [ ] **CI/CD Pipeline**
  - GitHub Actions workflow
  - Automated testing
  - Automated deployment
  - Environment management

#### Week 6: Polish & Documentation
- [ ] Architecture diagrams
- [ ] API documentation
- [ ] Learning story document
- [ ] Screenshots and demo video
- [ ] Update resume with achievements

### ⏳ Future Enhancements (Post-Job Hunt)

- [ ] Email verification for new accounts
- [ ] Password reset functionality  
- [ ] Recipe rating and review system
- [ ] User-generated recipe content
- [ ] Recipe collections/folders
- [ ] Meal planning features
- [ ] Shopping list management
- [ ] Social features (follow users, share recipes)
- [ ] Mobile app (React Native)
- [ ] Advanced search with filters
- [ ] Recipe recommendations (ML)
- [ ] Nutrition tracking
- [ ] Image upload for recipes
- [ ] Comments on recipes
- [ ] Print-friendly recipe view

## Getting Started

### Backend Tasks (Completed)
- ✅ Database schema design and implementation with indexes
- ✅ User authentication system (registration, login, JWT)
- ✅ Password security with Argon2 hashing
- ✅ Basic recipe and ingredient CRUD endpoints
- ✅ Recipe search by ingredients functionality
- ✅ Shopping list generation
- ✅ Search, filtering, and pagination
- ✅ Statistics and category endpoints
- ✅ CORS and security middleware
- ✅ User-specific features (liked recipes, user preferences)
- ✅ Custom error handling with HTTP status codes
- ✅ Clean service layer architecture
- ✅ Input validation and sanitization
- ✅ Optional authentication support for recipe endpoints

### Frontend Tasks (Completed)
- ✅ Next.js 15.5 project setup with TypeScript and React 19.1
- ✅ User authentication UI (register, login, profile) with validation
- ✅ Recipe and ingredient browsing pages
- ✅ Recipe search by ingredients interface
- ✅ Responsive design with Tailwind CSS 4
- ✅ Navigation component with auth dropdown
- ✅ Local storage for authentication tokens
- ✅ Protected routes and auth guards (ProtectedRoute component)
- ✅ User preferences and liked recipes with real-time updates
- ✅ LikeButton component for recipe likes
- ✅ AuthContext for global authentication state management
- ✅ Client-side form validation with error feedback
- ✅ Password visibility toggle
- ✅ Automatic redirect after login
- ✅ Anonymous browsing with optional authentication

### Pending Tasks
- [ ] Edit profile functionality (PUT /api/user/profile endpoint)
- [ ] Change password functionality (PUT /api/user/password endpoint)
- [ ] Shopping list management UI
- [ ] Rate limiting implementation
- [ ] Advanced search and filtering UI improvements
- [ ] Email verification for new accounts
- [ ] Password reset functionality
- [ ] User avatar upload
- [ ] Recipe ratings and reviews
- [ ] Social sharing features

### Security & Performance Tasks
- ✅ JWT token authentication
- ✅ Password hashing with Argon2
- ✅ Input validation and sanitization
- ✅ CORS configuration
- ✅ SQL injection prevention
- ⏳ Rate limiting implementation
- ⏳ XSS protection
- ⏳ Performance optimization and caching
- ⏳ API response compression
- ⏳ Database query optimization

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

**Note**: The current Dockerfile has `CGO_ENABLED=0` but should be `CGO_ENABLED=1` for SQLite support.

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
# Ensure frontend origin is whitelisted (http://localhost:3000)
```

**3. JWT Token Issues**
```bash
# Check JWT_SECRET environment variable
# Ensure .env file exists in backend directory
cd backend
cat .env  # Should show JWT_SECRET=your-secret-key

# Verify token is properly stored in localStorage
# Check browser console for token format in Authorization header
```

**4. Authentication Issues**
```bash
# Issue: User cannot login after registration
# Solution: Verify password hashing is working
# Check that Argon2 dependencies are installed: golang.org/x/crypto

# Issue: Token expired errors
# Solution: Tokens expire after 24 hours, user needs to login again

# Issue: Protected routes not working
# Solution: Check ProtectedRoute component and AuthContext
# Verify token is present in localStorage
```

**5. Docker Build Issues**
```bash
# Clear Docker cache
docker system prune

# Rebuild without cache
docker-compose build --no-cache

# Check container logs
docker-compose logs backend
docker-compose logs frontend

# Issue: Backend fails to start in Docker
# Solution: Ensure CGO_ENABLED=1 in Dockerfile for SQLite support
# Update backend/Dockerfile line 9:
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .
```

**6. AWS Deployment Issues**
```bash
# Check ECS service status
aws ecs describe-services --cluster recipe-cluster --services recipe-backend

# View CloudWatch logs
aws logs tail /ecs/recipe-backend --follow

# Verify security group rules
aws ec2 describe-security-groups --group-ids sg-xxxxx
```

**7. Frontend Build Issues**
```bash
# Issue: npm install fails
# Solution: Clear cache and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm cache clean --force
npm install

# Issue: TypeScript errors
# Solution: Check tsconfig.json and run type checking
npm run build
```

**8. Like Button Not Working**
```bash
# Issue: Likes not persisting
# Solution: Check JWT token is being sent in Authorization header
# Verify backend endpoint: POST /api/user/liked-recipes/add
# Check browser console for error messages

# Issue: Like status not updating
# Solution: Verify AuthContext is providing token
# Check that onLikeChange callback is updating parent state
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

## Backend Architecture

### Project Structure
```
backend/
├── internal/
│   ├── auth/
│   │   ├── handlers.go          # Authentication HTTP handlers
│   │   ├── models.go            # User and auth models
│   │   └── services.go          # Auth business logic (JWT, Argon2)
│   ├── errors/
│   │   └── http_errors.go       # Custom HTTP error types
│   ├── ingredients/
│   │   ├── handlers.go          # Ingredient HTTP handlers
│   │   ├── models.go            # Ingredient models
│   │   └── services.go          # Ingredient business logic
│   ├── recipes/
│   │   ├── handlers.go          # Recipe HTTP handlers
│   │   ├── models.go            # Recipe models
│   │   └── services.go          # Recipe business logic
│   └── users/
│       ├── handlers.go          # User profile HTTP handlers
│       ├── models.go            # User profile models
│       └── services.go          # User profile business logic
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── Dockerfile                   # Docker build configuration
└── .env                         # Environment variables (not committed)
```

### Middleware Stack
1. **Logging Middleware**: Logs all HTTP requests with timing
2. **CORS Middleware**: Handles cross-origin requests
3. **Auth Middleware**: Validates JWT tokens for protected routes
4. **Optional Auth Middleware**: Adds user context when token present

### Design Patterns
- **Service Layer Pattern**: Business logic separated from HTTP handlers
- **Repository Pattern**: Database operations encapsulated in services
- **Middleware Pattern**: Cross-cutting concerns handled via middleware
- **Error Handling Pattern**: Custom HTTPError types for consistent error responses

## Future Enhancements

### Planned Features
1. **User Preferences**: Save dietary restrictions and favorite categories
2. **Mobile App**: React Native mobile application
3. **Recipe Ratings**: User rating and review system with aggregated scores
4. **Social Features**: Share recipes and follow other users
5. **Meal Planning**: Weekly meal planning with automated shopping lists
6. **Nutrition Tracking**: Detailed nutritional analysis and daily tracking
7. **Recipe Collections**: Create custom recipe collections and folders
8. **Advanced Search**: Full-text search with filters and sorting
9. **Recipe Recommendations**: AI-powered recipe suggestions based on preferences
10. **Grocery Store Integration**: Direct shopping list export to grocery services

### Technical Improvements
1. **Caching**: Implement Redis for session management and API response caching
2. **Search**: Elasticsearch for advanced recipe search and autocomplete
3. **Testing**: Comprehensive unit tests and integration tests
   - Backend: Go testing package with table-driven tests
   - Frontend: Jest and React Testing Library
4. **CI/CD**: Automated testing and deployment pipeline with GitHub Actions
5. **Monitoring**: Application monitoring and logging with DataDog or Prometheus
6. **API Documentation**: Interactive API documentation with Swagger/OpenAPI
7. **Performance**: 
   - Database query optimization with prepared statements
   - CDN integration for static assets
   - Image optimization and lazy loading
8. **Security Enhancements**:
   - Rate limiting per user/IP
   - Email verification for new accounts
   - Password reset via email
   - Two-factor authentication (2FA)
   - HTTPS enforcement in production
9. **Database Improvements**:
   - Database connection pooling
   - Read replicas for scaling
   - Database migrations with tools like golang-migrate
10. **Code Quality**:
    - Linting with golangci-lint
    - Pre-commit hooks
    - Code coverage reports
    - API versioning

## Testing

### Quick Start
```bash
cd backend

# Run all tests
go test ./... -v

# Run with coverage report
go test ./... -cover
```

### Test Coverage
- **Overall Coverage**: ~40% (focusing on critical paths)
- **Authentication**: 45% - Password hashing, JWT validation, user registration
- **Recipe Service**: 35% - Database queries, recipe retrieval
- **Strategy**: Prioritize testing security-critical code and core business logic

### Test Structure
```
backend/
└── internal/
   ├── auth/
   │   ├── services.go
   │   └── services_test.go       # 4 auth security tests
   └── recipes/
       ├── services.go
       └── services_test.go       # 3 recipe logic tests
```

### Running Specific Tests
```bash
# Run only auth tests
go test ./internal/auth -v

# Run only recipe tests  
go test ./internal/recipes -v

# Run tests with coverage report
go test ./internal/auth -cover
go test ./internal/recipes -cover

# Generate HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API Rate Limiting (Planned)

### Rate Limit Strategy
- **Anonymous Users**: 100 requests per hour
- **Authenticated Users**: 1000 requests per hour
- **Sliding Window**: 1-hour window with gradual reset

### Implementation
```go
// Rate limiting middleware (to be implemented)
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check rate limit based on IP or user ID
        // Return 429 Too Many Requests if exceeded
        next(w, r)
    }
}
```
