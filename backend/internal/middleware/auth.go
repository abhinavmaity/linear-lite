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
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			unauthorized(c)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			unauthorized(c)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
		if tokenString == "" {
			unauthorized(c)
			return
		}

		claims := &AccessTokenClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			unauthorized(c)
			return
		}
		if claims.Subject == "" {
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
