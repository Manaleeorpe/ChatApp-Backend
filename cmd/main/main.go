package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asus/ChatApp/pkg/config"
	"github.com/asus/ChatApp/pkg/models"
	"github.com/asus/ChatApp/pkg/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load(".env")

	log.Println("GOOGLE_CLIENT_ID:", os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	config.Connect()
	db := config.GetDB()
	models.SetDB(db)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Friends{})
	db.AutoMigrate(&models.Messages{})

	// Use Gorilla Mux
	router := mux.NewRouter()
	routes.RegisterUserstoreRoutes(router)
	routes.RegisterFriendsstoreRoutes(router)
	routes.RegisterWebsocketstoreRoutes(router)
	routes.RegisterMessagestoreRoutes(router)
	routes.RegisterGoogleAuthstoreRoutes(router)

	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			switch origin {
			case "https://chatapp-frontend-production-ea8c.up.railway.app",
				"http://localhost:3000":
				return true
			default:
				return false
			}
		},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
	})
	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on 0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))

	//log.Fatal(http.ListenAndServe("localhost:8080", handler)) // Use mux router here
	//log.Fetal(http.ListenAndServe("0.0.0.0:8080", handler))

}
