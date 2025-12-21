package routes

import (
	"github.com/asus/ChatApp/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterMessagestoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/messages", controllers.SendAMessage).Methods("POST")
	router.HandleFunc("/messages/{Friend1ID}/{Friend2ID}", controllers.GetMessageHistory).Methods("Get")

}
