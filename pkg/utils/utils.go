package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asus/ChatApp/pkg/auth"
)

func ParseBody(r *http.Request, x interface{}) { // x will be the pointer to the struct which matches the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	// Unmarshal the JSON body into the provided interface
	if err := json.Unmarshal(body, x); err != nil {
		return
	}
}

func SessionUserID(r *http.Request) (uint, error) {
	//session, err := config.Store.Get(r, "session-name")
	session, err := auth.Store.Get(r, "session-name")
	if err != nil {
		return 0, err
	}

	userIDRaw, exists := session.Values["user_id"]
	if !exists {
		return 0, fmt.Errorf("user_id not in session")
	}

	userID, ok := userIDRaw.(uint)
	if !ok {
		return 0, fmt.Errorf("user_id is not a uint")
	}

	return userID, nil
}
