package db

import "github.com/go-pg/pg/v10"

func NewDb(url string) *pg.DB {
	opt, err := pg.ParseURL(url)
	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)
	return db
}
