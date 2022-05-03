package database

import (
	"DIVAYTHGRAM_BACKEND/internal/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func Init() *gorm.DB {

	err := godotenv.Load()
	if err != nil {
		panic("Faild .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open("postgres://mtagbysfcchslr:0e4c25be68e7f7bb1ed4b620972e38ca5c39bcce9b902e903a964c76de2344b4@ec2-34-247-72-29.eu-west-1.compute.amazonaws.com:5432/d3g2koorsvm6b4"), &gorm.Config{})
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
