package jwt

import (
	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	jwt.StandardClaims
	SessionID int64  `json:"session_id,string"`
	APPID     string `json:"appid,string"`
	UID       int64  `json:"uid,string"`

	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	PhotoURL    string `json:"photo_url"`
}

func NewUserClaim(
	userID, appID string,
	uid, sessionID, expiresAt int64,
	displayName, email, photoURL string) *UserClaims {
	return &UserClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        userID,
			ExpiresAt: expiresAt,
		},
		SessionID:   sessionID,
		APPID:       appID,
		UID:         uid,
		DisplayName: displayName,
		Email:       email,
		PhotoURL:    photoURL,
	}
}
