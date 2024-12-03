package jwt

import (
	"context"
	"github.com/golang-jwt/jwt"
)

type Issuer interface {
	Issuer(ctx context.Context, userClaims *UserClaims) (string, error)
}

type issuer struct {
	jwtSecret string
}

func NewIssuer(jwtSecret string) Issuer {
	return &issuer{jwtSecret: jwtSecret}
}

func (i *issuer) Issuer(ctx context.Context, userClaims *UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	return token.SignedString([]byte(i.jwtSecret))
}
