package config

import (
	"fmt"
	"log"
	"os"
	"time"

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
	log.Println("\n=== DATABASE CONNECTION DEBUGGING ===")

	// Get MySQL credentials
	dbUser := os.Getenv("MYSQLUSER")
	dbPass := os.Getenv("MYSQLPASSWORD")
	dbHost := os.Getenv("MYSQLHOST")
	dbPort := os.Getenv("MYSQLPORT")
	dbName := os.Getenv("MYSQLDATABASE")

	// Log available variables (without sensitive data)
	log.Printf("Database Config:")
	log.Printf("User: %s", dbUser)
	log.Printf("Host: %s", dbHost)
	log.Printf("Port: %s", dbPort)
	log.Printf("Database: %s", dbName)
	log.Printf("Password available: %v", dbPass != "")

	// Check if we have all required variables
	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing required MySQL environment variables")
	}

	// Create connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	log.Println("\n=== ATTEMPTING MYSQL CONNECTION ===")

	// Try to connect with retry logic
	var err error
	for i := 0; i < 3; i++ {
		db, err = gorm.Open("mysql", dsn)
		if err == nil {
			log.Printf("✅ Successfully connected to MySQL!")
			return
		}
		log.Printf("Connection attempt %d failed: %v", i+1, err)
		if i < 2 { // Don't sleep on last attempt
			time.Sleep(5 * time.Second)
		}
	}

	log.Fatal("Could not connect to MySQL:", err)
}
func GetDB() *gorm.DB {
	return db
}
