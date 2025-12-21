package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Initialize the db variable
func SetDB(database *gorm.DB) {
	db = database
}

type User struct {
	gorm.Model
	//ID           uint   `gorm:"primaryKey"`
	Name         string     `gorm:"" json:"name"`
	Email_id     string     `json:"email_id"`
	Phone_number string     `json:"phone_number"`
	Password     string     `json:"password"`
	GoogleID     string     `gorm:"uniqueIndex"`
	LastSeen     *time.Time `json:"last_seen"` //pointer allows null values
	// Friends     []User `gorm:"many2many:friends;"`
}

/*func (u *User) CreateUser() *User {
	db.NewRecord(u)
	db.Create(&u)
	log.Printf(u.Email_id)
	return u
}*/

func (u *User) CreateUser() (*User, error) {

	err := db.Create(&u).Error
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}
	return u, nil
}

func GetAllUsers() []User {
	var Users []User
	db.Find(&Users)
	return Users
}

func UpdateUserLastSeen(userID uint) error {
	now := time.Now()

	return db.Model(&User{}).
		Where("id = ?", userID).
		Update("last_seen", &now).
		Error
}

func GetUserLastSeen(UserID uint) (*time.Time, *gorm.DB) {
	var getUser User
	dbResult := db.Where("ID=?", UserID).Find(&getUser)
	return getUser.LastSeen, dbResult
}

func GetUserById(Id int64) (*User, *gorm.DB) {
	var getUser User
	db := db.Where("ID=?", Id).Find(&getUser)
	return &getUser, db
}

func GetUserByEmail(Email_id string) (*User, *gorm.DB) {
	var getUser User
	db := db.Where("Email_id=?", Email_id).Find(&getUser)
	return &getUser, db
}

func DeleteUser(Id int64) User {
	var user User
	db.Where("ID=?", Id).Delete(user)
	return user
}

/*
	func CreateUserFromOAuth(newUser User) (*User, error) {
		// Use the CreateUser method from your model

		user, db := GetUserByEmail(newUser.Email_id)

		if db.RecordNotFound() {
			createdUser, err := newUser.CreateUser()
			if err != nil {
				log.Println("Error creating user from OAuth:", err)
			}
			return createdUser, db.Error
		}

		return user, db.Error
	}
*/
func CreateUserFromOAuth(newUser User) (*User, error) {
	var user User
	err := db.Where("email_id = ?", newUser.Email_id).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		createdUser, err := newUser.CreateUser()
		if err != nil {
			log.Println("Error creating user from OAuth:", err)
			return nil, err
		}
		return createdUser, nil
	}

	if err != nil {
		log.Println("Database error while checking for user:", err)
		return nil, err
	}

	// User already exists
	return &user, nil
}

func GetSuggestedFriendsModel(ID int64) ([]User, error) {

	//get all the users i m friends with and compare with all the users except self
	myFriends, _ := GetAllUsersFriendsModel(ID, "Accepted") // friends struct
	myFriendsPending, _ := GetAllUsersFriendsModel(ID, "Pending")

	allFriends := append(myFriends, myFriendsPending...)

	allUsers := GetAllUsers() // User struct

	var suggestedFriends []User

	// Step 1: Make a map of friend IDs
	friendIDs := make(map[uint]bool)
	for _, f := range allFriends {
		friendIDs[f.UserID] = true
	}

	// Step 2: Filter users
	for _, user := range allUsers {
		if user.ID == uint(ID) {
			continue // skip yourself
		}
		isFriend := friendIDs[user.ID]
		if isFriend {
			continue // skip if already a friend
		}
		suggestedFriends = append(suggestedFriends, user)
	}

	return suggestedFriends, nil

}
