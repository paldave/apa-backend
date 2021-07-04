package jwt

import (
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type JWT interface {
	Generate(id string, email string, isAdmin bool) (string, error)
}

type jwt struct {
	ExpiresAt int64
	Issuer    string
	Signature string
}

type Claim struct {
	Id      string
	Email   string
	IsAdmin bool
	jwtgo.StandardClaims
}

func NewJWT(expire int, issuer string, signature string) JWT {
	return jwt{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(expire)).Unix(),
		Issuer:    issuer,
		Signature: signature,
	}
}

func (j jwt) Generate(id string, email string, isAdmin bool) (string, error) {
	claims := &Claim{
		Id:      id,
		Email:   email,
		IsAdmin: isAdmin,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: j.ExpiresAt,
			Issuer:    j.Issuer,
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.Signature))
	return signedToken, err
}
