package auth

import (
	"apa-backend/entity"
	"errors"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type Service interface {
	Authenticate(*loginDTO) (*entity.AuthToken, error)
}

type UserRepository interface {
	FindByEmail(string) (*entity.User, error)
}

type Securer interface {
	Hash(p string) string
	Compare(hash, p string) bool
}

type JWT interface {
	GenerateAccessToken(expiry int64, userId, email string, isAdmin bool) (string, string, error)
	GenerateRefreshToken(expiry int64, userId string) (string, string, error)
	Validate(token string) (*jwtgo.Token, jwtgo.MapClaims, error)
}

type service struct {
	ur  UserRepository
	r   Repository
	sec Securer
	jwt JWT
}

func NewService(ur UserRepository, r Repository, sec Securer, jwt JWT) *service {
	return &service{ur, r, sec, jwt}
}

func (s *service) Authenticate(req *loginDTO) (*entity.AuthToken, error) {
	var ent = &entity.AuthToken{}

	u, err := s.ur.FindByEmail(req.Email)
	if err != nil {
		return ent, err
	}

	if !s.sec.Compare(u.Password, req.Password) {
		return ent, errors.New("")
	}

	atExpiry := time.Now().Add(time.Minute * 60).Unix()
	rtExpiry := time.Now().Add(time.Hour * 24 * 7).Unix()

	/*
	 * TODO
	 * Implement proper ACL / Grouping
	 * Remove hardcoded boolean isAdmin pass
	 */
	aId, at, err := s.jwt.GenerateAccessToken(atExpiry, u.Id, u.Email, true)
	if err != nil {
		return ent, errors.New("")
	}

	if err = s.r.Create(&entity.RedisToken{
		Id:     aId,
		UserId: u.Id,
		Expiry: atExpiry,
	}); err != nil {
		return ent, errors.New("")
	}

	rId, rt, err := s.jwt.GenerateRefreshToken(rtExpiry, u.Id)
	if err != nil {
		return ent, errors.New("")
	}

	if err = s.r.Create(&entity.RedisToken{
		Id:     rId,
		UserId: u.Id,
		Expiry: rtExpiry,
	}); err != nil {
		return ent, errors.New("")
	}

	ent.AccessToken = at
	ent.RefreshToken = rt
	return ent, nil
}
