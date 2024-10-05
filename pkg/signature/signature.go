package signature

import (
	"books-api/pkg/exception"
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Signature struct {
	jwtSecretAccessToken string
	hmacSecretAcessToken []byte
}

type Signaturer interface {
	HashBscryptPassword(password string) (string, error)
	CheckBscryptPasswordHash(password, hash string) bool
	GenerateJWT(username string) (string, error)
	JWTCheck(token string) (*JwtAuthenticationRes, *exception.Exception)
	SignHMAC512(httpMethod, bodyJson, token string) (string, error)
	VerifyHMAC512(httpMethod, bodyJson, token, hash string) (bool, *exception.Exception)
}

func NewSignature(jwtToken, hmac string) Signaturer {
	return &Signature{
		jwtSecretAccessToken: jwtToken,
		hmacSecretAcessToken: []byte(hmac),
	}
}
func normalizeJson(bodyJson string) (string, error) {
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(bodyJson))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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

func (s *Signature) SignHMAC512(httpMethod, bodyJson, token string) (string, error) {
	var normalizedBodyJson string
	var err error
	if httpMethod == http.MethodPost || httpMethod == http.MethodPut {
		normalizedBodyJson, err = normalizeJson(bodyJson)
		if err != nil {
			return "", err
		}
	}
	stringToSign := httpMethod + ":" + normalizedBodyJson + ":" + token
	// --- hmac sha512 hasher with client secret as key
	h_sha := hmac.New(crypto.SHA512.New, s.hmacSecretAcessToken)
	h_sha.Write([]byte(stringToSign))
	// ---

	signature := hex.EncodeToString(h_sha.Sum(nil))
	return signature, nil
}

func (s *Signature) VerifyHMAC512(httpMethod, bodyJson, token, hash string) (bool, *exception.Exception) {
	var normalizedBodyJson string
	var err error
	if httpMethod == http.MethodPost || httpMethod == http.MethodPut {
		normalizedBodyJson, err = normalizeJson(bodyJson)
		if err != nil {
			return false, exception.Internal("error normalizing json", err)
		}
	}
	stringToSign := httpMethod + ":" + normalizedBodyJson + ":" + token
	mac := hmac.New(sha512.New, s.hmacSecretAcessToken)
	mac.Write([]byte(stringToSign))
	sig, err := hex.DecodeString(hash)
	if err != nil {
		return false, exception.PermissionDenied("error decoding signature")
	}

	if !hmac.Equal(sig, mac.Sum(nil)) {
		return false, exception.PermissionDenied("signature unmatched")
	}

	return hmac.Equal(sig, mac.Sum(nil)), nil
}
