package auth

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

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

	return db
}

func TestPasswordHashing(t *testing.T) {
	password := "SecurePass123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword() failed: %v", err)
	}

	if hash == "" {
		t.Error("hashPassword() returned empty hash")
	}

	if hash == password {
		t.Error("hashPassword() returned plaintext password!")
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword() failed on second call: %v", err)
	}

	if hash == hash2 {
		t.Error("Two hashes of same password should differ (random salt)")
	}
}

func TestRegistrationAndLogin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewAuthService(db, "test-secret")

	registerReq := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	response, err := service.registerUser(registerReq)
	if err != nil {
		t.Fatalf("registerUser() failed: %v", err)
	}

	if response == nil {
		t.Fatal("registerUser() returned nil response")
	}

	if response.Token == "" {
		t.Error("registerUser() returned empty token")
	}

	if response.User.Username != "testuser" {
		t.Errorf("Username = %v, want testuser", response.User.Username)
	}

	loginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResp, err := service.loginUser(loginReq)
	if err != nil {
		t.Fatalf("loginUser() with correct password failed: %v", err)
	}

	if loginResp.Token == "" {
		t.Error("loginUser() returned empty token")
	}

	wrongLoginReq := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err = service.loginUser(wrongLoginReq)
	if err == nil {
		t.Error("loginUser() should fail with wrong password")
	}

	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestJWTGeneration(t *testing.T) {
	service := &AuthService{jwtSecret: "test-secret-key"}

	testUser := &User{
		ID:       "test-user-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	token, err := service.generateJWT(testUser)
	if err != nil {
		t.Fatalf("generateJWT() failed: %v", err)
	}

	if token == "" {
		t.Error("generateJWT() returned empty token")
	}

	claims, err := service.validateJWT("Bearer " + token)
	if err != nil {
		t.Fatalf("validateJWT() failed: %v", err)
	}

	if claims == nil {
		t.Fatal("validateJWT() returned nil claims")
	}

	if claims.UserID != testUser.ID {
		t.Errorf("UserID = %v, want %v", claims.UserID, testUser.ID)
	}

	_, err = service.validateJWT("Bearer invalid.token.here")
	if err == nil {
		t.Error("validateJWT() should fail for invalid token")
	}

	tamperedToken := token[:len(token)-5] + "XXXXX"
	_, err = service.validateJWT("Bearer " + tamperedToken)
	if err == nil {
		t.Error("validateJWT() should fail for tampered token")
	}

	claims, err = service.validateJWT(token)
	if err != nil {
		t.Error("validateJWT() should work without Bearer prefix")
	}
}

func TestDuplicateUserRegistration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewAuthService(db, "test-secret")

	registerReq := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err := service.registerUser(registerReq)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	duplicateReq := RegisterRequest{
		Username: "differentuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	_, err = service.registerUser(duplicateReq)
	if err != ErrUserExists {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}

	duplicateReq2 := RegisterRequest{
		Username: "testuser",
		Email:    "different@example.com",
		Password: "password123",
	}

	_, err = service.registerUser(duplicateReq2)
	if err != ErrUserExists {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}
}
