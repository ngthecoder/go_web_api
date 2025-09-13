# Go Recipe Web API + Next.js
A web application using Go backend API and Next.js frontend

## Project Purpose
- Learn Web API development using Go language
- Build a cross-reference system for ingredients and recipes
- Gain team development experience

## Tech Stack
- **Backend**: Go 1.22.4 + SQLite3
- **Frontend**: Next.js 15.5.0 + TypeScript + Tailwind CSS
- **Database**: SQLite3 (for local development)

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
```

### API Endpoint Design

#### Endpoint List
| HTTP Method | Endpoint | Description | Required Parameters | Optional Parameters |
|-------------|----------|-------------|-------------------|-------------------|
| GET | `/ingredients` | Get ingredient list, search & filter | - | `search`, `category`, `sort`, `order`, `page`, `limit` |
| GET | `/ingredients/{id}` | Get ingredient details + related recipes | `id` | - |
| GET | `/recipes` | Get recipe list, search & filter | - | `search`, `category`, `max_time`, `difficulty`, `sort`, `order`, `page`, `limit` |
| GET | `/recipes/find-by-ingredients` | Find recipes by available ingredients | `ingredients` | `match_type`, `page`, `limit` |
| GET | `/recipes/{id}` | Get recipe details + required ingredients | `id` | - |
| GET | `/recipes/shopping-list/{id}` | Generate shopping list for recipe | `id` | `have_ingredients` |
| GET | `/categories` | Get category statistics | - | - |
| GET | `/stats` | Get overall statistics | - | - |

#### Detailed Specifications

**1: GET /api/recipes/find-by-ingredients**
```bash
GET /api/recipes/find-by-ingredients?ingredients=2,26&match_type=partial&limit=1
```
```json
{
  "matched_recipes": [
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
}
```

**2: GET /api/recipes/shopping-list/{id}**
```bash
# Example: Shopping list for summer vegetable curry (id=5), already have tomato and eggplant
GET /api/recipes/shopping-list/5?have_ingredients=1,4
```
```json
{
  "recipe": {
    "id": 5,
    "name": "Summer Vegetable Curry",
    "servings": 4
  },
  "need_to_buy": [
    {
      "id": 6,
      "name": "Potato",
      "quantity": 3,
      "unit": "pieces",
      "estimated_price": "¥150",
      "category": "vegetables"
    },
    {
      "id": 8,
      "name": "Curry Roux",
      "quantity": 1,
      "unit": "box",
      "estimated_price": "¥200",
      "category": "seasonings"
    }
  ],
  "already_have": [
    {"id": 1, "name": "Tomato", "quantity": 2, "unit": "pieces"},
    {"id": 4, "name": "Eggplant", "quantity": 1, "unit": "piece"}
  ],
  "total_estimated_cost": "¥350"
}
```

**3: GET /api/categories**
```json
{
  "ingredient_categories": {
    "vegetables": 15,
    "protein": 8,
    "grains": 5,
    "dairy": 4
  },
  "recipe_categories": {
    "breakfast": 12,
    "lunch": 18,
    "dinner": 25,
    "snacks": 8
  }
}
```

**4: GET /api/stats**
```json
{
  "total_ingredients": 32,
  "total_recipes": 63,
  "avg_prep_time": 15.5,
  "avg_cook_time": 22.3,
  "difficulty_distribution": {
    "easy": 45,
    "medium": 15,
    "hard": 3
  }
}
```

**5: GET /api/ingredients**
```bash
GET /api/ingredients?search=tomato&category=vegetables&sort=calories&order=desc&page=1&limit=10
```
```json
{
  "has_next": false,
  "ingredients": [
    {
      "id": 1,
      "name": "Tomato",
      "category": "vegetables",
      "calories_per_100g": 18,
      "description": "Fresh red tomatoes"
    }
  ],
  "page": 1,
  "page_size": 10,
  "total": 1,
  "total_pages": 1
}
```

**6: GET /api/ingredients/{id}**
```json
{
  "ingredient": {
    "id": 1,
    "name": "Tomato",
    "category": "vegetables",
    "calories_per_100g": 18,
    "description": "Fresh red tomatoes"
  },
  "recipes": [
    {
      "id": 1,
      "name": "Tomato Rice",
      "category": "dinner",
      "prep_time_minutes": 10,
      "cook_time_minutes": 25,
      "servings": 4,
      "difficulty": "easy",
      "instructions": "1. Heat oil in a frying pan...",
      "description": "Simple and delicious tomato rice"
    }
  ]
}
```

**7: GET /api/recipes**
```bash
GET /api/recipes?search=curry&category=dinner&difficulty=medium&max_time=60&sort=total_time&order=asc&page=1&limit=5
```
```json
{
  "has_next": false,
  "page": 1,
  "page_size": 5,
  "recipes": [
    {
      "id": 13,
      "name": "Summer Vegetable Curry",
      "category": "dinner",
      "prep_time_minutes": 20,
      "cook_time_minutes": 30,
      "servings": 4,
      "difficulty": "medium",
      "instructions": "1. Cut vegetables\n2. Stir-fry vegetables in pot\n3. Add water and simmer\n4. Add curry roux and dissolve\n5. Simmer more and complete",
      "description": "Healthy curry with plenty of summer vegetables"
    },
    {
      "id": 14,
      "name": "Chicken Curry",
      "category": "dinner",
      "prep_time_minutes": 15,
      "cook_time_minutes": 45,
      "servings": 4,
      "difficulty": "medium",
      "instructions": "1. Cut chicken\n2. Stir-fry onions\n3. Add chicken and stir-fry\n4. Add water and simmer\n5. Add curry roux",
      "description": "Authentic chicken curry"
    }
  ],
  "total": 2,
  "total_pages": 1
}
```

**8: GET /api/recipes/{id}**
```json
{
  "recipe": {
    "id": 1,
    "name": "Tomato Rice",
    "category": "dinner", 
    "prep_time_minutes": 10,
    "cook_time_minutes": 25,
    "servings": 4,
    "difficulty": "easy",
    "instructions": "1. Heat oil in a frying pan...",
    "description": "Simple and delicious tomato rice"
  },
  "ingredients": [
    {
      "ingredient_id": 1,
      "name": "Tomato",
      "quantity": 2,
      "unit": "pieces",
      "notes": "diced"
    },
    {
      "ingredient_id": 2,
      "name": "Onion",
      "quantity": 1,
      "unit": "piece", 
      "notes": "minced"
    }
  ]
}
```

## Frontend Design (Next.js)

### Page Details

#### Main Page (/)
- Toggle tabs for ingredients and recipes
- Search and filtering functionality
- Pagination
- Card-style data display

#### Ingredient Detail Page (/ingredients/[id])
- Display ingredient details (calories, category, description)
- List of recipes using this ingredient
- Click recipe cards to navigate to recipe details

#### Recipe Detail Page (/recipes/[id])
- Recipe details (cooking time, servings, difficulty, instructions)
- List of required ingredients (with quantities and units)
- Click ingredients to navigate to ingredient details
- Shopping list generation button

### User Flow

#### Flow 1: Recipe Search by Available Ingredients
1. Select "Search by Available Ingredients" tab on main page
2. Check available ingredients from the list (e.g., tomato, onion, eggplant)
3. Click "Search Recipes" button
4. View matched recipes displayed (sorted by score)
5. Click on a recipe of interest
6. Check instructions and all ingredients on recipe detail page
7. Use "Create Shopping List" button to check missing ingredients

#### Flow 2: Navigate from Recipe to Ingredient Details
1. Display ingredient list on recipe detail page
2. Click on an ingredient of interest (e.g., "Eggplant")
3. Navigate to eggplant detail page
4. Click on other recipes if interested to navigate

### Project Structure
```
frontend/
├── app/
│   ├── page.tsx              # Main page
│   ├── ingredients/
│   │   └── [id]/page.tsx     # Ingredient detail page
│   ├── recipes/
│   │   └── [id]/page.tsx     # Recipe detail page
│   └── globals.css           # Tailwind CSS
├── components/
│   ├── IngredientCard.tsx    # Ingredient card component
│   ├── RecipeCard.tsx        # Recipe card component
│   ├── SearchBar.tsx         # Search bar component
│   └── Pagination.tsx        # Pagination component
└── lib/
    └── api.ts                # API communication functions
```

## Features

### User Stories
1. Browse ingredient list: Users can search and filter ingredients
2. View ingredient details: Users can click ingredients to see details and recipes using them
3. Browse recipe list: Users can search and filter recipes
4. View recipe details: Users can click recipes to see details and required ingredients
5. Search recipes by available ingredients: Users can find recipes based on refrigerator contents
6. Generate shopping list: Users can select recipes and create shopping lists for missing ingredients
7. Bidirectional navigation: Navigate in cycles between ingredients → recipes → ingredients
8. View statistics: Users can browse category statistics and overall statistics
9. Paginated display: Efficiently browse large datasets with pagination

### Use Cases

**Scenario 1: Want to cook with refrigerator ingredients**
1. Select available ingredients (tomato, eggplant, potato)
2. System suggests recipes (with match scores)
3. Select a recipe and check details

**Scenario 2: Create shopping list for desired dish**
1. Select desired recipe (e.g., Summer Vegetable Curry)
2. Specify available ingredients
3. System auto-generates shopping list (with estimated prices)

## Team Development Workflow

### Backend Tasks
- Create database schema
- Implement basic GET endpoints
- Implement recipe search by available ingredients
- Implement shopping list generation
- Implement search and filtering functionality
- Implement pagination
- Add data aggregation endpoints
- Implement matching score calculation logic
- Enhance error handling
- Performance optimization (add indexes, etc.)
- Implement dedicated API documentation server

### Frontend Tasks (Next.js + TypeScript)
- Next.js project setup
- Implement API communication functions (lib/api.ts)
- Implement available ingredients selection UI
- Recipe suggestion display functionality
- Shopping list display functionality
- Implement search and filter UI (SearchBar.tsx)
- Implement pagination UI (Pagination.tsx)
- Implement ingredient and recipe card components
- Responsive design (Tailwind CSS)
- Bidirectional navigation functionality
- Handle loading and error states

#### Documentation Tasks
- Create API specification (static HTML site)
- Create endpoint usage examples and sample code
- Create developer API guide
- Create response format documentation

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
```

**3. SQLite3 Driver Error**
```bash
# Solution: Build with CGO enabled
CGO_ENABLED=1 go run main.go
```

## Production Environment (Incomplete)
```bash
docker compose up --build
```