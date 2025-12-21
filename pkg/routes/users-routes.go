package routes

import (
	"github.com/asus/ChatApp/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterUserstoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/users/me", controllers.GetMyProfile).Methods("GET")
	router.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{ID}", controllers.GetUserById).Methods("GET")
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users/{ID}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{ID}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/email/{email}", controllers.GetUserByEmailcontroller).Methods("GET") //get user by email

	router.HandleFunc("/users/suggestedfriends/{ID}", controllers.GetSuggestedFriends).Methods("GET")
}
