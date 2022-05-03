package database

import (
	"DIVAYTHGRAM_BACKEND/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() *gorm.DB {
	dsn := "host=ec2-34-247-72-29.eu-west-1.compute.amazonaws.com user=mtagbysfcchslr password=0e4c25be68e7f7bb1ed4b620972e38ca5c39bcce9b902e903a964c76de2344b4 dbname=d3g2koorsvm6b4 port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.UserAva{})
	db.AutoMigrate(&models.Stories{})
	db.AutoMigrate(&models.StoriesLikeDislike{})
	db.AutoMigrate(&models.StoriesEmoji{})
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		db = Init()
	}
	return db
}
