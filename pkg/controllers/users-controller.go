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

var NewUser models.User

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session end*/

	newUsers := models.GetAllUsers()
	res, _ := json.Marshal(newUsers)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	/*session
	session, _ := auth.Store.Get(r, "session-name")

	userIDRaw := session.Values["user_id"]
	userIDUint, ok := userIDRaw.(uint)
	if !ok {
		http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
		return
	}
	LoggedInID := int64(userIDUint)

	fmt.Fprintf(w, "Protected content for user ID: %d", LoggedInID)

	//session end*/
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
	ID, err := strconv.ParseInt(userID, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing user ID", err)
	}

	userDetails, _ := models.GetUserById(ID)
	res, _ := json.Marshal(userDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetUserByEmailcontroller(w http.ResponseWriter, r *http.Request) {
	// Session check
	_, err := utils.SessionUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	email := vars["email"]
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	user, db := models.GetUserByEmail(email)
	if db.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	res, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session end*/
	CreateUser := &models.User{}
	utils.ParseBody(r, CreateUser)
	user, _ := CreateUser.CreateUser()
	res, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session end*/

	UpdateUser := &models.User{}
	utils.ParseBody(r, UpdateUser)
	vars := mux.Vars(r)
	userID := vars["ID"]
	ID, err := strconv.ParseInt(userID, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing user ID", err)
	}
	userDetails, db := models.GetUserById(ID)
	if UpdateUser.Name != "" {
		userDetails.Name = UpdateUser.Name
	}
	if UpdateUser.Email_id != "" {
		userDetails.Email_id = UpdateUser.Email_id
	}
	if UpdateUser.Password != "" {
		userDetails.Password = UpdateUser.Password
	}
	db.Save(&userDetails)
	res, _ := json.Marshal(userDetails)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session end*/
	vars := mux.Vars(r)
	userID := vars["ID"]
	ID, err := strconv.ParseInt(userID, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing user ID", err)
	}
	user := models.DeleteUser(ID)
	res, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetMyProfile(w http.ResponseWriter, r *http.Request) {

	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userDetails, _ := models.GetUserById(int64(LoggedInID))
	res, _ := json.Marshal(userDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetSuggestedFriends(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//session end*/

	suggestedFriends, _ := models.GetSuggestedFriendsModel(int64(LoggedInID))
	res, _ := json.Marshal(suggestedFriends)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
