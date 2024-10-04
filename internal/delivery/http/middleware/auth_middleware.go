package api

import (
	"books-api/pkg/signature"
	"github.com/gin-gonic/gin"
	"strings"
)

type AuthMiddleware struct {
	Middleware
	signaturer signature.Signaturer
}

func NewAuthMiddleware(signaturer signature.Signaturer) *AuthMiddleware {
	return &AuthMiddleware{signaturer: signaturer}
}

func (m *AuthMiddleware) JWTAuthentication(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	authFields := strings.Fields(authHeader)
	if len(authFields) != 2 || strings.ToLower(authFields[0]) != "bearer" {
		m.UnauthorizedJSON(c, "Invalid token")
		return
	}
	token := authFields[1]

	res, exception := m.signaturer.JWTCheck(token)
	if exception != nil {
		m.ExceptionJSON(c, exception)
		return
	}

	c.Set("username", res.Username)
	c.Set("access_token", res.Token)

	c.Next()
}
