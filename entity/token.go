package entity

type AuthToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RedisToken struct {
	Id     string
	UserId string
	Expiry int64
}
