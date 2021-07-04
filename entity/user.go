package entity

type User struct {
	Id       string `json:"id" pg:",pk"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
