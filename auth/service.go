package auth

import (
	"apa-backend/entity"
	"errors"
	"net/http"
	"time"
)

type Service interface {
	Authenticate(*loginDTO) (*entity.AuthToken, error)
	AuthenticateRefresh(cookie string) (*entity.AuthToken, error)
	BuildCookie(name, value string) *http.Cookie
	Logout(tokenId, refreshId, userId string) error
}

type UserRepository interface {
	FindByEmail(string) (*entity.User, error)
	FindById(string) (*entity.User, error)
}

type Securer interface {
	Hash(p string) string
	Compare(hash, p string) bool
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

func (s *service) buildTokens(u *entity.User, AuthToken *entity.AuthToken) (*entity.AuthToken, error) {
	aId := entity.GenerateUuid()
	rId := entity.GenerateUuid()
	aExpiry := time.Now().Add(time.Minute * 15).Unix()
	rExpiry := time.Now().Add(time.Hour * 24 * 2).Unix()

	aDTO := &AccessDTO{
		Id:        aId,
		RefreshId: rId,
		UserId:    u.Id,
		Email:     u.Email,
		IsAdmin:   true, // TODO
		Expiry:    aExpiry,
	}

	rDTO := &RefreshDTO{
		Id:       rId,
		AccessId: aId,
		UserId:   u.Id,
		Expiry:   rExpiry,
	}

	/*
	 * TODO
	 * Implement proper ACL / Grouping
	 */
	at, err := s.jwt.GenerateAccessToken(aDTO)
	if err != nil {
		return AuthToken, errors.New("")
	}

	if err = s.r.Create(&entity.RedisToken{
		Id:     aId,
		UserId: u.Id,
		Expiry: aExpiry,
	}); err != nil {
		return AuthToken, errors.New("")
	}

	rt, err := s.jwt.GenerateRefreshToken(rDTO)
	if err != nil {
		return AuthToken, errors.New("")
	}

	if err = s.r.Create(&entity.RedisToken{
		Id:     rId,
		UserId: u.Id,
		Expiry: rExpiry,
	}); err != nil {
		return AuthToken, errors.New("")
	}

	AuthToken.AccessToken = at
	AuthToken.RefreshToken = rt
	return AuthToken, nil
}

func (s *service) Authenticate(req *loginDTO) (*entity.AuthToken, error) {
	var AuthToken = &entity.AuthToken{}

	u, err := s.ur.FindByEmail(req.Email)
	if err != nil {
		return AuthToken, err
	}

	if !s.sec.Compare(u.Password, req.Password) {
		return AuthToken, errors.New("")
	}

	return s.buildTokens(u, AuthToken)
}

func (s *service) AuthenticateRefresh(cookie string) (*entity.AuthToken, error) {
	var AuthToken = &entity.AuthToken{}

	token, claims, err := s.jwt.Validate(cookie)
	if err != nil || !token.Valid {
		return AuthToken, err
	}

	cId := claims["Id"].(string)
	cAID := claims["AccessId"].(string)
	cUID := claims["UserId"].(string)

	exists, err := s.r.Exists(cId, cUID)
	if err != nil || !exists {
		return AuthToken, err
	}

	u, err := s.ur.FindById(cUID)
	if err != nil {
		return AuthToken, err
	}

	if err := s.r.Delete(cId, cUID); err != nil {
		return AuthToken, err
	}

	if err := s.r.Delete(cAID, cUID); err != nil {
		return AuthToken, err
	}

	return s.buildTokens(u, AuthToken)
}

func (s *service) BuildCookie(name, value string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	// cookie.Expires =
	cookie.HttpOnly = true

	return cookie
}

func (s *service) Logout(tokenId, refreshId, userId string) error {
	if err := s.r.Delete(tokenId, userId); err != nil {
		return err
	}

	if err := s.r.Delete(refreshId, userId); err != nil {
		return err
	}

	return nil
}
