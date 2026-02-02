package users

import (
	"database/sql"
	"testing"

	"github.com/ngthecoder/go_web_api/internal/auth"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	// Create recipes table
	_, err = db.Exec(`
		CREATE TABLE recipes (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			category TEXT,
			prep_time_minutes INTEGER,
			cook_time_minutes INTEGER,
			servings INTEGER,
			difficulty TEXT,
			instructions TEXT,
			description TEXT
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create recipes table: %v", err)
	}

	// Create user_liked_recipes table
	_, err = db.Exec(`
		CREATE TABLE user_liked_recipes (
			user_id TEXT NOT NULL,
			recipe_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, recipe_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create user_liked_recipes table: %v", err)
	}

	return db
}

func seedTestUser(t *testing.T, db *sql.DB, userID, username, email, password string) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	_, err = db.Exec(
		"INSERT INTO users (id, username, email, password_hash) VALUES (?, ?, ?, ?)",
		userID, username, email, hashedPassword,
	)
	if err != nil {
		t.Fatalf("Failed to seed test user: %v", err)
	}
}

func seedTestRecipe(t *testing.T, db *sql.DB, id int, name string) {
	_, err := db.Exec(
		"INSERT INTO recipes (id, name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		id, name, "Dinner", 10, 20, 4, "easy", "Test instructions", "Test description",
	)
	if err != nil {
		t.Fatalf("Failed to seed test recipe: %v", err)
	}
}

func TestGetUserProfile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	// Seed a test user
	seedTestUser(t, db, "user-123", "testuser", "test@example.com", "password123")

	// Test successful profile retrieval
	profile, err := service.getUserProfile("user-123")
	if err != nil {
		t.Fatalf("getUserProfile() failed: %v", err)
	}

	if profile.ID != "user-123" {
		t.Errorf("Expected ID 'user-123', got '%s'", profile.ID)
	}
	if profile.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", profile.Username)
	}
	if profile.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", profile.Email)
	}

	// Test non-existent user
	_, err = service.getUserProfile("non-existent")
	if err == nil {
		t.Error("getUserProfile() should fail for non-existent user")
	}
}

func TestUpdateUserProfile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	seedTestUser(t, db, "user-123", "testuser", "test@example.com", "password123")

	// Test successful update
	profile, err := service.updateUserProfile("user-123", "newusername", "new@example.com")
	if err != nil {
		t.Fatalf("updateUserProfile() failed: %v", err)
	}

	if profile.Username != "newusername" {
		t.Errorf("Expected username 'newusername', got '%s'", profile.Username)
	}
	if profile.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got '%s'", profile.Email)
	}

	// Test duplicate username/email conflict
	seedTestUser(t, db, "user-456", "otheruser", "other@example.com", "password123")

	_, err = service.updateUserProfile("user-123", "otheruser", "new@example.com")
	if err == nil {
		t.Error("updateUserProfile() should fail when username is taken")
	}

	_, err = service.updateUserProfile("user-123", "newusername", "other@example.com")
	if err == nil {
		t.Error("updateUserProfile() should fail when email is taken")
	}
}

func TestLikedRecipes(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	seedTestUser(t, db, "user-123", "testuser", "test@example.com", "password123")
	seedTestRecipe(t, db, 1, "Recipe One")
	seedTestRecipe(t, db, 2, "Recipe Two")

	// Test getting empty liked recipes
	recipes, err := service.getLikedRecipes("user-123")
	if err != nil {
		t.Fatalf("getLikedRecipes() failed: %v", err)
	}
	if len(recipes) != 0 {
		t.Errorf("Expected 0 liked recipes, got %d", len(recipes))
	}

	// Test adding a liked recipe
	err = service.addLikedRecipe("user-123", 1)
	if err != nil {
		t.Fatalf("addLikedRecipe() failed: %v", err)
	}

	recipes, err = service.getLikedRecipes("user-123")
	if err != nil {
		t.Fatalf("getLikedRecipes() failed: %v", err)
	}
	if len(recipes) != 1 {
		t.Errorf("Expected 1 liked recipe, got %d", len(recipes))
	}
	if recipes[0].Name != "Recipe One" {
		t.Errorf("Expected recipe name 'Recipe One', got '%s'", recipes[0].Name)
	}

	// Test adding duplicate liked recipe (should fail)
	err = service.addLikedRecipe("user-123", 1)
	if err == nil {
		t.Error("addLikedRecipe() should fail for duplicate")
	}

	// Test adding non-existent recipe (should fail)
	err = service.addLikedRecipe("user-123", 999)
	if err == nil {
		t.Error("addLikedRecipe() should fail for non-existent recipe")
	}

	// Test removing liked recipe
	err = service.removeLikedRecipe("user-123", 1)
	if err != nil {
		t.Fatalf("removeLikedRecipe() failed: %v", err)
	}

	recipes, err = service.getLikedRecipes("user-123")
	if err != nil {
		t.Fatalf("getLikedRecipes() failed: %v", err)
	}
	if len(recipes) != 0 {
		t.Errorf("Expected 0 liked recipes after removal, got %d", len(recipes))
	}

	// Test removing non-existent liked recipe (should fail)
	err = service.removeLikedRecipe("user-123", 1)
	if err == nil {
		t.Error("removeLikedRecipe() should fail for recipe not in list")
	}
}

func TestChangePassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	seedTestUser(t, db, "user-123", "testuser", "test@example.com", "oldpassword")

	// Test with wrong current password
	err := service.changePassword("user-123", "wrongpassword", "newpassword")
	if err == nil {
		t.Error("changePassword() should fail with wrong current password")
	}

	// Test successful password change
	err = service.changePassword("user-123", "oldpassword", "newpassword")
	if err != nil {
		t.Fatalf("changePassword() failed: %v", err)
	}

	// Verify new password works (by trying to change again)
	err = service.changePassword("user-123", "newpassword", "anotherpassword")
	if err != nil {
		t.Error("changePassword() should work with new password")
	}

	// Verify old password no longer works
	err = service.changePassword("user-123", "oldpassword", "something")
	if err == nil {
		t.Error("changePassword() should fail with old password after change")
	}

	// Test non-existent user
	err = service.changePassword("non-existent", "password", "newpassword")
	if err == nil {
		t.Error("changePassword() should fail for non-existent user")
	}
}

func TestDeleteAccount(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewUserService(db)

	seedTestUser(t, db, "user-123", "testuser", "test@example.com", "password123")

	// Test with wrong password
	err := service.deleteAccount("user-123", "wrongpassword")
	if err == nil {
		t.Error("deleteAccount() should fail with wrong password")
	}

	// Verify user still exists
	_, err = service.getUserProfile("user-123")
	if err != nil {
		t.Error("User should still exist after failed deletion")
	}

	// Test successful deletion
	err = service.deleteAccount("user-123", "password123")
	if err != nil {
		t.Fatalf("deleteAccount() failed: %v", err)
	}

	// Verify user no longer exists
	_, err = service.getUserProfile("user-123")
	if err == nil {
		t.Error("User should not exist after deletion")
	}

	// Test deleting non-existent user
	err = service.deleteAccount("non-existent", "password")
	if err == nil {
		t.Error("deleteAccount() should fail for non-existent user")
	}
}
