package main

import (
	"github.com/chema0/url-shortener/config"
	"github.com/chema0/url-shortener/pkg/db"
)

func main() {
	conf := config.LoadConfig()

	db, err := db.NewDatabase(&conf.Database)
	if err != nil {
		panic("Run!")
	}

	println(db)
}
