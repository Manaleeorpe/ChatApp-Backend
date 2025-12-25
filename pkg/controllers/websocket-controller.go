package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/asus/ChatApp/pkg/models"
	"github.com/asus/ChatApp/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   uint
	Conn *websocket.Conn
}

type UserStatusResponse struct {
	Online   bool    `json:"online"`
	LastSeen *string `json:"last_seen,omitempty"`
}

type IncomingWSMessage struct {
	Type     string `json:"type"`
	Content  string `json:"content,omitempty"`
	IsTyping bool   `json:"isTyping,omitempty"`
}

var clientsMu sync.RWMutex           //used for creating mutext
var clients = make(map[uint]*Client) //This creates a map (dictionary) to store all currently connected WebSocket clients

var upgrader = websocket.Upgrader{ //upgrades http connection to websocket
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println("Origin:", origin) // Optional: log the origin
		return origin == "http://localhost:8080" || origin == "ws://localhost:8080" || origin == "http://localhost:3000" || origin == "ws://localhost:3000"
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
	clientsMu.Lock()
	clients[uint(senderID)] = client
	clientsMu.Unlock()

	log.Println("New client connected:", uint(senderID))

	clientsMu.Lock()
	for id, client := range clients {
		log.Printf("UserID: %d, Connected: %v\n", id, client.Conn != nil)
	}
	clientsMu.Unlock()

	//used for deleting the user from map
	defer func() {
		models.UpdateUserLastSeen(uint(senderID))
		conn.Close()
		clientsMu.Lock()
		delete(clients, uint(senderID))
		clientsMu.Unlock()

		log.Println("Client disconnected:", uint(senderID))
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// 👇 Parse incoming JSON
		var incoming IncomingWSMessage
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Println("Invalid JSON from client:", err)
			continue
		}

		// 👇 Handle by message type
		switch incoming.Type {

		case "message":
			// Send to receiver if online
			clientsMu.RLock()
			receiverClient := clients[uint(receiverID)]
			clientsMu.RUnlock()

			if receiverClient != nil {
				err := receiverClient.Conn.WriteJSON(map[string]interface{}{
					"type":    "message",
					"content": incoming.Content,
				})
				if err != nil {
					log.Println("Write error:", err)
				}
			}

			// Save to DB
			_, err := SendAMessageFunc(
				uint(senderID),
				uint(receiverID),
				incoming.Content,
			)
			if err != nil {
				log.Println("DB error:", err)
			}

		case "typing":
			// Forward typing event ONLY
			clientsMu.RLock()
			receiverClient := clients[uint(receiverID)]
			clientsMu.RUnlock()

			if receiverClient != nil {
				receiverClient.Conn.WriteJSON(map[string]interface{}{
					"type":     "typing",
					"isTyping": incoming.IsTyping,
				})
			}
			//log.Println("typing received from front end", incoming.Type, incoming.IsTyping)

		default:
			log.Println("Unknown WS message type:", incoming.Type)
		}
	}
}

func GetIsUserOnline(w http.ResponseWriter, r *http.Request) {
	// ✅ Session check
	_, err := utils.SessionUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ✅ Get userID from URL
	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userID", http.StatusBadRequest)
		return
	}

	// ✅ Check online (in-memory)
	clientsMu.RLock()
	_, online := clients[uint(userID)]
	clientsMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// ✅ Online case
	if online {
		json.NewEncoder(w).Encode(UserStatusResponse{
			Online: true,
		})
		return
	}

	// ✅ Offline case → fetch last seen using YOUR function
	lastSeen, dbResult := models.GetUserLastSeen(uint(userID))

	var lastSeenStr *string
	if dbResult.Error == nil && lastSeen != nil {
		formatted := lastSeen.Format(time.RFC3339)
		lastSeenStr = &formatted
	}

	json.NewEncoder(w).Encode(UserStatusResponse{
		Online:   false,
		LastSeen: lastSeenStr,
	})
}
