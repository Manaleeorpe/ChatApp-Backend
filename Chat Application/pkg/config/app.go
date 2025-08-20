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
	// Try MYSQL_URL first as it contains all connection details
	mysqlURL := os.Getenv("MYSQL_URL")
	if mysqlURL != "" {
		log.Printf("Found MYSQL_URL, attempting to connect...")
		var err error
		db, err = gorm.Open("mysql", mysqlURL)
		if err != nil {
			log.Printf("Failed to connect using MYSQL_URL: %v, falling back to individual credentials", err)
		} else {
			log.Println("Connected to MySQL using URL!")
			return
		}
	}

	// Fallback to individual variables if URL doesn't work
	dbUser := os.Getenv("MYSQL_USER") // Try with underscore
	if dbUser == "" {
		dbUser = os.Getenv("MYSQLUSER") // Try without underscore
	}

	dbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if dbPass == "" {
		dbPass = os.Getenv("MYSQLPASSWORD")
	}

	dbHost := os.Getenv("MYSQL_HOST") // Try with underscore
	if dbHost == "" {
		dbHost = os.Getenv("MYSQLHOST") // Try without underscore
	}

	dbPort := os.Getenv("MYSQL_PORT") // Try with underscore
	if dbPort == "" {
		dbPort = os.Getenv("MYSQLPORT") // Try without underscore
	}

	dbName := os.Getenv("MYSQL_DATABASE") // Try with underscore
	if dbName == "" {
		dbName = os.Getenv("MYSQLDATABASE") // Try without underscore
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
