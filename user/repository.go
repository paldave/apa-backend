package user

import (
	"apa-backend/entity"

	"github.com/go-pg/pg/v10"
)

type Repository interface {
	Create(*entity.User) error
	Exists(email string) (bool, error)
	FindByEmail(string) (*entity.User, error)
	FindById(string) (*entity.User, error)
}

type repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(user *entity.User) error {
	_, err := r.db.Model(user).Insert()
	return err
}

func (r *repository) Exists(email string) (bool, error) {
	var user entity.User
	count, err := r.db.Model(&user).
		Where("email = ?", email).
		Limit(1).
		SelectAndCount()
	return count > 0, err
}

func (r *repository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Model(&user).
		Where("email = ?", email).
		Limit(1).
		Select()
	return &user, err
}

func (r *repository) FindById(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Model(&user).
		Where("id = ?", id).
		Limit(1).
		Select()
	return &user, err
}
