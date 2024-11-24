package repositories

import (
	"fmt"
	"inventory-management/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	// env config
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// connect DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// DB migration
	err = DB.AutoMigrate(&models.Product{}, &models.Inventory{}, &models.Order{})
	if err != nil {
		log.Fatal("Database migration failed:", err)
	}

	fmt.Println("Database connected and migrated")
}
