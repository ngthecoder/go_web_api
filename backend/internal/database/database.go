package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	URL      string
}

func InitDB(config Config) error {
	var err error
	var dbURL string

	if config.URL != "" {
		dbURL = config.URL
		if !strings.Contains(dbURL, "sslmode=") {
			if strings.Contains(dbURL, "?") {
				dbURL += "&sslmode=require"
			} else {
				dbURL += "?sslmode=require"
			}
		}
	} else {
		if config.Host == "" {
			config.Host = "postgres"
		}
		if config.Port == "" {
			config.Port = "5432"
		}

		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.User, config.Password, config.Host, config.Port, config.DBName)
	}

	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if err := createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS ingredients (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			calories_per_100g INTEGER NOT NULL,
			description TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS recipes (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			prep_time_minutes INTEGER NOT NULL,
			cook_time_minutes INTEGER NOT NULL,
			servings INTEGER NOT NULL,
			difficulty TEXT NOT NULL,
			instructions TEXT NOT NULL,
			description TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS recipe_ingredients (
			recipe_id INTEGER NOT NULL,
			ingredient_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit TEXT NOT NULL,
			notes TEXT,
			PRIMARY KEY (recipe_id, ingredient_id),
			FOREIGN KEY (recipe_id) REFERENCES recipes (id),
			FOREIGN KEY (ingredient_id) REFERENCES ingredients (id)
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS user_liked_recipes (
			user_id TEXT NOT NULL,
			recipe_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, recipe_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		)`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("All tables created successfully")
	return nil
}

func createIndexes() error {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_ingredients_category ON ingredients(category)",
		"CREATE INDEX IF NOT EXISTS idx_ingredients_name ON ingredients(name)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_category ON recipes(category)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_difficulty ON recipes(difficulty)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_category_difficulty ON recipes(category, difficulty)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_ingredient_id ON recipe_ingredients(ingredient_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
	}

	for _, index := range indexes {
		if _, err := DB.Exec(index); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	log.Println("All indexes created successfully")
	return nil
}
