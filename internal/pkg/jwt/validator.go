package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type Validator interface {
	Validator(ctx context.Context, token string) (*UserClaims, error)
}

type validatorImpl struct {
	jwtSecret      string
	sessionChecker SessionChecker
}

type SessionChecker func(ctx context.Context, userID, sessionID int64) (bool, error)

func NewValidator(jwtSecret string, sessionChecker SessionChecker) Validator {
	return &validatorImpl{
		jwtSecret:      jwtSecret,
		sessionChecker: sessionChecker,
	}
}

func (v *validatorImpl) Validator(ctx context.Context, jwtToken string) (*UserClaims, error) {
	claim, err := v.getClaim(jwtToken)
	if err != nil {
		return nil, err
	}
	if claim.ExpiresAt < time.Now().Unix() {
		return claim, fmt.Errorf("token expired")
	}

	if claim.SessionID == 0 {
		return claim, fmt.Errorf("token is empty")
	}

	if valid, err := v.sessionChecker(ctx, claim.UID, claim.SessionID); err != nil {
		return claim, err
	} else if !valid {
		return claim, fmt.Errorf("session is invalid")
	}
	return claim, nil
}

func (v *validatorImpl) getClaim(jwtToken string) (*UserClaims, error) {
	claims := new(UserClaims)
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(v.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*UserClaims), nil
}
