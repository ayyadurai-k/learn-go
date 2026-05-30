package initializers

import (
	"log"
  	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(){
	var err error

	dsn := "host=localhost user=postgres password=password dbname=go_crud port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database")
	}
}