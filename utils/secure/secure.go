package secure

import "golang.org/x/crypto/bcrypt"

type Securer interface {
	Hash(p string) string
	Compare(hash string, p string) bool
}

type securer struct{}

func NewSecurer() Securer {
	return securer{}
}

func (s securer) Hash(p string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (s securer) Compare(hash string, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}
