package api

import (
	"books-api/pkg/signature"
	"github.com/gin-gonic/gin"
	"log/slog"
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

func (m *AuthMiddleware) SignatureAuthentication(c *gin.Context) {
	signatureHeader := c.GetHeader("X-Req-Signature")
	if signatureHeader == "" {
		m.UnauthorizedJSON(c, "signature empty")
		return
	}
	httpMethod, bodyJson, err := m.ParseSignatureHTTPMethod(c)
	if err != nil {
		m.BadRequestJSON(c, err.Error())
		return
	}
	token := m.GetToken(c)
	ok, exception := m.signaturer.VerifyHMAC512(httpMethod, bodyJson, token, signatureHeader)
	if exception != nil {
		m.ExceptionJSON(c, exception)
		return
	}
	if !ok {
		m.UnauthorizedJSON(c, "signature invalid")
	}

	c.Next()
}

func (m *AuthMiddleware) ErrorHandler(c *gin.Context) {

	defer func() {
		if err0 := recover(); err0 != nil {
			slog.Any("error", err0)
			m.InternalErrorJSON(c, "Request is halted unexpectedly, please contact the administrator.", err0)
		}
	}()
	c.Next()
}
