package middleware

import (
	"strings"

	apperrors "github.com/abhinavmaity/linear-lite/backend/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextKeyAuthUserID = "auth_user_id"
	ContextKeyAuthEmail  = "auth_email"
)

type AccessTokenClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func RequireAuth(secret string) gin.HandlerFunc {
	secret = strings.TrimSpace(secret)

	return func(c *gin.Context) {
		if secret == "" {
			unauthorized(c)
			return
		}

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			unauthorized(c)
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != "Bearer" {
			unauthorized(c)
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			unauthorized(c)
			return
		}

		claims := &AccessTokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
		if err != nil || !token.Valid {
			unauthorized(c)
			return
		}
		if claims.Subject == "" || claims.Email == "" || claims.ExpiresAt == nil || claims.IssuedAt == nil {
			unauthorized(c)
			return
		}

		c.Set(ContextKeyAuthUserID, claims.Subject)
		c.Set(ContextKeyAuthEmail, claims.Email)

		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	apperrors.Write(c, apperrors.Unauthorized("authentication is required"), GetRequestID(c))
}
