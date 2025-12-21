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
)

// session store
var Store = sessions.NewCookieStore([]byte("your-secret-key"))

var GoogleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback", // "http://localhost:8080/auth/google/callback",
	ClientID:     config.GoogleClientID,                        //"272806201443-0hot4ej2vc1u8vof49gvmjgvh3m3f01d.apps.googleusercontent.com", //
	ClientSecret: config.GoogleClientSecret,                    //.Getenv("GOOGLE_CLIENT_SECRET"), //"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",                                      //
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

var oauthStateString = "random" // you can generate a secure random string here

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := GoogleOauthConfig.AuthCodeURL(oauthStateString)
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

	userIDRaw, exists := session.Values["user_id"]
	if !exists {
		log.Println("user_id not in session")
	}
	userID, _ := userIDRaw.(uint)

	// Clear session values
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1 // This deletes the cookie
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	// ✅ Update last seen
	if err := models.UpdateUserLastSeen(userID); err != nil {
		log.Println("Failed to update last seen:", err)
	}

	// Redirect or respond
	http.Redirect(w, r, "/auth/google/login", http.StatusSeeOther) // Or return a JSON response
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	token, err := GoogleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Code exchange failed", http.StatusInternalServerError)
		return
	}

	client := GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
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
		return
	}

	// ✅ Here you would: find/create user in your DB
	// ✅ Then generate JWT and return it (example below)

	newUser := models.User{
		Name: addUTCTimestamp(userInfo.Name),
		//Name:     userInfo.Name,
		Email_id: userInfo.Email,
		GoogleID: userInfo.Id,
	}

	CreatedUser, _ := models.CreateUserFromOAuth(newUser)

	log.Println("User info:", CreatedUser)

	/*session
	session, _ := Store.Get(r, "session-name")
	session.Values["user_id"] = CreatedUser.ID
	session.Save(r, w)*/
	oldSession, _ := Store.Get(r, "session-name")
	oldSession.Options.MaxAge = -1 // Mark for deletion
	_ = oldSession.Save(r, w)

	// Create a fresh new session
	newSession, _ := Store.New(r, "session-name")
	newSession.Values["user_id"] = CreatedUser.ID
	newSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	_ = newSession.Save(r, w)

	log.Printf("New session set with user ID: %v", CreatedUser.ID)

	log.Printf("User ID stored in session: %v\n", newSession.Values["user_id"])
	log.Printf("User ID: %d", CreatedUser.ID)

	// Example: generate a JWT token (pseudo code)
	// tokenString := generateJWT(userInfo.Email)
	//w.Write([]byte("Login successful for user: " + CreatedUser.GoogleID))
	//http.Redirect(w, r, "http://localhost:8080/me", http.StatusSeeOther)

	//redirecting to the react app
	http.Redirect(w, r, "http://localhost:3000/dashboard", http.StatusSeeOther)
}
