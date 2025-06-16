package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hursty1/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)
func HashPassword(password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), len([]byte(password)))
	if err != nil {
		return "", fmt.Errorf("failed to hash password %s", err)
	}
	return string(hashed_password), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return err //true
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})
	ss, err := jwt.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(tokenSecret), nil
    })
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid Token")
	}
	if claims.Issuer != "chirpy" {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}
	userID, err := uuid.Parse(claims.Subject)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid UUID in token subject: %w", err)
    }

    return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("authorization header was missing")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}
	return parts[1], nil

}

func GetPolkaApiKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("authorization header was missing")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "apikey" {
		return "", fmt.Errorf("authorization header format must be ApiKey {token}")
	}
	return parts[1], nil

}

func CheckRefreshToken(refresh_token database.RefreshToken) (bool, error) {
	// fmt.Printf("Refresh token: %s\n", refresh_token)
	if refresh_token.ExpiresAt.Before(time.Now()) {
		//expired
		return false, fmt.Errorf("Token Expired")
	}
	null := sql.NullTime{}
	if refresh_token.RevokedAt != null {
		return false, fmt.Errorf("Token Revoked")
	}

	return true, nil
}

