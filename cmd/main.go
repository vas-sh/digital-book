package main

import (
	"digital-book/internal/config"
	"digital-book/internal/handlers"
	"digital-book/internal/repo"
	"digital-book/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	var err error
	db, err := gorm.Open(postgres.Open(config.Config.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	rep := repo.New(db)
	srv := service.New(rep)
	server := handlers.New(rep, srv)
	server.Run()
}
