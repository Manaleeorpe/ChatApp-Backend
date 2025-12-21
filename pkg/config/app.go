package config

import (
	"fmt"
	"log"
	"os"

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
	fmt.Println("Attempting to connect to the database...") // Debug log

	//local
	d, err := gorm.Open("mysql", "root:Alohomora9*@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local")
	//root:aDnJDdqULNszDrrzltxgCqGVSMHSNoJc@mysql.railway.internal:3306/railway

	//"root:aDnJDdqULNszDrrzltxgCqGVSMHSNoJc@tcp(mysql.railway.internal:3306)/railway?charset=utf8&parseTime=True&loc=Local"

	//d, err := gorm.Open("mysql", "root:aDnJDdqULNszDrrzltxgCqGVSMHSNoJc@tcp(mysql.railway.internal:3306)/railway?charset=utf8&parseTime=True&loc=Local")
	//dsn := "root:aDnJDdqULNszDrrzltxgCqGVSMHSNoJc@tcp(gondola.proxy.rlwy.net:33106)/railway?charset=utf8&parseTime=True&loc=Local"
	//d, err := gorm.Open("mysql", dsn)

	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}
	fmt.Println("Connected to database")
	db = d
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
