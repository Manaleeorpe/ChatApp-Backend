package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var (
	db *gorm.DB
)
var GoogleClientID string
var GoogleClientSecret string

func init() {
	_ = godotenv.Load() // Load .env file
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	//GoogleClientID = `272806201443-0hot4ej2vc1u8vof49gvmjgvh3m3f01d.apps.googleusercontent.com`
	//GoogleClientSecret = `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`
	log.Println("Loaded Google Client ID:", GoogleClientID)
	log.Println("Loaded Google Client Secret:", GoogleClientSecret)

}

// Connect to the database
func Connect() {
	log.Println("Connecting to database...")

	_ = godotenv.Load() // works locally, ignored on Railway

	mysqlURL := os.Getenv("MYSQL_URL")
	//sqlURL := os.Getenv("SQL_URL")

	var dsn string

	if mysqlURL != "" {
		mysqlURL = strings.TrimPrefix(mysqlURL, "mysql://")
		parts := strings.Split(mysqlURL, "@")
		if len(parts) == 2 {
			userPass := parts[0]
			hostPortDB := strings.Split(parts[1], "/")
			if len(hostPortDB) == 2 {
				dsn = fmt.Sprintf(
					"%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					userPass,
					hostPortDB[0],
					hostPortDB[1],
				)
			}
		}
	} else {
		// Local / fallback
		dbUser := os.Getenv("MYSQLUSER")
		dbPass := os.Getenv("MYSQLPASSWORD")
		dbHost := os.Getenv("MYSQLHOST")
		dbPort := os.Getenv("MYSQLPORT")
		dbName := os.Getenv("MYSQLDATABASE")

		if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
			log.Fatal("Missing MySQL environment variables")
		}

		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser, dbPass, dbHost, dbPort, dbName,
		)
	}

	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("✅ Database connected")
}

func GetDB() *gorm.DB {
	return db
}

/*package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

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
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Connected to MySQL!")
}
func GetDB() *gorm.DB {
	return db
}

*/
