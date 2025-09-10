package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	db        *sql.DB
	jwtSecret string
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) registerUser(registerRequest RegisterRequest) (*AuthResponse, error) {
	exists, err := s.userExists(registerRequest.Email, registerRequest.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	uuid := uuid.New()
	encodedPasswordHash, err := s.hashPassword(registerRequest.Password)
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'));"
	_, err = s.db.Exec(query, uuid, registerRequest.Username, registerRequest.Email, encodedPasswordHash)
	if err != nil {
		return nil, err
	}

	user, err := s.getUserByEmail(registerRequest.Email)
	if err != nil {
		return nil, err
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s *AuthService) userExists(email, username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = ? OR username = ?`
	err := s.db.QueryRow(query, email, username).Scan(&count)
	return count > 0, err
}

func (s *AuthService) loginUser(loginRequest LoginRequest) (*AuthResponse, error) {
	verified, err := s.verifyPassword(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return nil, err
	}

	if !verified {
		return nil, ErrInvalidCredentials
	}

	user, err := s.getUserByEmail(loginRequest.Email)
	if err != nil {
		return nil, err
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  *user,
		Token: token,
	}, nil
}

func (s *AuthService) getUserByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE email = ?`

	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

// still not complete (simplied version)
func (s *AuthService) generateJWT(user *User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	token := base64.StdEncoding.EncodeToString(claimsJSON)

	return token, nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 16)

	passwordHash := append(salt, hash...)

	return base64.StdEncoding.EncodeToString(passwordHash), nil
}

func (s *AuthService) verifyPassword(email string, password string) (bool, error) {
	storedEncodedPasswordHash := ""
	query := "SELECT password_hash FROM users WHERE email = ?;"
	err := s.db.QueryRow(query, email).Scan(&storedEncodedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	storedDecodedPasswordHashBytes, err := base64.StdEncoding.DecodeString(storedEncodedPasswordHash)
	if err != nil {
		return false, err
	}

	salt := storedDecodedPasswordHashBytes[:16]
	correctHash := storedDecodedPasswordHashBytes[16:]

	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 16)

	return subtle.ConstantTimeCompare(correctHash, hash) == 1, nil
}
