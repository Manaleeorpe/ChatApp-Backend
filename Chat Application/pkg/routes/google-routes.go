package routes

import (
	"github.com/asus/ChatApp/pkg/auth"
	"github.com/gorilla/mux"
)

var RegisterGoogleAuthstoreRoutes = func(router *mux.Router) {

	router.HandleFunc("/auth/google/login", auth.HandleGoogleLogin)
	router.HandleFunc("/auth/google/callback", auth.HandleGoogleCallback)
	router.HandleFunc("/auth/google/logout", auth.HandleGoogleLogout).Methods("GET")

}
