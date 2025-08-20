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
	log.Println("\n=== DATABASE CONNECTION DEBUGGING ===")

	// Check all possible MySQL URL variables
	log.Println("\n1. Checking URL variables:")
	log.Printf("MYSQL_URL: '%s'", os.Getenv("MYSQL_URL"))
	log.Printf("MYSQL_PUBLIC_URL: '%s'", os.Getenv("MYSQL_PUBLIC_URL"))

	// Check Railway system variables
	log.Println("\n2. Checking Railway system variables:")
	log.Printf("RAILWAY_TCP_PROXY_DOMAIN: '%s'", os.Getenv("RAILWAY_TCP_PROXY_DOMAIN"))
	log.Printf("RAILWAY_TCP_PROXY_PORT: '%s'", os.Getenv("RAILWAY_TCP_PROXY_PORT"))
	log.Printf("RAILWAY_PRIVATE_DOMAIN: '%s'", os.Getenv("RAILWAY_PRIVATE_DOMAIN"))

	// Check all possible user/password variables
	log.Println("\n3. Checking authentication variables:")
	log.Printf("MYSQLUSER: '%s'", os.Getenv("MYSQLUSER"))
	log.Printf("MYSQL_ROOT_PASSWORD exists: %v", os.Getenv("MYSQL_ROOT_PASSWORD") != "")
	log.Printf("MYSQLPASSWORD exists: %v", os.Getenv("MYSQLPASSWORD") != "")

	// Check all possible host/port variables
	log.Println("\n4. Checking connection variables:")
	log.Printf("MYSQLHOST: '%s'", os.Getenv("MYSQLHOST"))
	log.Printf("MYSQLPORT: '%s'", os.Getenv("MYSQLPORT"))

	// Check all possible database name variables
	log.Println("\n5. Checking database name variables:")
	log.Printf("MYSQL_DATABASE: '%s'", os.Getenv("MYSQL_DATABASE"))
	log.Printf("MYSQLDATABASE: '%s'", os.Getenv("MYSQLDATABASE"))

	log.Println("\n=== STARTING CONNECTION ATTEMPT ===")

	// Get MySQL credentials from Railway
	dbUser := os.Getenv("MYSQLUSER")
	if dbUser == "" {
		log.Println("❌ No MYSQLUSER found")
	} else {
		log.Printf("✅ Found MYSQLUSER: %s", dbUser)
	}

	dbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if dbPass == "" {
		log.Println("❌ No MYSQL_ROOT_PASSWORD found, trying MYSQLPASSWORD")
		dbPass = os.Getenv("MYSQLPASSWORD")
		if dbPass == "" {
			log.Println("❌ No MYSQLPASSWORD found either")
		} else {
			log.Println("✅ Found MYSQLPASSWORD")
		}
	} else {
		log.Println("✅ Found MYSQL_ROOT_PASSWORD")
	}

	// Use Railway's TCP Proxy for connection
	dbHost := os.Getenv("RAILWAY_TCP_PROXY_DOMAIN")
	if dbHost == "" {
		log.Println("❌ No RAILWAY_TCP_PROXY_DOMAIN found, trying MYSQLHOST")
		dbHost = os.Getenv("MYSQLHOST")
		if dbHost == "" {
			log.Println("❌ No MYSQLHOST found either")
		} else {
			log.Printf("✅ Found MYSQLHOST: %s", dbHost)
		}
	} else {
		log.Printf("✅ Found RAILWAY_TCP_PROXY_DOMAIN: %s", dbHost)
	}

	dbPort := os.Getenv("RAILWAY_TCP_PROXY_PORT")
	if dbPort == "" {
		log.Println("❌ No RAILWAY_TCP_PROXY_PORT found, trying MYSQLPORT")
		dbPort = os.Getenv("MYSQLPORT")
		if dbPort == "" {
			log.Println("❌ No MYSQLPORT found either")
		} else {
			log.Printf("✅ Found MYSQLPORT: %s", dbPort)
		}
	} else {
		log.Printf("✅ Found RAILWAY_TCP_PROXY_PORT: %s", dbPort)
	}

	dbName := os.Getenv("MYSQL_DATABASE")
	if dbName == "" {
		log.Println("❌ No MYSQL_DATABASE found, trying MYSQLDATABASE")
		dbName = os.Getenv("MYSQLDATABASE")
		if dbName == "" {
			log.Println("❌ No MYSQLDATABASE found either")
		} else {
			log.Printf("✅ Found MYSQLDATABASE: %s", dbName)
		}
	} else {
		log.Printf("✅ Found MYSQL_DATABASE: %s", dbName)
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
