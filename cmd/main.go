package main

import (
	"digital-book/internal/handlers"
	"digital-book/internal/repo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const dsn = "host=localhost user=user password=1111 dbname=test port=5432 sslmode=disable TimeZone=Europe/Kiev"

func main() {
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	rep := repo.New(db)
	server := handlers.New(rep)
	server.Run()
}
