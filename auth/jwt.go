package auth

import (
	jwtgo "github.com/dgrijalva/jwt-go"
)

type AccessClaims struct {
	Id        string
	RefreshId string
	UserId    string
	Email     string
	IsAdmin   bool
	jwtgo.StandardClaims
}

type AccessDTO struct {
	Id        string
	RefreshId string
	UserId    string
	Email     string
	IsAdmin   bool
	Expiry    int64
}

type RefreshClaims struct {
	Id       string
	AccessId string
	UserId   string
	jwtgo.StandardClaims
}

type RefreshDTO struct {
	Id       string
	AccessId string
	UserId   string
	Expiry   int64
}

type JWT interface {
	GenerateAccessToken(data *AccessDTO) (string, error)
	GenerateRefreshToken(data *RefreshDTO) (string, error)
	Validate(token string) (*jwtgo.Token, jwtgo.MapClaims, error)
}

type jwt struct {
	ExpiresAt int64
	Issuer    string
	Signature string
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

func (j *jwt) GenerateAccessToken(data *AccessDTO) (string, error) {
	claims := &AccessClaims{
		Id:        data.Id,
		RefreshId: data.RefreshId,
		UserId:    data.UserId,
		Email:     data.Email,
		IsAdmin:   data.IsAdmin,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: data.Expiry,
			Issuer:    j.Issuer,
		},
	}

	token, err := j.Sign(claims)
	return token, err
}

func (j *jwt) GenerateRefreshToken(data *RefreshDTO) (string, error) {
	claims := &RefreshClaims{
		Id:       data.Id,
		AccessId: data.AccessId,
		UserId:   data.UserId,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: data.Expiry,
			Issuer:    j.Issuer,
		},
	}

	token, err := j.Sign(claims)
	return token, err
}

func (j *jwt) Validate(token string) (*jwtgo.Token, jwtgo.MapClaims, error) {
	claims := jwtgo.MapClaims{}
	t, err := jwtgo.ParseWithClaims(token, claims, func(token *jwtgo.Token) (interface{}, error) {
		return []byte(j.Signature), nil
	})

	return t, claims, err
}
