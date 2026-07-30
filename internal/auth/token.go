package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"project-manager-go/internal/domain"
)

type TokenManager struct {
	AccessSecret  []byte
	RefreshSecret []byte
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	Now           func() time.Time
}

type Claims struct {
	ID    int             `json:"id"`
	Email string          `json:"email"`
	Role  domain.UserRole `json:"role,omitempty"`
	Type  string          `json:"type"`
	jwt.RegisteredClaims
}

func (m TokenManager) GenerateAccessToken(user domain.User) (string, error) {
	return m.generate(user, "access", m.AccessTTL, m.AccessSecret)
}

func (m TokenManager) GenerateRefreshToken(user domain.User) (string, error) {
	return m.generate(user, "refresh", m.RefreshTTL, m.RefreshSecret)
}

func (m TokenManager) VerifyAccessToken(token string) (domain.AuthUser, error) {
	claims, err := m.verify(token, "access", m.AccessSecret)
	if err != nil {
		return domain.AuthUser{}, err
	}

	return domain.AuthUser{
		ID:    claims.ID,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}

func (m TokenManager) VerifyRefreshToken(token string) (domain.AuthUser, error) {
	claims, err := m.verify(token, "refresh", m.RefreshSecret)
	if err != nil {
		return domain.AuthUser{}, err
	}

	return domain.AuthUser{
		ID:    claims.ID,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (m TokenManager) now() time.Time {
	if m.Now != nil {
		return m.Now()
	}
	return time.Now()
}

func (m TokenManager) generate(user domain.User, tokenType string, ttl time.Duration, secret []byte) (string, error) {
	now := m.now()
	claims := Claims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("signing %s token: %w", tokenType, err)
	}
	return signed, nil
}

func (m TokenManager) verify(tokenString string, tokenType string, secret []byte) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secret, nil
	}, jwt.WithTimeFunc(m.now))
	if err != nil {
		return Claims{}, domain.ErrInvalidRefresh
	}
	if !token.Valid || claims.Type != tokenType {
		return Claims{}, domain.ErrInvalidRefresh
	}

	return claims, nil
}
