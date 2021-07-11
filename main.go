package main

import (
	"apa-backend/auth"
	"apa-backend/config"
	"apa-backend/config/db"
	"apa-backend/router"
	"apa-backend/user"
	"log"

	"apa-backend/utils/jwt"
	"apa-backend/utils/secure"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	r := router.New()
	v1 := r.Group("/v1")

	rdb := db.NewRedis(config.RedisUrl)
	db := db.NewDb(config.DBUrl)

	sec := secure.NewSecurer()
	jwt := jwt.NewJWT(config.JWTIssuer, config.JWTTokenSignature)

	userRepo := user.NewRepository(db)
	user.RegisterHandlers(v1, user.NewService(userRepo, sec))

	authRepo := auth.NewRepository(rdb)
	auth.RegisterHandlers(v1.Group("/auth"), auth.NewService(userRepo, authRepo, sec, jwt))

	r.Logger.Fatal(r.Start(":8001"))
}
