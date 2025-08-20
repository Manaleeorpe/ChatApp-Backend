package config

import (
	"fmt"
	"log"
	"os"
	"strings"
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

	// First try using MYSQL_URL or SQL_URL if available
	mysqlURL := os.Getenv("MYSQL_URL")
	sqlURL := os.Getenv("SQL_URL")
	var dsn string

	if mysqlURL != "" {
		log.Printf("MYSQL_URL available, converting to DSN format")
		// Remove mysql:// if present
		mysqlURL = strings.TrimPrefix(mysqlURL, "mysql://")

		// Split into user:pass@host:port/dbname
		parts := strings.Split(mysqlURL, "@")
		if len(parts) == 2 {
			userPass := parts[0]
			hostPortDB := strings.Split(parts[1], "/")
			if len(hostPortDB) == 2 {
				dsn = fmt.Sprintf("%s@tcp(%s)/%s", userPass, hostPortDB[0], hostPortDB[1])
			}
		}
		log.Printf("Formatted DSN (hiding credentials): ...@tcp(%s)/...", strings.Split(parts[1], "/")[0])
	} else if sqlURL != "" {
		log.Printf("SQL_URL available, converting to DSN format")
		// Remove mysql:// if present
		sqlURL = strings.TrimPrefix(sqlURL, "mysql://")

		// Split into user:pass@host:port/dbname
		parts := strings.Split(sqlURL, "@")
		if len(parts) == 2 {
			userPass := parts[0]
			hostPortDB := strings.Split(parts[1], "/")
			if len(hostPortDB) == 2 {
				dsn = fmt.Sprintf("%s@tcp(%s)/%s", userPass, hostPortDB[0], hostPortDB[1])
			}
		}
		log.Printf("Formatted DSN (hiding credentials): ...@tcp(%s)/...", strings.Split(parts[1], "/")[0])
	} else {
		// Fall back to individual credentials
		dbUser := os.Getenv("MYSQLUSER")
		dbPass := os.Getenv("MYSQLPASSWORD")
		dbHost := os.Getenv("MYSQLHOST")
		dbPort := os.Getenv("MYSQLPORT")
		dbName := os.Getenv("MYSQLDATABASE")

		// Log available variables and their status
		log.Printf("Database Config (Individual Variables):")
		log.Printf("MYSQLUSER available: %v (Value: %s)", dbUser != "", dbUser)
		log.Printf("MYSQLHOST available: %v (Value: %s)", dbHost != "", dbHost)
		log.Printf("MYSQLPORT available: %v (Value: %s)", dbPort != "", dbPort)
		log.Printf("MYSQLDATABASE available: %v (Value: %s)", dbName != "", dbName)
		log.Printf("MYSQLPASSWORD available: %v", dbPass != "")

		// Check if we have all required variables
		if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
			log.Fatal("Missing required MySQL environment variables and MYSQL_URL not available")
		}

		// Create connection string
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPass, dbHost, dbPort, dbName)
	}

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
