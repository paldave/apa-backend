package user

import (
	"apa-backend/entity"
)

type Service interface {
	Create(*UserDTO) (entity.User, error)
	Exists(email string) (bool, error)
	FindByEmail(email string) (entity.User, error)
}

type Securer interface {
	Hash(p string) string
}

type service struct {
	repo Repository
	sec  Securer
}

func NewService(repo Repository, sec Securer) Service {
	return service{repo, sec}
}

func (s service) Create(req *UserDTO) (entity.User, error) {
	new := entity.User{
		Id:       entity.GenerateBaseId(),
		Name:     req.Name,
		Email:    req.Email,
		Password: s.sec.Hash(req.Password),
	}

	err := s.repo.Create(new)
	if err != nil {
		return entity.User{}, err
	}

	return new, nil
}

func (s service) Exists(email string) (bool, error) {
	return s.repo.Exists(email)
}

func (s service) FindByEmail(email string) (entity.User, error) {
	return s.repo.FindByEmail(email)
}
