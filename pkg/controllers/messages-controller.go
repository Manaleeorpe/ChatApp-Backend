package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asus/ChatApp/pkg/models"
	"github.com/asus/ChatApp/pkg/utils"
	"github.com/gorilla/mux"
)

var Messages models.Messages

func SendAMessage(w http.ResponseWriter, r *http.Request) {

	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session

	SendMessage := &models.Messages{}
	utils.ParseBody(r, SendMessage)

	if SendMessage.SenderID == 0 || SendMessage.RecieverID == 0 {
		http.Error(w, "SenderID and RecieverID are required", http.StatusBadRequest)
		return
	}
	_, result1 := models.GetUserById(int64(SendMessage.SenderID))

	//var Friend2 models.User
	//result2 := db.Where("ID = ?", SendMessage.Friend2UserID).First(&Friend2)
	_, result2 := models.GetUserById(int64(SendMessage.RecieverID))

	if result1.RecordNotFound() || result2.RecordNotFound() {
		http.Error(w, "Sender or Reciever does not exist", http.StatusBadRequest)
		return
	}

	ID := models.GetAFriendRequestByOFriendIDs(int64(SendMessage.SenderID), int64(SendMessage.RecieverID))

	SenderReciverFriendrequest, _ := models.GetFriendRequestById(int64(ID))

	if ID == 0 || SenderReciverFriendrequest.Status != "Accepted" {
		http.Error(w, "Sender and Reciever are not friends", http.StatusBadRequest)
		return
	}

	Message, _ := SendMessage.SendAMessageModel() //calling the SendMessage method from models package SendMessage() method is of struct type Message
	res, _ := json.Marshal(Message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetMessageHistory(w http.ResponseWriter, r *http.Request) {

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

	IntID1, err := strconv.ParseInt(ID1, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing ID1", err)
	}

	IntID2, err := strconv.ParseInt(ID2, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing ID2", err)
	}

	if IntID1 == 0 || IntID2 == 0 {
		http.Error(w, "SenderID and RecieverID are required", http.StatusBadRequest)
		return
	}
	_, result1 := models.GetUserById(int64(IntID1))

	//var Friend2 models.User
	//result2 := db.Where("ID = ?", SendMessage.Friend2UserID).First(&Friend2)
	_, result2 := models.GetUserById(int64(IntID2))

	if result1.RecordNotFound() || result2.RecordNotFound() {
		http.Error(w, "Sender or Reciever does not exist", http.StatusBadRequest)
		return
	}

	ID := models.GetAFriendRequestByOFriendIDs(int64(IntID1), int64(IntID2))

	SenderReciverFriendrequest, _ := models.GetFriendRequestById(int64(ID))

	if ID == 0 || SenderReciverFriendrequest.Status != "Accepted" {
		http.Error(w, "Sender and Reciever are not friends", http.StatusBadRequest)
		return
	}

	AllMessages, _ := models.GetMessageHistoryModel(int64(IntID1), int64(IntID2))

	if len(AllMessages) == 0 {
		http.Error(w, "There are no messages", http.StatusBadRequest)
		return
	}

	res, _ := json.Marshal(AllMessages)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func SendAMessageFunc(senderID, receiverID uint, content string) (*models.Messages, error) {

	//SendMessage := &models.Messages{}

	if senderID == 0 || receiverID == 0 {
		return nil, errors.New("sender or receiver does not exist")
	}
	_, result1 := models.GetUserById(int64(senderID))

	//var Friend2 models.User
	//result2 := db.Where("ID = ?", SendMessage.Friend2UserID).First(&Friend2)
	_, result2 := models.GetUserById(int64(receiverID))

	if result1.RecordNotFound() || result2.RecordNotFound() {
		return nil, errors.New("sender or receiver does not exist")
	}

	ID := models.GetAFriendRequestByOFriendIDs(int64(senderID), int64(receiverID))

	SenderReciverFriendrequest, _ := models.GetFriendRequestById(int64(ID))

	if ID == 0 || SenderReciverFriendrequest.Status != "Accepted" {

		return nil, errors.New("sender and Reciever are not friends")
	}

	msg := &models.Messages{
		SenderID:   senderID,
		RecieverID: receiverID,
		Content:    content,
	}

	Message, err := msg.SendAMessageModel()
	if err != nil {
		return nil, err
	}

	return Message, nil
}
