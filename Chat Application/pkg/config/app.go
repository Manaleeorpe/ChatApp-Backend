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
	_ = godotenv.Load(".env") // Load .env file
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	log.Println("Loaded Google Client ID:", GoogleClientID)
}

// Connect to the database
func Connect() {
	fmt.Println("Attempting to connect to the database...") // Debug log

	//local
	//d, err := gorm.Open("mysql", "root:Alohomora9*@tcp(127.0.0.1:3306)/testdb?charset=utf8&parseTime=True&loc=Local")
	//root:aDnJDdqULNszDrrzltxgCqGVSMHSNoJc@mysql.railway.internal:3306/railway
	d, err := gorm.Open("mysql", "root:aDnJDdqULNszDrrzltxgCqGVSMHSNoJc@tcp(mysql.railway.internal:3306)/railway?charset=utf8&parseTime=True&loc=Local")
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
