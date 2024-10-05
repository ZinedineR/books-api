package signature

import (
	"books-api/pkg/exception"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Signature struct {
	jwtSecretAccessToken string
}

type Signaturer interface {
	HashBscryptPassword(password string) (string, error)
	CheckBscryptPasswordHash(password, hash string) bool
	GenerateJWT(username string) (string, error)
	JWTCheck(token string) (*JwtAuthenticationRes, *exception.Exception)
}

func NewSignature(jwtToken string) Signaturer {
	return &Signature{
		jwtSecretAccessToken: jwtToken,
	}
}

func (s *Signature) HashBscryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Signature) CheckBscryptPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Username string `json:"Username"`
}

type JwtAuthenticationRes struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (s *Signature) GenerateJWT(username string) (string, error) {
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "books-api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		Username: username,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	signedToken, err := token.SignedString([]byte(s.jwtSecretAccessToken))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *Signature) JWTCheck(token string) (*JwtAuthenticationRes, *exception.Exception) {
	fmt.Println("JWTAuthentication", token)
	fmt.Println("secret", s.jwtSecretAccessToken)

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecretAccessToken), nil
	})
	if err != nil {
		return nil, exception.Unauthenticated("Invalid token, " + err.Error())
	}

	var username string
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok || jwtToken.Valid {
		username = fmt.Sprintf("%v", claims["name"])
	} else {
		return nil, exception.Unauthenticated("Invalid token")
	}

	return &JwtAuthenticationRes{
		Username: username,
		Token:    token,
	}, nil
}
