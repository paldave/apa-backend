package secure

import "golang.org/x/crypto/bcrypt"

type securer struct{}

func NewSecurer() *securer {
	return &securer{}
}

func (s *securer) Hash(p string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (s *securer) Compare(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}
