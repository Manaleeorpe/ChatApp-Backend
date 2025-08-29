package auth

//auth needs controller, controller needs auth
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/asus/ChatApp/pkg/config"
	"github.com/asus/ChatApp/pkg/models"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"crypto/rand"
    "encoding/base64"
    "fmt"
)

// session store
var Store = sessions.NewCookieStore([]byte("your-secret-key"))

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "https://chatapp-backend-production-5b08.up.railway.app/auth/google/callback",
	ClientID:     config.GoogleClientID,     //"272806201443-0hot4ej2vc1u8vof49gvmjgvh3m3f01d.apps.googleusercontent.com", //
	ClientSecret: config.GoogleClientSecret, //.Getenv("GOOGLE_CLIENT_SECRET"), //"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",                                      //
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

//var oauthStateString = "random" // you can generate a secure random string here
func generateRandomString() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", fmt.Errorf("failed to generate random string: %w", err)
    }
    return base64.URLEncoding.EncodeToString(b), nil
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "session-name")
    if err != nil {
        http.Error(w, "Failed to get session", http.StatusInternalServerError)
        return
    }

    state, err := generateRandomString()
    if err != nil {
        http.Error(w, "Failed to generate state", http.StatusInternalServerError)
        return
    }

    session.Values["oauth_state"] = state
    if err := session.Save(r, w); err != nil {
        http.Error(w, "Failed to save session", http.StatusInternalServerError)
        return
    }

    url := GoogleOauthConfig.AuthCodeURL(state)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func addUTCTimestamp(name string) string {
	timestamp := time.Now().UTC().Format("20060102150405") // e.g. 20250708120000
	return name + "-" + timestamp
}
func HandleGoogleLogout(w http.ResponseWriter, r *http.Request) {
	// Get session
	session, err := Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	// Clear session values
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1 // This deletes the cookie
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	// Redirect or respond
	http.Redirect(w, r, "/auth/google/login", http.StatusSeeOther) // Or return a JSON response
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Retrieve session (needed to verify the oauth state)
	session, err := Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		log.Printf("session get error: %v", err)
		return
	}

	// Verify state stored in session matches state returned by Google
	expectedState, _ := session.Values["oauth_state"].(string)
	receivedState := r.FormValue("state")
	if expectedState == "" || receivedState != expectedState {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		log.Printf("state mismatch: expected='%s' received='%s'", expectedState, receivedState)
		return
	}

	// Clear the oauth_state to prevent replay attacks
	delete(session.Values, "oauth_state")
	if err := session.Save(r, w); err != nil {
		// Log the error but continue if verification passed; you may choose to fail here instead
		log.Printf("warning: failed to clear oauth_state from session: %v", err)
	}

	// Exchange the authorization code for a token
	token, err := GoogleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		log.Printf("code exchange error: %v", err)
		return
	}

	// Create an HTTP client using the token and fetch user info
	client := GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		log.Printf("userinfo request error: %v", err)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Id            string `json:"id"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		log.Printf("userinfo decode error: %v", err)
		return
	}

	// Create or find user in DB (handle error)
	newUser := models.User{
		Name:     userInfo.Name,
		Email_id: userInfo.Email,
		GoogleID: userInfo.Id,
	}

	CreatedUser, err := models.CreateUserFromOAuth(newUser)
	if err != nil {
		http.Error(w, "Failed to create/find user", http.StatusInternalServerError)
		log.Printf("CreateUserFromOAuth error: %v", err)
		return
	}

	log.Printf("User info: %+v", CreatedUser)

	// Update session with authenticated user id
	session.Values["user_id"] = CreatedUser.ID
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		// Secure: true, // enable in production when using HTTPS
	}
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save user session", http.StatusInternalServerError)
		log.Printf("session save error: %v", err)
		return
	}

	log.Printf("New session set with user ID: %v", CreatedUser.ID)

	// Redirect to frontend (development URL shown; change to production URL as needed)
	http.Redirect(w, r, "http://localhost:3000/dashboard", http.StatusSeeOther)
}












