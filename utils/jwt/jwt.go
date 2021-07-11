package jwt

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type jwt struct {
	ExpiresAt int64
	Issuer    string
	Signature string
}

type AccessClaim struct {
	Id      string
	UserId  string
	Email   string
	IsAdmin bool
	jwtgo.StandardClaims
}

type RefreshClaim struct {
	Id     string
	UserId string
	jwtgo.StandardClaims
}

func NewJWT(issuer string, signature string) *jwt {
	return &jwt{
		Issuer:    issuer,
		Signature: signature,
	}
}

func (j *jwt) Sign(claims jwtgo.Claims) (string, error) {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.Signature))
	return signedToken, err
}

func (j *jwt) GenerateAccessToken(expiry int64, userId, email string, isAdmin bool) (string, string, error) {
	id := uuid.NewString()
	claims := &AccessClaim{
		Id:      id,
		UserId:  userId,
		Email:   email,
		IsAdmin: isAdmin,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: expiry,
			Issuer:    j.Issuer,
		},
	}

	token, err := j.Sign(claims)
	return id, token, err
}

func (j *jwt) GenerateRefreshToken(expiry int64, userId string) (string, string, error) {
	id := uuid.NewString()
	claims := &AccessClaim{
		Id:     id,
		UserId: userId,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: expiry,
			Issuer:    j.Issuer,
		},
	}

	token, err := j.Sign(claims)
	return id, token, err
}
