package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/asus/ChatApp/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   uint
	Conn *websocket.Conn
}

var clients = make(map[uint]*Client) //This creates a map (dictionary) to store all currently connected WebSocket clients

var upgrader = websocket.Upgrader{ //upgrades http connection to websocket
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println("Origin:", origin) // Optional: log the origin
		allowedOrigins := []string{
			"http://localhost:8080",
			"ws://localhost:8080",
			"http://localhost:3000",
			"ws://localhost:3000",
			"https://chatapp-backend-production-5b08.up.railway.app",
			"ws://chatapp-backend-production-5b08.up.railway.app",
			"wss://chatapp-backend-production-5b08.up.railway.app",
		}

		for _, o := range allowedOrigins {
			if origin == o {
				return true
			}
		}
		return false
		//return true
	}, //allows all the origins
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	//session
	LoggedInID, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println(LoggedInID)
	//session

	vars := mux.Vars(r)
	senderIDStr := vars["SenderID"]
	receiverIDStr := vars["ReceiverID"]

	senderID, err := strconv.ParseUint(senderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid SenderID", http.StatusBadRequest)
		return
	}

	receiverID, err := strconv.ParseUint(receiverIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid RecieverID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		ID:   uint(senderID),
		Conn: conn,
	}
	clients[uint(senderID)] = client

	log.Println("New client connected:", uint(senderID))
	for id, client := range clients {
		log.Printf("UserID: %d, Connected: %v\n", id, client.Conn != nil)
	}

	defer func() {
		conn.Close()
		delete(clients, uint(senderID))
		log.Println("Client disconnected:", uint(senderID))
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		//log.Printf("Message from SenderID %d to RecieverID %d: %s\n", senderID, receiverID, string(msg))

		//sending the message

		receiverClient := clients[uint(receiverID)]

		if receiverClient != nil {
			//err := receiverClient.Conn.WriteMessage(websocket.TextMessage, msg)
			err := receiverClient.Conn.WriteMessage(websocket.TextMessage, []byte(msg))

			if err != nil {
				log.Println("Could not send message to receiver:", err)
			}
		}
		//db message addition
		message, error := SendAMessageFunc(uint(senderID), uint(receiverID), string(msg))

		if error != nil {
			log.Println("Could not send message to receiver:", err)
		} else {
			log.Println(message.Sender.Name, message.Content)
		}

		// Optionally echo back to sender
		/*senderClient := clients[uint(senderID)]

		if senderClient != nil {
			err := senderClient.Conn.WriteMessage(websocket.TextMessage, []byte("Message delivered"))
			if err != nil {
				log.Println("Could not send confirmation to sender:", err)
			}
		}*/

		//send the messages to everyone excepts the sender
		/*for id, cl := range clients {
			if id != userID {
				err := cl.Conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Write error:", err)
				}
			}
		}*/

	}
}

func GetIsUserOnline(w http.ResponseWriter, r *http.Request) {
	//session
	_, error := utils.SessionUserID(r)

	if error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//session

	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	//log.Println("userIDStr", userIDStr)

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// Check if the user is online

	/*if _, exists := clients[uint(userID)]; exists {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("true"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("false"))
	}*/

	found := false
	for id := range clients {
		if id == uint(userID) {
			found = true
			break
		}
	}
	w.WriteHeader(http.StatusOK)
	if found {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}

}



