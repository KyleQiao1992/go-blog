package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the plain password using bcrypt.
func HashPassword(plain string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares hashed password and plain password.
func CheckPassword(hashed string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

func GenerateToken(userId uint, userName string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not configured")
	}

	// Set token claims.
	claims := jwt.MapClaims{
		"sub": userId,                                // subject: user ID
		"usr": userName,                              // username
		"exp": time.Now().Add(24 * time.Hour).Unix(), // expiration time (24 hours)
		"iat": time.Now().Unix(),                     // issued at
	}

	//create token with HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// sign token with secret
	// the signature ensures the token cannot be modified without the secret key
	return token.SignedString([]byte(secret))
}

func ParseAndValidateToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET not configured")
	}

	// Parse token.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
