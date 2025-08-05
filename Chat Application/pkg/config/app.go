package config

import (
    "fmt"
    "log"
    "net/url"
    "os"
    "strings"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
    rawURL := os.Getenv("SQl_URL")
    if rawURL == "" {
        log.Fatal("SQl_URL not set")
    }

    u, err := url.Parse(rawURL)
    if err != nil {
        log.Fatal("Invalid SQl_URL:", err)
    }

    user := u.User.Username()
    pass, _ := u.User.Password()
    hostPort := u.Host
    dbname := strings.TrimPrefix(u.Path, "/")

    // Build DSN for MySQL driver
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        user, pass, hostPort, dbname)

    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }
    log.Println("Connected to MySQL!")
}
