package routes

import (
	"github.com/asus/ChatApp/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterFriendsstoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/friends/friendRequestStatus/{ID}/{status}", controllers.GetAllUsersFriends).Methods("GET") //get users all friends only users json and can choose status

	router.HandleFunc("/friends/friendRequestUserCanAccept/{ID}", controllers.FriendRequestUserCanAccept).Methods("GET")

	router.HandleFunc("/friends", controllers.SendAFriendReuest).Methods("POST") //send a friend request
	//router.HandleFunc("/friends/{ID}/{status}", controllers.UpdateFriendRequest).Methods("PUT") //update a friend request accept or rejector blocka and unblock
	router.HandleFunc("/friends/{Friend1ID}/{Friend2ID}/{status}", controllers.UpdateFriendRequest).Methods("PUT") //new

	router.HandleFunc("/friends/{ID}", controllers.DeleteFriendRequest).Methods("DELETE") //delete a friend request
	//router.HandleFunc("/friends/{ID}", controllers.GetFriendRequestStatus).Methods("GET")       //get a friend request status
	router.HandleFunc("/friends/{Friend1ID}/{Friend2ID}", controllers.GetFriendRequestID).Methods("GET") //API to get id of friend request from 2 users id

	//get all the friend requests sent to the user
}
