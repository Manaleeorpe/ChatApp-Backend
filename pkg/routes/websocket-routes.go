package routes

import (
	"github.com/asus/ChatApp/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterWebsocketstoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/ws/isOnline/{userID}", controllers.GetIsUserOnline).Methods("GET")
	router.HandleFunc("/ws/{SenderID}/{ReceiverID}", controllers.HandleWebSocket)

}
