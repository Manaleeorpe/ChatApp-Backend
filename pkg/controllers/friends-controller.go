package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asus/ChatApp/pkg/models"
	"github.com/asus/ChatApp/pkg/utils"
	"github.com/gorilla/mux"
)

var Friend models.Friends

func SendAFriendReuest(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session
	CreateFriend := &models.Friends{}
	CreateFriend.Status = "Pending"
	utils.ParseBody(r, CreateFriend)

	if CreateFriend.Friend1UserID == 0 || CreateFriend.Friend2UserID == 0 {
		http.Error(w, "Friend1UserID and Friend2UserID are required", http.StatusBadRequest)

		return
	}

	//get a user should be in model
	//var Friend1 models.User
	//result1 := db.Where("ID = ?", CreateFriend.Friend1UserID).First(&Friend1)
	_, result1 := models.GetUserById(int64(CreateFriend.Friend1UserID))

	//var Friend2 models.User
	//result2 := db.Where("ID = ?", CreateFriend.Friend2UserID).First(&Friend2)
	_, result2 := models.GetUserById(int64(CreateFriend.Friend2UserID))

	if result1.RecordNotFound() || result2.RecordNotFound() {
		http.Error(w, "Friend 1 or Friend 2 does not exist", http.StatusBadRequest)
		return
	}

	ID := models.GetAFriendRequestByOFriendIDs(int64(CreateFriend.Friend1UserID), int64(CreateFriend.Friend2UserID))

	if ID != 0 {
		http.Error(w, "Relationship between Friend 1 and Friend 2 is already created", http.StatusBadRequest)
		return
	}

	Friend, _ := CreateFriend.CreateFriend() //calling the CreateFriend method from models package CreateFriend() method is of struct type Friend
	res, _ := json.Marshal(Friend)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
func UpdateFriendRequest(w http.ResponseWriter, r *http.Request) {
	// Session
	LoggedInID, error := utils.SessionUserID(r)
	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("Logged in ID:", LoggedInID)

	// Parse vars
	vars := mux.Vars(r)
	ID1 := vars["Friend1ID"]
	ID2 := vars["Friend2ID"]
	FriendRequestStatus := vars["status"]

	Friend1ID, err1 := strconv.ParseInt(ID1, 0, 0)
	Friend2ID, err2 := strconv.ParseInt(ID2, 0, 0)

	if err1 != nil {
		fmt.Println("Error while parsing User ID1", err1)
	}

	if err2 != nil {
		fmt.Println("Error while parsing User ID2", err2)
	}

	ID := models.GetAFriendRequestByOFriendIDs(int64(Friend1ID), int64(Friend2ID))

	/*ID, err := strconv.ParseInt(FriendRequestID, 10, 64)
	if err != nil || ID == 0 {
		http.Error(w, "Invalid or missing Friend Request ID", http.StatusBadRequest)
		return
	}*/

	// Get request from DB
	friendRequestDetails, db := models.GetFriendRequestById(ID)
	if db.RecordNotFound() {
		http.Error(w, "Friend Request ID does not exist", http.StatusBadRequest)
		return
	}

	// Validate status
	if FriendRequestStatus != "Accepted" && FriendRequestStatus != "Rejected" {
		http.Error(w, "Invalid status. Only 'Accepted' or 'Rejected' allowed.", http.StatusBadRequest)
		return
	}

	// Update & return
	friendRequestDetails.Status = FriendRequestStatus
	db.Save(&friendRequestDetails)

	res, _ := json.Marshal(friendRequestDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetAllUsersFriends(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session
	vars := mux.Vars(r)
	userID := vars["ID"]
	status := vars["status"]

	ID, err := strconv.ParseInt(userID, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing user ID", err)
	}

	_, result1 := models.GetUserById(int64(ID))

	if result1.RecordNotFound() {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	AllFriends, _ := models.GetAllUsersFriendsModel(ID, status)

	if len(AllFriends) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]")) // return an empty JSON array
		return
	}

	res, _ := json.Marshal(AllFriends)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteFriendRequest(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session

	vars := mux.Vars(r)
	FriendRequestID := vars["ID"]

	ID, err := strconv.ParseInt(FriendRequestID, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing Friend Request ID", err)
	}

	_, result1 := models.GetFriendRequestById(int64(ID))

	if result1.RecordNotFound() {
		http.Error(w, "Friend Request ID does not exist", http.StatusBadRequest)
		return
	}

	result := models.DeleteFriendRequestModel(int64(ID))

	if result.Error != nil {
		fmt.Println("Error deleting friend request:", result.Error)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Friend Request Deleted successfully"))

}

func GetFriendRequestID(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session

	vars := mux.Vars(r)
	ID1 := vars["Friend1ID"]
	ID2 := vars["Friend2ID"]

	Friend1ID, err1 := strconv.ParseInt(ID1, 0, 0)
	Friend2ID, err2 := strconv.ParseInt(ID2, 0, 0)

	if err1 != nil {
		fmt.Println("Error while parsing User ID1", err1)
	}

	if err2 != nil {
		fmt.Println("Error while parsing User ID2", err2)
	}

	ID := models.GetAFriendRequestByOFriendIDs(int64(Friend1ID), int64(Friend2ID))

	if ID == 0 {
		http.Error(w, "Relationship between Friend 1 and Friend 2 is not created", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	idStr := strconv.FormatInt(ID, 10)

	// Write response
	w.Write([]byte(idStr))

}

func FriendRequestUserCanAccept(w http.ResponseWriter, r *http.Request) {
	//sessionLoggedInID, error := utils.SessionUserID(r)

	LoggedInID, error := utils.SessionUserID(r)
	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session
	vars := mux.Vars(r)
	userID := vars["ID"]

	ID, err := strconv.ParseInt(userID, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing user ID", err)
	}

	_, result1 := models.GetUserById(int64(ID))

	if result1.RecordNotFound() {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}

	AllFriends, _ := models.FriendRequestUserCanAcceptModel(ID)

	if len(AllFriends) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]")) // return an empty JSON array
		return
	}

	res, _ := json.Marshal(AllFriends)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
