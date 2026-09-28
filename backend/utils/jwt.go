package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalid = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	Type     string `json:"type"` // access | refresh
	jwt.RegisteredClaims
}

// SignToken 生成 JWT
func SignToken(secret string, c *Claims, ttl int) (string, error) {
	now := time.Now()
	c.Issuer = "go-vue-admin"
	c.IssuedAt = jwt.NewNumericDate(now)
	c.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(ttl) * time.Second))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT
func ParseToken(secret, tokenStr string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, ErrTokenExpired
		}
		return claims, ErrTokenInvalid
	}
	if !t.Valid {
		return claims, ErrTokenInvalid
	}
	return claims, nil
}
