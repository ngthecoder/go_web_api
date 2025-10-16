package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
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

	encodedPasswordHash, err := HashPassword(registerRequest.Password)
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);"
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
	query := `SELECT COUNT(*) FROM users WHERE email = $1 OR username = $2`
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

func (s *AuthService) verifyPassword(email string, password string) (bool, error) {
	storedEncodedPasswordHash := ""
	query := "SELECT password_hash FROM users WHERE email = $1;"
	err := s.db.QueryRow(query, email).Scan(&storedEncodedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return VerifyPasswordHash(password, storedEncodedPasswordHash)
}

func (s *AuthService) getUserByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE email = $1`

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

func (s *AuthService) generateJWT(user *User) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	encodedHeader := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(headerJSON)

	claims := JWTClaims{
		UserID: user.ID,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encodedClaims := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(claimsJSON)

	message := encodedHeader + "." + encodedClaims

	signature, err := s.createHMACSignature(message)
	if err != nil {
		return "", err
	}

	token := message + "." + signature

	return token, nil
}

func (s *AuthService) createHMACSignature(message string) (string, error) {
	h := hmac.New(sha256.New, []byte(s.jwtSecret))

	_, err := h.Write([]byte(message))
	if err != nil {
		return "", err
	}

	signature := h.Sum(nil)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(signature), nil
}

func (s *AuthService) validateJWT(tokenString string) (*JWTClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	tokenStringParts := strings.Split(tokenString, ".")
	if len(tokenStringParts) != 3 {
		return nil, errors.New("Invalid JWT format")
	}

	headerString := tokenStringParts[0]
	claimsString := tokenStringParts[1]
	signatureString := tokenStringParts[2]

	message := headerString + "." + claimsString

	expectedTokenString, err := s.createHMACSignature(message)

	if signatureString != expectedTokenString {
		return nil, errors.New("Invalid token signature")
	}

	decodedClaims, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(claimsString)
	if err != nil {
		return nil, errors.New("Invalid claims encoding")
	}

	var claims JWTClaims
	err = json.Unmarshal(decodedClaims, &claims)
	if err != nil {
		return nil, errors.New("Invalid claims format")
	}

	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("Token expired")
	}

	return &claims, nil
}
