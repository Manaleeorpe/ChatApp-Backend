package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var GoogleClientID string
var GoogleClientSecret string

func init() {
	_ = godotenv.Load(".env") // Load .env file
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	log.Println("Loaded Google Client ID:", GoogleClientID)
}

func Connect() {
	rawURL := os.Getenv("SQL_URL")
	if rawURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Parse URL
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal("Invalid DATABASE_URL:", err)
	}

	user := u.User.Username()
	pass, _ := u.User.Password()
	hostPort := u.Host
	dbname := strings.TrimPrefix(u.Path, "/")

	// Build DSN for jinzhu gorm
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, hostPort, dbname)

	// Open connection
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Connected to MySQL!")
}

