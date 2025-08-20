package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var GoogleClientID string
var GoogleClientSecret string

func init() {
	_ = godotenv.Load(".env") // Load .env file
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	log.Println("Loaded Google Client ID:", GoogleClientID)
}

func Connect() {
	// Using exact variable names from Railway configuration
	dbUser := os.Getenv("MYSQLUSER")
	dbPass := os.Getenv("MYSQL_ROOT_PASSWORD") // Changed to use root password
	dbHost := os.Getenv("MYSQLHOST")
	dbPort := os.Getenv("MYSQLPORT")

	// Try both database name variables that exist in Railway
	dbName := os.Getenv("MYSQLDATABASE")
	if dbName == "" {
		dbName = os.Getenv("MYSQL_DATABASE")
	}

	// Debug logging
	log.Printf("Database Config - User: %v, Host: %v, Port: %v, DBName: %v",
		dbUser != "", dbHost != "", dbPort != "", dbName != "")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing variables - User:", dbUser == "",
			"Pass:", dbPass == "",
			"Host:", dbHost == "",
			"Port:", dbPort == "",
			"DB:", dbName == "")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Connected to MySQL!")
}
func GetDB() *gorm.DB {
	return db
}
