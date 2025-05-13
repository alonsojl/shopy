package token

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const prefix = "Bearer "

var (
	ErrInvalidAuthorization = errors.New("invalid authorization")
	ErrInvalidAccessToken   = errors.New("invalid access token")
)

type JWT struct {
	key []byte
}

func NewJWT(key string) *JWT {
	return &JWT{
		key: []byte(key),
	}
}

func (j *JWT) Generate(ctx context.Context, expiresAt time.Duration) (string, error) {
	var (
		now = time.Now()
		iat = now.Unix()
		eat = now.Add(expiresAt).Unix()
	)

	payload := jwt.MapClaims{
		"foo": "bar",
		"iat": iat,
		"exp": eat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	accessToken, err := token.SignedString(j.key)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func (j *JWT) Validate(authorization string) error {
	if authorization == "" || !strings.HasPrefix(authorization, prefix) {
		return ErrInvalidAuthorization
	}

	accessToken := authorization[len(prefix):]
	token, err := jwt.Parse(accessToken, j.validateMethod)
	if err != nil {
		return fmt.Errorf("failed to parse access token: %w", err)
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return ErrInvalidAccessToken
	}

	return nil
}

func (j *JWT) validateMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return j.key, nil
}
