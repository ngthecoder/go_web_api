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

### User Stories
1. **Browse Recipes**: Users can search and filter recipes by category, difficulty, cooking time
2. **Find Recipes by Ingredients**: Users can select available ingredients and find matching recipes
3. **View Details**: Users can view detailed recipe instructions and ingredient requirements
4. **User Registration**: New users can create accounts securely
5. **User Login**: Existing users can log in to access their profiles
6. **Generate Shopping Lists**: Users can create shopping lists for recipes they want to make
7. **Bidirectional Navigation**: Navigate between ingredients and recipes seamlessly

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

**4. JWT Token Issues**
```bash
# Check JWT_SECRET environment variable
# Ensure token is properly stored in localStorage
# Verify token format in Authorization header
```

## Production Environment (Incomplete)
```bash
# Docker setup is configured but not fully tested
docker compose up --build
```

## Future Enhancements

### Planned Features
1. **User Preferences**: Save liked recipes and dietary restrictions
2. **Mobile App**: React Native mobile application
3. **Recipe Ratings**: User rating and review system

### Technical Improvements
1. **Database Migration**: Move to PostgreSQL for production
2. **Caching**: Implement Redis for session management and caching
3. **Testing**: Unit tests and integration tests
4. **CI/CD**: Automated testing and deployment pipeline
5. **Monitoring**: Application monitoring and logging
6. **API Documentation**: Interactive API documentation with Swagger
