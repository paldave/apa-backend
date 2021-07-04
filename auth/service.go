package auth

import (
	"apa-backend/entity"
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	Authenticate(email string, passwd string) (entity.AuthToken, error)
}

type Repository interface {
	FindByEmail(email string) (entity.User, error)
}

type Securer interface {
	Hash(p string) string
	Compare(hash string, p string) bool
}

type JWT interface {
	Generate(id string, email string, isAdmin bool) (string, error)
}

type service struct {
	repo Repository
	sec  Securer
	jwt  JWT
}

func NewService(repo Repository, sec Securer, jwt JWT) Service {
	return service{repo, sec, jwt}
}

func (s service) Authenticate(email string, passwd string) (entity.AuthToken, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return entity.AuthToken{}, err
	}

	if !s.sec.Compare(u.Password, passwd) {
		return entity.AuthToken{}, errors.New("")
	}

	/*
	 * TODO
	 * Implement proper ACL / Grouping
	 * Remove hardcoded boolean isAdmin pass
	 */
	t, err := s.jwt.Generate(uuid.New().String(), u.Email, true)
	if err != nil {
		return entity.AuthToken{}, errors.New("")
	}

	return entity.AuthToken{Token: t}, nil
}
